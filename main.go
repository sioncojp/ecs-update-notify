package ecsupdatenotify

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sioncojp/ecs-update-notify/internal/pidfile"
	"go.uber.org/zap"
)

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	s := logger.Sugar()
	log = Logger{sugar: s}
}

// Run ... run ecs-update-notify
func Run(file string) int {
	exitCh := make(chan int)

	// create PIDFILE and defer remove
	if err := pidfile.Create(pid); err != nil {
		log.sugar.Fatalf("failed to remove the pidfile: %s: %s", pid, err)
	}
	defer pidfile.Remove(pid)

	ctx, cancel := context.WithCancel(context.Background())

	// start ecs-update-notify
	go Start(ctx, file, exitCh)

	// receive syscall
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go signalHandler(sig, cancel)

	return <-exitCh
}

// Start ... check and notify update of ecs service
func Start(ctx context.Context, file string, exitCh chan int) {
	log.sugar.Info("start ecs-update-notify...")

	// load config
	config, err := LoadToml(file)
	if err != nil {
		log.sugar.Fatalf("failed to load toml file: %s\n", file)
	}

	// check update
	go func() {
		for {
			config.CheckUpdate()
			time.Sleep(time.Duration(config.Interval) * time.Second)
		}
	}()

	// receive context
	select {
	case <-ctx.Done():
		log.sugar.Info("received done, exiting in 500 milliseconds")
		time.Sleep(500 * time.Millisecond)
		exitCh <- 0
		return
	}
}

// signalHandler ... Receive signal handler and do context.cancel
func signalHandler(sig chan os.Signal, cancel context.CancelFunc) {
	for {
		select {
		case s := <-sig:
			switch s {
			case syscall.SIGINT:
				log.sugar.Info("received SIGINT signal")
				log.sugar.Info("shutdown...")
				cancel()
			case syscall.SIGTERM:
				log.sugar.Info("received SIGTERM signal")
				log.sugar.Info("shutdown...")
				cancel()
			}
		}
	}
}
