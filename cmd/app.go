package main

import (
	"fmt"
	"math"
	"time"

	"../config"
	"../internal/audio"
	"../internal/leds"
	"../utils"
)

var (
	ratio      float64
	tinted     uint32
	configData config.Config
	ledColors  = make([]uint32, ledCount)
	avgPitch   float64
	buff       = make([]float64, *audio.Bufsize)
	eq         = make([]int, ledCount)
	ledCount   = 144 * 4
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
		var (
			channelLedCount = ledCount / 2
			ledDivider      = float64(channelLedCount) / 40.0
			channelStart    = 0
			channelEnd      = channelLedCount
		)

		for {
			if !audio.IsReady {
				time.Sleep(1 * time.Second)
				continue
			}
			energies := audio.GetMelEnergies(buff)
			pitchVal := audio.GetPitchVal(buff)

			if pitchVal < 9500 {
				ratio := configData.Note_ratio
				avgPitch = avgPitch*ratio + pitchVal*(1-ratio)
			}
			noteIndex := utils.GetNoteIndex(avgPitch)
			color := utils.GetFloatColor(configData.NoteColors, noteIndex)

			// Channel 1
			for ledIndex := channelStart; ledIndex < channelEnd; ledIndex++ {
				ratio = utils.GetEnergy(energies, float64(ledIndex)/ledDivider)
				ratio = math.Pow(ratio, configData.Pre_power)
				ratio *= configData.Multi
				ratio = math.Pow(ratio, configData.Post_power)
				if ratio > configData.Max_level {
					ratio = configData.Max_level
				} else if ratio < configData.Min_level {
					ratio = 0
				}

				eq[ledIndex] = int(ratio * 100)
				ledColors[ledIndex] = utils.AddColor(int(ledColors[ledIndex]), int(color), ratio)
				ledColors[ledIndex] = utils.AvgColor(int(ledColors[ledIndex]), int(color), configData.Tint_alpha*ratio)
				ledColors[ledIndex] = utils.FadeColor(int(ledColors[ledIndex]), configData.Fade_ratio)
				leds.SetMirror(ledIndex, ledCount, ledColors[ledIndex])
			}

			leds.Render()
		}
	}()

	go func() {
		var (
			configData config.Config
			err        error
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
