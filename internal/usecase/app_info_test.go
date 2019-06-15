package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	usecase "github.com/vterdunov/janna/internal/usecase"
)

func TestUsecase_AppInfo(t *testing.T) {
	tests := map[string]struct {
		want    *usecase.AppInfoResponse
		wantErr bool
		prepare func(*AppInfoRepository)
	}{
		"success": {
			want:    &usecase.AppInfoResponse{Commit: "test", BuildTime: "2000-01-01"},
			wantErr: false,
			prepare: func(m *AppInfoRepository) {
				resp := &usecase.AppInfoResponse{
					Commit:    "test",
					BuildTime: "2000-01-01",
				}
				m.On("AppInfo").Return(resp, nil)
			},
		},
		"withError": {
			want:    &usecase.AppInfoResponse{},
			wantErr: true,
			prepare: func(m *AppInfoRepository) {
				m.On("AppInfo").Return(&usecase.AppInfoResponse{}, errors.New("smthg"))
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &AppInfoRepository{}
			defer m.AssertExpectations(t)

			if tc.prepare != nil {
				tc.prepare(m)
			}

			u := usecase.NewUsecase(m, nil)
			got, err := u.AppInfo()

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.want, got)
			assert.NoError(t, err)
		})
	}
}
