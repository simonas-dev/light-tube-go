package main

import (
    "github.com/cocoonlife/goalsa"
    "fmt"
    "flag"
    "time"
)

var (
    srcPath    = flag.String("src", "", "Path to source file. Required")
    Samplerate = flag.Int("samplerate", 44100, "Sample rate to use for the audio file")
    Blocksize  = flag.Int("blocksize", 1024, "Blocksize use for the audio file")
    Bufsize    = flag.Int("bufsize", 1024, "Bufsize use for the audio file")
    Silence    = flag.Float64("silence", -90.0, "Threshold to use when detecting silence")
    Threshold  = flag.Float64("threshold", 0.0, "Detection threshold")
    Verbose    = flag.Bool("verbose", false, "Print verbose output")
    help       = flag.Bool("help", false, "Print this help")
)

func main() {
    fmt.Println("Go!")

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
    
    buff := make([]float64, uint(*Bufsize))
    
    for {
        start := time.Now()    
        c.Read(buff)
        p.Write(buff)

        elapsed := time.Since(start)
        fmt.Println("Audio")
        fmt.Println(elapsed)
    }
}
