package appinfo

// Repository describes abstract methods to receive information about the application
type Repository interface {
	GetAppInfo() Response
}
