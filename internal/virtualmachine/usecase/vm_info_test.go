package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	usecase "github.com/vterdunov/janna/internal/virtualmachine/usecase"
)

func TestUsecase_VMInfo(t *testing.T) {
	tests := map[string]struct {
		uuid    string
		want    usecase.VMInfoResponse
		wantErr bool
		prepare func(*VMRepository)
	}{
		"success": {
			uuid:    "ddd",
			want:    usecase.VMInfoResponse{},
			wantErr: false,
			prepare: func(m *VMRepository) {
				m.On("VMInfo", mock.AnythingOfType("string")).Return(usecase.VMInfoResponse{}, nil)
			},
		},
		"withError": {
			uuid:    "dddd",
			want:    usecase.VMInfoResponse{},
			wantErr: true,
			prepare: func(m *VMRepository) {
				m.On("VMInfo", mock.AnythingOfType("string")).Return(usecase.VMInfoResponse{}, errors.New("smthg"))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &VMRepository{}
			defer m.AssertExpectations(t)

			if tc.prepare != nil {
				tc.prepare(m)
			}

			u := usecase.VMInfo{m}

			got, err := u.VMInfo(tc.uuid)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.want, got)
			assert.NoError(t, err)
		})
	}
}
