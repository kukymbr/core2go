package logtools

import (
	"fmt"

	"go.uber.org/zap"
)

// CatchPanic catches the panic and writes it to the log.
func CatchPanic(log *zap.Logger, onPanic ...func(recovered any)) {
	recovered := recover()
	if recovered == nil {
		return
	}

	log.Error(fmt.Sprintf("%v", recovered), zap.Bool("is_panic", true))

	if len(onPanic) == 0 {
		return
	}

	for _, fn := range onPanic {
		fn(recovered)
	}
}
