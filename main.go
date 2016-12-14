package main
import (
    "fmt"
    "net"
    "reflect"
    "unsafe"
	"strings"
    "io/ioutil"
	"encoding/json"
    "path/filepath"
    "os"
)
var httpret  []byte
var httpretstr string
func b2s(buf []byte) string {
    return *(*string)(unsafe.Pointer(&buf))
}
func getCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}
type config struct{
    Port string
    TargetSSH string
    TargetHTTP string
    HTTPReturn200 string
}
var myconfig config
func s2b(s *string) []byte {
    return *(*[]byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(s))))
}
func handleReConn(c net.Conn,s net.Conn) {
    defer c.Close()
    defer s.Close()
    buf := make([]byte, 4096)
    for {
        len1,err1:=s.Read(buf)
        if err1!=nil || len1<=0{
            return
        }
        len2,err2:=c.Write(buf[0:len1])
        if err2!=nil || len2<=0{
            return
        }
    }
}
func handleConn(c net.Conn) {
    var target net.Conn
    var err error
    defer func(){
        if target!=nil{
            target.Close()
        }
        c.Close()
    }()
    buf:=make([]byte,4096)
    var len=0
    var flag=0
    var targetIP string
    for{
        len,err=c.Read(buf)
        if len<=0||err!=nil {
            return
        }
        if flag==0{
            if b2s(buf[0:3])=="SSH"{
                targetIP=myconfig.TargetSSH
            }else{
                if myconfig.HTTPReturn200==""||strings.Index(b2s(buf),myconfig.HTTPReturn200)==-1{
                targetIP=myconfig.TargetHTTP
                }else{
                    c.Write(httpret)
                    continue
                }
            }
            target,err=net.Dial("tcp",targetIP)
            if err!=nil{
                return
            }
            flag=1
            go handleReConn(c,target)
        }
        len2,err2:=target.Write(buf[0:len])
        if err2!=nil || len2<=0{
            return
        }
    }
}

func main() {
    version:="1.0"
    fmt.Println("CrabProxy Ver",version)
    httpretstr="HTTP/1.1 200 OK\x0d\x0aDate: Sun, 11 Dec 2016 09:26:11 GMT\x0d\x0aConnection: keep-alive\x0d\x0a\x0d\x0a"
    httpret=s2b(&httpretstr)
    jsonbuf,jsonerr:=ioutil.ReadFile(getCurrentDirectory()+"/config.json")
    if jsonerr != nil {
        fmt.Println("failed to open config.json:", jsonerr)
        return
    }
    jsonerr1:=json.Unmarshal(jsonbuf,&myconfig)
    if jsonerr1 != nil {
        fmt.Println("failed to Unmarshal config.json:", jsonerr1)
        return
    }
    if myconfig.Port==""{
        fmt.Println("port is null,please edit config.json")
        return
    }
    if myconfig.TargetSSH==""&&myconfig.TargetHTTP==""{
        fmt.Println("both targetSSH and targetHTTP cannot be null,please edit config.json")
        return
    }
    l, err := net.Listen("tcp", ":"+myconfig.Port)
    if err != nil {
        fmt.Println("listen error:", err)
        return
    }
    fmt.Println("server running at port",myconfig.Port)
    for {
        c, err := l.Accept()
        if err != nil {
            fmt.Println("accept error:", err)
            break
        }
        go handleConn(c)
    }
}