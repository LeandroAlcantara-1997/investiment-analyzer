package customcontext

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	hostKey      contextKey = "host"
	ipKey        contextKey = "ip"
	methodKey    contextKey = "method"
	userAgentKey contextKey = "user-agent"
	headerKey    contextKey = "header"
	paramKey     contextKey = "param"
)

func AddHost(ctx *gin.Context) {
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), hostKey, ctx.Request.Host))
}

func AddIP(ctx *gin.Context) {
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), ipKey, ctx.RemoteIP()))
}

func AddMethod(ctx *gin.Context) {
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), methodKey, ctx.Request.Method))
}

func AddEndpoint(ctx *gin.Context) {
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), methodKey, ctx.Request.URL.Path))
}

func AddUserAgent(ctx *gin.Context) {
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), userAgentKey, ctx.Request.UserAgent()))
}

func AddHeader(ctx *gin.Context) {
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), headerKey, ctx.Request.Header))
}

func AddParams(ctx *gin.Context) {
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), paramKey, ctx.Params))
}

func GetHost(ctx context.Context) string {
	if v, ok := ctx.Value(hostKey).(string); ok {
		return v
	}

	return ""
}

func GetIP(ctx context.Context) string {
	if v, ok := ctx.Value(ipKey).(string); ok {
		return v
	}

	return ""
}

func GetMethod(ctx context.Context) string {
	if v, ok := ctx.Value(methodKey).(string); ok {
		return v
	}

	return ""
}

func GetEndpoint(ctx context.Context) string {
	if v, ok := ctx.Value(methodKey).(string); ok {
		return v
	}

	return ""
}

func GetUserAgent(ctx context.Context) string {
	if v, ok := ctx.Value(userAgentKey).(string); ok {
		return v
	}

	return ""
}

func GetHeader(ctx context.Context) http.Header {
	if v, ok := ctx.Value(headerKey).(http.Header); ok {
		return v
	}

	return nil
}

func GetParam(ctx context.Context) *gin.Params {
	if v, ok := ctx.Value(paramKey).(*gin.Params); ok {
		return v
	}

	return nil
}
