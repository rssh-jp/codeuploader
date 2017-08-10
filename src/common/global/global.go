package global
import(
    "fmt"
    "flag"
    "log"
    "os"
    "sync"

    "common/config"
)
type global struct{
    appPath string
    config config.Config

}
func (g *global)setAppPath(appPath string){
    g.appPath = appPath
    fmt.Println(g.appPath)
}
func (g *global) AppPath()string{
    return g.appPath
}

func (g *global)setConfig(conf config.Config){
    g.config = conf
    fmt.Println(g.config)
}
func (g *global) Config()config.Config{
    return g.config
}

var (
    gInstance *global = &global{}
    once sync.Once
)

// --------------------------------------------------
// 外部から参照できる部分
// --------------------------------------------------
func GetInstance() *global{
    return gInstance
}
func Initialize(){
    once.Do(func(){
        g := GetInstance()
        p, _ := os.Getwd()
        g.setAppPath(p)

        var err error
        var c string
        var conf config.Config
        flag.StringVar(&c, "c", "", "config.json path")
        flag.Parse()
        if c == ""{
            conf, err = config.GetConfig("config.json")
            if err != nil{
                log.Fatal(err)
            }
        }else{
            conf, err = config.GetConfig(c)
            if err != nil{
                log.Fatal(err)
            }
        }
        g.setConfig(conf)
    })
}
