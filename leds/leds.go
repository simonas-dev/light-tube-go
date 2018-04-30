package leds

import (
    "../pkg/go_ws2811"
    "fmt"
)

var (
    c *ws2811.WS281x
)

func Init() {
    config := ws2811.DefaultConfig
    config.Brightness = 255
    config.Invert = false

    controller, err := ws2811.New(144*4, &config); if (err != nil) {
    fmt.Println(err)
    }

    err = controller.Init(); if (err != nil) {
    fmt.Println(err)
    }
    fmt.Println("Done")

    c = controller

    c.SetLed(100, uint32(0xff0000))
    c.Render()
}

func SetMirror(index int, total int, color uint32) {
    pivot := total / 2
    c.SetLed(pivot-index, color)
    c.SetLed(pivot+index, color)
}

func Render() {
    c.Render()
}


func Close() {
    c.Close()
}
