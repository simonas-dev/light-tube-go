package main

import (
	"log"
	"github.com/kaey/pulse"
	"fmt"
	"flag"
	"github.com/simonassank/aubio-go"
	"github.com/simonassank/go_ws2811"
	"math"
	"time"
)

// Configs
var (
	srcPath    = flag.String("src", "", "Path to source file. Required")
	Samplerate = flag.Int("samplerate", 44100, "Sample rate to use for the audio file")
	Blocksize  = flag.Int("blocksize", 512, "Blocksize use for the audio file")
	Bufsize    = flag.Int("bufsize", 1024, "Bufsize use for the audio file")
	Silence    = flag.Float64("silence", -90.0, "Threshold to use when detecting silence")
	Threshold  = flag.Float64("threshold", 0.0, "Detection threshold")
	Verbose    = flag.Bool("verbose", false, "Print verbose output")
	help       = flag.Bool("help", false, "Print this help")
	LedCount   = 144
)

// Final vars
var (
	A4 = 440.0
	C0 = A4 * math.Pow(2, -4.75)
	colors = []uint32 {
		uint32(0xFF0000),
		uint32(0xFF0000),
		uint32(0xFFFF00),
		uint32(0xFFFF00),
		uint32(0xC3F2FF),
		uint32(0x7F8BFD),
		uint32(0x7F8BFD),
		uint32(0xF37900),
		uint32(0xF37900),
		uint32(0x33CC33),
		uint32(0x33CC33),
		uint32(0x8EC9FF),
	}
)

// Non final vars
var (
ratio float64
energies []float64
fftgrain *aubio.ComplexBuffer
inputBuffer *aubio.SimpleBuffer
led_index int
led_colors = make([]uint32, LedCount)
buff = make([]float64, int(*Bufsize))
)

func main() {
	fmt.Println("Go!")

	ws2811.Init(18, LedCount, 255)

	s := pulse.Sample{
		Format:   pulse.FormatS32le,
		Rate:     44100,
		Channels: 2,
	}

	r, err := pulse.NewReader(&s, "echo", "bluez_source.4A:79:DF:CA:CB:6F")
	if err != nil {
		log.Fatalln("Blue: %s", err)
	}
	defer r.Close()
	defer r.Drain()

	rl, err := r.Latency()
	if err != nil {
		log.Fatalln(err)
	}

	w, err := pulse.NewWriter(&s, "echo", "alsa_output.0.analog-stereo")
	if err != nil {
		log.Fatalln("Out: %s", err)
	}
	defer w.Close()
	defer w.Drain()

	wl, err := w.Latency()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Reader latency", rl)
	log.Println("Writer latency", wl)

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

	_ = phVoc


	_ = ratio
	_ = energies
	_ = fftgrain
	_ = inputBuffer
	_ = led_index
	_ = led_colors
	_ = buff

	buf := make([]byte, 2048)
	go func() {
		for {

			err := r.Read(buf)
			if err != nil {
				log.Fatalln(err)
			}

			err = w.Write(buf)
			if err != nil {
				log.Fatalln(err)
			}
		}	
	}()

	for {
		start := time.Now()
		inputBuffer = aubio.NewSimpleBufferData(uint(*Bufsize), buff)
		pitch.Do(inputBuffer)
		pitch_val := pitch.Buffer().Slice()[0]
		color := colors[GetNoteIndex(pitch_val)]

		fmt.Println(color)

		phVoc.Do(inputBuffer)
		fftgrain = phVoc.Grain()
		fb.Do(fftgrain)
		energies = fb.Buffer().Slice()

		for led_index = 0; led_index < LedCount; led_index++ {
			ratio = energies[int(Round(float64(led_index)/3.6))]
			ratio = math.Pow(ratio, 1)
			// ratio *= 15
			// ratio = math.Pow(ratio, 4)
			if ratio > 0.8 {
				ratio = 0.8
			}
			
			tinted := AvgColor(int(color), int(0x00ffff), 0.5)
			led_colors[led_index] = AvgColor(int(led_colors[led_index]), int(tinted), ratio)
			led_colors[led_index] = AvgColor(int(led_colors[led_index]), 0, 0.05)
			ws2811.SetLed(led_index, led_colors[led_index])
		}
		
		// fmt.Println(led_colors)
		ws2811.Render()
		fmt.Println(time.Since(start))
	}
}

func GetColor(ratio float64) uint32 {
	num := uint32(255*ratio)
	return (num << 16)+(num << 8)+num
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func AvgColor(a int, b int, ratio float64) uint32{
	inv_ratio := 1-ratio
	c_a := GetColorNum(a)
	c_b := GetColorNum(b)
	return uint32(int((c_a[0]*inv_ratio)+(c_b[0]*ratio)) << 16 + int((c_a[1]*inv_ratio)+(c_b[1]*ratio)) << 8 + int((c_a[2]*inv_ratio)+(c_b[2]*ratio)))
}

func GetColorNum(color int) []float64 {
	return []float64{float64(color & 0xFF0000 >> 16), float64(color & 0xFF00 >> 8), float64(color & 0xFF)}
}

func GetNoteIndex(freqHz float64) int {
	div := freqHz/C0
	if div == 0 { return 0 }
	h := int(Round(12 * math.Log2(div)))
	n := h % 12
	return int(n)
}
