package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"main/http"
	"main/storeage"
)

var Dest string
var Port int
var Bind string

func init() {
	flag.Set("logtostderr", "true")
	flag.StringVar(&Dest, "storeto", "./temp", "dir to store states")
	flag.StringVar(&Bind, "bind", "0.0.0.0", "bind to ip")
	flag.IntVar(&Port,"port", 8080, "bind to port")
	flag.Parse()
}

func main() {
	glog.Info("run state service")
	glog.Info(fmt.Sprintf("store dir: '%s'", Dest))
	fmt.Println("Start")
	storeageProvider := storeage.FileStorageBuilder{
		Directory: Dest,
	}
    s := http.InitServer(Port, Bind, &storeageProvider)
	s.Run()
}

