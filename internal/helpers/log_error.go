package helpers

import (
	"log"
	"net/http"
	"runtime"
)

func LogErrorWithRequest(err error, r *http.Request, message string) {
	if err == nil {
		return
	}

	pc, file, line, ok := runtime.Caller(1)
	var fnName string
	if ok {
		fn := runtime.FuncForPC(pc)
		fnName = fn.Name()
	} else {
		file = "unknown"
		line = 0
		fnName = "unknown"
	}

	log.Printf(
		"[ERROR] %s: %v | func=%s, file=%s:%d | method=%s, url=%s\n",
		message, err, fnName, file, line, r.Method, r.URL.Path,
	)
}
