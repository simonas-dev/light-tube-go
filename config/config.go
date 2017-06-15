package config

import (
    "io/ioutil"
    "encoding/json"
)

var (
    CONFIG_FILE_NAME = "./configs.txt"
)


type Config struct {
    Max_level float32 `json:"max_level"`
    Min_level float32 `json:"min_level"`
    Pre_power float32 `json:"pre_power"`
    Multi float32 `json:"multi"`
    Post_power float32 `json:"post_power"`
    Tint_alpha float32 `json:"tint_alpha"`
    Fade_ratio float32 `json:"fade_ratio"`
    Tint_color string `json:"tint_color"`
}

func Load() Config {
    file, _ := ioutil.ReadFile(CONFIG_FILE_NAME)
    var jsontype Config
    _ = json.Unmarshal([]byte(file), &jsontype)
    return jsontype
}