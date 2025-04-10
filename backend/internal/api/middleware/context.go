package middleware

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
)

type httpReqRespKey struct{}

type HttpReqResp struct {
	W http.ResponseWriter
	R *http.Request
}

func WithHttpContext(ctx huma.Context, next func(huma.Context)) {
	r, w := humachi.Unwrap(ctx)

	newCtx := huma.WithValue(ctx, httpReqRespKey{}, &HttpReqResp{
		W: w,
		R: r,
	})

	next(newCtx)
}

func GetHttpContext(ctx context.Context) (*http.Request, http.ResponseWriter, bool) {
	if httpObj, ok := ctx.Value(httpReqRespKey{}).(*HttpReqResp); ok {
		return httpObj.R, httpObj.W, true
	}
	return nil, nil, false
}
