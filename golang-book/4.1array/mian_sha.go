package main

import(
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var alg=flag.String("alg","sha256","Hash algorithm(sha256, sha384, sha512)")

func main(){
	flag.Parse()

	input,err:=io.ReadAll(os.Stdin)
	if err!=nil{
		log.Fatal("Failed to read from stdin:%v",err)
	}

	//根据flag的值选择hash算法
	switch *alg{
	case "sha256":
		hash:=sha256.Sum256(input)
		fmt.Printf("sha256:%x\n",hash)
	case "sha384":
		hash:=sha512.Sum384(input)
		fmt.Printf("sha384:%x\n",hash)
	case "sha512":
		hash:=sha512.Sum512(input)
		fmt.Printf("sha512:%x\n",hash)
	default:
		fmt.Fprintf(os.Stderr,"Unkonw algo '%s'",*alg)
		os.Exit(1)
	}
}