package middleware

import (
	"bytes"
	"io"
	"log"
	"log/slog"
	"os"
	"road2ca/internal/entity"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"
	"time"

	"github.com/google/uuid"
)

type Logger interface {
	SettingLogger(c *minigin.Context)
}

type logger struct{
	accessLogger *slog.Logger
	errorLogger  *slog.Logger
}

func NewLogger() (Logger, error) {
	accessFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	errorFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	accessLogger := slog.New(slog.NewTextHandler(accessFile, nil))
	errorLogger := slog.New(slog.NewTextHandler(errorFile, nil))
	return &logger{
		accessLogger: accessLogger,
		errorLogger:  errorLogger,
	}, nil
}

// SettingLogger ログを出力するミドルウェア
func (l *logger) SettingLogger(c *minigin.Context) {
	startTime := time.Now()

	// リクエストボディは一度しか読み取れないため、後続で再度読み取れるようにする
	var reqBody []byte
	if c.Request.Body != nil {
		reqBody, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	c.Next()
	duration := time.Since(startTime)
	statusCode := c.Writer.Status()

	var userId string
	if user, ok := c.Request.Context().Value(service.ContextKey).(*entity.User); ok {
		uuid, err := uuid.FromBytes(user.ID)
		if err != nil {
			userId = "unknown"
		} else {
			userId = uuid.String()
		}
	} else {
		userId = "unknown"
	}

	accessAttr := []slog.Attr{
		slog.Time("time", startTime),
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
		slog.Int("status", statusCode),
		slog.Duration("duration", duration),
		slog.String("ip", c.Request.RemoteAddr),
		slog.String("userId", userId),
		slog.String("requestBody", string(reqBody)),
	}
	accessArgs := make([]any, len(accessAttr))
	for i, attr := range accessAttr {
		accessArgs[i] = attr
	}

	l.accessLogger.Info("access", accessArgs...)

	log.Printf("Method: %s, Path: %s, Status: %d, Duration: %s",
		c.Request.Method, c.Request.URL.Path, statusCode, duration)
}
