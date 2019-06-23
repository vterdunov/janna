package appinfo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vterdunov/janna/internal/appinfo"
	"github.com/vterdunov/janna/internal/appinfo/mocks"
)

func TestUsecase_AppInfo(t *testing.T) {
	tests := map[string]struct {
		want    appinfo.Response
		prepare func(*mocks.Repository)
	}{
		"success": {
			want:    appinfo.Response{Commit: "test", BuildTime: "2000-01-01"},
			prepare: func(m *mocks.Repository) {
				resp := appinfo.Response{
					Commit:    "test",
					BuildTime: "2000-01-01",
				}
				m.On("GetAppInfo").Return(resp, nil)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocks.Repository{}
			defer m.AssertExpectations(t)

			if tc.prepare != nil {
				tc.prepare(m)
			}

			command := appinfo.NewAppInfo(m)
			got := command.Execute()

			assert.Equal(t, tc.want, got)
		})
	}
}
