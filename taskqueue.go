package cerberus

import (
	"io"
	"time"
)

type Task func()

type TaskQueue interface {
	Enqueue(task Task)
	Dequeue() (task Task, ok bool)
	Peek() (task Task, ok bool)
	RunEvery(d time.Duration) error
	Empty() bool
	Clear()

	io.Closer
}
