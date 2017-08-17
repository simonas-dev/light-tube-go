package main

import (
    "fmt"
    "github.com/simonassank/go_ws2811"
    "math"
    "time"
    "./config"
    "./utils"
    "./audio"
)

var (
    ratio           float64
    led_index       int
    tinted          uint32
    config_data     config.Config
    led_colors      = make([]uint32, led_count)
    buff            = make([]float64, int(1024))
    eq              = make([]int, led_count)
    led_count       = 144
    correct_gamma   = false
)

func main() {
    fmt.Println("Go!")

    ws2811.Init(18, led_count, 255)

    c, p := audio.NewAudio()

    go func() {
        for {
            samples, err := c.Read(buff)
            if (samples == 0) {
                break
            }
            if (err != nil) {
                 break
            }

            samples, err = p.Write(buff)

            for (err != nil) {
                fmt.Println(err)
                samples, err = p.Write(buff)
            }
            if (err != nil) {
                 break
            }
        }
    }()

    go func() {
        var (
            channel_led_count       = led_count/2
            led_divider             = float64(channel_led_count)/40.0
            channel_1_start         = 0
            channel_1_end           = channel_led_count
            channel_2_start         = channel_led_count
            channel_2_end           = led_count-1
        )

        for {
            energies := audio.GetMelEnergies(buff)
            pitch_val := audio.GetPitchVal(buff)

            color := utils.GetFloatColor(config_data.Note_Colors, utils.GetNoteIndex(pitch_val))

            // Channel 1
            for led_index := channel_1_start; led_index < channel_1_end; led_index++ {
                ratio = utils.GetEnergy(energies, float64(led_index)/led_divider)
                ratio = math.Pow(ratio, config_data.Pre_power)
                ratio *= config_data.Multi
                ratio = math.Pow(ratio, config_data.Post_power)
                if ratio > config_data.Max_level {
                    ratio = config_data.Max_level
                } else if (ratio < config_data.Min_level) {
                    ratio = 0
                }

                inv_led_index := (channel_led_count)-led_index
                eq[inv_led_index] = int(ratio * 100)

                tinted = utils.AvgColor(int(color), int(config_data.Tint_color), config_data.Tint_alpha)
                led_colors[inv_led_index] = utils.AvgColor(int(led_colors[inv_led_index]), int(tinted), ratio)
                led_colors[inv_led_index] = utils.AvgColor(int(led_colors[inv_led_index]), 0, config_data.Fade_ratio)
                ws2811.SetLed(inv_led_index, led_colors[inv_led_index], correct_gamma)
            }

            // Channel 2
            for led_index := channel_2_start; led_index < channel_2_end; led_index++ {
                ratio = utils.GetEnergy(energies, float64(led_index-channel_2_start)/led_divider)
                ratio = math.Pow(ratio, config_data.Pre_power)
                ratio *= config_data.Multi
                ratio = math.Pow(ratio, config_data.Post_power)
                if ratio > config_data.Max_level {
                    ratio = config_data.Max_level
                } else if (ratio < config_data.Min_level) {
                    ratio = 0
                }
                eq[led_index] = int(ratio * 100)

                tinted = utils.AvgColor(int(color), int(config_data.Tint_color), config_data.Tint_alpha)
                led_colors[led_index] = utils.AvgColor(int(led_colors[led_index]), int(tinted), ratio)
                led_colors[led_index] = utils.AvgColor(int(led_colors[led_index]), 0, config_data.Fade_ratio)
                ws2811.SetLed(led_index, led_colors[led_index], correct_gamma)
            }

            ws2811.Render()
        }
    }()

    go func() {
        var (
            new_config config.Config
            err error
        )
        for {
            new_config, err = config.Load()
            if (err == nil) {
                config_data = new_config
            }
            time.Sleep(1000)
        }
    }()

    for {
        time.Sleep(100)
    }
}
