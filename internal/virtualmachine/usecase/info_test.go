package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vterdunov/janna/internal/virtualmachine/mocks"
	"github.com/vterdunov/janna/internal/virtualmachine/usecase"
)

func TestUsecase_VMInfo(t *testing.T) {
	tests := map[string]struct {
		uuid    string
		want    usecase.VMInfoResponse
		wantErr bool
		prepare func(*mocks.VMRepository)
	}{
		"success": {
			uuid:    "ddd",
			want:    usecase.VMInfoResponse{},
			wantErr: false,
			prepare: func(m *mocks.VMRepository) {
				m.On("VMInfo", mock.AnythingOfType("string")).Return(usecase.VMInfoResponse{}, nil)
			},
		},
		"withError": {
			uuid:    "dddd",
			want:    usecase.VMInfoResponse{},
			wantErr: true,
			prepare: func(m *mocks.VMRepository) {
				m.On("VMInfo", mock.AnythingOfType("string")).Return(usecase.VMInfoResponse{}, errors.New("smthg"))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocks.VMRepository{}
			defer m.AssertExpectations(t)

			params := usecase.VMInfoRequest{
				UUID: tc.uuid,
			}

			if tc.prepare != nil {
				tc.prepare(m)
			}

			c := usecase.NewVMInfo(m, params)

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
