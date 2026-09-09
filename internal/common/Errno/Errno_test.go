package errno

import (
	"errors"
	"testing"
)

func TestWithMsg_ReplacesMessageWithoutMutatingOriginal(t *testing.T) {
	origMsg := ParmError.Msg

	got := ParmError.WithMsg("id 必填")

	if got.Code != ParmError.Code {
		t.Errorf("code = %d, want %d", got.Code, ParmError.Code)
	}
	if got.Msg != "id 必填" {
		t.Errorf("msg = %q, want %q", got.Msg, "id 必填")
	}
	if ParmError.Msg != origMsg {
		t.Errorf("mutated global ParmError.Msg = %q, want %q", ParmError.Msg, origMsg)
	}
	if got == ParmError {
		t.Error("WithMsg should return a copy, not the original pointer")
	}
}

func TestWithMsg_DoesNotConcatenateOK(t *testing.T) {
	orig := OK.Msg
	got := OK.WithMsg("ping ok")
	if got.Msg != "ping ok" {
		t.Errorf("msg = %q, want %q", got.Msg, "ping ok")
	}
	if OK.Msg != orig {
		t.Errorf("mutated global OK.Msg = %q, want %q", OK.Msg, orig)
	}
}

func TestWithErrMsg_SetsRawErrorWithoutMutatingOriginal(t *testing.T) {
	origErrMsg := ServerError.ErrMsg
	raw := errors.New("connection refused")

	got := ServerError.WithErrMsg(raw)

	if got.Msg != ServerError.Msg {
		t.Errorf("Msg = %q, want %q", got.Msg, ServerError.Msg)
	}
	if got.ErrMsg != raw.Error() {
		t.Errorf("ErrMsg = %q, want %q", got.ErrMsg, raw.Error())
	}
	if ServerError.ErrMsg != origErrMsg {
		t.Errorf("mutated global ServerError.ErrMsg = %q, want %q", ServerError.ErrMsg, origErrMsg)
	}
}
