package main

import (
	"fmt"
	"time"

	"../config"
	"../internal/anims"
	"../internal/audio"
	"../internal/leds"
)

var (
	configData config.Config
	avgPitch   float64
	buff       = make([]float64, *audio.Bufsize)
	ledCount   = 144 * 4
	ledColors  = make([]uint32, ledCount)
)

func main() {
	defer leds.Close()
	fmt.Println("Go!")
	leds.Init()
	configData, _ = config.Load()

	channel := make(chan []float64, int(*audio.Bufsize))
	go audio.AudioPassThrough(channel)

	go func() {
		for {
			buff = <-channel
		}
	}()

	go func() {
		for {
			if !audio.IsReady {
				time.Sleep(1 * time.Second)
				continue
			}
			energies := audio.GetMelEnergies(buff)
			pitchVal := audio.GetPitchVal(buff)

			anims.ReduceAubioAnim(ledColors, energies, pitchVal, avgPitch, configData)
			anims.ReduceWithFade(ledColors, configData)
			//anims.ReduceWithRipple(ledColors, configData)

			leds.SetArray(ledColors)
			leds.Render()
		}
	}()

	go func() {
		var (
			newConfig config.Config
			err       error
		)
		for {
			newConfig, err = config.Load()
			if err == nil {
				configData = newConfig
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		time.Sleep(1 * time.Second)
	}
}
