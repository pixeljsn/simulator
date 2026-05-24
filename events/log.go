package events

import "fmt"

type Log struct {
	MaxEntries int
	entries    []string
}

func NewLog(max int) *Log {
	return &Log{MaxEntries: max, entries: make([]string, 0, max)}
}

func (l *Log) Add(tick int, format string, args ...any) {
	msg := fmt.Sprintf("[T%03d] %s", tick, fmt.Sprintf(format, args...))
	l.entries = append([]string{msg}, l.entries...)
	if len(l.entries) > l.MaxEntries {
		l.entries = l.entries[:l.MaxEntries]
	}
}

func (l *Log) Entries() []string { return l.entries }
