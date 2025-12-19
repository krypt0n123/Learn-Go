package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	defaultWidth, defaultHeight = 600, 320
	cells                       = 100
	xyrange                     = 30.0
	angle                       = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

type SurfaceFunction func(x, y float64) float64

type Params struct {
	Width       int
	Height      int
	PeakColor   string
	ValleyColor string
	Shape       string
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("服务器已启动")
	log.Fatal(http.ListenAndServe("localhost:8000",nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")

	params := parseRequestParameters(r)
	surfaceFunc := getSurfaceFunction(params.Shape)

	generateSVG(w, params, surfaceFunc)
}

// 解析URL查询参数以实现自定义
func parseRequestParameters(r *http.Request) Params {
	p := Params{
		Width:       defaultWidth,
		Height:      defaultHeight,
		PeakColor:   "#ff0000",
		ValleyColor: "#0000ff",
		Shape:       "sinc", //默认图形
	}
	if val, err := strconv.Atoi(r.URL.Query().Get("width")); err == nil && val > 0 {
		p.Width = val
	}
	if val, err := strconv.Atoi(r.URL.Query().Get("height")); err == nil && val > 0 {
		p.Height = val
	}
	if color := r.URL.Query().Get("peak"); color != "" {
		p.PeakColor = "#" + color
	}
	if color := r.URL.Query().Get("valley"); color != "" {
		p.ValleyColor = "#" + color
	}
	if shape := r.URL.Query().Get("shape"); shape != "" {
		p.Shape = shape
	}
	return p
}

// 根据shape参数选择要渲染的函数
func getSurfaceFunction(shape string) SurfaceFunction {
	switch shape {
	case "eggbox":
		return eggBox
	case "moguls":
		return moguls
	case "saddle":
		return saddle
	default:
		return sinc
	}
}

// 原始函数：sin(r)/r
func sinc(x, y float64) float64 {
	r := math.Hypot(x, y) //到(0，0)的距离
	if r == 0 {
		return 1.0
	}
	return math.Sin(r) / r
}

// 蛋盒函数
func eggBox(x, y float64) float64 {
	return 0.1 * (math.Cos(x) + math.Cos(y))
}

// 雪堆函数
func moguls(x, y float64) float64 {
	return 0.1 * (math.Sin(x*math.Pi/5) - math.Cos(y*math.Pi/5))
}

// 马鞍函数(双曲抛物面)
func saddle(x, y float64) float64 {
	a, b := 25.0, 15.0
	return (y*y)/(b*b) - (x*x)/(a*a)
}

func generateSVG(out io.Writer, p Params, f SurfaceFunction) {
	minZ, maxZ := getMinMaxZ(f)

	//输出SVG标签头
	fmt.Fprintf(out,"<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke:grey;fill:white; stroke-width:0.7' "+
		"width='%d' height='%d'>", p.Width, p.Height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az := corner(i+1, j, p, f)
			bx, by, bz := corner(i, j, p, f)
			cx, cy, cz := corner(i, j+1, p, f)
			dx, dy, dz := corner(i+1, j+1, p, f)

			if isInvalid(ax, ay, bx, by, cx, cy, dx, dy) {
				continue
			}

			avgZ := (az + bz + cz + dz) / 4.0
			color := getColor(avgZ, minZ, maxZ, p.PeakColor, p.ValleyColor)

			fmt.Fprintf(out, "<polygon style='fill: %s' points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				color, ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintln(out,"</svg>")
}

func corner(i, j int, p Params, f SurfaceFunction) (sx, sy, z float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z = f(x, y)

	xyscale := float64(p.Width) / 2 / xyrange
	zscale := float64(p.Height) * 0.4
	sx = float64(p.Width)/2 + (x-y)*cos30*xyscale
	sy = float64(p.Height)/2 + (x+y)*sin30*xyscale - z*zscale
	return
}

// isInvalid 检查传入的浮点数是否有NaN或无穷大的值
func isInvalid(values ...float64) bool {
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return true
		}
	}
	return false
}

// getMinmaxZ遍历所有点，找出给定函数的z值的最大和最小值，用于颜色缩放
func getMinMaxZ(f SurfaceFunction) (min, max float64) {
	min = math.Inf(1)
	max = math.Inf(-1)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			x := xyrange * (float64(i)/cells - 0.5)
			y := xyrange * (float64(j)/cells - 0.5)
			z := f(x, y)
			if !math.IsNaN(z) && !math.IsInf(z, 0) {
				if z < min {
					min = z
				}
				if z > max {
					max = z
				}
			}
		}
	}
	return
}

// getColor根据一个值在最大最小之间的位置，计算出两种十六进制颜色之间的插值颜色
func getColor(val, min, max float64, peakColor, valleyColor string) string {
	//将val标准化到[0,1]区间
	ratio := (val - min) / (max - min)
	if max-min == 0 {
		ratio = 0.5
	}
	pr, pg, pb := parseHexColor(peakColor)
	vr, vg, vb := parseHexColor(valleyColor)

	r := uint8(float64(vr) + ratio*(float64(pr)-float64(vr)))
	g := uint8(float64(vg) + ratio*(float64(pg)-float64(vg)))
	b := uint8(float64(vb) + ratio*(float64(pb)-float64(vb)))

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// 将#RRGGBB格式的颜色字符串解析为R，G，B分量
func parseHexColor(s string) (r, g, b uint8) {
	s = s[1:] //移除‘#’
	if len(s) != 6 {
		return 0, 0, 0
	}
	rVal, _ := strconv.ParseInt(s[0:2], 16, 0)
	gVal, _ := strconv.ParseInt(s[2:4], 16, 0)
	bVal, _ := strconv.ParseInt(s[4:6], 16, 0)
	return uint8(rVal), uint8(gVal), uint8(bVal)
}