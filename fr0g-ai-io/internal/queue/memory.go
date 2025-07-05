package queue

import (
	"fmt"
	"sync"
)

// MemoryQueue implements an in-memory message queue
type MemoryQueue struct {
	messages []*Message
	maxSize  int
	mu       sync.RWMutex
}

// NewMemoryQueue creates a new in-memory queue
func NewMemoryQueue(maxSize int) *MemoryQueue {
	return &MemoryQueue{
		messages: make([]*Message, 0),
		maxSize:  maxSize,
	}
}

// Enqueue adds a message to the queue
func (q *MemoryQueue) Enqueue(message *Message) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.messages) >= q.maxSize {
		return fmt.Errorf("queue is full (max size: %d)", q.maxSize)
	}

	q.messages = append(q.messages, message)
	return nil
}

// Dequeue removes and returns the first message from the queue
func (q *MemoryQueue) Dequeue() (*Message, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.messages) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	message := q.messages[0]
	q.messages = q.messages[1:]
	return message, nil
}

// Size returns the current size of the queue
func (q *MemoryQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.messages)
}

// IsEmpty returns true if the queue is empty
func (q *MemoryQueue) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.messages) == 0
}

// Clear removes all messages from the queue
func (q *MemoryQueue) Clear() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.messages = q.messages[:0]
	return nil
}
