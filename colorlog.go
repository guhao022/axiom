package axiom

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	Gray = uint8(iota + 90)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	//NRed      = uint8(31) // Normal
	EndColor = "\033[0m"

	INFO = "INFO"
	TRAC = "TRAC"
	ERRO = "ERRO"
	WARN = "WARN"
	SUCC = "SUCC"
	SKIP = "SKIP"
)

func CLog(format string, a ...interface{}) {
	fmt.Println(colorLogS(format, a...))
}

func colorLogS(format string, a ...interface{}) string {
	log := fmt.Sprintf(format, a...)

	var clog string

	if runtime.GOOS != "windows" {
		// Level.
		i := strings.Index(log, "]")
		if log[0] == '[' && i > -1 {
			clog += "[" + getColorLevel(log[1:i]) + "]"
		}

		log = log[i+1:]

		// Error.
		log = strings.Replace(log, "[ ", fmt.Sprintf("\033[%dm", Red), -1)
		log = strings.Replace(log, " ]", EndColor+"", -1)

		// Path.
		log = strings.Replace(log, "( ", fmt.Sprintf("\033[%dm", Yellow), -1)
		log = strings.Replace(log, " )", EndColor+"", -1)

		log = strings.Replace(log, "< ", fmt.Sprintf("\033[%dm", Cyan), -1)
		log = strings.Replace(log, " >", EndColor+"", -1)

		// Highlights.
		log = strings.Replace(log, "# ", fmt.Sprintf("\033[%dm", Magenta), -1)
		log = strings.Replace(log, " #", EndColor+"", -1)

		log = strings.Replace(log, "@@ ", fmt.Sprintf("\033[%dm", Green), -1)
		log = strings.Replace(log, " @@", EndColor+"", -1)

		log = clog + log

	} else {
		// Level.
		i := strings.Index(log, "]")
		if log[0] == '[' && i > -1 {
			clog += "[" + log[1:i] + "]"
		}

		log = log[i+1:]

		// Error.
		log = strings.Replace(log, "[ ", "", -1)
		log = strings.Replace(log, " ]", "", -1)

		// Path.
		log = strings.Replace(log, "( ", "", -1)
		log = strings.Replace(log, " )", "", -1)

		// Highlights.
		log = strings.Replace(log, "# ", "", -1)
		log = strings.Replace(log, " #", "", -1)

		log = strings.Replace(log, "@@ ", "", -1)
		log = strings.Replace(log, " @@", "", -1)

		log = clog + log
	}

	return time.Now().Format("2006/01/02 15:04:05 ") + log
}

// getColorLevel returns colored level string by given level.
func getColorLevel(level string) string {
	level = strings.ToUpper(level)
	switch level {
	case INFO:
		return fmt.Sprintf("\033[%dm%s\033[0m", Blue, level)
	case TRAC:
		return fmt.Sprintf("\033[%dm%s\033[0m", Cyan, level)
	case ERRO:
		return fmt.Sprintf("\033[%dm%s\033[0m", Red, level)
	case WARN:
		return fmt.Sprintf("\033[%dm%s\033[0m", Magenta, level)
	case SUCC:
		return fmt.Sprintf("\033[%dm%s\033[0m", Green, level)
	case SKIP:
		return fmt.Sprintf("\033[%dm%s\033[0m", Yellow, level)
	default:
		return level
	}
	return level
}