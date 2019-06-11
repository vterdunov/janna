package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	usecase "github.com/vterdunov/janna/internal/usecase"
)

func TestUsecase_VMInfo(t *testing.T) {
	tests := []struct {
		name      string
		arg       string
		want      usecase.VMInfoResponse
		wantError bool
		prepare   func(*VMWareRepository)
	}{
		{
			name:      "test",
			arg:       "ddd",
			want:      usecase.VMInfoResponse{},
			wantError: false,
			prepare: func(c *VMWareRepository) {
				c.On("VMInfo", "ddd").Return(usecase.VMInfoResponse{}, nil).Once()
			},
		},
		{
			name:      "test2",
			arg:       "dddd",
			want:      usecase.VMInfoResponse{},
			wantError: true,
			prepare: func(mock *VMWareRepository) {
				mock.On("VMInfo", "dddd").Return(usecase.VMInfoResponse{}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := new(VMWareRepository)

			if tt.prepare != nil {
				tt.prepare(mock)
			}

			u := usecase.NewUsecase(nil, mock)

			// run the tested function
			got, err := u.VMInfo(tt.arg)

			if tt.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got) // notice the go convention: want is first argument, got is the second argument
			}

			// assert that the mocks were called correctly.
			mock.AssertExpectations(t)
		})
	}

	// TODO: switch to table test
	// vmWareRepositoryMock := new(VMWareRepository)
	// vmWareRepositoryMock.On("VMInfo", "ddd").Return(usecase.VMInfoResponse{}, nil)
	// u := usecase.NewUsecase(nil, vmWareRepositoryMock)
	// _, _ = u.VMInfo("ddd")

	// vmWareRepositoryMock.AssertExpectations(t)
}
