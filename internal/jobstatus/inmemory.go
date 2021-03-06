// status used to add and get information
// such as current deploy progress, deploy error messages, etc.

package jobstatus

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/vterdunov/janna/internal/virtualmachine"
)

// Storage stores information about something
type Storage struct {
	sync.RWMutex
	cleanInterval     time.Duration
	defaultExpiration time.Duration

	tasks map[string]*TaskStatus
}

// TaskStatus keep status messages and other metadata of task
type TaskStatus struct {
	sync.RWMutex
	id         string
	Status     map[string]interface{}
	Created    time.Time
	expiration int64
}

// NewStorage creates a new in-memory storage
func NewStorage() *Storage {
	cleanInterval := time.Second * 10
	expirationTime := time.Hour * 24
	tasks := make(map[string]*TaskStatus)
	s := Storage{
		cleanInterval:     cleanInterval,
		defaultExpiration: expirationTime,
		tasks:             tasks,
	}

	go s.gc()

	return &s
}

// NewTask creates a new unique status for a task
func (s *Storage) NewTask() virtualmachine.TaskStatuser {

	expiration := time.Now().Add(s.defaultExpiration).UnixNano()
	uuid := uuid.New().String()
	status := make(map[string]interface{})
	r := TaskStatus{
		id:         uuid,
		Created:    time.Now(),
		expiration: expiration,
		Status:     status,
	}
	s.Lock()
	s.tasks[uuid] = &r
	s.Unlock()

	return &r
}

func (s *Storage) FindByID(id string) virtualmachine.TaskStatuser {
	for _, task := range s.tasks {
		if task.id == id {
			return task
		}
	}
	return nil
}

// Id returns task Id
func (t *TaskStatus) ID() string {
	return t.id
}

// Str a key-value pairs to a task status message
func (t *TaskStatus) Str(keyvals ...string) virtualmachine.TaskStatuser {
	t.Lock()
	defer t.Unlock()

	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			t.Status[fmt.Sprint(keyvals[i])] = keyvals[i+1]
		} else {
			t.Status[fmt.Sprint(keyvals[i])] = "(MISSING)"
		}
	}

	return t
}

// StrArr a key-value pairs to a task status message
func (t *TaskStatus) StrArr(key string, arr []string) virtualmachine.TaskStatuser {
	t.Lock()
	defer t.Unlock()

	t.Status[key] = arr
	return t
}

// Get status messages from a task
func (t *TaskStatus) Get() (statuses map[string]interface{}) {
	t.Lock()
	defer t.Unlock()

	return t.Status
}

// gc search and clean expired tasks from in-memory storage
func (s *Storage) gc() {
	ticker := time.NewTicker(s.cleanInterval)

	for range ticker.C {
		if s.tasks == nil {
			return
		}

		s.RLock()
		for _, task := range s.tasks {
			isTaskExpired := time.Now().UnixNano() > task.expiration
			if isTaskExpired {
				delete(s.tasks, task.id)
			}
		}
		s.RUnlock()
	}
}
