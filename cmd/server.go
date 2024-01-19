package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dungnh3/skool-mn/config"
	"github.com/dungnh3/skool-mn/internal/repositories/mysql"
	"github.com/dungnh3/skool-mn/internal/services"
	pkgdb "github.com/dungnh3/skool-mn/pkg/db"
	l "github.com/dungnh3/skool-mn/pkg/log"
	"github.com/spf13/cobra"
)

// there will be 2 signals, one is from err chan, e.g. when you can run the job
// 1 is interrupt signal
// this one is used for code that doesn't have health/live check
func handleExitSignal(errChan <-chan error) {
	quit := catchInterruptSignal()
	for {
		select {
		case <-quit:
			return
		case <-errChan:
			return
		}
	}
}

func catchInterruptSignal() <-chan os.Signal {
	// wait for interrupt signal here
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

func newServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "start server",
		RunE: func(*cobra.Command, []string) error {
			var err error
			conf := config.Load()
			db := pkgdb.ConnectMySQL(conf.MySQL)
			repo := mysql.New(db, conf)
			logger := l.New()
			ctx := context.Background()
			svr := services.New(conf, repo)

			errChan := make(chan error)
			// handle server error when running
			// when the server returns, we can get stuck forever waiting for interrupt signal
			go func() {
				if err = svr.Run(); err != nil {
					logger.Error("running application error", l.Error(err))
					errChan <- err
				}
			}()
			handleExitSignal(errChan)

			// graceful shutdown, give 20 seconds to finish ongoing requests
			ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
			defer cancel()
			if err = svr.Close(ctx); err != nil {
				logger.Error("exception error when shutting down server", l.Error(err))
			}
			return nil
		},
	}
}
