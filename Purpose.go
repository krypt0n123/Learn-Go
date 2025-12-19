package main

import "fmt"

type User struct {
	ID   int
	Name string
}

func main() {
	num := 100
	pi := 3.14159
	str := "Go语言"
	user := User{ID: 1, Name: "Gemini"}
	ptr := &user

	fmt.Println("--- 通用 ---")
	fmt.Printf("默认值 %%v: %v\n", user)
	fmt.Printf("带字段名 %%+v: %+v\n", user)
	fmt.Printf("Go语法 %%#v: %#v\n", str)
	fmt.Printf("类型 %%T: %T\n", user)

	fmt.Println("\n--- 整数 ---")
	fmt.Printf("十进制 %%d: %d\n", num)
	fmt.Printf("二进制 %%b: %b\n", num)
	fmt.Printf("八进制 %%o: %o\n", num)
	fmt.Printf("十六进制 %%x: %x\n", num)
	fmt.Printf("字符 %%c: %c\n", num) // 100 对应的 ASCII 字符是 'd'

	fmt.Println("\n--- 浮点数 ---")
	fmt.Printf("标准小数 %%f: %f\n", pi)
	fmt.Printf("两位小数 %%.2f: %.2f\n", pi)
	fmt.Printf("科学计数法 %%e: %e\n", pi)

	fmt.Println("\n--- 字符串 ---")
	fmt.Printf("普通字符串 %%s: %s\n", str)
	fmt.Printf("带引号字符串 %%q: %q\n", str)

	fmt.Println("\n--- 指针 ---")
	fmt.Printf("指针地址 %%p: %p\n", ptr)
}