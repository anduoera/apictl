package utils

import (
	"apictl/dto"
	"fmt"
	"github.com/pterm/pterm"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	ApiFileName     = "api"
	ServiceFileName = "service"
	ModelFileName   = "model"
)

func CheckDirectoryStructure(path string) (err error, project dto.ProjectStruct) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err, project
	}
	for _, file := range files {
		if file.IsDir() {
			switch file.Name() {
			case ApiFileName:
				err, project.ApiSelect = getApiSelect(filepath.Join(path, file.Name()))
			case ServiceFileName:
				err, project.ServiceSelect = getServiceSelect(filepath.Join(path, file.Name()))
			case ModelFileName:
				err, project.DtoSelect = getDtoSelect(filepath.Join(path, file.Name()))
			}
		}
	}
	if project.ApiSelect == nil || len(project.ApiSelect) == 0 {
		return fmt.Errorf("get project api select error"), project
	}
	if project.ServiceSelect == nil || len(project.ServiceSelect) == 0 {
		return fmt.Errorf("get project service select error"), project
	}
	if project.DtoSelect == nil || len(project.DtoSelect) == 0 {
		return fmt.Errorf("get project dto select error"), project
	}
	return nil, project
}

func getDtoSelect(path string) (err error, dtos map[string]dto.ProjectStructDtoSelect) {
	dtos = make(map[string]dto.ProjectStructDtoSelect)
	fullPath := filepath.Join(path, "dto")
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return err, dtos
	}
	for _, file := range files {
		subPath := filepath.Join(fullPath, file.Name())
		subPath = strings.ReplaceAll(subPath, "test1", "test2")
		dtos[file.Name()] = dto.ProjectStructDtoSelect{
			SelectName: file.Name(),
			Path:       subPath,
		}
	}
	return
}

func getServiceSelect(path string) (err error, services map[string]dto.ProjectStructServiceSelect) {
	services = make(map[string]dto.ProjectStructServiceSelect)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err, services
	}
	for _, file := range files {
		subPath := filepath.Join(path, file.Name())
		subPath = strings.ReplaceAll(subPath, "test1", "test2")
		if fileExists(filepath.Join(subPath, "service.go")) {
			services[file.Name()] = dto.ProjectStructServiceSelect{
				SelectName:      file.Name(),
				Path:            subPath,
				ServiceFilePath: filepath.Join(subPath, "service.go"),
			}
		}
	}
	return
}

func getApiSelect(path string) (err error, apis map[string]dto.ProjectStructApiSelect) {
	apis = make(map[string]dto.ProjectStructApiSelect)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err, apis
	}
	for _, file := range files {
		subPath := filepath.Join(path, file.Name())
		subPath = strings.ReplaceAll(subPath, "test1", "test2")
		if fileExists(filepath.Join(subPath, "api.go")) &&
			fileExists(filepath.Join(subPath, "v1")) &&
			fileExists(filepath.Join(subPath, "v1/api.go")) {
			apis[file.Name()] = dto.ProjectStructApiSelect{
				SelectName:     file.Name(),
				Path:           filepath.Join(subPath, "v1"),
				ApiFilePath:    filepath.Join(subPath, "api.go"),
				ApiVersionPath: filepath.Join(subPath, "v1/api.go"),
			}
		}

	}
	return
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func SuccessProjectStructPrint() {
	tree := pterm.TreeNode{
		Text: "project",
		Children: []pterm.TreeNode{
			{
				Text: "api-demo",
				Children: []pterm.TreeNode{
					{
						Text: "api",
						Children: []pterm.TreeNode{
							// Grandchildren nodes
							{Text: "api.go"},
							{
								Text: "v1",
								Children: []pterm.TreeNode{
									{Text: "api.go"},
								},
							},
						},
					},
				},
			},
			{
				Text: "service",
				Children: []pterm.TreeNode{
					{
						Text: "service-demo",
						Children: []pterm.TreeNode{
							{Text: "service.go"},
						},
					},
				},
			},
			{
				Text: "model",
				Children: []pterm.TreeNode{
					{
						Text: "dto",
						Children: []pterm.TreeNode{
							{Text: "dto-demo.go"},
						},
					},
				},
			},
		},
	}
	pterm.DefaultTree.WithRoot(tree).Render()
}
