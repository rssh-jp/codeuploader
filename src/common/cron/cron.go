package cron
import(
    "fmt"
    "io/ioutil"
    "strings"
    //"strconv"
    "time"

    "common/global"
    "github.com/robfig/cron"
)
type CronData struct{
    Path string
    Name string
}
func GetCronFileList()(ret []CronData){
    g := global.GetInstance()
    conf := g.Config()
    files, err := ioutil.ReadDir(conf.CronDir)
    if err != nil{
        return
    }
    for _, val := range files{
        name := val.Name()
        data := CronData{conf.CronDir + name, name}
        ret = append(ret, data)
    }
    return
}
func GetCronScheduleTime(cronStr string)(ret []time.Time){
    splitList := strings.Split(cronStr, " ")
    if len(splitList) < 6{
        return
    }
    timeStr := strings.Join(splitList[:5], " ")
    schedule ,err := cron.Parse("0 " + timeStr)
    if err != nil{
        fmt.Println("error : ", err)
        return
    }
    now := time.Now()
    t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(-1 * time.Minute)
    for;;{
        wk := schedule.Next(t)
        if wk.Year() != now.Year() || wk.Month() != now.Month() || wk.Day() != now.Day(){
            break
        }
        ret = append(ret, wk)
        t = wk
    }

    return
}
