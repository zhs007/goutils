package goutils

import "errors"

var (
	// ErrInvalidArrayLength - invalid array length
	ErrInvalidArrayLength = errors.New("invalid array length")
	// ErrInvalidJsonString - invalid json string
	ErrInvalidJsonString = errors.New("invalid json string")
	// ErrInvalidJsonInt - invalid json int
	ErrInvalidJsonInt = errors.New("invalid json int")
	// ErrInvalidJsonBool - invalid json bool
	ErrInvalidJsonBool = errors.New("invalid json bool")
	// ErrInvalidVersion - invalid Version
	ErrInvalidVersion = errors.New("invalid Version")
	// ErrDuplicateMsgCtx - duplicate msg ctx
	ErrDuplicateMsgCtx = errors.New("duplicate msg ctx")
	// ErrNoMsgCtx - no msg ctx
	ErrNoMsgCtx = errors.New("no msg ctx")
	// ErrInvalidMsgName - invalid MsgName
	ErrInvalidMsgName = errors.New("invalid MsgName")
	// ErrNoMsgName - no MsgName
	ErrNoMsgName = errors.New("no MsgName")
)
