package usecase

func (u *Usecase) AppInfo() (*AppInfoResponse, error) {
	return u.appInfoRepository.AppInfo()
}

type AppInfoResponse struct {
	Commit    string
	BuildTime string
}
