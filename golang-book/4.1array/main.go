package main

import(
	"fmt"
	"crypto/sha256"
	"encoding/binary"
	"array/popcount"
)

func DiffBits(h1,h2 [32]byte)int{
	count:=0
	for i:=0;i<4;i++{//32字节=4*8字节(uint64)
		offset:=i*8

		u1:=binary.BigEndian.Uint64(h1[offset:offset+8])
		u2:=binary.BigEndian.Uint64(h2[offset:offset+8])
		xorVal:=u1^u2
		count+=popcount.PopCount(xorVal)
	}
	return count
}

func main(){
	h1:=sha256.Sum256([]byte("x"))
	h2:=sha256.Sum256([]byte("X"))

	fmt.Printf("Hash 1 (%s): %x\n","x",h1)
	fmt.Printf("Hash 2 (%s): %x\n","X",h2)

	diff:=DiffBits(h1,h2)
	fmt.Printf("Number of different bits: %d\n",diff)

	//测试相同输入
	h3 := sha256.Sum256([]byte("hello"))
	h4 := sha256.Sum256([]byte("hello"))

	fmt.Printf("Hash 3 (%s): %x\n", "hello", h3)
	fmt.Printf("Hash 4 (%s): %x\n", "hello", h4)
	diff2 := DiffBits(h3, h4)
	fmt.Printf("Number of different bits: %d\n", diff2)
}