package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("テスト")
		if buf.String() != "テスト\n" {
			t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
		}
	}
}