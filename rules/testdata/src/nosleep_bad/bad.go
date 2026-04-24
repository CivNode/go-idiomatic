package nosleepbad

import (
	"context"
	"time"
)

// Sleep alongside a channel receive: coordination by timing.
func WaitForWorker(done <-chan struct{}) {
	time.Sleep(10 * time.Millisecond) // want `time.Sleep alongside channel operations`
	<-done
}

// Sleep with ctx in scope: should use ctx.Done.
func WaitForCtx(ctx context.Context) {
	time.Sleep(time.Second) // want `time.Sleep alongside context.Context`
	_ = ctx
}
