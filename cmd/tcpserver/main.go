package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    // replace the module path below with whatever you have in your go.mod
    "go-tcp-server/internal/config"
    "go-tcp-server/internal/server"
)

func main() {
    // 1) Load YAML configuration (config.yaml by default or via -config flag)
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // 2) Set up a structured logger
    logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

    // 3) Create a context that cancels on SIGINT/SIGTERM for graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigCh
        logger.Println("shutdown signal received, terminating...")
        cancel()
    }()

    // 4) Instantiate and run the TCP server
    srv := server.New(cfg, logger)
    if err := srv.Run(ctx); err != nil {
        logger.Fatalf("server error: %v", err)
    }

    logger.Println("server exited cleanly")
}
