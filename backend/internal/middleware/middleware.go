package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogrusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// เริ่มจับเวลา
		startTime := time.Now()

		// ส่งไม้ต่อให้ Handler ตัวถัดไปทำงาน
		c.Next()

		// คำนวณเวลาที่ใช้ไป
		latencyTime := time.Since(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// บันทึก Log ผ่าน Logrus
		entry := logrus.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"method":       reqMethod,
			"uri":          reqUri,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			if statusCode >= 500 {
				entry.Error("Server Error")
			} else if statusCode >= 400 {
				entry.Warn("Client Error")
			} else {
				entry.Info("Success")
			}
		}
	}
}
