package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/auth"
	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/config"
)

func main() {
	os.Exit(run())
}

func run() int {
	settings, err := config.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, config.UserFacingError(err))
		return 1
	}

	addr := fmt.Sprintf("127.0.0.1:%d", settings.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Port conflict: %s is already in use. Use --port to select another port or update ~/.sql-proxy/config.yaml.\n",
			addr,
		)
		return 1
	}
	defer listener.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// OAuth credentials are optional in local/test environments; if configured,
	// perform auth bootstrap for downstream dialer initialization.
	if client, err := auth.GetClient(ctx); err == nil {
		_ = client
	} else if !errors.Is(err, auth.ErrMissingCredentials) {
		fmt.Fprintf(os.Stderr, "Authentication unavailable: %v\n", err)
		return 1
	}

	go func() {
		<-ctx.Done()
		fmt.Fprintln(os.Stdout, "Shutting down proxy...")
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) || ctx.Err() != nil {
				break
			}
			fmt.Fprintf(os.Stderr, "Listener error: %v\n", err)
			continue
		}
		_ = conn.Close()
	}

	return 0
}
