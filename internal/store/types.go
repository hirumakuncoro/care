package store

import (
	"sync"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusDone    Status = "done"
	StatusMissed  Status = "missed"
)

type Reminder struct {
	ID         string    // "2026-10-06#breakfast"
	SlotID     string
	Date       string    // "2006-01-02" WIB
	DueMs      int64     // epoch milis
	Status     Status
	Attempts   int
	LastSentMs int64     // 0 = belum pernah
	MessageIDs []int     // untuk edit jadi ✅
	AckedMs    int64     // 0 = belum
}

type State struct {
	mu        sync.Mutex
	Today     string               // "2006-01-02" WIB
	Reminders map[string]*Reminder // key = Reminder.ID
	Offset    int                  // offset getUpdates
}

func NewState() *State {
	return &State{
		Reminders: make(map[string]*Reminder),
		Offset:    -1,
	}
}

func (s *State) GetReminder(id string) *Reminder {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Reminders[id]
}

func (s *State) SetReminder(rem *Reminder) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Reminders[rem.ID] = rem
}

func (s *State) AllReminders() []*Reminder {
	s.mu.Lock()
	defer s.mu.Unlock()
	var list []*Reminder
	for _, r := range s.Reminders {
		list = append(list, r)
	}
	return list
}

func (s *State) ClearReminders() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Reminders = make(map[string]*Reminder)
}

func (s *State) SetToday(date string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Today = date
}

func (s *State) GetToday() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Today
}

func (s *State) SetOffset(offset int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Offset = offset
}

func (s *State) GetOffset() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Offset
}