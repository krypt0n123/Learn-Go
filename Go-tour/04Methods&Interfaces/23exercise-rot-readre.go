package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct{
	r io.Reader
}

func (rot *rot13Reader)Read(p []byte)(n int, err error){
	n,err=rot.r.Read(p)

	for i:=0;i<n;i++{
		b:=p[i]
		if(b>='a'&& b<='m')||(b>='A'&&b<='M'){
			p[i]+=13
		}else if (b>='n'&& b<='z')||(b>='N'&&b<='Z'){
			p[i]-=13
		}
	}
	return n,err
}

func main(){
	s:=strings.NewReader("Lbh penpxrq gur pbqr!")
	r:=rot13Reader{s}
	io.Copy(os.Stdout,&r)
}