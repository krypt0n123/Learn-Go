package main

import (
	"strings"
	"golang.org/x/tour/wc")

func WordCount(s string)map[string]int{
	wordmap:=make(map[string]int)

	words:=strings.Fields(s)

	for _,word := range words{
		wordmap[word]++
	}

	return wordmap
}

func main(){
	wc.Test(WordCount)
}