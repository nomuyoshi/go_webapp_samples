package trace

import (
	"bytes"
	"testing"
)

func TestTrace(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	tracer.Trace("テスト")
	if buf.String() != "テスト\n" {
		t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
	}
}
