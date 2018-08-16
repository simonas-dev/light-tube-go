package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
)

var (
	CONFIG_FILE_NAME = "./config.json"
)

type ConfigJSON struct {
	NoteRatio           float64  `json:"note_ratio"`
	MinEnergySubraction bool     `json:"min_energy_subraction"`
	MaxLevel            float64  `json:"max_level"`
	MinLevel            float64  `json:"min_level"`
	PrePower            float64  `json:"pre_power"`
	Multi               float64  `json:"multi"`
	PostPower           float64  `json:"post_power"`
	TintAlpha           float64  `json:"tint_alpha"`
	FadeRatio           float64  `json:"fade_ratio"`
	TintColor           string   `json:"tint_color"`
	NoteColors          []string `json:"note_colors"`
}

type Config struct {
	NoteRatio           float64
	MinEnergySubraction bool
	MaxLevel            float64
	MinLevel            float64
	PrePower            float64
	Multi               float64
	PostPower           float64
	TintAlpha           float64
	FadeRatio           float64
	TintColor           uint32
	NoteColors          []uint32
}

func Load() (Config, error) {
	file, _ := ioutil.ReadFile(CONFIG_FILE_NAME)
	var jsonData ConfigJSON
	var config Config
	json.Unmarshal([]byte(file), &jsonData)

	if len(jsonData.NoteColors) > 0 {
		config = Config{
			NoteRatio:           jsonData.NoteRatio,
			MinEnergySubraction: jsonData.MinEnergySubraction,
			MaxLevel:            jsonData.MaxLevel,
			MinLevel:            jsonData.MinLevel,
			PrePower:            jsonData.PrePower,
			Multi:               jsonData.Multi,
			PostPower:           jsonData.PostPower,
			TintAlpha:           jsonData.TintAlpha,
			FadeRatio:           jsonData.FadeRatio,
			TintColor:           hexStrToInt(jsonData.TintColor),
			NoteColors:          hexStrArrToIntArr(jsonData.NoteColors, len(jsonData.NoteColors))}
		return config, nil
	} else {
		return config, errors.New("I/O fuckup")
	}
}

func hexStrToInt(hex string) uint32 {
	num, _ := strconv.ParseUint(hex, 16, 32)
	return uint32(num)
}

func hexStrArrToIntArr(hexArr []string, size int) []uint32 {
	numArr := make([]uint32, size)
	for i := 0; i < size; i++ {
		numArr[i] = hexStrToInt(hexArr[i])
	}
	return numArr
}
