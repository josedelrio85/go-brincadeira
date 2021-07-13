package nivoriacomp

import (
	"log"
	"runtime"

	"github.com/josedelrio85/voalarm"
)

// ResponseError an alarm when an error occurs
func ResponseError(message string, err error) {
	fancyHandleError(err)

	alarm := voalarm.NewClient("")
	alarm.SendAlarm("nivoriacomp", voalarm.Acknowledgement, err)
}

// fancyHandleError logs the error and indicates the line and function
func fancyHandleError(err error) (b bool) {
	if err != nil {
		// using 1 => it will actually log where the error happened, 0 = this function.
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
		b = true
	}
	return
}
