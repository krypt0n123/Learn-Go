package main

import "fmt"

func main(){
	const freezingF,boilingF=32.0,212.0
	fmt.Printf("%g 华氏度=%g 度\n",freezingF,fToc(freezingF))
	fmt.Printf("%g 华氏度=%g 度",boilingF,fToc(boilingF))
}

func fToc(f float32)float32{
	return (f-32)*5/9
}