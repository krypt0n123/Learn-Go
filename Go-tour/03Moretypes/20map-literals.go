package main

import "fmt"

type Vertex struct{
	Lat, Long float64
}

var m=map[string]Vertex{
	"Bell":Vertex{
		40,-74,
	},
	"Google":Vertex{
		37.1212,-122,
	},
}

func main(){
	fmt.Println(m)
}