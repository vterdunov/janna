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
		"taskExist": {st, id, task},
		"taskNotFound": {st, "empty_id", nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.FindByID(tc.id)
			assert.Equal(t, tc.want, got)
		})
	}
}

// func TestTaskStatus_ID(t *testing.T) {
// 	type fields struct {
// 		RWMutex    sync.RWMutex
// 		id         string
// 		Status     map[string]interface{}
// 		Created    time.Time
// 		expiration int64
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t := &TaskStatus{
// 				RWMutex:    tt.fields.RWMutex,
// 				id:         tt.fields.id,
// 				Status:     tt.fields.Status,
// 				Created:    tt.fields.Created,
// 				expiration: tt.fields.expiration,
// 			}
// 			if got := t.ID(); got != tt.want {
// 				t.Errorf("TaskStatus.ID() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestTaskStatus_Str(t *testing.T) {
// 	type fields struct {
// 		RWMutex    sync.RWMutex
// 		id         string
// 		Status     map[string]interface{}
// 		Created    time.Time
// 		expiration int64
// 	}
// 	type args struct {
// 		keyvals []string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   virtualmachine.TaskStatuser
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t := &TaskStatus{
// 				RWMutex:    tt.fields.RWMutex,
// 				id:         tt.fields.id,
// 				Status:     tt.fields.Status,
// 				Created:    tt.fields.Created,
// 				expiration: tt.fields.expiration,
// 			}
// 			if got := t.Str(tt.args.keyvals...); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("TaskStatus.Str() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestTaskStatus_StrArr(t *testing.T) {
// 	type fields struct {
// 		RWMutex    sync.RWMutex
// 		id         string
// 		Status     map[string]interface{}
// 		Created    time.Time
// 		expiration int64
// 	}
// 	type args struct {
// 		key string
// 		arr []string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   virtualmachine.TaskStatuser
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t := &TaskStatus{
// 				RWMutex:    tt.fields.RWMutex,
// 				id:         tt.fields.id,
// 				Status:     tt.fields.Status,
// 				Created:    tt.fields.Created,
// 				expiration: tt.fields.expiration,
// 			}
// 			if got := t.StrArr(tt.args.key, tt.args.arr); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("TaskStatus.StrArr() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestTaskStatus_Get(t *testing.T) {
// 	type fields struct {
// 		RWMutex    sync.RWMutex
// 		id         string
// 		Status     map[string]interface{}
// 		Created    time.Time
// 		expiration int64
// 	}
// 	tests := []struct {
// 		name         string
// 		fields       fields
// 		wantStatuses map[string]interface{}
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t := &TaskStatus{
// 				RWMutex:    tt.fields.RWMutex,
// 				id:         tt.fields.id,
// 				Status:     tt.fields.Status,
// 				Created:    tt.fields.Created,
// 				expiration: tt.fields.expiration,
// 			}
// 			if gotStatuses := t.Get(); !reflect.DeepEqual(gotStatuses, tt.wantStatuses) {
// 				t.Errorf("TaskStatus.Get() = %v, want %v", gotStatuses, tt.wantStatuses)
// 			}
// 		})
// 	}
// }

// func TestStorage_gc(t *testing.T) {
// 	type fields struct {
// 		RWMutex           sync.RWMutex
// 		cleanInterval     time.Duration
// 		defaultExpiration time.Duration
// 		tasks             map[string]*TaskStatus
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &Storage{
// 				RWMutex:           tt.fields.RWMutex,
// 				cleanInterval:     tt.fields.cleanInterval,
// 				defaultExpiration: tt.fields.defaultExpiration,
// 				tasks:             tt.fields.tasks,
// 			}
// 			s.gc()
// 		})
// 	}
// }
