package asserts

import (
	"fmt"
	"github.com/go-test/deep"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}, message string) {
	if diff := deep.Equal(exp, act); diff != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(
			"\t\033[31m✖ %s:%d:%s\n\n\t\033[36mexp\033[39m: %#v\n\n\t\033[36mgot\033[39m: %#v\n\n\t\033[36mdiff\033[39m:\n\t\t%s\033[39m\n\n",
			filepath.Base(file),
			line,
			message,
			exp,
			act,
			strings.Join(diff, "\n\t\t"),
		)
		tb.FailNow()
	} else {
		fmt.Printf("\t\033[32;1m✔ \033[37;0m%s\033[39m\n", message)
	}
}

// Success fails the test if err != nil
func Success(tb testing.TB, err error) {
	if nil != err {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(
			"\t\033[31m✖ %s:%d\n\n\t\033[36mget unexpected error\u001B[39m:\n\t\t%s\u001B[39m\n\n",
			filepath.Base(file),
			line,
			err,
		)
		tb.FailNow()
	}
}
