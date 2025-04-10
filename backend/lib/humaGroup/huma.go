package humagroup

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// huma lacks route group so I created one

type HumaGroup struct {
	basePath    string
	tags        []string
	api         huma.API
	middlewares huma.Middlewares
}

type HumaGroupOptions struct {
	Summary     string
	Description string
}

func NewHumaGroup(
	api huma.API,
	basePath string,
	tags []string,
	middlewares ...func(ctx huma.Context, next func(huma.Context)),
) *HumaGroup {
	return &HumaGroup{
		api:         api,
		basePath:    basePath,
		tags:        tags,
		middlewares: middlewares,
	}
}

func (g *HumaGroup) Use(middlewares ...func(ctx huma.Context, next func(huma.Context))) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func Post[I, O any](
	g *HumaGroup,
	path string,
	handler func(context.Context, *I) (*O, error),
	operationName string,
	options *HumaGroupOptions,
	middlewares ...func(ctx huma.Context, next func(huma.Context)),
) {
	operation := huma.Operation{
		OperationID: operationName,
		Method:      http.MethodPost,
		Path:        g.basePath + path,
		Tags:        g.tags,
		Middlewares: append(g.middlewares, middlewares...),
	}
	if options != nil {
		operation.Summary = options.Summary
		operation.Description = options.Description
	}
	huma.Register(g.api, operation, handler)
}

func Get[I, O any](
	g *HumaGroup,
	path string,
	handler func(context.Context, *I) (*O, error),
	operationName string,
	options *HumaGroupOptions,
	middlewares ...func(ctx huma.Context, next func(huma.Context)),
) {
	operation := huma.Operation{
		OperationID: operationName,
		Method:      http.MethodGet,
		Path:        g.basePath + path,
		Tags:        g.tags,
		Middlewares: append(g.middlewares, middlewares...),
	}
	if options != nil {
		operation.Summary = options.Summary
		operation.Description = options.Description
	}
	huma.Register(g.api, operation, handler)
}

func Patch[I, O any](
	g *HumaGroup,
	path string,
	handler func(context.Context, *I) (*O, error),
	operationName string,
	middlewares ...func(ctx huma.Context, next func(huma.Context)),
) {
	operation := huma.Operation{
		OperationID: operationName,
		Method:      http.MethodPatch,
		Path:        g.basePath + path,
		Tags:        g.tags,
		Middlewares: append(g.middlewares, middlewares...),
	}
	huma.Register(g.api, operation, handler)
}

func Put[I, O any](
	g *HumaGroup,
	path string,
	handler func(context.Context, *I) (*O, error),
	operationName string,
	middlewares ...func(ctx huma.Context, next func(huma.Context)),
) {
	operation := huma.Operation{
		OperationID: operationName,
		Method:      http.MethodPut,
		Path:        g.basePath + path,
		Tags:        g.tags,
		Middlewares: append(g.middlewares, middlewares...),
	}
	huma.Register(g.api, operation, handler)
}

func Delete[I, O any](
	g *HumaGroup,
	path string,
	handler func(context.Context, *I) (*O, error),
	operationName string,
	middlewares ...func(ctx huma.Context, next func(huma.Context)),
) {
	operation := huma.Operation{
		OperationID: operationName,
		Method:      http.MethodDelete,
		Path:        g.basePath + path,
		Tags:        g.tags,
		Middlewares: append(g.middlewares, middlewares...),
	}
	huma.Register(g.api, operation, handler)
}
