package config

import (
    "io/ioutil"
    "encoding/json"
)

var (
    CONFIG_FILE_NAME = "./configs.txt"
)


type ConfigJSON struct {
    Max_level float32 `json:"max_level"`
    Min_level float32 `json:"min_level"`
    Pre_power float32 `json:"pre_power"`
    Multi float32 `json:"multi"`
    Post_power float32 `json:"post_power"`
    Tint_alpha float32 `json:"tint_alpha"`
    Fade_ratio float32 `json:"fade_ratio"`
    Tint_color string `json:"tint_color"`
    Note_Colors []string `json:"note_colors"`
}

type Config struct {
    Max_level float32
    Min_level float32
    Pre_power float32
    Multi float32
    Post_power float32
    Tint_alpha float32
    Fade_ratio float32
    Tint_color uint32
    Note_Colors []uint32
}

func Load() Config {
    file, _ := ioutil.ReadFile(CONFIG_FILE_NAME)
    var jsonData ConfigJSON
    _ = json.Unmarshal([]byte(file), &jsonData)
    config = Config {
        Max_level = jsonData.Max_level,
        Min_level = jsonData.Min_level,
        Pre_power = jsonData.Pre_power,
        Multi = jsonData.Multi,
        Post_power = jsonData.Post_power,
        Tint_alpha = jsonData.Tint_alpha,
        Fade_ratio = jsonData.Fade_ratio,
        Tint_color = hexStrToInt(jsonData.Tint_color)
        Note_Colors = hexStrArrToIntArr(jsonData.Note_Colors)
    }
    return config
}

func hexStrToInt(hex string) uint32 {
    num, _ := strconv.ParseInt(hex, 16, 32)
    return num
}

func hexStrArrToIntArr(hexArr []string) uint32[] {
    var numArr [len(hexArr)]uint32
    for i := 0; i < len(hexArr); i++ {
        numArr[i] = hexStrToInt(hexArr[i])
    }
    return numArr
}


