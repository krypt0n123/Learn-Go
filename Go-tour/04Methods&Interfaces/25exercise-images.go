package main

import (
	"golang.org/x/tour/pic"
	"image"
	"image/color"
)
type Image struct{
	Width,Height int
}

func (i *Image)Bounds() image.Rectangle{
	return image.Rect(0,0,i.Width,i.Height)
}

func (i *Image)ColorModel() color.Model{
	return color.RGBAModel
}

func (i *Image)At(x,y int) color.Color{
	v:=uint8(x^y)
	return color.RGBA{v,v,255,255}
}

func main(){
	m:=&Image{Width: 256,Height: 256}
	pic.ShowImage(m)
}