package usecase

// AppInfo implements usecase to ask the programm information.
// Such as build time and commit hash
type AppInfo struct {
	AppInfoRepository
}

func NewAppInfo(r AppInfoRepository) *AppInfo {
	return &AppInfo{
		AppInfoRepository: r,
	}
}

func (a *AppInfo) Execute() (AppInfoResponse, error) {
	return a.AppInfo()
}

type AppInfoResponse struct {
	Commit    string
	BuildTime string
}
