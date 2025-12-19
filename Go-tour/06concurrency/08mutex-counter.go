package main

import (
	"fmt"
	"sync"
	"time"
)

//SafeCounter是安全并发的
type SafeCounter struct{
	mu  sync.Mutex
	v 	map[string]int
}

//Inc对给定键的计数加一
func (c *SafeCounter) Inc(key string){
	c.mu.Lock()
	//锁定使得一次只有一个Go协程可以访问映射成c，v
	c.v[key]++
	c.mu.Unlock()
}

//Value返回给定键的计数的当前值
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	//锁定使得一次只有一个Go协程可以访问映射成c，v
	defer c.mu.Unlock()
	return c.v[key]
}

func main(){
	c:=SafeCounter{v:make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go c.Inc("somekey")
	}

	time.Sleep(time.Second)
	fmt.Println(c.Value("somekey"))
}