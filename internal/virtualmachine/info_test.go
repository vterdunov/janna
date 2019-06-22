package virtualmachine_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vterdunov/janna/internal/virtualmachine"
	"github.com/vterdunov/janna/internal/virtualmachine/mocks"
)

func TestUsecase_VMInfo(t *testing.T) {
	tests := map[string]struct {
		uuid    string
		want    virtualmachine.VMInfoResponse
		wantErr bool
		prepare func(*mocks.VMRepository)
	}{
		"success": {
			uuid:    "ddd",
			want:    virtualmachine.VMInfoResponse{},
			wantErr: false,
			prepare: func(m *mocks.VMRepository) {
				m.On("VMInfo", mock.AnythingOfType("string")).Return(virtualmachine.VMInfoResponse{}, nil)
			},
		},
		"withError": {
			uuid:    "dddd",
			want:    virtualmachine.VMInfoResponse{},
			wantErr: true,
			prepare: func(m *mocks.VMRepository) {
				m.On("VMInfo", mock.AnythingOfType("string")).Return(virtualmachine.VMInfoResponse{}, errors.New("smthg"))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocks.VMRepository{}
			defer m.AssertExpectations(t)

			params := virtualmachine.VMInfoRequest{
				UUID: tc.uuid,
			}

			if tc.prepare != nil {
				tc.prepare(m)
			}

			c := virtualmachine.NewVMInfo(m, params)

			got, err := c.Execute()

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.want, got)
			assert.NoError(t, err)
		})
	}
}
