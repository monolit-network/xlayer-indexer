package ctxlog

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey string

const (
	RequestIDKey ctxKey = "requestID"
	UserIDKey    ctxKey = "userID"
	UserInfoKey  ctxKey = "userInfo"
	IPKey        ctxKey = "ip"
	APIKeyKey    ctxKey = "apiKey"

	BillingInfoKey ctxKey = "billingInfo"
)

func Info(ctx context.Context, logger *zap.Logger, msg string, fields ...zap.Field) {
	fields = populateFields(ctx, fields...)
	logger.Info(msg, fields...)
}

func Error(ctx context.Context, logger *zap.Logger, msg string, fields ...zap.Field) {
	fields = populateFields(ctx, fields...)
	logger.Error(msg, fields...)
}

func Warn(ctx context.Context, logger *zap.Logger, msg string, fields ...zap.Field) {
	fields = populateFields(ctx, fields...)
	logger.Warn(msg, fields...)
}

func Debug(ctx context.Context, logger *zap.Logger, msg string, fields ...zap.Field) {
	fields = populateFields(ctx, fields...)
	logger.Debug(msg, fields...)
}

func ExtractRequestID(ctx context.Context) string {
	reqID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return reqID
}

func ExtractUserID(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func ExtractIP(ctx context.Context) string {
	ip, ok := ctx.Value(IPKey).(string)
	if !ok {
		return ""
	}
	return ip
}

func ExtractAPIKey(ctx context.Context) string {
	apiKey, ok := ctx.Value(APIKeyKey).(string)
	if !ok {
		return ""
	}
	return apiKey
}

func populateFields(ctx context.Context, fields ...zap.Field) []zap.Field {
	if requestID := ExtractRequestID(ctx); requestID != "" {
		fields = append(fields, zap.String("requestID", requestID))
	}
	if userID := ExtractUserID(ctx); userID != "" {
		fields = append(fields, zap.String("userID", userID))
	}
	if ip := ExtractIP(ctx); ip != "" {
		fields = append(fields, zap.String("ip", ip))
	}
	// Api Keys should not be logged in production
	// if apiKey := ExtractAPIKey(ctx); apiKey != "" {
	// 	fields = append(fields, zap.String("apiKey", apiKey))
	// }
	return fields
}
