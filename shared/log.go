package shared

import "fmt"

var logLevel = 0

func PushLog(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
	logLevel++
}

func PopLog() {
	if logLevel > 0 {
		fmt.Print("\x1b[2K\x1b[F\x1b[2K")
		logLevel--
	}
}

func SwapLog(format string, args ...any) {
	PopLog()
	PushLog(format, args...)
}
