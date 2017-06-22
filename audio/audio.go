package audio

import (
    "github.com/cocoonlife/goalsa"
    "github.com/simonassank/aubio-go"
    "fmt"
    "flag"
)

var (
    SrcPath     = flag.String("src", "", "Path to source file. Required")
    Samplerate  = flag.Int("samplerate", 44100, "Sample rate to use for the audio file")
    Blocksize   = flag.Int("blocksize", 512, "Blocksize use for the audio file")
    Bufsize     = flag.Int("bufsize", 1024, "Bufsize use for the audio file")
)

var (
    pitch       *aubio.Pitch
    phVoc       *aubio.PhaseVoc
    fb          *aubio.FilterBank
    inBuff      *aubio.SimpleBuffer
)

func NewAudio() (c *alsa.CaptureDevice, p *alsa.PlaybackDevice) {
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


    pitch = aubio.NewPitch(
        aubio.PitchDefault,
        uint(*Bufsize),
        uint(*Blocksize),
        uint(*Samplerate),
    )
    pitch.SetUnit(aubio.PitchOutFreq)
    pitch.SetTolerance(0.9)

    phVoc, _ = aubio.NewPhaseVoc(uint(*Bufsize), uint(*Blocksize))
    fb = aubio.NewFilterBank(40, uint(*Bufsize))
    fb.SetMelCoeffsSlaney(uint(*Samplerate))

    return c, p
}

func GetAnalaysis(buffer []float64) ([]float64, float64) {
    inBuff = aubio.NewSimpleBufferData(uint(*Bufsize), buffer)
    pitch.Do(inBuff)
    phVoc.Do(inBuff)
    fb.Do(phVoc.Grain())
    return fb.Buffer().Slice(), pitch.Buffer().Slice()[0]
}