package atreugo

import (
	"errors"
	"testing"

	"github.com/valyala/fasthttp"
)

func Test_viewToHandler(t *testing.T) {
	called := false
	err := errors.New("error")

	view := func(ctx *RequestCtx) error {
		called = true

		return err
	}

	ctx := new(fasthttp.RequestCtx)

	handler := viewToHandler(view, defaultErrorView)
	handler(ctx)

	if !called {
		t.Error("View is not called")
	}

	if ctx.Response.StatusCode() != fasthttp.StatusInternalServerError {
		t.Errorf("Status code == %d, want %d", ctx.Response.StatusCode(), fasthttp.StatusInternalServerError)
	}

	if string(ctx.Response.Body()) != err.Error() {
		t.Errorf("Response body == %s, want %s", ctx.Response.Body(), err.Error())
	}
}

func Test_isEqual(t *testing.T) {
	v1 := func() {} // nolint:ifshort
	v2 := func() {} // nolint:ifshort

	if !isEqual(v1, v1) {
		t.Errorf("Values are equals")
	}

	if isEqual(v1, v2) {
		t.Errorf("Values are not equals")
	}
}

func Test_middlewaresInclude(t *testing.T) {
	fnIncluded := func(ctx *RequestCtx) error { return nil }
	fnNotIncluded := func(ctx *RequestCtx) error { return nil }
	ms := []Middleware{fnIncluded}

	if !middlewaresInclude(ms, fnIncluded) {
		t.Errorf("The middleware '%p' is included in '%p'", fnIncluded, ms)
	}

	if middlewaresInclude(ms, fnNotIncluded) {
		t.Errorf("The middleware '%p' is not included in '%p'", fnNotIncluded, ms)
	}
}

func Test_appendMiddlewares(t *testing.T) {
	fn := func(ctx *RequestCtx) error { return nil }
	fnSkip := func(ctx *RequestCtx) error { return nil }

	dst := []Middleware{}
	src := []Middleware{fn, fnSkip}
	skip := []Middleware{fnSkip}

	dst = appendMiddlewares(dst, src, skip...)

	if middlewaresInclude(dst, fnSkip) {
		t.Errorf("The middleware '%p' must not be appended in '%p'", fnSkip, dst)
	}

	if !middlewaresInclude(dst, fn) {
		t.Errorf("The middleware '%p' must be appended in '%p'", fn, dst)
	}
}
