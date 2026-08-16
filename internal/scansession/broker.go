// Package scansession implements a short-lived in-memory pub/sub broker that
// pairs a browser tab waiting for a barcode with a phone that scans it. A
// session is created by the desktop tab, shown to the phone as a QR code, and
// discarded a few minutes later — nothing here is persisted to MongoDB since
// none of it needs to survive a server restart.
package scansession

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const sessionTTL = 5 * time.Minute

type session struct {
	subscribers []chan string
	expiresAt   time.Time
}

type Broker struct {
	mu       sync.Mutex
	sessions map[string]*session
}

func NewBroker() *Broker {
	b := &Broker{sessions: make(map[string]*session)}
	go b.reapExpired()
	return b
}

// Create reserves a new random session ID for a desktop tab to hand to a phone via QR code.
func (b *Broker) Create() string {
	id := randomID()

	b.mu.Lock()
	defer b.mu.Unlock()
	b.sessions[id] = &session{expiresAt: time.Now().Add(sessionTTL)}
	return id
}

// Subscribe returns a channel that receives every barcode submitted for this
// session from now on. The returned cleanup func must be called when the
// subscriber (the GraphQL subscription resolver) stops listening.
func (b *Broker) Subscribe(sessionID string) (ch chan string, cleanup func(), ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	s, exists := b.sessions[sessionID]
	if !exists {
		return nil, nil, false
	}

	ch = make(chan string, 1)
	s.subscribers = append(s.subscribers, ch)

	cleanup = func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if s, exists := b.sessions[sessionID]; exists {
			for i, c := range s.subscribers {
				if c == ch {
					s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
					break
				}
			}
		}
		close(ch)
	}
	return ch, cleanup, true
}

// Publish delivers a scanned barcode to every subscriber of the session.
// Returns false if the session doesn't exist (expired or never created).
func (b *Broker) Publish(sessionID, code string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	s, exists := b.sessions[sessionID]
	if !exists {
		return false
	}

	for _, ch := range s.subscribers {
		select {
		case ch <- code:
		default:
			// Subscriber isn't ready for another value yet; drop rather than block.
		}
	}
	return true
}

func (b *Broker) reapExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		b.mu.Lock()
		now := time.Now()
		for id, s := range b.sessions {
			if now.After(s.expiresAt) {
				for _, ch := range s.subscribers {
					close(ch)
				}
				delete(b.sessions, id)
			}
		}
		b.mu.Unlock()
	}
}

func randomID() string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
