package main

import (
	"github.com/cocoonlife/goalsa"
	"fmt"
	"flag"
	"github.com/simonassank/aubio-go"
	"github.com/simonassank/go_ws2811"
	"math"
	"time"
	"io/ioutil"
	"encoding/json"
	// ui "github.com/gizak/termui"
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
	ConfigsFileName = "./configs.txt"
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
	eq = make([]int, LedCount)
	tinted uint32
	inv_led_index int
)

type SettingsObject struct {
	max_level float32 `json:"max_level"`
	min_level float32 `json:"min_level"`
	pre_power float32 `json:"pre_power"`
	multi float32 `json:"multi"`
	post_power float32 `json:"post_power"`
	tint_alpha float32 `json:"tint_alpha"`
	fade_ratio float32 `json:"fade_ratio"`
	tint_color string `json:"tint_color"`
}

var (
	max_level = 0.95
	min_level = 0.0
	pre_power = 1.0 // 1.15
	multi = 1.125
	post_power = 1.0
	tint_alpha = 0.125 // 0.5
	fade_ratio = 0.05
	tint_color = int(0x00ff77)
)

func main() {
	fmt.Println("Go!")

	ws2811.Init(18, LedCount, 255)

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

	go func() {
		for {
			c.Read(buff)
			p.Write(buff)
		}	
	}()

	go func() {
		for {
			// start := time.Now()
			inputBuffer = aubio.NewSimpleBufferData(uint(*Bufsize), buff)
			pitch.Do(inputBuffer)
			pitch_val := pitch.Buffer().Slice()[0]
			color := GetFloatColor(colors, GetNoteIndex(pitch_val))

			phVoc.Do(inputBuffer)
			fftgrain = phVoc.Grain()
			fb.Do(fftgrain)
			energies = fb.Buffer().Slice()

			channel_led_count := LedCount/2
			led_divider := float64(channel_led_count)/40.0
			channel_1_start := 0
			channel_1_end := channel_led_count
			channel_2_start := channel_led_count
			channel_2_end := LedCount-1

			// Channel 1
			for led_index = channel_1_start; led_index < channel_1_end; led_index++ {
				ratio = GetEnergy(energies, float64(led_index)/led_divider)
				ratio = math.Pow(ratio, pre_power)
				ratio *= multi
				ratio = math.Pow(ratio, post_power)
				if ratio > max_level {
					ratio = max_level
				} else if (ratio < min_level) {
					ratio = 0
				}

				inv_led_index := (channel_led_count)-led_index
				eq[inv_led_index] = int(ratio * 100)

				tinted = AvgColor(int(color), tint_color, tint_alpha)
				led_colors[inv_led_index] = AvgColor(int(led_colors[inv_led_index]), int(tinted), ratio)
				led_colors[inv_led_index] = AvgColor(int(led_colors[inv_led_index]), 0, fade_ratio)
				ws2811.SetLed(inv_led_index, led_colors[inv_led_index])
			}

			// Channel 2
			for led_index = channel_2_start; led_index < channel_2_end; led_index++ {
				ratio = GetEnergy(energies, float64(led_index-channel_2_start)/led_divider)
				ratio = math.Pow(ratio, pre_power)
				ratio *= multi
				ratio = math.Pow(ratio, post_power)
				if ratio > max_level {
					ratio = max_level
				} else if (ratio < min_level) {
					ratio = 0
				}
				eq[led_index] = int(ratio * 100)

				tinted = AvgColor(int(color), tint_color, tint_alpha)
				led_colors[led_index] = AvgColor(int(led_colors[led_index]), int(tinted), ratio)
				led_colors[led_index] = AvgColor(int(led_colors[led_index]), 0, fade_ratio)
				ws2811.SetLed(led_index, led_colors[led_index])
			}

			// fmt.Println(led_colors)
			ws2811.Render()
			// _ = start
			// fmt.Println(time.Since(start))
		}
	}()

	go func() {
		for {
			file, e := ioutil.ReadFile(ConfigsFileName)
			if e != nil {
			    fmt.Printf("File error: %v\n", e)
			}
			var jsontype SettingsObject
			err := json.Unmarshal(file, &jsontype)
			fmt.Printf("File: %s, Data: %s Err: %s\n", file, jsontype.max_level, err)
			time.Sleep(100)
		}
	}()
	
	for {
		time.Sleep(100)
	}

	// ui.Init()
	// go func() {
	// 	for {
	// 		spl0 := ui.NewSparkline()
	// 		spl0.Data = eq
	// 		spl0.Title = "Eq"
	// 		spl0.Height = 10
	// 		spl0.LineColor = ui.ColorMagenta

	// 		// single
	// 		spls0 := ui.NewSparklines(spl0)
	// 		spls0.Height = 30
	// 		spls0.Width = 144
	// 		spls0.Border = true

	// 		ui.Render(spls0)
	// 	}
	// }()

	// ui.Handle("/sys/kbd/C-c", func(ui.Event) {
	// 	ui.StopLoop()
	// })

	// ui.Loop()
}

func GetColor(ratio float64) uint32 {
	num := uint32(255*ratio)
	return (num << 16)+(num << 8)+num
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func GetEnergy(arr []float64, f_index float64) float64 {
	index := int(f_index)
	ratio := f_index - float64(index)
	return arr[index] * (1-ratio) + arr[index+1] * ratio
}

func GetFloatColor(arr []uint32, f_index float64) uint32 {
	index := int(f_index)
	if len(arr)-1 == index {
		return arr[index]
	}
	ratio := f_index - float64(index)
	return AvgColor(int(arr[index]), int(arr[index+1]), ratio)
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

func GetNoteIndex(freqHz float64) float64 {
	div := freqHz/C0
	if div == 0 { return 0 }
	h := 12 * math.Log2(div)
	rem := h - float64(int(h))
	n := int(h) % 12.0
	return float64(n)+rem
}
