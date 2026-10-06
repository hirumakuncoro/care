package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"care/internal/clock"
	"care/internal/config"
	"care/internal/handler"
	"care/internal/scheduler"
	"care/internal/telegram"
)

func main() {
	// Load config dari env
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	log.Printf("Bot starting. ChatID: %d", cfg.ChatID)

	// Init Telegram client
	client := telegram.NewClient(cfg.Token)

	// Init scheduler
	sch := scheduler.New(cfg, client)

	// Init handler
	h := handler.New(cfg, client, sch)

	// Boot: buang update lama (offset=-1)
	_, err = client.GetUpdates(-1, 0)
	if err != nil {
		log.Printf("warning: boot getUpdates failed: %v", err)
	}
	sch.State().SetOffset(0)

	// Boot: set Today
	now := clock.NowMs()
	today := clock.DateKey(now)
	sch.State().SetToday(today)
	log.Printf("Boot complete. Today: %s", today)

	// Context untuk graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		log.Printf("Signal received: %v", sig)
		cancel()
	}()

	// Ticker (setiap 30 detik)
	ticker := time.NewTicker(time.Duration(cfg.TickMs) * time.Millisecond)
	defer ticker.Stop()

	// Goroutine 1: tick loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				sch.Tick()
			}
		}
	}()

	// Goroutine 2: poll loop (long-poll, blocking)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				h.Poll()
			}
		}
	}()

	// Main goroutine: tunggu signal / context cancel
	<-ctx.Done()
	log.Println("Shutting down...")
	time.Sleep(1 * time.Second) // grace period
	log.Println("Bot stopped")
}
