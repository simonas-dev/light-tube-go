package anims

import (
	"math"

	"../../config"
	"../../utils"
)

// ReduceAubioAnim TODO
func ReduceAubioAnim(ledColors []uint32, energies []float64, pitchVal float64, avgPitch float64, configData config.Config) {
	var (
		ledCount        = len(ledColors)
		channelLedCount = ledCount / 2
		ledDivider      = float64(channelLedCount) / 40.0
		channelStart    = 0
		channelEnd      = channelLedCount
	)

	if pitchVal < 9500 {
		ratio := configData.NoteRatio
		avgPitch = avgPitch*ratio + pitchVal*(1-ratio)
	}
	noteIndex := utils.GetNoteIndex(avgPitch)
	color := utils.GetFloatColor(configData.NoteColors, noteIndex)

	var ratio float64

	pivot := ledCount / 2

	for ledIndex := channelStart; ledIndex < channelEnd; ledIndex++ {
		ratio = utils.GetEnergy(energies, float64(ledIndex)/ledDivider)
		ratio = math.Pow(ratio, configData.PrePower)
		ratio = math.Pow(ratio, configData.PostPower)
		ratio *= configData.Multi
		if ratio > configData.MaxLevel {
			ratio = configData.MaxLevel
		} else if ratio < configData.MinLevel {
			ratio = 0
		}

		leftIndex := pivot - ledIndex
		ledColors[leftIndex] = calculateMixColor(leftIndex, ledColors, color, ratio)
		ledColors[leftIndex] = calculateMixColor(leftIndex, ledColors, color, configData.TintAlpha*ratio)

		rightIndex := pivot + ledIndex
		ledColors[rightIndex] = calculateMixColor(rightIndex, ledColors, color, ratio)
		ledColors[rightIndex] = calculateMixColor(rightIndex, ledColors, color, configData.TintAlpha*ratio)
	}
}

func calculateMixColor(index int, src []uint32, color uint32, ratio float64) uint32 {
	return utils.AddColor(int(src[index]), int(color), ratio)
}

// ReduceWithFade TODO
func ReduceWithFade(src []uint32, configData config.Config) {
	for index, item := range src {
		src[index] = utils.FadeColor(int(item), configData.FadeRatio)
	}
}
