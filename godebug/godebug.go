package godebug

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

var debugFlag map[string]bool

func init() {
	debugFlag = make(map[string]bool)
}

func SetDebugFlags(s string) {
	for _, vv := range strings.Split(s, ",") {
		debugFlag[vv] = true
	}
}

func DebugOn(s string) bool {
	return debugFlag[s]
}

var ColorRed = "\033[31;40m"
var ColorYellow = "\033[33;40m"
var ColorGreen = "\033[32;40m"
var ColorCyan = "\033[36;40m"
var ColorReset = "\033[0m"

// LF()  Return the File name and Line no as a string.
func LF(d ...int) string {
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	_, file, line, ok := runtime.Caller(depth)
	if ok {
		return fmt.Sprintf("File: %s LineNo:%d", file, line)
	} else {
		return fmt.Sprintf("File: Unk LineNo:Unk")
	}
}

// SVar returns 'v' in JSON format
func SVar(v interface{}) string {
	s, err := json.Marshal(v)
	// s, err := json.MarshalIndent ( v, "", "\t" )
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}

// SVarI returns 'v' in indented JSON format
func SVarI(v interface{}) string {
	// s, err := json.Marshal ( v )
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}

// InArrayString searches 'arr' for 's', returning -1 if not found or index.  Use Dijkstra L algorythm.
func InArrayString(s string, arr []string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return -1
}

// InArrayInt searches int/'arr' for int/'s', returning -1 if not found or index.  Use Dijkstra L algorythm.
func InArrayInt(s int, arr []int) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return -1
}
