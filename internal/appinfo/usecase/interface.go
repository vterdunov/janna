package usecase

// AppInfoRepository abstract methods to receive information
// about the application
type AppInfoRepository interface {
	GetAppInfo() (AppInfoResponse, error)
}
