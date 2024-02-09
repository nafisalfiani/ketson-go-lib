package log

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/nafisalfiani/ketson-go-lib/appcontext"
	"github.com/nafisalfiani/ketson-go-lib/errors"
)

const (
	Zerolog = "zerolog"
	Logrus  = "logrus"
)

type Interface interface {
	Trace(ctx context.Context, obj any)
	Debug(ctx context.Context, obj any)
	Info(ctx context.Context, obj any)
	Warn(ctx context.Context, obj any)
	Error(ctx context.Context, obj any)
	Fatal(ctx context.Context, obj any)
	Panic(obj any)
}

type Config struct {
	Level string `env:"LEVEL"`
}

var now = time.Now

func Init(cfg Config, logType string) Interface {
	switch logType {
	case Zerolog:
		return initZerolog(cfg)
	case Logrus:
		return initLogrus(cfg)
	default:
		return initZerolog(cfg)
	}
}

func getContextFields(ctx context.Context) map[string]any {
	reqstart := appcontext.GetRequestStartTime(ctx)
	apprespcode := appcontext.GetAppResponseCode(ctx)
	appErrMsg := appcontext.GetAppErrorMessage(ctx)
	timeElapsed := "0ms"
	if !time.Time.IsZero(reqstart) {
		timeElapsed = fmt.Sprintf("%dms", int64(now().Sub(reqstart)/time.Millisecond))
	}

	cf := map[string]interface{}{
		"request_id":      appcontext.GetRequestId(ctx),
		"user_agent":      appcontext.GetUserAgent(ctx),
		"user_id":         appcontext.GetUserId(ctx),
		"service_version": appcontext.GetServiceVersion(ctx),
		"time_elapsed":    timeElapsed,
	}

	if apprespcode > 0 {
		cf["app_resp_code"] = apprespcode
	}

	if appErrMsg != "" {
		cf["app_err_msg"] = appErrMsg
	}

	return cf
}

func getCaller(obj any) any {
	switch tr := obj.(type) {
	case error:
		file, line, msg, err := errors.GetCaller(tr)
		if err == nil {
			obj = fmt.Sprintf("%s:%#v --- %s", file, line, msg)
		}
	case string:
		obj = tr
	default:
		obj = fmt.Sprintf("%#v", tr)
	}

	return obj
}

func getPanicStacktrace() map[string]any {
	errStack := strings.Split(strings.ReplaceAll(string(debug.Stack()), "\t", ""), "\n")
	return map[string]any{
		"stracktrace": errStack,
	}
}
