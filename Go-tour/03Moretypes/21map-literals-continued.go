package main

import "fmt"

type Vertex struct{
	Lat,Long float64
}

var m=map[string]Vertex{
	"Bell"	:	{40,-234},
	"Google":	{37,-234},
}

func main(){
	fmt.Println(m)
}