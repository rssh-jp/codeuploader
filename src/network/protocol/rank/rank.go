package rank
import(
    "bytes"
    "fmt"
    "net/http"
    "html/template"

    "common"
    "common/global"
    "network/protocol"
    "network/protocol/notfound"
    "network/mysql"
)

// ------------------------------------------------------------------
// protocol
// ------------------------------------------------------------------
type Protocol struct{
    common.IProtocol
}

type link struct{
    Url string
    Name string
}
type cronData struct{
    startTime string
    requireMinutes int
}
type codeData struct{
    Count int
    Name string
}
func (p *Protocol) Execute(w http.ResponseWriter, r *http.Request){
    r.ParseForm()
    g := global.GetInstance()
    conf := g.Config()

    db, err := mysql.GetConnection()
    if err != nil{
        fmt.Println(err)
        nf := &notfound.Protocol{}
        nf.Execute(w, r)
        return
    }
    defer db.Close()

    query := "select count, language_type, name, question_type from code order by count"
    rows, err := db.Query(query)
    if err != nil{
        fmt.Println(err)
        nf := &notfound.Protocol{}
        nf.Execute(w, r)
        return
    }
    var count, language_type, question_type int
    var name string
    codeDataHash := make(map[int]map[int][]codeData)
    for rows.Next(){
        rows.Scan(&count, &language_type, &name, &question_type)
        if _, ok := codeDataHash[question_type]; !ok{
            codeDataHash[question_type] = make(map[int][]codeData)
        }
        if _, ok := codeDataHash[question_type][language_type]; !ok{
            codeDataHash[question_type][language_type] = []codeData{}
        }
        codeDataHash[question_type][language_type] = append(codeDataHash[question_type][language_type], codeData{count, name})
    }
    fmt.Println(codeDataHash)


    templateFiles := protocol.GetTemplateFiles([]string{conf.PublicDir + "public/views/rank.html"})
    tmpl, err := template.ParseFiles(templateFiles...)
    if err != nil{
        fmt.Println(err)
        
        nf := &notfound.Protocol{}
        nf.Execute(w, r)
        return
    }
    data := struct{
        Active string
        Hash map[int]map[int][]codeData
    }{
        Active: "rank",
        Hash : codeDataHash,
    }
    var wkW bytes.Buffer
    err = tmpl.Execute(&wkW, data)
    if err != nil{
        fmt.Println(err)
        nf := &notfound.Protocol{}
        nf.Execute(w, r)
        return
    }else{
        w.Write(wkW.Bytes())
    }
}

