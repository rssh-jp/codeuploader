package main
import(
    "fmt"
    "log"
    "net/http"
    "net"

    "common"
    "common/global"
    "network/protocol/top"
    "network/protocol/rank"
)

func test(){
    fmt.Println("test")
    fmt.Println("test2")
}

// ------------------------------------------------------------------
// server
// ------------------------------------------------------------------
type Server struct{}
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request){
    redirectUrl := ""
    var protocol common.IProtocol
    switch r.URL.Path{
    case "/top":
        protocol = &top.Protocol{}
    case "/rank":
        protocol = &rank.Protocol{}
    default:
        redirectUrl = "/top"
    }

    if redirectUrl == ""{
        protocol.Execute(w, r)
    }else{
        http.Redirect(w, r, "/top", http.StatusFound)
    }

}

func serve(){
    var listener net.Listener
    listener, err := net.Listen("tcp", ":2000")
    if err != nil{
        log.Fatal(err)
    }
    s := Server{}
    http.Handle("/", s)
    http.Serve(listener, nil)
}

func initialize(){
    global.Initialize()
}

func main(){
    initialize()
    serve()
}
