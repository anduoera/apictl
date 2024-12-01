package dto

type ProjectStruct struct {
	ServiceSelect map[string]ProjectStructServiceSelect
	ApiSelect     map[string]ProjectStructApiSelect
	DtoSelect     map[string]ProjectStructDtoSelect
}
type ProjectStructServiceSelect struct {
	SelectName      string
	Path            string
	ServiceFilePath string
}

type ProjectStructApiSelect struct {
	SelectName     string
	Path           string
	ApiFilePath    string
	ApiVersionPath string
}

type ProjectStructDtoSelect struct {
	SelectName string
	Path       string
}
