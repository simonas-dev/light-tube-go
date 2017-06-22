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

// Non final vars
var (
    ratio           float64
    led_index       int
    tinted          uint32
    inv_led_index   int
    configData      config.Config
    led_colors      = make([]uint32, LedCount)
    buff            = make([]float64, int(1024))
    eq              = make([]int, LedCount)
    LedCount        = 144
)

func main() {
    fmt.Println("Go!")

    ws2811.Init(18, LedCount, 255)

    c, p := audio.NewAudio()

    go func() {
        for {
            c.Read(buff)
            p.Write(buff)
        }
    }()

    go func() {
        var (
            channel_led_count = LedCount/2
            led_divider             = float64(channel_led_count)/40.0
            channel_1_start         = 0
            channel_1_end           = channel_led_count
            channel_2_start         = channel_led_count
            channel_2_end           = LedCount-1
            energies                = make([]float64, int(40))
            pitch_val               float64
            color                   uint32
            inv_led_index           int
        )

        for {
            energies, pitch_val = audio.GetAnalaysis(buff)

            color = utils.GetFloatColor(configData.Note_Colors, utils.GetNoteIndex(pitch_val))

            fmt.Println(energies)

            // Channel 1
            for led_index = channel_1_start; led_index < channel_1_end; led_index++ {
                ratio = utils.GetEnergy(energies, float64(led_index)/led_divider)
                ratio = math.Pow(ratio, configData.Pre_power)
                ratio *= configData.Multi
                ratio = math.Pow(ratio, configData.Post_power)
                if ratio > configData.Max_level {
                    ratio = configData.Max_level
                } else if (ratio < configData.Min_level) {
                    ratio = 0
                }

                inv_led_index = (channel_led_count)-led_index
                eq[inv_led_index] = int(ratio * 100)

                tinted = utils.AvgColor(int(color), int(configData.Tint_color), configData.Tint_alpha)
                led_colors[inv_led_index] = utils.AvgColor(int(led_colors[inv_led_index]), int(tinted), ratio)
                led_colors[inv_led_index] = utils.AvgColor(int(led_colors[inv_led_index]), 0, configData.Fade_ratio)
                ws2811.SetLed(inv_led_index, led_colors[inv_led_index])
            }

            // Channel 2
            for led_index = channel_2_start; led_index < channel_2_end; led_index++ {
                ratio = utils.GetEnergy(energies, float64(led_index-channel_2_start)/led_divider)
                ratio = math.Pow(ratio, configData.Pre_power)
                ratio *= configData.Multi
                ratio = math.Pow(ratio, configData.Post_power)
                if ratio > configData.Max_level {
                    ratio = configData.Max_level
                } else if (ratio < configData.Min_level) {
                    ratio = 0
                }
                eq[led_index] = int(ratio * 100)

                tinted = utils.AvgColor(int(color), int(configData.Tint_color), configData.Tint_alpha)
                led_colors[led_index] = utils.AvgColor(int(led_colors[led_index]), int(tinted), ratio)
                led_colors[led_index] = utils.AvgColor(int(led_colors[led_index]), 0, configData.Fade_ratio)
                ws2811.SetLed(led_index, led_colors[led_index])
            }

            ws2811.Render()
        }
    }()

    go func() {
        var (
            newConfig config.Config
            err error
        )
        for {
            newConfig, err = config.Load()
            if (err == nil) {
                configData = newConfig
            }
            time.Sleep(1000)
        }
    }()
    
    for {
        time.Sleep(100)
    }
}
