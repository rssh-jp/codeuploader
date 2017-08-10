package protocol
import(
    "common/global"
)
func GetTemplateFiles(list []string)(ret []string){
    g := global.GetInstance()
    conf := g.Config()

    ret = make([]string, len(list) + 3)
    index := 0
    for index=0;index<len(list);index++{
        ret[index] = list[index]
    }
    ret[index + 0] = conf.PublicDir + "public/views/header.html"
    ret[index + 1] = conf.PublicDir + "public/views/template.html"
    ret[index + 2] = conf.PublicDir + "public/views/common.css"
    
    return
}
