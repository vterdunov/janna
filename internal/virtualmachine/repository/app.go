package repository

import (
	"github.com/vterdunov/janna/internal/version"
	"github.com/vterdunov/janna/internal/virtualmachine/usecase"
)

type AppRepository struct{}

func NewAppRepository() *AppRepository {
	return &AppRepository{}
}

func (a *AppRepository) AppInfo() (usecase.AppInfoResponse, error) {
	buildTime, commit := version.GetBuildInfo()

	appInfo := usecase.AppInfoResponse{
		BuildTime: buildTime,
		Commit:    commit,
	}

	return appInfo, nil
}