package main

import (
	"fmt"
	"popcount/popcount"
)

func main() {
	num := uint64(11) //1011
	count := popcount.PopCountShift(num)
	fmt.Printf("%d的二进制数中有%d个'1'\n", num, count)

	num = uint64(1234567890) //1011
	count = popcount.PopCountShift(num)
	fmt.Printf("%d的二进制数中有%d个'1'", num, count)
}
