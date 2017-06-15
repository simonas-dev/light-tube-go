package main

import (
    _ "fmt"
    "github.com/simonassank/go_ws2811"
)

func main() {
    ws2811.Init(18, 144, 100)
    
    for {
        led_colors := make([]uint32, 144)

        for led_index := 0; led_index < 144; led_index++ {
            led_colors[led_index] = 0x0000ff
            ws2811.SetBitmap(led_colors)
            ws2811.Render()
        }
    }
}