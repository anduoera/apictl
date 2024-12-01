package logic

import (
	"apictl/core"
	"apictl/dto"
	"apictl/utils"
	"apictl/utils/generate_tmpl"
	"github.com/pterm/pterm"
	"os"
	"reflect"
	"runtime"
	"strings"
)

type addCommand struct {
	ctx               *core.Context
	args              []string
	project           dto.ProjectStruct
	apiName           string
	projectName       string
	api               dto.ProjectStructApiSelect
	service           dto.ProjectStructServiceSelect
	dto               dto.ProjectStructDtoSelect
	generateFileInput dto.GenerateFileValue
}

func init() {
	AddCommandTag(AddCommandUse, newAddCommand)
}

func newAddCommand(ctx *core.Context, args []string) Command {
	return &addCommand{
		ctx:  ctx,
		args: args,
	}
}

func (a *addCommand) Run() {
	taskList := make([]func() bool, 0)
	taskList = append(taskList, a.checkProjectStructTask)
	taskList = append(taskList, a.getUserInputApiNameTask)
	taskList = append(taskList, a.getUserSelectApiTask)
	taskList = append(taskList, a.getUserSelectServiceTask)
	taskList = append(taskList, a.getUserSelectDtoTask)
	taskList = append(taskList, a.formatParamTask)
	taskList = append(taskList, a.generateFileTask)
	taskList = append(taskList, a.generateParamTask)
	taskList = append(taskList, a.generateApiTask)
	taskList = append(taskList, a.generateServiceTask)

	for i := 0; i < len(taskList); i++ {
		a.ctx.Logger(core.DebugLevel, "do task:", runtime.FuncForPC(reflect.ValueOf(taskList[i]).Pointer()).Name())
		if taskList[i]() {
			return
		}
	}
}

func (a *addCommand) checkProjectStructTask() (errorTag bool) {
	cwd, err := os.Getwd()
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "get project struct error:", err.Error())
		return true
	}
	a.ctx.Logger(core.InfoLevel, "project struct:", cwd)
	err, a.project = utils.CheckDirectoryStructure(cwd)
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "check project struct error:", err.Error())
		utils.SuccessProjectStructPrint()
		return true
	}
	split := strings.Split(cwd, "\\")
	a.projectName = split[len(split)-1]
	return false
}

func (a *addCommand) getUserInputApiNameTask() (errorTag bool) {
	apiName := pterm.DefaultInteractiveTextInput
	apiName.DefaultText = "Please Enter your api name[example:GetList]"
	show, err := apiName.Show()
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "get api name error:", err.Error())
		return true
	}
	show = strings.TrimSpace(show)
	err = utils.ContainsOnlyLetters(show)
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "get api name error:", err.Error())
		return true
	}
	a.apiName = show
	return false
}

func (a *addCommand) getUserSelectApiTask() (errorTag bool) {
	apisOption := make([]string, 0)
	for _, v := range a.project.ApiSelect {
		apisOption = append(apisOption, v.SelectName)
	}
	selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(apisOption).Show()
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "show api select error:", err.Error())
		return true
	}
	a.ctx.Logger(core.DebugLevel, "Selected api option: ", pterm.Green(selectedOption))
	a.api = a.project.ApiSelect[selectedOption]
	return false
}

func (a *addCommand) getUserSelectServiceTask() (errorTag bool) {
	servicesOption := make([]string, 0)
	for _, v := range a.project.ServiceSelect {
		servicesOption = append(servicesOption, v.SelectName)
	}
	selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(servicesOption).Show()
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "show services select error:", err.Error())
		return true
	}
	a.ctx.Logger(core.DebugLevel, "Selected services option: ", pterm.Green(selectedOption))
	a.service = a.project.ServiceSelect[selectedOption]
	return false
}

func (a *addCommand) getUserSelectDtoTask() (errorTag bool) {
	dtosOption := make([]string, 0)
	for _, v := range a.project.DtoSelect {
		dtosOption = append(dtosOption, v.SelectName)
	}
	selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(dtosOption).Show()
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "show dto select error:", err.Error())
		return true
	}
	a.ctx.Logger(core.DebugLevel, "Selected dto option: ", pterm.Green(selectedOption))
	a.dto = a.project.DtoSelect[selectedOption]
	return false
}

func (a *addCommand) formatParamTask() (errorTag bool) {
	a.generateFileInput = dto.GenerateFileValue{
		ProjectName:            a.projectName,
		ApiNameLowerCase:       strings.ToLower(string(a.apiName[0])) + a.apiName[1:],
		ApiNameCapital:         a.apiName,
		StructNameAbbreviation: strings.ToLower(string(a.apiName[0])),
		PackName:               a.service.SelectName,
		ServiceName:            a.service.SelectName,
	}
	return false
}

func (a *addCommand) generateFileTask() (errorTag bool) {
	err := utils.GenerateFile(a.generateFileInput, utils.NameToFileNamePath(a.api.Path, a.apiName), generate_tmpl.GenerateApiDemo)
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "generate api file error:", err.Error())
		return true
	}

	err = utils.GenerateFile(a.generateFileInput, utils.NameToFileNamePath(a.service.Path, a.apiName), generate_tmpl.GenerateServiceDemo)
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "generate service file error:", err.Error())
		return true
	}

	return false
}

func (a *addCommand) generateParamTask() (errorTag bool) {
	err := utils.GenerateDto(a.dto, a.generateFileInput, "Input")
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "generate dto param error:", err.Error())
		return true
	}
	err = utils.GenerateDto(a.dto, a.generateFileInput, "Output")
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "generate dto param error:", err.Error())
		return true
	}
	return false
}

func (a *addCommand) generateApiTask() (errorTag bool) {
	err := utils.GenerateApi(a.api.ApiFilePath, a.generateFileInput)
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "generate api interface error:", err.Error())
		return true
	}
	err = utils.GenerateApi(a.api.ApiVersionPath, a.generateFileInput)
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "generate version api interface error:", err.Error())
		return true
	}
	err = utils.GenerateApiFunc(a.api, a.generateFileInput)
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "generate api interface func error:", err.Error())
		return true
	}
	return false
}

func (a *addCommand) generateServiceTask() (errorTag bool) {
	err := utils.GenerateService(a.service.ServiceFilePath, a.generateFileInput)
	if err != nil {
		a.ctx.Logger(core.ErrorLevel, "generate service interface func error:", err.Error())
		return true
	}
	return false
}
