package notfound
import(
    "bytes"
    "fmt"
    "net/http"
    "html/template"
    "common"
    "common/global"
    "network/protocol"
)

// ------------------------------------------------------------------
// protocol
// ------------------------------------------------------------------
type Protocol struct{
    common.IProtocol
}
func (p *Protocol) Execute(w http.ResponseWriter, r *http.Request){
    g := global.GetInstance()
    conf := g.Config()
    
    data := struct{
        Active string
    }{
        Active: "top",
    }
    templateFiles := protocol.GetTemplateFiles([]string{conf.PublicDir + "public/views/notfound.html"})
    tmpl, err := template.ParseFiles(templateFiles...)
    if err != nil{
        fmt.Println(err)
        return
    }
    var wkW bytes.Buffer
    err = tmpl.Execute(&wkW, data)
    if err != nil{
        fmt.Println(err)
        return
    }else{
        w.Write(wkW.Bytes())
    }
}
