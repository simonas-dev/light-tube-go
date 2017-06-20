package config

import (
    "io/ioutil"
    "encoding/json"
    "strconv"
    "errors"
)

var (
    CONFIG_FILE_NAME = "./config.json"
)


type ConfigJSON struct {
    Max_level float64 `json:"max_level"`
    Min_level float64 `json:"min_level"`
    Pre_power float64 `json:"pre_power"`
    Multi float64 `json:"multi"`
    Post_power float64 `json:"post_power"`
    Tint_alpha float64 `json:"tint_alpha"`
    Fade_ratio float64 `json:"fade_ratio"`
    Tint_color string `json:"tint_color"`
    Note_Colors []string `json:"note_colors"`
}

type Config struct {
    Max_level float64
    Min_level float64
    Pre_power float64
    Multi float64
    Post_power float64
    Tint_alpha float64
    Fade_ratio float64
    Tint_color uint32
    Note_Colors []uint32
}

func Load() (Config, error) {
    file, _ := ioutil.ReadFile(CONFIG_FILE_NAME)
    var jsonData ConfigJSON
    var config Config
    json.Unmarshal([]byte(file), &jsonData)
    
    if (len(jsonData.Note_Colors) > 0) {
        config = Config {
            Max_level: jsonData.Max_level,
            Min_level: jsonData.Min_level,
            Pre_power: jsonData.Pre_power,
            Multi: jsonData.Multi,
            Post_power: jsonData.Post_power,
            Tint_alpha: jsonData.Tint_alpha,
            Fade_ratio: jsonData.Fade_ratio,
            Tint_color: hexStrToInt(jsonData.Tint_color),
            Note_Colors: hexStrArrToIntArr(jsonData.Note_Colors)}
        return config, nil
    } else {
        return config, errors.New("I/O fuckup")
    }
}

func hexStrToInt(hex string) uint32 {
    num, _ := strconv.ParseUint(hex, 16, 32)
    return uint32(num)
}

func hexStrArrToIntArr(hexArr []string) []uint32 {
    numArr := make([]uint32, 12)
    for i := 0; i < 12; i++ {
        numArr[i] = hexStrToInt(hexArr[i])
    }
    return numArr
}


