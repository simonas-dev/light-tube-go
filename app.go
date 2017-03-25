package main

import (
	"github.com/cocoonlife/goalsa"
	"fmt"
	"flag"
	"github.com/simonassank/aubio-go"
	"github.com/simonassank/go_ws2811"
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

	ws2811.Init(18, 144, 255)
	ws2811.SetLed(0, uint32(234542))
	ws2811.Render()

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

	b4 := make([]float64, 256)

	for {
		c.Read(b4)
		p.Write(b4)

		go func() {
			// sBuffer := aubio.NewSimpleBufferData(256, b4)
			// pitch.Do(sBuffer)
			// pitch_val := pitch.Buffer().Slice()[0]
			// sBuffer.Free()
			// fmt.Println(pitch_val)
			inputBuffer := aubio.NewSimpleBufferData(256, b4)
			phVoc.Do(inputBuffer)
			fftgrain := phVoc.Grain()
			fb.Do(fftgrain)
			energies := fb.Buffer().Slice()
			
			fmt.Println(energies)
		}()
	}
}



