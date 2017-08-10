package common
import(
    "net/http"
    "os"
)

// テキストファイルの取得
func GetFileString(path string)(ret string, err error){
    file, err := os.Open(path)
    if err != nil{
        return
    }
    defer file.Close()
    buf := make([]byte, 256)
    for{
        n, err := file.Read(buf)
        if n == 0{
            break
        }
        if err != nil{
            break
        }
        ret += string(buf[:n])
    }
    return
}

// ------------------------------------------------------------------
// interface
// ------------------------------------------------------------------
type IProtocol interface{
    Execute(w http.ResponseWriter, r *http.Request)
}

