package watchdog

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Start launches the parent-process watchdog and signal handler.
// When the parent process dies or a termination signal is received,
// the returned context's cancel function is called.
func Start(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	// Signal handler
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	// Platform-specific parent monitor
	go monitorParent(ctx, cancel)

	return ctx, cancel
}
