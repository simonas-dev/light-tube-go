package audio

import (
	"github.com/cocoonlife/goalsa"
	"fmt"
	"flag"
)

var (
	srcPath     = flag.String("src", "", "Path to source file. Required")
	samplerate  = flag.Int("samplerate", 44100, "Sample rate to use for the audio file")
	Blocksize   = flag.Int("blocksize", 512, "Blocksize use for the audio file")
	Bufsize     = flag.Int("bufsize", 1024, "Bufsize use for the audio file")
	Silence     = flag.Float64("silence", -90.0, "Threshold to use when detecting silence")
	Threshold   = flag.Float64("threshold", 0.0, "Detection threshold")
	Verbose     = flag.Bool("verbose", false, "Print verbose output")
	help        = flag.Bool("help", false, "Print this help")
	LedCount    = 144
)

func NewAudio() (c *CaptureDevice, p *PlaybackDevice) {
	c, errC := alsa.NewCaptureDevice(
		"plughw:CARD=Device,DEV=0",
		2,
		alsa.FormatFloat64LE,
		*Samplerate,
		alsa.BufferParams{
			*Samplerate,
			1,
			1,
		 },
	)

	fmt.Println(errC)

	p, errP := alsa.NewPlaybackDevice(
		"plughw:CARD=Device,DEV=0",
		2,
		alsa.FormatFloat64LE,
		*Samplerate,
		alsa.BufferParams{
			*Samplerate,
			1,
			1,
		},
	)

	fmt.Println(errP)

	return c, p
}

func GetEnergies(buffer []float64) []float64 {
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

	inputBuffer = aubio.NewSimpleBufferData(uint(*Bufsize), buffer)
	pitch.Do(inputBuffer)
	pitch_val := pitch.Buffer().Slice()[0]
	color := utils.GetFloatColor(configData.Note_Colors, utils.GetNoteIndex(pitch_val))

	phVoc.Do(inputBuffer)
	fftgrain = phVoc.Grain()
	fb.Do(fftgrain)
	return fb.Buffer().Slice()
}