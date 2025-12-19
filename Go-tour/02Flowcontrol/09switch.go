package main

import(
	"fmt"
	"runtime"
)
func main(){
	fmt.Printf("Go的运行环境:")
	switch os:=runtime.GOOS;os {
	case "darwin":
		fmt.Println("MacOS")
	case "linux":
		fmt.Println("Linux")
	default:
		fmt.Printf("%s.\n",os)
	}
} 