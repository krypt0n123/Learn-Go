package main

import (
	"fmt"
	"gotempconv/tempconv" // 导入你修改后的包
)

func main() {
	boilingK := tempconv.CToK(tempconv.BoilingC) // 100°C
	fmt.Printf("沸点 (Boiling point) 是 %v, 或者 %v\n", tempconv.BoilingC, boilingK)
	// 输出: 沸点 (Boiling point) 是 100°C, 或者 373.15K

	// 测试开尔文转换为摄氏度
	absoluteZeroInC := tempconv.KToC(0) // 0K
	fmt.Printf("绝对零度 (Absolute zero) 是 %v, 或者 %v\n", tempconv.Kelvin(0), absoluteZeroInC)
	// 输出: 绝对零度 (Absolute zero) 是 0K, 或者 -273.15°C
}
