package middleware

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yishanzhilu/api-template/pkg/common"
)

var timeFormat = "02/Jan/2006:15:04:05 -0700"

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}

// GinLogger is the logrus logger middleware
func GinLogger() gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknow"
	}
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()

		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)
		if dataLength < 0 {
			dataLength = 0
		}
		if viper.GetString("runmode") == "debug" {
			msg := fmt.Sprintf("[GIN] %s |%s  %s %3d| %s  %s %-7s %s | size: %d | %s | (%dms)\n",
				clientIP,
				statusColor, reset, statusCode,
				methodColor, reset, method,
				path,
				dataLength,
				clientUserAgent,
				latency)
			if len(c.Errors) > 0 {
				common.Logger.Error(msg, c.Errors.ByType(gin.ErrorTypePrivate).String())
			} else {
				if statusCode > 499 {
					common.Logger.Error(msg)
				} else if statusCode > 399 {
					common.Logger.Warn(msg)
				} else {
					common.Logger.Info(msg)
				}
			}
			return
		}

		entry := common.Logger.WithFields(logrus.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Errorln(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			if statusCode > 499 {
				entry.Error("\n")
			} else if statusCode > 399 {
				entry.Warn("\n")
			} else {
				entry.Info("\n")
			}
		}
	}
}
