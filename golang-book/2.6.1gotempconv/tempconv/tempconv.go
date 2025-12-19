package tempconv

import "fmt"

type Celsius float64
type Fahrenheit float64
type Kelvin float64 // 练习 2.1: 添加 Kelvin 类型

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

func (c Celsius) String() string    { return fmt.Sprintf("%g°C", c) }
func (f Fahrenheit) String() string { return fmt.Sprintf("%g°F", f) }
// 练习 2.1: 为 Kelvin 添加 String() 方法
func (k Kelvin) String() string    { return fmt.Sprintf("%gK", k) }