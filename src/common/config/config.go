package config

import(
    "encoding/json"
    "io"
    "strings"

    "common"
)

// コンフィグ
type Config struct{
    CronDir string
    PublicDir string
    Ex string
}
func getConfigByString(str string)(ret Config){
    dec := json.NewDecoder(strings.NewReader(str))
    for{
        if err := dec.Decode(&ret); err == io.EOF{
            break
        } else if err != nil{
            return
        }
    }
    return
}
func GetConfig(path string)(ret Config, err error){
    str, err := common.GetFileString(path)
    if err != nil{
        return
    }
    ret = getConfigByString(str)
    return
}
