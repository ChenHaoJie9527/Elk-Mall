package logger

import (
	"fmt"
	"math"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var bufPool = buffer.NewPool()

// kvEncoder 本地控制台用：时间、级别、caller、消息，字段写成 key=value，不再甩一段 JSON。
type kvEncoder struct {
	zapcore.EncoderConfig
	fields []zapcore.Field
}

func newKVEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	if cfg.LineEnding == "" {
		cfg.LineEnding = zapcore.DefaultLineEnding
	}
	return &kvEncoder{EncoderConfig: cfg}
}

func (e *kvEncoder) Clone() zapcore.Encoder {
	c := *e
	if len(e.fields) > 0 {
		c.fields = append([]zapcore.Field(nil), e.fields...)
	}
	return &c
}

// arrayBuf 让 *buffer.Buffer 满足 zapcore.PrimitiveArrayEncoder（EncodeTime/Level/Caller 需要）。
type arrayBuf struct{ *buffer.Buffer }

func (a arrayBuf) AppendByteString(v []byte)     { a.AppendString(string(v)) }
func (a arrayBuf) AppendComplex128(v complex128) { a.AppendString(fmt.Sprint(v)) }
func (a arrayBuf) AppendComplex64(v complex64)   { a.AppendString(fmt.Sprint(v)) }
func (a arrayBuf) AppendFloat64(v float64)       { a.AppendFloat(v, 64) }
func (a arrayBuf) AppendFloat32(v float32)       { a.AppendFloat(float64(v), 32) }
func (a arrayBuf) AppendInt(v int)               { a.Buffer.AppendInt(int64(v)) }
func (a arrayBuf) AppendInt64(v int64)           { a.Buffer.AppendInt(v) }
func (a arrayBuf) AppendInt32(v int32)           { a.Buffer.AppendInt(int64(v)) }
func (a arrayBuf) AppendInt16(v int16)           { a.Buffer.AppendInt(int64(v)) }
func (a arrayBuf) AppendInt8(v int8)             { a.Buffer.AppendInt(int64(v)) }
func (a arrayBuf) AppendUint(v uint)             { a.Buffer.AppendUint(uint64(v)) }
func (a arrayBuf) AppendUint64(v uint64)         { a.Buffer.AppendUint(v) }
func (a arrayBuf) AppendUint32(v uint32)         { a.Buffer.AppendUint(uint64(v)) }
func (a arrayBuf) AppendUint16(v uint16)         { a.Buffer.AppendUint(uint64(v)) }
func (a arrayBuf) AppendUint8(v uint8)           { a.Buffer.AppendUint(uint64(v)) }
func (a arrayBuf) AppendUintptr(v uintptr)       { a.Buffer.AppendUint(uint64(v)) }

func (e *kvEncoder) EncodeEntry(ent zapcore.Entry, extra []zapcore.Field) (*buffer.Buffer, error) {
	buf := bufPool.Get()
	out := arrayBuf{buf}
	sep := e.ConsoleSeparator
	if sep == "" {
		sep = "  "
	}
	n := 0
	kv := func(key string, write func()) {
		if n > 0 {
			buf.AppendString(sep)
		}
		n++
		buf.AppendString(key)
		buf.AppendByte('=')
		write()
	}

	if e.TimeKey != "" && e.EncodeTime != nil {
		kv(e.TimeKey, func() { e.EncodeTime(ent.Time, out) })
	}
	if e.LevelKey != "" && e.EncodeLevel != nil {
		kv(e.LevelKey, func() { e.EncodeLevel(ent.Level, out) })
	}
	if ent.Caller.Defined && e.CallerKey != "" {
		kv(e.CallerKey, func() {
			if e.EncodeCaller != nil {
				e.EncodeCaller(ent.Caller, out)
			} else {
				buf.AppendString(ent.Caller.TrimmedPath())
			}
		})
	}
	if e.MessageKey != "" && ent.Message != "" {
		kv(e.MessageKey, func() { appendLogfmtValue(buf, ent.Message) })
	}

	for _, f := range e.fields {
		appendKV(buf, sep, f)
	}
	for _, f := range extra {
		appendKV(buf, sep, f)
	}

	if ent.Stack != "" && e.StacktraceKey != "" {
		buf.AppendByte('\n')
		buf.AppendString(ent.Stack)
	}
	buf.AppendString(e.LineEnding)
	return buf, nil
}

func appendLogfmtValue(buf *buffer.Buffer, s string) {
	if strings.ContainsAny(s, " \t=\"") {
		buf.AppendByte('"')
		buf.AppendString(s)
		buf.AppendByte('"')
		return
	}
	buf.AppendString(s)
}

func appendKV(buf *buffer.Buffer, sep string, f zapcore.Field) {
	if f.Type == zapcore.SkipType {
		return
	}
	buf.AppendString(sep)
	buf.AppendString(f.Key)
	buf.AppendByte('=')
	switch f.Type {
	case zapcore.StringType:
		appendLogfmtValue(buf, f.String)
	case zapcore.ByteStringType:
		if b, ok := f.Interface.([]byte); ok {
			buf.AppendBytes(b)
		}
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
		if prettyConsole && f.Key == "status" {
			buf.AppendString(colorStatus(int(f.Integer)))
			return
		}
		buf.AppendInt(f.Integer)
	case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type, zapcore.UintptrType:
		buf.AppendUint(uint64(f.Integer))
	case zapcore.BoolType:
		buf.AppendBool(f.Integer == 1)
	case zapcore.DurationType:
		buf.AppendString(formatLatency(time.Duration(f.Integer)))
	case zapcore.TimeType:
		t := time.Unix(0, f.Integer)
		if loc, ok := f.Interface.(*time.Location); ok && loc != nil {
			t = t.In(loc)
		}
		buf.AppendString(t.Format("15:04:05.000"))
	case zapcore.Float64Type:
		buf.AppendFloat(math.Float64frombits(uint64(f.Integer)), 64)
	case zapcore.Float32Type:
		buf.AppendFloat(float64(math.Float32frombits(uint32(f.Integer))), 32)
	case zapcore.ErrorType:
		if err, ok := f.Interface.(error); ok && err != nil {
			buf.AppendString(err.Error())
		} else {
			buf.AppendString("<nil>")
		}
	case zapcore.StringerType:
		if s, ok := f.Interface.(fmt.Stringer); ok && s != nil {
			buf.AppendString(s.String())
		}
	default:
		if f.Interface != nil {
			fmt.Fprint(buf, f.Interface)
		} else if f.String != "" {
			buf.AppendString(f.String)
		} else {
			buf.AppendInt(f.Integer)
		}
	}
}

func (e *kvEncoder) add(f zapcore.Field) { e.fields = append(e.fields, f) }

func (e *kvEncoder) AddArray(k string, v zapcore.ArrayMarshaler) error {
	e.add(zap.Array(k, v))
	return nil
}
func (e *kvEncoder) AddObject(k string, v zapcore.ObjectMarshaler) error {
	e.add(zap.Object(k, v))
	return nil
}
func (e *kvEncoder) AddBinary(k string, v []byte)            { e.add(zap.Binary(k, v)) }
func (e *kvEncoder) AddByteString(k string, v []byte)        { e.add(zap.ByteString(k, v)) }
func (e *kvEncoder) AddBool(k string, v bool)                { e.add(zap.Bool(k, v)) }
func (e *kvEncoder) AddComplex128(k string, v complex128)    { e.add(zap.Complex128(k, v)) }
func (e *kvEncoder) AddComplex64(k string, v complex64)      { e.add(zap.Complex64(k, v)) }
func (e *kvEncoder) AddDuration(k string, v time.Duration)   { e.add(zap.Duration(k, v)) }
func (e *kvEncoder) AddFloat64(k string, v float64)          { e.add(zap.Float64(k, v)) }
func (e *kvEncoder) AddFloat32(k string, v float32)          { e.add(zap.Float32(k, v)) }
func (e *kvEncoder) AddInt(k string, v int)                  { e.add(zap.Int(k, v)) }
func (e *kvEncoder) AddInt64(k string, v int64)              { e.add(zap.Int64(k, v)) }
func (e *kvEncoder) AddInt32(k string, v int32)              { e.add(zap.Int32(k, v)) }
func (e *kvEncoder) AddInt16(k string, v int16)              { e.add(zap.Int16(k, v)) }
func (e *kvEncoder) AddInt8(k string, v int8)                { e.add(zap.Int8(k, v)) }
func (e *kvEncoder) AddString(k, v string)                   { e.add(zap.String(k, v)) }
func (e *kvEncoder) AddTime(k string, v time.Time)           { e.add(zap.Time(k, v)) }
func (e *kvEncoder) AddUint(k string, v uint)                { e.add(zap.Uint(k, v)) }
func (e *kvEncoder) AddUint64(k string, v uint64)            { e.add(zap.Uint64(k, v)) }
func (e *kvEncoder) AddUint32(k string, v uint32)            { e.add(zap.Uint32(k, v)) }
func (e *kvEncoder) AddUint16(k string, v uint16)            { e.add(zap.Uint16(k, v)) }
func (e *kvEncoder) AddUint8(k string, v uint8)              { e.add(zap.Uint8(k, v)) }
func (e *kvEncoder) AddUintptr(k string, v uintptr)          { e.add(zap.Uintptr(k, v)) }
func (e *kvEncoder) AddReflected(k string, v any) error {
	e.add(zap.Reflect(k, v))
	return nil
}
func (e *kvEncoder) OpenNamespace(string) {}
