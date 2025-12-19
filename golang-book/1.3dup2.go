package main

import (
	"bufio"
	"fmt"
	"os"
)

func main(){
	counts:=make(map[string]int)
	files:=os.Args[1:]
	if len(files)==0{
		countLines(os.Stdin,counts,"terminal input")
	}else{
		for _,arg:=range files{
			f,err:=os.Open(arg)
			if err !=nil{
				fmt.Fprint(os.Stderr,"dup2:%v\n",err)
				continue
			}
			countLines(f,counts,arg)
			f.Close()
		}
	}
}

func countLines(f *os.File,counts map[string]int,filename string){
	input:=bufio.NewScanner(f)
	for input.Scan(){
		line:=input.Text()
		counts[line]++
		if counts[line]==2{
			fmt.Printf("Duplicate line:'%s'\tfound in :%s\n",line,filename)
		}
	}
}