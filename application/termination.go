package application

import (
	"context"
	"github.com/maxmanuylov/utils"
	"os"
	"os/signal"
)

type TerminationWatcher struct {
	C <-chan os.Signal
	c chan os.Signal
}

func NewTerminationWatcher() *TerminationWatcher {
	return NewTerminationWatcherN(1)
}

func NewTerminationWatcherN(n int) *TerminationWatcher {
	signalsChannel := make(chan os.Signal, n)
	notifyOnTermination(signalsChannel)
	return &TerminationWatcher{
		C: signalsChannel,
		c: signalsChannel,
	}
}

func (tw *TerminationWatcher) Stop() {
	signal.Stop(tw.c)
	close(tw.c)
}

func WaitForTermination() {
	<-NewTerminationWatcher().C
}

func notifyOnTermination(c chan<- os.Signal) {
	signal.Notify(c, sigInt, sigTerm)
}

func Context(ctx context.Context) context.Context {
	tw := NewTerminationWatcher()
	utils.Defer(ctx, tw.Stop)

	appCtx, cancel := context.WithCancel(ctx)

	go func() {
		if _, ok := <-tw.C; ok {
			cancel()
		}
	}()

	return appCtx
}
