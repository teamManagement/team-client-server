package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"team-client-server/entrance"
)

func main() {

	go func() {
		host := "0.0.0.0:6061"
		if err := http.ListenAndServe(host, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", host)
			os.Exit(1)
		}
	}()

	entrance.Run()
}
