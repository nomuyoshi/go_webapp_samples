// ログ出力パッケージ
package trace

import (
	"fmt"
	"io"
)

// Tracer はログ出力できるオブジェクトを表すインターフェース
type Tracer struct {
	out io.Writer
}

func (t Tracer) Trace(args ...interface{}) {
	// ゼロ値の場合何もしない
	if t.out == nil {
		return
	}
	fmt.Fprintln(t.out, args...)
}

func New(w io.Writer) Tracer {
	return Tracer{out: w}
}
