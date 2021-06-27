// ログ出力パッケージ
package trace

import (
	"fmt"
	"io"
)

// Tracer はログ出力できるオブジェクトを表すインターフェース
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(args ...interface{}) {
	fmt.Fprint(t.out, args...)
	fmt.Fprint(t.out, "\n")
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
