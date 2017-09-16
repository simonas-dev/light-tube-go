package leds

import (
    "github.com/simonassank/go_ws2811"
)

func SetMirror(index int, total int, color uint32) {
    pivot := total / 2
    ws2811.SetLed(pivot-index, color, false)
    ws2811.SetLed(pivot+index, color, false)
}

func TurnOn(total int) {
    // Initial Center Clustering
    SetMirror(0, total, 0xff0000)
    ws2811.Render()
}
