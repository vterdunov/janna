package appinfo

// AppInfo is a command that implements a usecase that requests
// information about the program. Such as build time and commit hash.
type AppInfo struct {
	Repository
}

func NewAppInfo(r Repository) *AppInfo {
	return &AppInfo{
		Repository: r,
	}
}

// Execute returns the application information
func (a *AppInfo) Execute() Response {
	return a.GetAppInfo()
}

type Response struct {
	Commit    string
	BuildTime string
}
