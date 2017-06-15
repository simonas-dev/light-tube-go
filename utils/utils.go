package utils

import "math"

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