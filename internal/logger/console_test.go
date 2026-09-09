package logger

import (
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestFormatLatency(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{500 * time.Microsecond, "500µs"},
		{65*time.Millisecond + 661*time.Microsecond, "65.7ms"},
		{1500 * time.Millisecond, "1.50s"},
	}
	for _, tt := range tests {
		if got := formatLatency(tt.d); got != tt.want {
			t.Errorf("formatLatency(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestKVEncoder_NoJSONBlob(t *testing.T) {
	encCfg := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		MessageKey:       "msg",
		LineEnding:       "\n",
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeTime:       zapcore.TimeEncoderOfLayout("15:04:05.000"),
		ConsoleSeparator: "  ",
	}
	enc := newKVEncoder(encCfg)
	ent := zapcore.Entry{
		Level:   zapcore.InfoLevel,
		Time:    time.Date(2026, 9, 9, 13, 47, 26, 142000000, time.Local),
		Message: "user login",
	}
	out, err := enc.EncodeEntry(ent, []zapcore.Field{
		zap.String("method", "POST"),
		zap.String("uri", "/users/login"),
		zap.String("username", "elk"),
	})
	if err != nil {
		t.Fatal(err)
	}
	line := out.String()
	if strings.Contains(line, "{") {
		t.Fatalf("still looks like JSON: %q", line)
	}
	for _, want := range []string{
		"time=13:47:26.142",
		"level=INFO",
		`msg="user login"`,
		"method=POST",
		"uri=/users/login",
		"username=elk",
	} {
		if !strings.Contains(line, want) {
			t.Fatalf("missing %q in %q", want, line)
		}
	}
}
