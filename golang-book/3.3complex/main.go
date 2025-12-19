package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"math/cmplx"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", fractalHandler)
	log.Println("Server started on http://localhost:8000")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// 处理HTTP请求
func fractalHandler(w http.ResponseWriter, r *http.Request) {
	centerX, centerY, zoom := 0.0, 0.0, 1.0
	var err error

	if xStr := r.FormValue("x"); xStr != "" {
		centerX, err = strconv.ParseFloat(xStr, 64)
		if err != nil {
			http.Error(w, "invalid 'x' parameter", http.StatusBadRequest)
			return
		}
	}
	if yStr := r.FormValue("y"); yStr != "" {
		centerY, err = strconv.ParseFloat(yStr, 64)
		if err != nil {
			http.Error(w, "invalid 'y' parameter", http.StatusBadRequest)
			return
		}
	}
	if zoomStr := r.FormValue("zoom"); zoomStr != "" {
		zoom, err = strconv.ParseFloat(zoomStr, 64)
		if err != nil {
			http.Error(w, "invalid 'zoom' parameter", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "image/png")

	generateFractal(w, centerX, centerY, zoom)
}

//根据中心点和缩放级别生成图像
func generateFractal(out io.Writer,centerX,centerY,zoom float64){
	const(
		width,height=1024,1024
	)
	viewSize:=4.0/zoom
	xmin:=centerX-viewSize/2
	xmax:=centerX+viewSize/2
	ymin:=centerY-viewSize/2
	ymax:=centerY+viewSize/2
	
	img:=image.NewRGBA(image.Rect(0,0,width,height))
	for py:=0;py<height;py++{
		y:=float64(py)/height*(ymax-ymin)+ymin
		for px:=0;px<width;px++{
			x:=float64(px)/width*(xmax-xmin)+xmin
			z:=complex(x,y)
			img.Set(px,py,mandlbrot(z))
		}
	}
	png.Encode(out,img)
}

func mandlbrot(z complex128) color.Color {
	const iterations = 200

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			r := uint8(math.Sin(float64(n)*0.15+0)*127 + 128)
			g := uint8(math.Sin(float64(n)*0.15+2)*127 + 128)
			b := uint8(math.Sin(float64(n)*0.15+4)*127 + 128)
			return color.RGBA{r, g, b, 255}
		}
	}
	return color.Black
}

