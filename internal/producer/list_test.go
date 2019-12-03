package producer_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vterdunov/janna/internal/virtualmachine"
	"github.com/vterdunov/janna/internal/virtualmachine/mocks"
)

func TestNewVMList(t *testing.T) {
	tests := map[string]struct {
		want    []virtualmachine.VMListResponse
		wantErr bool
		prepare func(*mocks.VMRepository)
	}{
		"success": {
			want:    []virtualmachine.VMListResponse{},
			wantErr: false,
			prepare: func(m *mocks.VMRepository) {
				params := virtualmachine.VMListRequest{
					Datacenter:   "DC1",
					Folder:       "folder",
					ResourcePool: "rp",
				}

				m.On("VMList", params).Return([]virtualmachine.VMListResponse{}, nil)
			},
		},
		"withError": {
			want:    []virtualmachine.VMListResponse{},
			wantErr: true,
			prepare: func(m *mocks.VMRepository) {
				params := virtualmachine.VMListRequest{
					Datacenter:   "DC1",
					Folder:       "folder",
					ResourcePool: "rp",
				}

				m.On("VMList", params).Return([]virtualmachine.VMListResponse{}, errors.New("smthg"))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocks.VMRepository{}
			defer m.AssertExpectations(t)

			params := virtualmachine.VMListRequest{
				Datacenter:   "DC1",
				Folder:       "folder",
				ResourcePool: "rp",
			}

			if tc.prepare != nil {
				tc.prepare(m)
			}

			c := virtualmachine.NewVMList(m, params)

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
