// status used to add and get information
// such as current deploy progress, deploy error messages, etc.

package jobstatus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vterdunov/janna/internal/jobstatus"
	"github.com/vterdunov/janna/internal/virtualmachine"
)

func TestStorage_FindByID(t *testing.T) {
	st := jobstatus.NewStorage()
	task := st.NewTask()
	id := task.ID()

	tests := map[string]struct {
		s    *jobstatus.Storage
		id   string
		want virtualmachine.TaskStatuser
	}{
		"taskExist":    {st, id, task},
		"taskNotFound": {st, "not_exist_id", nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.FindByID(tc.id)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTaskStatus_Str(t *testing.T) {
	st := jobstatus.NewStorage()

	tests := map[string]struct {
		storage *jobstatus.Storage
		keyvals []string
		want    map[string]interface{}
	}{
		"ok": {
			storage: st,
			keyvals: []string{"key", "value"},
			want:    map[string]interface{}{"key": "value"},
		},
		"missing value": {
			storage: st,
			keyvals: []string{"key2"},
			want:    map[string]interface{}{"key2": "(MISSING)"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			task := tc.storage.NewTask()
			taskWithValues := task.Str(tc.keyvals...)
			got := taskWithValues.Get()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTaskStatus_StrArr(t *testing.T) {
	st := jobstatus.NewStorage()

	tests := map[string]struct {
		storage *jobstatus.Storage
		key     string
		arr     []string
		want    map[string]interface{}
	}{
		"ok": {
			storage: st,
			key:     "key",
			arr:     []string{"value1", "value2"},
			want:    map[string]interface{}{"key": []string{"value1", "value2"}},
		},
		"missing arr": {
			storage: st,
			key:     "key",
			arr:     nil,
			want:    map[string]interface{}{"key": []string(nil)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			task := tc.storage.NewTask()
			taskWithValues := task.StrArr(tc.key, tc.arr)
			got := taskWithValues.Get()
			assert.Equal(t, tc.want, got)
		})
	}
}
