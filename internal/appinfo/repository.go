package appinfo

import (
	"github.com/vterdunov/janna/internal/version"
)

type AppRepository struct{}

func NewAppRepository() AppRepository {
	return AppRepository{}
}

func (a AppRepository) GetAppInfo() Response {
	buildTime, commit := version.GetBuildInfo()

	appInfo := Response{
		BuildTime: buildTime,
		Commit:    commit,
	}

	return appInfo
}
