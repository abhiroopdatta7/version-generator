package filetype

type FileType interface {
	WriteVersion(filePath string, version string) error
}
