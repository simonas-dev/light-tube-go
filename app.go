package main

import (
	"github.com/cocoonlife/goalsa"
	"fmt"
	"flag"
	"github.com/simonassank/aubio-go"
	"github.com/simonassank/go_ws2811"
	"math"
)

var (
	srcPath    = flag.String("src", "", "Path to source file. Required")
	Samplerate = flag.Int("samplerate", 44100, "Sample rate to use for the audio file")
	Blocksize  = flag.Int("blocksize", 256, "Blocksize use for the audio file")
	Bufsize    = flag.Int("bufsize", 512, "Bufsize use for the audio file")
	Silence    = flag.Float64("silence", -90.0, "Threshold to use when detecting silence")
	Threshold  = flag.Float64("threshold", 0.0, "Detection threshold")
	Verbose    = flag.Bool("verbose", false, "Print verbose output")
	help       = flag.Bool("help", false, "Print this help")
)


func main() {
	fmt.Println("Go!")

	ws2811.Init(18, 144, 255)

	c, _ := alsa.NewCaptureDevice(
		"plughw:CARD=Set,DEV=0",
		 2,
		 alsa.FormatFloat64LE,
		 *Samplerate,
		 alsa.BufferParams{},
	)

	p, _ := alsa.NewPlaybackDevice(
		"plughw:CARD=Set,DEV=0",
		 2,
		 alsa.FormatFloat64LE,
		 *Samplerate,
		 alsa.BufferParams{},
	)
	
	pitch := aubio.NewPitch(
		aubio.PitchDefault,
		uint(*Bufsize),
		uint(*Blocksize),
		uint(*Samplerate),
	)
	pitch.SetUnit(aubio.PitchOutFreq)
	pitch.SetTolerance(0.7)

	phVoc, _ := aubio.NewPhaseVoc(uint(*Bufsize), uint(*Blocksize))
	fb := aubio.NewFilterBank(40, uint(*Bufsize))
	fb.SetMelCoeffsSlaney(uint(*Samplerate))

	b4 := make([]float64, 512)

	calc_outs := 0
	audio_outs := 0

	_ = phVoc
	_ = calc_outs
	
	go func() {
		var (
			ratio float64
			energies []float64
			fftgrain *aubio.ComplexBuffer
			inputBuffer *aubio.SimpleBuffer
			led_index int
			led_colors = make([]uint32, 144)
		)
		
		for {
			inputBuffer = aubio.NewSimpleBufferData(512, b4)
			pitch.Do(inputBuffer)
			_ = pitch.Buffer().Slice()[0]

			phVoc.Do(inputBuffer)
			fftgrain = phVoc.Grain()
			fb.Do(fftgrain)
			energies = fb.Buffer().Slice()

			for led_index = 0; led_index < 144; led_index++ {
				ratio = energies[int(Round(float64(led_index)/3.6))]
				// ratio = math.Pow(ratio, 1.5)
				ratio = ratio * 30.0
				// ratio = math.Pow(ratio, 1.5)
				// fmt.Println(ratio)
				if ratio > 1 {
					ratio = 1
				}

				led_colors[led_index] = AvgColor(int(led_colors[led_index]), int(GetColor(ratio)))
				led_colors[led_index] = AvgColor(int(led_colors[led_index]), 0)
			}
			ws2811.SetBitmap(led_colors)
			ws2811.Render()
		}
	}()

	for {
		_, _ = c.Read(b4)
		_, _ = p.Write(b4)
		audio_outs += 1
	}
}

func GetColor(ratio float64) uint32 {
	num := uint32(255*ratio)
	return (num << 16)+(num << 8)+num
}

func Round(f float64) float64 {
    return math.Floor(f + .5)
}

func AvgColor(a int, b int) uint32{
	return uint32(float64(GetRed(a)+GetRed(b)) / 2.0 + float64(GetGreen(a)+GetGreen(b)) / 2.0 + float64(GetBlue(a)+GetBlue(b)) / 2)
}

func GetRed(color int) int {
	return color & 0xFF0000 >> 16
}

func GetBlue(color int) int {
	return color & 0xFF00 >> 8
}

func GetGreen(color int) int {
	return color & 0xFF
}

