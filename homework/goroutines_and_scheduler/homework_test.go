package main

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	tasks *list.List
}

type SchedulerTask struct {
	Value    Task
	Priority int
}

func NewScheduler() Scheduler {
	return Scheduler{
		tasks: list.New(),
	}
}

func NewSchedulerTask(task Task) SchedulerTask {
	return SchedulerTask{
		Value:    task,
		Priority: task.Priority,
	}
}

func (s *Scheduler) AddSchedulerTask(task SchedulerTask) {
	if s.tasks.Len() == 0 {
		s.tasks.PushFront(task)
		return
	}

	elem := s.tasks.Front()
	for {
		if task.Priority <= elem.Value.(SchedulerTask).Priority {
			s.tasks.InsertBefore(task, elem)
			return
		}

		next := elem.Next()

		if next == nil {
			s.tasks.InsertAfter(task, elem)
			return
		}

		elem = next
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.AddSchedulerTask(NewSchedulerTask(task))
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	elem := s.tasks.Front()
	for {
		if elem.Value.(SchedulerTask).Value.Identifier == taskID {
			task := s.tasks.Remove(elem).(SchedulerTask)
			task.Priority = newPriority
			s.AddSchedulerTask(task)

			return
		}

		elem = elem.Next()
	}
}

func (s *Scheduler) GetTask() Task {
	elem := s.tasks.Back()
	return s.tasks.Remove(elem).(SchedulerTask).Value
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task31 := Task{Identifier: 31, Priority: 30}
	task32 := Task{Identifier: 32, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}
	task6 := Task{Identifier: 6, Priority: 150}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task31)
	scheduler.AddTask(task32)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)
	scheduler.AddTask(task6)

	task = scheduler.GetTask()
	assert.Equal(t, task6, task)

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)

	task = scheduler.GetTask()
	assert.Equal(t, task31, task)

	task = scheduler.GetTask()
	assert.Equal(t, task32, task)
}
