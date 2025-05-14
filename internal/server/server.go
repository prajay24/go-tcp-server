package server

import (
    "context"
    "fmt"
    "log"
    "net"
    "runtime"
    "sync/atomic"
    "time"

    "go-tcp-server/internal/config"
)

type Server struct {
    cfg         *config.Config
    listener    net.Listener
    connections uint64
    startTime   time.Time
    logger      *log.Logger
}

func New(cfg *config.Config, logger *log.Logger) *Server {
    // allow Go to use all CPUs
    runtime.GOMAXPROCS(runtime.NumCPU())
    return &Server{
        cfg:       cfg,
        startTime: time.Now(),
        logger:    logger,
    }
}

func (s *Server) Run(ctx context.Context) error {
    addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return fmt.Errorf("listen error: %w", err)
    }
    s.listener = ln
    s.logger.Printf("Listening on %s", addr)

    // start stats reporter
    ticker := time.NewTicker(s.cfg.StatsInterval)
    go s.reportStats(ctx, ticker)

    // accept loop
    for {
        select {
        case <-ctx.Done():
            ln.Close()
            ticker.Stop()
            return nil
        default:
        }

        conn, err := ln.Accept()
        if err != nil {
            s.logger.Printf("Accept error: %v", err)
            continue
        }

        atomic.AddUint64(&s.connections, 1)
        if tcp, ok := conn.(*net.TCPConn); ok {
            tcp.SetNoDelay(true)
        }

        go s.handleConn(conn)
    }
}

func (s *Server) handleConn(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    resp := []byte{0x18, 0x00, 0x00, 0x00}

    for {
        if _, err := conn.Read(buf); err != nil {
            return
        }
        if _, err := conn.Write(resp); err != nil {
            return
        }
    }
}

func (s *Server) reportStats(ctx context.Context, ticker *time.Ticker) {
    for {
        select {
        case <-ticker.C:
            total := atomic.LoadUint64(&s.connections)
            elapsed := time.Since(s.startTime).Seconds()
            tps := float64(total) / elapsed
            s.logger.Printf("Stats: connections=%d, TPS=%.2f, uptime=%.0fs",
                total, tps, elapsed)
        case <-ctx.Done():
            return
        }
    }
}
