package scheduler

import (
	"fmt"
	"strings"

	"care/internal/clock"
	"care/internal/config"
	"care/internal/schedule"
	"care/internal/store"
	"care/internal/telegram"
)

type Scheduler struct {
	config *config.Config
	client *telegram.Client
	state  *store.State
}

func New(cfg *config.Config, tg *telegram.Client) *Scheduler {
	return &Scheduler{
		config: cfg,
		client: tg,
		state:  store.NewState(),
	}
}

func (s *Scheduler) State() *store.State {
	return s.state
}

// formatMessage buat text reminder dari slot
func (s *Scheduler) formatMessage(slot *schedule.Slot) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("<b>%s %02d:%02d WIB</b>\n", slot.Title, slot.At.Hour, slot.At.Min))
	for _, item := range slot.Items {
		if item.Note != "" {
			buf.WriteString(fmt.Sprintf("• %s (%s)\n", item.Name, item.Note))
		} else {
			buf.WriteString(fmt.Sprintf("• %s\n", item.Name))
		}
	}
	return buf.String()
}

// makeKeyboard buat inline button - semua slot cuma "Selesai"
func (s *Scheduler) makeKeyboard(remID string) *telegram.InlineKeyboard {
	return &telegram.InlineKeyboard{
		Rows: [][]telegram.InlineButton{
			{
				{Text: "✅ Selesai", CallbackData: fmt.Sprintf("ack:%s", remID)},
			},
		},
	}
}

// sendReminder kirim pesan reminder baru
func (s *Scheduler) sendReminder(rem *store.Reminder, slot *schedule.Slot) error {
	text := s.formatMessage(slot)
	keyboard := s.makeKeyboard(rem.ID)

	msgID, err := s.client.SendMessage(s.config.ChatID, text, keyboard)
	if err != nil {
		return err
	}

	rem.MessageIDs = append(rem.MessageIDs, msgID)
	rem.LastSentMs = clock.NowMs()
	rem.Attempts++
	s.state.SetReminder(rem)

	return nil
}

// editReminderDone edit pesan jadi ✅ setelah di-ack
func (s *Scheduler) editReminderDone(rem *store.Reminder, slot *schedule.Slot) error {
	text := s.formatMessage(slot) + "\n✅ Selesai"

	// Edit pesan terakhir
	if len(rem.MessageIDs) == 0 {
		return nil // tidak ada pesan, skip
	}

	lastMsgID := rem.MessageIDs[len(rem.MessageIDs)-1]
	return s.client.EditMessage(s.config.ChatID, lastMsgID, text)
}
