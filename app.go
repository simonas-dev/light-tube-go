package main

import (
    "fmt"
    "math"
    "time"
    "./config"
    "./utils"
    "./audio"
    "./leds"
    // "os"
    // "bufio"
    // "strconv"
)


type TimelineItem struct {
    BufferId uint64
    RatioArr []float64
    Color    uint32
}

var (
    ratio           float64
    led_index       int
    tinted          uint32
    config_data     config.Config
    led_colors      = make([]uint32, led_count)
    avg_pitch       float64
    buff            = make([]float64, *audio.Bufsize)
    eq              = make([]int, led_count)
    led_count       = 144 * 4
    correct_gamma   = false
    timelineArr     []TimelineItem
)

var (
    channel_led_count       = led_count/2
    led_divider             = float64(channel_led_count)/40.0
    channel_1_start         = 0
    channel_1_end           = channel_led_count
)

func main() {
    defer leds.Close()

    fmt.Println("Go!")

    leds.Init()

    config_data, _ = config.Load()

    bufferChannel := make(chan []float64, int(*audio.Bufsize))
    playbackChannel := make(chan uint64)
    go audio.AudioPassThrough(bufferChannel, playbackChannel)

    go func() {
        for {
            if (!audio.IsReady) {
                time.Sleep(1 * time.Second)
                continue
            }
            bufferId := uint64(audio.BufferId)
            energies := audio.GetMelEnergies(buff)
            pitch_val := audio.GetPitchVal(buff)
            if (pitch_val < 9500) {
                ratio := 0.5
                avg_pitch = avg_pitch * ratio + pitch_val * (1-ratio)
            }
            note_index := utils.GetNoteIndex(avg_pitch)
            color := utils.GetFloatColor(config_data.Note_Colors, note_index)

            min, _ := MinMax(energies)
            for index, value := range energies {
                 energies[index] = value - min
            }

            ratioArr := make([]float64, channel_1_end)
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
                ratioArr[led_index] = ratio
            }
            
            timelineArr = append(timelineArr, TimelineItem {
                bufferId,
                ratioArr,
                color,
            })

            // fmt.Println("Add To Timeline")
        }
    }()

    go func() {
        _ = <-playbackChannel
        fmt.Println("Start Render")

        // f, _ := os.Create("./timeline.csv")

        lastRenderIndex := -1

        for {
            bufferId := uint64(audio.WriterId)
            index := TimelineIndexOf(timelineArr, bufferId)
            
            if (lastRenderIndex == index) {
                continue
            }

            timelineItem := timelineArr[index]
            
            // logLine := ""

            for led_index := channel_1_start; led_index < channel_1_end; led_index++ {
                ratio := ReturnTimelineEnergy(timelineArr, index, led_index, 0)
                // if (led_index == 20) {
                //     ratioStr := strconv.FormatFloat(ratio, 'E', 5, 64)
                //     logLine += ratioStr
                // }
                led_colors[led_index] = utils.AddColor(int(led_colors[led_index]), int(timelineItem.Color), ratio)               
                led_colors[led_index] = utils.AvgColor(int(led_colors[led_index]), int(timelineItem.Color), config_data.Tint_alpha * ratio)
                led_colors[led_index] = utils.FadeColor(int(led_colors[led_index]), config_data.Fade_ratio)
                leds.SetMirror(led_index, led_count, led_colors[led_index])
            }

            
            // w := bufio.NewWriter(f)
            // w.WriteString(logLine + "\n")
            // w.Flush()
            // fmt.Println("Save Log File")

            leds.Render()
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
            time.Sleep(1 * time.Second)
        }
    }()

    for {
        buff = <-bufferChannel
    }
}

func ReturnTimelineEnergy(array []TimelineItem, centerIndex int, energyIndex int, avgWeight int) float64 {
    sum := float64(0)
    couter := 0
    for index := centerIndex - avgWeight; index <= centerIndex + avgWeight; index++ {
        if (index > 0 && index < len(array)) {
            sum += array[index].RatioArr[energyIndex]
            couter += 1
        }
    }
    return sum / float64(couter)
}

func TimelineIndexOf(array []TimelineItem, bufferId uint64) int {
    for index := 0; index < len(array); index++ {
        if (array[index].BufferId >= bufferId) {
            return index
        }
    }
    return -1
}

func MinMax(array []float64) (float64, float64) {
    var max float64 = array[0]
    var min float64 = array[0]
    for _, value := range array {
        if max < value {
            max = value
        }
        if min > value {
            min = value
        }
    }
    return min, max
}
