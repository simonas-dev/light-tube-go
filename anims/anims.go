package anims

import (
    "fmt"
    _"math"
    "../utils"
    "../leds"
    "math/rand"
)

var (
    is_running = false
    led_cnt = 144 * 4
//  colors = []uint32{0xff0000, 0x0000ff, 0x00ff00 }
    colors = []uint32{0xFF0000, 0xFF7F00, 0xFFFF00, 0x7FFF00, 0x00FF00, 0x00FF7F, 0x00FFFF, 0x007FFF, 0x0000FF, 0x7F00FF, 0xFF00FF, 0xFF007F}
)

func RunLava() {
    if (!is_running) {
        is_running = true
        led_colors := make([]uint32, led_cnt)
        for (is_running) {
            color := colors[rand.Intn(3)]
            PushToFront(color, led_colors)
            leds.Set(led_colors)
            leds.Render()
            fmt.Println("cyc done")
        }
    }
}

func StopLava() {
    if (is_running) {
        is_running = false
        fmt.Println("Alive")
    }
}

func PushToFront(item uint32, array []uint32) {
    for i := len(array)-2; i >= 0; i-- {
        array[i+1] = array[i]
    }
    array[0] = item
    fmt.Println(array)
}

func wheel(pos int) uint32 {
    if (pos < 85) {
        return utils.CreateColor(pos * 3, 255 - pos * 3, 0)
    } else if (pos < 170) {
        pos -= 85
        return utils.CreateColor(255 - pos * 3, 0, pos * 3)
    } else {
        pos -= 170
        return utils.CreateColor(0, 3 * pos, 255 - pos * 3)
    }
}
