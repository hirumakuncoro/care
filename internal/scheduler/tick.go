package scheduler

import (
	"fmt"
	"log"

	"care/internal/clock"
	"care/internal/schedule"
	"care/internal/store"
)

func (s *Scheduler) Tick() {
	now := clock.NowMs()
	today := clock.DateKey(now)

	// Cek pergantian hari: purge & reset
	currentToday := s.state.GetToday()
	if currentToday != today {
		s.state.SetToday(today)
		s.state.ClearReminders()
	}

	// Jika Today baru pertama kali di-set, jangan buat reminder untuk slot yang sudah lewat
	if currentToday == "" {
		s.state.SetToday(today)
	}

	// Tiap slot: cek apakah perlu buat reminder
	for _, slot := range schedule.Slots {
		remID := fmt.Sprintf("%s#%s", today, slot.ID)
		rem := s.state.GetReminder(remID)

		// Belum ada reminder untuk slot ini hari ini
		if rem == nil {
			dueMs := clock.DueMs(today, slot.At.Hour, slot.At.Min)

			// Buat reminder kalau sudah waktunya
			if now >= dueMs {
				rem = &store.Reminder{
					ID:     remID,
					SlotID: slot.ID,
					Date:   today,
					DueMs:  dueMs,
					Status: store.StatusPending,
				}
				s.state.SetReminder(rem)
			} else {
				continue
			}
		}

		// Jika sudah done/missed, skip
		if rem.Status != store.StatusPending {
			continue
		}

		// Cek window: sudah melebihi DueMs + WindowMs?
		if now > rem.DueMs+slot.WindowMs {
			rem.Status = store.StatusMissed
			s.state.SetReminder(rem)
			s.editReminderDone(rem, slot)
			continue
		}

		// Cek retry: sudah waktunya kirim atau kirim ulang?
		if rem.LastSentMs == 0 || (now-rem.LastSentMs >= s.config.RetryEveryMs) {
			// Cek MaxAttempts
			if rem.Attempts >= s.config.MaxAttempts {
				rem.Status = store.StatusMissed
				s.state.SetReminder(rem)
				s.editReminderDone(rem, slot)
				continue
			}

			// Kirim atau kirim ulang
			if err := s.sendReminder(rem, slot); err != nil {
				log.Printf("error sending reminder %s: %v", remID, err)
			}
		}
	}
}
