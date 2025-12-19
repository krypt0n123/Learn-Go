package main

import (
	"fmt"
	"time"
)

func say(s string){
	currentTime:=time.Now()
	TimewithMS := currentTime.Format("15:04:05.000")
	for i:=0;i<5;i++{
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s,TimewithMS)
	}
}

func main(){
	go say("world")
	say("hello")
}