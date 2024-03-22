package goutils

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseLevel(t *testing.T) {
	lv := parseLevel("debug")
	assert.True(t, lv == slog.LevelDebug)

	lv = parseLevel("info")
	assert.Equal(t, lv, slog.LevelInfo)

	lv = parseLevel("warn")
	assert.Equal(t, lv, slog.LevelWarn)

	lv = parseLevel("error")
	assert.Equal(t, lv, slog.LevelError)

	lv = parseLevel("Error")
	assert.Equal(t, lv, slog.LevelError)

	lv = parseLevel("other")
	assert.Equal(t, lv, slog.LevelInfo)

	t.Logf("Test_parseLevel OK")
}

func Test_InitLogger2(t *testing.T) {
	type Obj struct {
		A int    `json:"a"`
		B string `json:"b"`
	}

	obj := &Obj{A: 100, B: "abc"}

	InitLogger2("app", "v1.0.1", "debug", true, "./")

	InitLogger2("app", "v1.0.1", "debug", false, "./")

	slog.LogAttrs(context.Background(), slog.LevelError, "test", Err(ErrDuplicateMsgCtx), slog.Any("obj", obj))

	t.Logf("Test_InitLogger2 OK")
}
