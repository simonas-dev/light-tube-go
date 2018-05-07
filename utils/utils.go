package utils

import (
	"math"
)

var (
	A4 = 440.0
	C0 = A4 * math.Pow(2, -4.75)
)

func GetColor(ratio float64) uint32 {
	num := uint32(255 * ratio)
	return (num << 16) + (num << 8) + num
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func GetEnergy(arr []float64, f_index float64) float64 {
	index := int(f_index)
	ratio := f_index - float64(index)
	return arr[index]*(1-ratio) + arr[index+1]*ratio
}

func GetFloatColor(arr []uint32, f_index float64) uint32 {
	index := int(f_index)
	if len(arr)-1 == index {
		return arr[index]
	}
	ratio := f_index - float64(index)
	return AvgColor(int(arr[index]), int(arr[index+1]), ratio)
}

func AvgColor(a int, b int, ratio float64) uint32 {
	inv_ratio := 1 - ratio
	c_a := GetColorNum(a)
	c_b := GetColorNum(b)
	return CreateColor(int((c_a[0]*inv_ratio)+(c_b[0]*ratio)), int((c_a[1]*inv_ratio)+(c_b[1]*ratio)), int((c_a[2]*inv_ratio)+(c_b[2]*ratio)))
}

func FadeColor(color int, fade_ratio float64) uint32 {
	arr := GetColorNum(color)
	res := make([]int, 3)
	for index, color := range arr {
		ratio := math.Pow(color/255, 1.0)
		res[index] = int(color - (color*fade_ratio)*ratio)
	}
	return CreateColor(res[0], res[1], res[2])
}

func AddColor(a int, b int, ratio float64) uint32 {
	c_a := GetColorNum(a)
	c_b := GetColorNum(b)
	return CreateColor(int(c_a[0]+c_b[0]*ratio), int(c_a[1]+c_b[1]*ratio), int(c_a[2]+c_b[2]*ratio))
}

func MinusColor(a int, b int) uint32 {
	c_a := GetColorNum(a)
	c_b := GetColorNum(b)
	red := int(c_a[0] - c_b[0])
	if red < 0 {
		red = 0
	}
	blue := int(c_a[1] - c_b[1])
	if blue < 0 {
		blue = 0
	}
	green := int(c_a[2] - c_b[2])
	if green < 0 {
		green = 0
	}
	return CreateColor(red, green, blue)
}

func FadeColorChannels(a int, ratio float64) uint32 {
	color := GetColorNum(a)
	red := (1 - ratio*math.Pow(color[0]/255, 5)) * color[0]
	green := (1 - ratio*math.Pow(color[1]/255, 5)) * color[1]
	blue := (1 - ratio*math.Pow(color[2]/255, 5)) * color[2]

	return uint32(int(red)<<16 + int(green)<<8 + int(blue))
}

func CreateColor(r int, g int, b int) uint32 {
	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if b > 255 {
		b = 255
	}
	return uint32(r<<16 + g<<8 + b)
}

func GetColorNum(color int) []float64 {
	return []float64{float64(color & 0xFF0000 >> 16), float64(color & 0xFF00 >> 8), float64(color & 0xFF)}
}

func GetNoteIndex(freqHz float64) float64 {
	div := freqHz / C0
	if div == 0 {
		return 0
	}
	h := 12 * math.Log2(div)
	if h < 0 {
		return 0
	}
	rem := h - float64(int(h))
	n := int(h) % 12.0
	return float64(n) + rem
}

func GetNoteIndex2(freqHz float64) float64 {
	if freqHz < 200 && freqHz > 42000 {
		return 0
	}
	if freqHz > 800 {
		return 11
	}
	return ((freqHz + 200) / 1000) * 11
}
