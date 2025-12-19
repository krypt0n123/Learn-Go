package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
			fmt.Fprintf(os.Stderr, "自动添加 http://前缀: %s\n", url)
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprint(os.Stderr, "fetch:%v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		fmt.Printf("HTTP Status: %s ---\n",resp.Status)

		_,err=io.Copy(os.Stdout,resp.Body)
		if err != nil {
			fmt.Fprint(os.Stderr, "fetch:reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}
