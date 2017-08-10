package top
import(
    "bytes"
    "database/sql"
    "fmt"
    "net/http"
    "html/template"
    "os"
    "os/exec"
    "io"
    "strconv"
    "strings"

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
func (p *Protocol) Execute(w http.ResponseWriter, r *http.Request){
    g := global.GetInstance()
    conf := g.Config()
    switch r.Method{
    case "GET":
        fmt.Println("GET")
    case "POST":
        fmt.Println("POST")
    }
    
    data := struct{
        Active string
        CountStr string
        ErrorStr string
    }{
        Active : "top",
        CountStr : "",
        ErrorStr : "",
    }


    file, h, err := r.FormFile("Input")
    if err == nil{
        defer file.Close()
        name    := r.FormValue("name")
        version := r.FormValue("version")
        language := r.FormValue("language")
        fmt.Println("lang : ", language)
        if name != ""{
            filename := "/opt/cron-manager/codegolf/" + version + "/" + name + "_" + h.Filename
            output := "/opt/cron-manager/codegolf/" + version + "/" + name + "_" + h.Filename + ".out"
            answerFile := "/opt/cron-manager/codegolf/" + version + "/answer"

            f, err := os.Create(filename)
            if err != nil{
                fmt.Println(err)
                execNotFound(w, r)
                return
            }
            defer f.Close()
            io.Copy(f, file)

            var cmd *exec.Cmd
            var language_type int
            switch(language){
            case "golang":
                cmd = exec.Command("/usr/local/go/bin/go", "run", filename)
                language_type = 1
            case "php":
                cmd = exec.Command("/usr/bin/php", filename)
                language_type = 2
            case "node":
                cmd = exec.Command("node", filename)
                language_type = 3
            }
            var question_type int
            switch(version){
            case "v1":
                question_type = 1
            case "v2":
                question_type = 2
            case "v3":
                question_type = 3
            }
            var out bytes.Buffer
            cmd.Stdout = &out
            err = cmd.Run()
            if err != nil{
                fmt.Println("err : ", err)
            }
            outfile, err := os.Create(output)
            if err != nil{
                fmt.Println(err)
                execNotFound(w, r)
                return
            }
            defer outfile.Close()
            io.Copy(outfile, &out)

            o, _ := exec.Command("diff", output, answerFile).Output()
            if string(o) == ""{
                str, _ := common.GetFileString(filename)
                count := strconv.Itoa(len(str))
                data.CountStr = count
                db, err := mysql.GetConnection()
                if err != nil{
                    fmt.Println(err)
                    execNotFound(w, r)
                    return
                }
                defer db.Close()

                fmt.Println("+++++++++++++++++++++")
                err = insertOnUpdateCode(db, language_type, question_type, len(str), name)
                if err != nil{
                    fmt.Println(err)
                }
            }else{
                data.ErrorStr = string(o)
            }
        }
    }



    templateFiles := protocol.GetTemplateFiles([]string{conf.PublicDir + "public/views/top.html"})
    tmpl, err := template.ParseFiles(templateFiles...)
    if err != nil{
        fmt.Println(err)
        
        nf := &notfound.Protocol{}
        nf.Execute(w, r)
        return
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
func execNotFound(w http.ResponseWriter, r *http.Request){
    nf := &notfound.Protocol{}
    nf.Execute(w, r)
}
func insertOnUpdateCode(db *sql.DB, language_type, question_type, count int, name string)(err error){
    // 検索
    query_args := []interface{}{language_type, name, question_type}
    query := "select id, count from code where language_type = ? and name = ? and question_type = ?"
    var id, oldCount int
    fmt.Println(query)
    fmt.Println(query_args)
    rows, err := db.Query(query, query_args...)
    if err != nil{
        return
    }
    defer rows.Close()

    for rows.Next(){
        rows.Scan(&id, &oldCount)
    }

    fmt.Println(id)
    fmt.Println(oldCount)

    // 既にデータがあるかどうか
    // データがあった場合に以前のカウントよりも今のが少ないか
    if id == 0 && oldCount == 0{
        itemList := make([]string, 4)
        itemList[0] = "\"" + strconv.Itoa(count)         + "\""
        itemList[1] = "\"" + strconv.Itoa(language_type) + "\""
        itemList[2] = "\"" + name                        + "\""
        itemList[3] = "\"" + strconv.Itoa(question_type) + "\""
        query := "insert into code (count, language_type, name, question_type) values ("
        query += strings.Join(itemList, ", ")
        query += ")"
        fmt.Println(query)
        _, err = db.Query(query)
        if err != nil{
            fmt.Println(err)
        }
    }
    if oldCount > count{
        query_args = []interface{}{strconv.Itoa(count), id}
        query := "update code set count = ? where id = ?"
        fmt.Println(query)
        fmt.Println(query_args)
        _, err = db.Query(query, query_args...)
        if err != nil{
            fmt.Println(err)
        }
    }
    return
}

