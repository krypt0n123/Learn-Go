package main

import (
	"fmt"
	"strings"
)

type List[T any] struct{
	next *List[T]
	val	  T
}

func (l *List[T])Push(v T) *List[T]{
	return &List[T]{next: l,val: v}
}

func (l *List[T])Len() int{
	count:=0
	for node:=l; node != nil; node = node.next{
		count++
	}
	return count
}

func (l *List[T]) String() string{
	var sb strings.Builder
	sb.WriteString("[")
	for node:=l;node!=nil;node=node.next{
		sb.WriteString(fmt.Sprintf("%v",node.val))
		if node.next!=nil{
			sb.WriteString(" ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func main(){
	fmt.Println("----整数链表----")
	var intHead *List[int]
	fmt.Printf("初始状态：%s, 长度：%d\n",intHead,intHead.Len())

	intHead=intHead.Push(10)
	intHead=intHead.Push(20)
	intHead=intHead.Push(30)

	fmt.Printf("After add elements:%s, lenth: %d\n",intHead,intHead.Len())
	fmt.Println()

	//-----------------------------------------------------------
	fmt.Println("----字符串链表----")
	var stringHead *List[string]
	fmt.Printf("初始状态：%s, 长度：%d\n",stringHead,stringHead.Len())
	stringHead=stringHead.Push("Go")
	stringHead=stringHead.Push("is")
	stringHead=stringHead.Push("awesome")

	fmt.Printf("After add elements:%s, lenth: %d\n",stringHead,stringHead.Len())
}