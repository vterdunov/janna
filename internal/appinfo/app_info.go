package appinfo

// AppInfo is a command that implements a usecase that requests
// information about the program. Such as build time and commit hash.
type AppInfo struct {
	AppInfoRepository
}

func NewAppInfo(r AppInfoRepository) *AppInfo {
	return &AppInfo{
		AppInfoRepository: r,
	}
}

// Execute returns the application information
func (a *AppInfo) Execute() (AppInfoResponse, error) {
	return a.getAppInfo()
}

type AppInfoResponse struct {
	Commit    string
	BuildTime string
}
