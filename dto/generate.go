package dto

import "html/template"

type GenerateFileValue struct {
	ProjectName            string
	ApiNameLowerCase       string
	ApiNameCapital         string
	StructNameAbbreviation string
	PackName               string
	ServiceName            string
	AngleBracket           template.HTML
}
