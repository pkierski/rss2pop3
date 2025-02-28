package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pkierski/pop3srv"
)

func main() {
	listenAddrs := multistringFlag{}
	flag.Var(&listenAddrs, "p", "listen address, maybe specified multiple times. If not specified default is ':pop3'")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	srv := pop3srv.NewServer(pop3srv.AllowAllAuthorizer{}, RssMboxProvider{})

	var gr sync.WaitGroup
	startListen := func(addr string) {
		defer gr.Done()
		slog.Info("starting ListenAndServe", "addr", addr)
		err := srv.ListenAndServe(addr)
		if errors.Is(err, pop3srv.ErrServerClosed) {
			slog.Info("ListenAndServe exited normally", "addr", addr)
		} else {
			slog.Error("ListenAndServe exited", "addr", addr, "error", err)
		}
	}

	if len(listenAddrs.values) == 0 {
		listenAddrs.values = append(listenAddrs.values, ":pop3")
	}
	for _, addr := range listenAddrs.values {
		gr.Add(1)
		go startListen(addr)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("POP3 shutdown", "error", err)
	} else {
		slog.Info("server shut down gracefully")
	}

	gr.Wait()
}

type multistringFlag struct {
	values []string
}

func (m *multistringFlag) String() string {
	return fmt.Sprintf("[%v]", strings.Join(m.values, ", "))
}

func (m *multistringFlag) Set(v string) error {
	m.values = append(m.values, v)
	return nil
}
