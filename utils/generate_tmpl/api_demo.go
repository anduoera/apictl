package generate_tmpl

const GenerateApiDemo = `
// -----------------------------------------------------------------------
// Copyright (c) 2022 Akuvox Corporation and Akubela-Eevee Contributors.
// All rights reserved.
// -----------------------------------------------------------------------

package v1

import (
	"{{ .ProjectName }}/api/common"
	"{{ .ProjectName }}/core"
	"{{ .ProjectName }}/model"
	"{{ .ProjectName }}/model/dto"
	"{{ .ProjectName }}/model/exception"
	"{{ .ProjectName }}/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
)
type {{ .ApiNameLowerCase }}Api struct {
	ctx     *core.Context
	gin     *gin.Context
	service *service
	throw   dto.Throw
	header  {{ .ApiNameCapital }}ReqHeader
	body    {{ .ApiNameCapital }}ReqBody
}

type {{ .ApiNameCapital }}ReqHeader struct {}

type {{ .ApiNameCapital }}ReqBody struct {}

type {{ .ApiNameCapital }}ResBody struct {}

// new{{ .ApiNameCapital }}Api
// @description
// @parameter ctx
// @parameter gin
// @parameter service
// @return *{{ .ApiNameLowerCase }}Api
func new{{ .ApiNameCapital }}Api(ctx *core.Context, gin *gin.Context, service *service) *{{ .ApiNameLowerCase }}Api {
	return &{{ .ApiNameLowerCase }}Api{
		ctx:     ctx,
		gin:     gin,
		service: service,
	}
}

// responseCode
// @description
// @receiver {{ .StructNameAbbreviation }}
// @parameter throw
// @return code
func ({{ .StructNameAbbreviation }} *{{ .ApiNameLowerCase }}Api) responseCode(throw dto.Throw) (code int) {
	switch throw.ErrorCode {
	default:
		code = http.StatusInternalServerError
	}
	return code
}

// check
// @description
// @receiver {{ .StructNameAbbreviation }}
// @return errorTag
func ({{ .StructNameAbbreviation }} *{{ .ApiNameLowerCase }}Api) check() (errorTag bool) {
	err := {{ .StructNameAbbreviation }}.gin.ShouldBindHeader(&{{ .StructNameAbbreviation }}.header)
	if err != nil {
		{{ .StructNameAbbreviation }}.ctx.Logger.Warn(utils.TraceLog({{ .StructNameAbbreviation }}.gin, common.RequestCheckHeaderError+err.Error()))
		{{ .StructNameAbbreviation }}.throw.ErrorCode = exception.ErrorCode{{ .ApiNameCapital }}ApiCheckBindHeader
		{{ .StructNameAbbreviation }}.throw.ErrorMessage = exception.ErrorMessage{{ .ApiNameCapital }}ApiCheckBindHeader
		return true
	}
	err = {{ .StructNameAbbreviation }}.gin.ShouldBindJSON(&{{ .StructNameAbbreviation }}.body)
	if err != nil {
		{{ .StructNameAbbreviation }}.ctx.Logger.Warn(utils.TraceLog({{ .StructNameAbbreviation }}.gin, common.RequestCheckBodyError+err.Error()))
		{{ .StructNameAbbreviation }}.throw.ErrorCode = exception.ErrorCode{{ .ApiNameCapital }}ApiCheckBindBody
		{{ .StructNameAbbreviation }}.throw.ErrorMessage = exception.ErrorMessage{{ .ApiNameCapital }}ApiCheckBindBody
		return true
	}

	return false
}

// do
// @description
// @receiver {{ .StructNameAbbreviation }}
func ({{ .StructNameAbbreviation }} *{{ .ApiNameLowerCase }}Api) do() {
	errorTag, throw, out := {{ .StructNameAbbreviation }}.service.{{ .PackName }}.{{ .ApiNameCapital }}({{ .StructNameAbbreviation }}.gin, dto.{{ .ApiNameCapital }}Input{})
	if errorTag {
		model.ResponseFail({{ .StructNameAbbreviation }}.ctx, {{ .StructNameAbbreviation }}.gin, {{ .StructNameAbbreviation }}.responseCode(throw), throw)
		return
	}
	
	var result {{ .ApiNameCapital }}ResBody
    	if err := copier.Copy(&result, &out); err != nil {
    		{{ .StructNameAbbreviation }}.ctx.Logger.Warn(utils.TraceLog({{ .StructNameAbbreviation }}.gin, "copier copy err, error is: "+err.Error()))
    		{{ .StructNameAbbreviation }}.throw.ErrorCode = exception.ErrorCodeCommandHandler{{ .ApiNameCapital }}ApiCopier
    		{{ .StructNameAbbreviation }}.throw.ErrorMessage = exception.ErrorMessageCommandHandler{{ .ApiNameCapital }}ApiCopier
    		model.ResponseFail({{ .StructNameAbbreviation }}.ctx, {{ .StructNameAbbreviation }}.gin, {{ .StructNameAbbreviation }}.responseCode({{ .StructNameAbbreviation }}.throw), {{ .StructNameAbbreviation }}.throw)
    		return
    	}
	
	model.ResponseSuccess({{ .StructNameAbbreviation }}.ctx, {{ .StructNameAbbreviation }}.gin, http.StatusOK, result)
}

// {{ .ApiNameCapital }}
// @description
// @receiver a
// @parameter c
func (a *api) {{ .ApiNameCapital }}(c *gin.Context) {
	handler := new{{ .ApiNameCapital }}Api(a.ctx, c, a.service)
	if errorTag := handler.check(); errorTag {
		model.ResponseFail(a.ctx, c, http.StatusBadRequest, handler.throw)
		return
	}
	handler.do()
}
`
