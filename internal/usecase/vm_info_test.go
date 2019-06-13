package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	usecase "github.com/vterdunov/janna/internal/usecase"
)

func TestUsecase_VMInfo(t *testing.T) {
	tests := map[string]struct {
		uuid      string
		want      usecase.VMInfoResponse
		wantError bool
		prepare   func(*VMWareRepository)
	}{
		"success": {
			uuid:      "ddd",
			want:      usecase.VMInfoResponse{},
			wantError: false,
			prepare: func(m *VMWareRepository) {
				m.On("VMInfo", mock.AnythingOfType("string")).Return(usecase.VMInfoResponse{}, nil)
			},
		},
		"withError": {
			uuid:      "dddd",
			want:      usecase.VMInfoResponse{},
			wantError: true,
			prepare: func(m *VMWareRepository) {
				m.On("VMInfo", mock.AnythingOfType("string")).Return(usecase.VMInfoResponse{}, errors.New("smthg"))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &VMWareRepository{}
			defer m.AssertExpectations(t)

			if tc.prepare != nil {
				tc.prepare(m)
			}

			u := usecase.NewUsecase(nil, m)

			got, err := u.VMInfo(tc.uuid)

			if tc.wantError {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.want, got)
			assert.NoError(t, err)
		})
	}
}
