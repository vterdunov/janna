package appinfo

import (
	"github.com/vterdunov/janna/internal/version"
)

type AppRepository struct{}

func NewAppRepository() AppRepository {
	return AppRepository{}
}

func (a AppRepository) getAppInfo() (AppInfoResponse, error) {
	buildTime, commit := version.GetBuildInfo()

	appInfo := AppInfoResponse{
		BuildTime: buildTime,
		Commit:    commit,
	}

	return appInfo, nil
}
