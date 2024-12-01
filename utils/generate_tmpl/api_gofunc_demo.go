package generate_tmpl

const GenerateApiGoFuncDemo = `
// {{ .ApiNameCapital }}
// @description
// @receiver a
// @parameter c
func (a api) {{ .ApiNameCapital }}(c *gin.Context) {
	switch common.GetApiVersion(a.ctx, c) {
	case ApiVersionV1:
		v1.New(a.ctx).{{ .ApiNameCapital }}(c)
	default:
		common.UseDefaultApi(a.ctx, c, v1.New(a.ctx).{{ .ApiNameCapital }})
	}
}
`
