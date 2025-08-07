package log

import "fmt"

const (
	Default = iota
	Black
	DarkGray
	Red
	Green
	Yellow
	Blue
	Purple
	Cyan
	White
	BoldPurple
	BoldRed
	BoldGreen
)

var colorCode = []string{"0m", "0;30m", "1;30m", "0;31m", "0;32m", "0;33m", "0;34m", "0;35m", "0;36m", "0;37m", "1;35m", "1;31m", "1;32m"}

// 预设颜色代码
func ColorTxt(color int, str string) string {
	return fmt.Sprintf("\x1b[%s%s\x1b[0m", colorCode[color], str)
}

// 自定义颜色代码
func ColorTxtForCode(color string, str string) string {
	return fmt.Sprintf("\x1b[%s%s\x1b[0m", color, str)
}
