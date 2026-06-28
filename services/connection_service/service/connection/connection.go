package connection

import (
	"sync"
	"sync/atomic"

	"so-many-v2/realtime_comments/pkg/logg"
)

const sendBuffer = 16

type Subscriber struct {
	id   uint64
	send chan []byte
}

func (s *Subscriber) Send() <-chan []byte {
	return s.send
}

type ConnectionService struct {
	logger *logg.Logger

	mu     sync.RWMutex
	repo   map[string][]*Subscriber
	nextID atomic.Uint64

	done chan struct{}
	once sync.Once
}

func NewConnectionService(logger *logg.Logger) *ConnectionService {
	return &ConnectionService{
		logger: logger,
		repo:   make(map[string][]*Subscriber),
		done:   make(chan struct{}),
	}
}

func (s *ConnectionService) Add(channel string) *Subscriber {
	sub := &Subscriber{
		id:   s.nextID.Add(1),
		send: make(chan []byte, sendBuffer),
	}

	s.mu.Lock()
	s.repo[channel] = append(s.repo[channel], sub)
	s.mu.Unlock()

	return sub
}

func (s *ConnectionService) Remove(channel string, sub *Subscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()

	subs := s.repo[channel]
	for i, c := range subs {
		if c == sub {
			s.repo[channel] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
	if len(s.repo[channel]) == 0 {
		delete(s.repo, channel)
	}
}

func (s *ConnectionService) Broadcast(channel string, msg []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, sub := range s.repo[channel] {
		select {
		case sub.send <- msg:
		default:
			s.logger.WithField("subscriber", sub.id).
				Warn("slow subscriber, message dropped")
		}
	}
}

func (s *ConnectionService) Done() <-chan struct{} {
	return s.done
}

func (s *ConnectionService) Close() {
	s.once.Do(func() { close(s.done) })
}