package logging

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type log struct {
	TimeStamp  string `json:"date"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	ClientIP   string `json:"clientIP"`
	StatusCode int    `json:"statusCode"`
	Latency    string `json:"latency"`
	UserAgent  string `json:"userAgent"`
	Error      string `json:"error,omitempty"`
}

func Logger(p gin.LogFormatterParams) string {
	var l log
	l.TimeStamp = p.TimeStamp.Format(time.RFC3339)
	l.Path = p.Path
	l.Method = p.Method
	l.ClientIP = p.ClientIP
	l.StatusCode = p.StatusCode
	l.Latency = p.Latency.String()
	l.UserAgent = p.Request.UserAgent()
	l.Error = p.ErrorMessage

	out, _ := json.Marshal(l)
	return fmt.Sprintln(string(out))
}
