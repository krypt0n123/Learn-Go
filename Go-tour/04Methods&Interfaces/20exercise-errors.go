package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64)(float64,error){
	if x<0{
		return 0,ErrNegativeSqrt(x)
	}
	z:=math.Sqrt(x)
	return z,nil
}

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string{
	return fmt.Sprintf("Can't sqrt negative num:%v",float64(e))
}

func main(){
	result,err:=Sqrt(2)
	if err!=nil{
		fmt.Println(err)
	}else{
		fmt.Println(result)
	}

	result,err = Sqrt(-2)
	if err!=nil{
		fmt.Println(err)
	}else{
		fmt.Println(result)
	}
}