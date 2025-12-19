package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)

func walkRecursive(t *tree.Tree, ch chan int){
	if t==nil{
		return
	}
		//中序遍历（左-根-右）
		walkRecursive(t.Left,ch)
		ch<-t.Value
		walkRecursive(t.Right,ch)
}

func Walk(t *tree.Tree,ch chan int){
	walkRecursive(t,ch)
	close(ch)
}

func Same(t1,t2 *tree.Tree) bool{
	ch1:= make(chan int)
	ch2:=make(chan int)

	go Walk(t1,ch1)
	go Walk(t2,ch2)

	for{
		v1,ok1:=<-ch1
		v2,ok2:=<-ch2

		if ok1!=ok2{
			return false
		}
		if !ok1{
			break
		}
		if v1!=v2{
			return false
		}
	}
	return true
}

func main(){
	fmt.Println("test Walk function")
	ch:=make(chan int)
	go Walk(tree.New(1),ch)

	fmt.Println("elements form tree.New(1):")
	for v:=range ch{
			fmt.Printf("%d ",v)
	}
	fmt.Println("\n")

	fmt.Println("test Same function")
	fmt.Printf("Same(tree.New(1),tree.New(1))-->%v\n",Same(tree.New(1),tree.New(1)))
	fmt.Printf("Same(tree.New(1),tree.New(2))-->%v\n",Same(tree.New(1),tree.New(2)))
}