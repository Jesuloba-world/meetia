package middleware

import (
	"context"
	"net/http"
)

type httpReqRespKey struct{}

type HttpReqResp struct {
	W http.ResponseWriter
	R *http.Request
}

func WithHttpContextChi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), httpReqRespKey{}, &HttpReqResp{
			W: w,
			R: r,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetHttpContext(ctx context.Context) (*http.Request, http.ResponseWriter, bool) {
	if httpObj, ok := ctx.Value(httpReqRespKey{}).(*HttpReqResp); ok {
		return httpObj.R, httpObj.W, true
	}
	return nil, nil, false
}
