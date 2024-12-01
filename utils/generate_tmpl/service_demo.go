package generate_tmpl

const GenerateServiceDemo = `
// -----------------------------------------------------------------------
// Copyright (c) 2022 Akuvox Corporation and Akubela-Eevee Contributors.
// All rights reserved.
// -----------------------------------------------------------------------

package {{ .ServiceName }}

import (
	"{{ .ProjectName }}/core"
	"{{ .ProjectName }}/model/dto"
	"{{ .ProjectName }}/utils"
	"github.com/gin-gonic/gin"
)

type {{ .ApiNameLowerCase }}Service struct {
	ctx   *core.Context
	gin   *gin.Context
	throw dto.Throw
	in    dto.{{ .ApiNameCapital }}Input
	out   dto.{{ .ApiNameCapital }}Output
}

// new{{ .ApiNameCapital }}Service
// @description
// @parameter ctx
// @parameter c
// @parameter in
// @return *{{ .ApiNameLowerCase }}Service
func new{{ .ApiNameCapital }}Service(ctx *core.Context, gin *gin.Context, in dto.{{ .ApiNameCapital }}Input) *{{ .ApiNameLowerCase }}Service {
	return &{{ .ApiNameLowerCase }}Service{
		ctx: ctx,
		gin: gin,
		in:  in,
	}
}

// {{ .ApiNameLowerCase }}
// @description
// @receiver {{ .StructNameAbbreviation }}
// @return errorTag
func ({{ .StructNameAbbreviation }} *{{ .ApiNameLowerCase }}Service) {{ .ApiNameLowerCase }}() (errorTag bool) {
	taskList := make([]func() bool, 0)
	for i := 0; i {{ .AngleBracket }} len(taskList); i++ {
		{{ .StructNameAbbreviation }}.ctx.Logger.Info(utils.GetTraceTag({{ .StructNameAbbreviation }}.gin) + utils.TraceFuncName(taskList[i]))
		if taskList[i]() {
			return true
		}
	}
	return false
}

// {{ .ApiNameCapital }}
// @description
// @receiver s
// @parameter ctx
// @parameter in
// @return errorTag
// @return throw
// @return out
func (s *service) {{ .ApiNameCapital }}(c *gin.Context, in dto.{{ .ApiNameCapital }}Input) (errorTag bool, throw dto.Throw, out dto.{{ .ApiNameCapital }}Output) {
	handler := new{{ .ApiNameCapital }}Service(s.ctx, c, in)
	return handler.{{ .ApiNameLowerCase }}(), handler.throw, handler.out
}
`
