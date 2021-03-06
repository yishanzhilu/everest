package middleware

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	ginLogger := logrus.New()
	if viper.GetString("runmode") != "debug" {
		ginLogger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		ginLogger.SetFormatter(&logrus.TextFormatter{})
		ginLogger.SetOutput(os.Stdout)
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
		dataBytes := c.Writer.Size()

		if dataBytes < 0 {
			dataBytes = 0
		}
		if viper.GetString("runmode") == "debug" {
			statusColor := colorForStatus(statusCode)
			methodColor := colorForMethod(method)
			msg := fmt.Sprintf("[GIN] %s |%s  %s %3d| %s  %s %-7s %s | size: %d bytes | %s | (%dms)\n",
				clientIP,
				statusColor, reset, statusCode,
				methodColor, reset, method,
				path,
				dataBytes,
				clientUserAgent,
				latency)
			if len(c.Errors) > 0 {
				ginLogger.Error(msg, c.Errors.String())
			} else {
				if statusCode > 499 {
					ginLogger.Error(msg)
				} else if statusCode > 399 {
					ginLogger.Warn(msg)
				} else {
					ginLogger.Info(msg)
				}
			}
			return
		}

		entry := ginLogger.WithFields(logrus.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataBytes":  dataBytes,
			"userAgent":  clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Errorln(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			if statusCode > 499 {
				entry.Error()
			} else if statusCode > 399 {
				entry.Warn()
			} else {
				entry.Info()
			}
		}
	}
}
