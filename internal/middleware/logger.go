package middleware

import (
	"bytes"
	"context"
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

type logger struct {
	accessLogger *slog.Logger
	errorLogger  *slog.Logger
}

func NewLogger() (Logger, error) {
	accessFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	errorFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	accessLogger := slog.New(slog.NewJSONHandler(accessFile, nil))
	errorLogger := slog.New(slog.NewJSONHandler(errorFile, nil))
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

	var reqId string
	reqUuid, err := uuid.NewV7()
	if err != nil {
		reqId = "unknown"
	} else {
		reqId = reqUuid.String()
	}

	// 一応後続ハンドラでもリクエストIDを参照できるようにcontextにセット
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), service.ReqIdContextKey, reqId))

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
		slog.String("requestId", reqId),
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

	if statusCode >= 100 && statusCode < 200 {
		l.accessLogger.Info("informational", accessArgs...)
	} else if statusCode >= 200 && statusCode < 300 {
		l.accessLogger.Info("access", accessArgs...)
	} else if statusCode >= 300 && statusCode < 400 {
		l.accessLogger.Info("redirection", accessArgs...)
	} else if statusCode >= 400 && statusCode < 500 {
		l.accessLogger.Warn("client error", accessArgs...)
	} else if statusCode >= 500 && statusCode < 600 {
		l.accessLogger.Error("server error", accessArgs...)
	} else {
		l.accessLogger.Error("unknown status", accessArgs...)
	}

	if len(c.Errors) > 0 {
		errorAttr := []slog.Attr{
			slog.String("requestId", reqId),
			slog.Time("time", startTime),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", statusCode),
			slog.Duration("duration", duration),
			slog.String("ip", c.Request.RemoteAddr),
			slog.String("userId", userId),
			slog.String("requestBody", string(reqBody)),
		}
		for _, err := range c.Errors {
			errorAttr = append(errorAttr, slog.String("error", err.Error()))
		}
		errorArgs := make([]any, len(errorAttr))
		for i, attr := range errorAttr {
			errorArgs[i] = attr
		}
		l.errorLogger.Error("error", errorArgs...)
	}

	log.Printf("Method: %s, Path: %s, Status: %d, Duration: %s",
		c.Request.Method, c.Request.URL.Path, statusCode, duration)
}
