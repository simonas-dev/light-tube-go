package main

import (
	Alsa "github.com/cocoonlife/goalsa"	
)
import "fmt"
import "flag"
import "github.com/simonassank/aubio-go"
import "C"

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
	// Alsa
	c, err := Alsa.NewCaptureDevice("plughw:CARD=Set,DEV=0",
		 2, Alsa.FormatFloat64LE, 44100, Alsa.BufferParams{})
	p, err := Alsa.NewPlaybackDevice("plughw:CARD=Set,DEV=0",
		 2, Alsa.FormatFloat64LE, 44100, Alsa.BufferParams{})
	
	fmt.Println(c)
	fmt.Println(err)
	
	pitch := aubio.NewPitch(aubio.PitchDefault, uint(*Bufsize), uint(*Blocksize), uint(*Samplerate))
	pitch.SetUnit(aubio.PitchOutFreq)
	pitch.SetTolerance(0.7)

	b4 := make([]float64, 256)

	for {
		c.Read(b4)
		p.Write(b4)

		go func() {
			buffer := aubio.NewSimpleBufferData(256, b4)
			pitch.Do(buffer)
			pitch_val := pitch.Buffer().Slice()[0]
			buffer.Free()
			fmt.Println(pitch_val)
		}()
	}
}


