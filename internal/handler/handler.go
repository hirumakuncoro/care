package handler

import (
	"log"
	"strings"

	"care/internal/clock"
	"care/internal/config"
	"care/internal/scheduler"
	"care/internal/store"
	"care/internal/telegram"
)

type Handler struct {
	config    *config.Config
	client    *telegram.Client
	scheduler *scheduler.Scheduler
}

func New(cfg *config.Config, tg *telegram.Client, sch *scheduler.Scheduler) *Handler {
	return &Handler{
		config:    cfg,
		client:    tg,
		scheduler: sch,
	}
}

// Poll long-poll getUpdates, tangani callback
func (h *Handler) Poll() error {
	offset := h.scheduler.State().GetOffset()

	updates, err := h.client.GetUpdates(offset, 30)
	if err != nil {
		log.Printf("getUpdates error: %v", err)
		return nil // jangan return error, coba lagi
	}

	for _, upd := range updates {
		offset = upd.UpdateID + 1

		// Handle callback query (tombol ✅)
		if upd.CallbackQuery != nil {
			h.handleCallback(upd.CallbackQuery)
		}

		// Handle message text (opsional, untuk logging)
		if upd.Message != nil && upd.Message.From != nil {
			// Filter: hanya penerima yang diizinkan
			if upd.Message.From.ID != h.config.ChatID {
				continue
			}
		}
	}

	h.scheduler.State().SetOffset(offset)
	return nil
}

// handleCallback proses tombol "Selesai"
func (h *Handler) handleCallback(cb *telegram.CallbackQuery) {
	// Filter: hanya ChatID yang diizinkan
	if cb.From.ID != h.config.ChatID {
		h.client.AnswerCallback(cb.ID, "❌ Tidak diizinkan")
		return
	}

	// Parse callback data: "ack:<remID>"
	if !strings.HasPrefix(cb.Data, "ack:") {
		return
	}

	remID := strings.TrimPrefix(cb.Data, "ack:")
	rem := h.scheduler.State().GetReminder(remID)
	if rem == nil {
		h.client.AnswerCallback(cb.ID, "❌ Reminder tidak ditemukan")
		return
	}

	// Hanya terima jika masih pending
	if rem.Status != store.StatusPending {
		h.client.AnswerCallback(cb.ID, "✅ Sudah dikonfirmasi sebelumnya")
		return
	}

	// Mark as done
	rem.Status = store.StatusDone
	rem.AckedMs = h.getCurrentMs()
	h.scheduler.State().SetReminder(rem)

	// Answer callback query (notif di Telegram)
	h.client.AnswerCallback(cb.ID, "✅ Tercatat")

	// Edit pesan jadi ✅ (opsional, tapi bagus untuk visual)
	if len(rem.MessageIDs) > 0 {
		lastMsgID := rem.MessageIDs[len(rem.MessageIDs)-1]
		if cb.Message != nil {
			// Edit pesan dengan status done
			text := cb.Message.Text + "\n\n✅ <i>Selesai</i>"
			h.client.EditMessage(h.config.ChatID, lastMsgID, text)
		}
	}
}

func (h *Handler) getCurrentMs() int64 {
	return clock.NowMs()
}