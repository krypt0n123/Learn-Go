package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

var mu sync.Mutex
var count int
var palette = []color.Color{
	color.White,                        // 背景色
	color.RGBA{0x00, 0xff, 0x00, 0xff}, // 绿
	color.RGBA{0xff, 0x00, 0x00, 0xff}, // 红
	color.RGBA{0x00, 0x00, 0xff, 0xff}, // 蓝
	color.RGBA{0xff, 0xff, 0x00, 0xff}, // 黄
}

const (
	backgroundIndex = 0 //背景色
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cycles := 5
		if cyclesStr := r.URL.Query().Get("cycles"); cyclesStr != "" {
			if parsedCycles, err := strconv.Atoi(cyclesStr); err == nil {
				cycles = parsedCycles
			}
		}
		lissajous(w, cycles)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
	}
	fmt.Fprintf(w, "Host=%q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr=%q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q]=%q\n", k, v)
	}
}

func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

func lissajous(out io.Writer, cycles int) {
	const (
		// cycles  = 5
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 8
	)

	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)

			colorIndx := uint8(int(t/(2*math.Pi))%(len(palette)-1)) + 1
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				colorIndx)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
