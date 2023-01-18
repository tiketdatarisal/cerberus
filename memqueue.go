package cerberus

import (
	"context"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"sync"
	"time"
)

type MemQueue struct {
	taskQueue *linkedlistqueue.Queue
	mutex     *sync.Mutex
	context   context.Context
	cancel    context.CancelFunc
}

// NewMemQueue return an instance of MemQueue object.
func NewMemQueue() *MemQueue {
	ctx, cancel := context.WithCancel(context.Background())
	return &MemQueue{
		taskQueue: linkedlistqueue.New(),
		mutex:     &sync.Mutex{},
		context:   ctx,
		cancel:    cancel,
	}
}

// Enqueue enqueue a new Task to MemQueue.
func (q *MemQueue) Enqueue(task Task) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.taskQueue.Enqueue(task)
}

// Dequeue return and remove next Task from MemQueue.
func (q *MemQueue) Dequeue() (task Task, ok bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	res, ok := q.taskQueue.Dequeue()
	if !ok {
		return nil, ok
	}

	task, ok = res.(Task)
	return task, ok
}

// Peek return next Task from MemQueue.
func (q *MemQueue) Peek() (task Task, ok bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	res, ok := q.taskQueue.Peek()
	if !ok {
		return nil, ok
	}

	task, ok = res.(Task)
	return task, ok
}

// RunEvery dequeue and run Task every d duration.
func (q *MemQueue) RunEvery(d time.Duration) error {
	ticker := time.Tick(d)
	for {
		select {
		case <-q.context.Done():
			return q.context.Err()
		case <-ticker:
			task, ok := q.Dequeue()
			if !ok {
				continue
			}

			if task != nil {
				task()
			}
		}
	}
}

// Empty return true when MemQueue is empty.
func (q *MemQueue) Empty() bool {
	return q.taskQueue.Empty()
}

// Clear remove all Task from MemQueue.
func (q *MemQueue) Clear() {
	q.taskQueue.Clear()
}

// Close cancel all pending (not running) Task.
func (q *MemQueue) Close() error {
	if q.cancel != nil {
		q.cancel()
	}

	return nil
}
