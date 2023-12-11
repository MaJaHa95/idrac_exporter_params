package logging

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var mx sync.Mutex
var verbose bool
var datefmt string = "2006-01-02T15:04:05.000"

func print(l string, a ...interface{}) {
	mx.Lock()
	dt := time.Now()
	f := dt.Format(datefmt) + fmt.Sprintf(" %-5s ", l)
	fmt.Print(f)
	f = strings.TrimSpace(fmt.Sprint(a...))
	fmt.Println(f)
	mx.Unlock()
}

func printf(l string, f string, a ...interface{}) {
	mx.Lock()
	dt := time.Now()
	f = dt.Format(datefmt) + fmt.Sprintf(" %-5s ", l) + strings.TrimSpace(f) + "\n"
	fmt.Printf(f, a...)
	mx.Unlock()
}

func SetVerbose(v bool) {
	mx.Lock()
	verbose = v
	mx.Unlock()
}

func Info(a ...interface{}) {
	print("INFO", a...)
}

func Infof(fmt string, a ...interface{}) {
	printf("INFO", fmt, a...)
}

func Error(err error, a ...interface{}) {
	print("ERROR", a...)
	printf("%s", err.Error())
}

func Errorf(err error, fmt string, a ...interface{}) {
	printf("ERROR", fmt, a...)
	printf("%s", err.Error())
}

func Fatal(err error, a ...interface{}) {
	print("FATAL", a...)
	printf("%s", err.Error())
	os.Exit(1)
}

func Fatalf(err error, fmt string, a ...interface{}) {
	printf("FATAL", fmt, a...)
	printf("%s", err.Error())
	os.Exit(1)
}

func Debug(a ...interface{}) {
	if verbose {
		print("DEBUG", a...)
	}
}

func Debugf(fmt string, a ...interface{}) {
	if verbose {
		printf("DEBUG", fmt, a...)
	}
}
