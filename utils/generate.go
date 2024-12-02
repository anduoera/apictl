package utils

import (
	"apictl/dto"
	"apictl/utils/generate_tmpl"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func GenerateFile(in dto.GenerateFileValue, savePath, tmplPath string) (err error) {
	if len(in.ProjectName) == 0 ||
		len(in.ApiNameLowerCase) == 0 ||
		len(in.ApiNameCapital) == 0 ||
		len(in.StructNameAbbreviation) == 0 ||
		len(in.PackName) == 0 ||
		len(in.ServiceName) == 0 {
		return fmt.Errorf("generatr file fail input")
	}
	in.AngleBracket = template.HTML("<")
	file, err := os.Create(savePath)
	if err != nil {
		return
	}
	defer file.Close()
	tmpl := template.New("generate")
	tmpl = template.Must(tmpl.Parse(tmplPath))
	if err != nil {
		return
	}
	err = tmpl.Execute(file, in)
	if err != nil {
		return
	}
	err = file.Sync()
	if err != nil {
		return
	}
	return nil
}

func GenerateDto(in dto.ProjectStructDtoSelect, info dto.GenerateFileValue, dto string) (err error) {
	file, err := getFileToString(in.Path)
	if err != nil {
		return
	}
	print(file)
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", file, parser.ParseComments)
	if err != nil {
		return
	}
	newStruct := &ast.TypeSpec{
		Name: ast.NewIdent(info.ApiNameCapital + dto),
		Type: &ast.StructType{
			Struct: token.NoPos,
			Fields: &ast.FieldList{},
		},
	}
	genDecl := &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{newStruct},
	}
	node.Decls = append(node.Decls, genDecl)
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}
	err = ioutil.WriteFile(in.Path, buf.Bytes(), 0644)
	if err != nil {
		return
	}
	return nil
}

func GenerateApi(path string, info dto.GenerateFileValue) (err error) {
	file, err := getFileToString(path)
	if err != nil {
		return
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", file, parser.ParseComments)
	if err != nil {
		return
	}
	ast.Walk(newApiMethodAdder{
		file: node,
		info: info,
	}, node)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

type newApiMethodAdder struct {
	file *ast.File
	info dto.GenerateFileValue
}

func (v newApiMethodAdder) Visit(node ast.Node) ast.Visitor {
	if decl, ok := node.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
		for _, spec := range decl.Specs {
			if tspec, ok := spec.(*ast.TypeSpec); ok {
				if tspec.Name.Name == "Api" {
					newMethod := &ast.Field{
						Names: []*ast.Ident{ast.NewIdent(v.info.ApiNameCapital)},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{ast.NewIdent("c")},
										Type:  ast.NewIdent("*gin.Context"),
									},
								},
							},
						},
					}
					tspec.Type.(*ast.InterfaceType).Methods.List = append(
						tspec.Type.(*ast.InterfaceType).Methods.List,
						newMethod,
					)
				}
			}

		}
	}
	return v
}

func GenerateApiFunc(in dto.ProjectStructApiSelect, info dto.GenerateFileValue) (err error) {
	filePath := in.ApiFilePath

	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	tmpl, err := template.New("method").Parse(generate_tmpl.GenerateApiGoFuncDemo)
	if err != nil {
		return err
	}
	var method strings.Builder
	err = tmpl.Execute(&method, info)
	if err != nil {
		return err
	}
	newContent := string(content) + method.String()

	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GenerateService(path string, info dto.GenerateFileValue) (err error) {
	file, err := getFileToString(path)
	if err != nil {
		return
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", file, parser.ParseComments)
	if err != nil {
		return
	}
	ast.Walk(newServiceMethodAdder{
		file: node,
		info: info,
	}, node)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

type newServiceMethodAdder struct {
	file *ast.File
	info dto.GenerateFileValue
}

func (v newServiceMethodAdder) Visit(node ast.Node) ast.Visitor {
	if decl, ok := node.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
		for _, spec := range decl.Specs {
			if tspec, ok := spec.(*ast.TypeSpec); ok {
				if tspec.Name.Name == "Service" {
					newMethod := &ast.Field{
						Names: []*ast.Ident{ast.NewIdent(v.info.ApiNameCapital)},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{ast.NewIdent("c")},
										Type:  ast.NewIdent("*gin.Context"),
									},
									{
										Names: []*ast.Ident{ast.NewIdent("in")},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("dto"),
											Sel: ast.NewIdent(v.info.ApiNameCapital + "Input"),
										},
									},
								},
							},
							Results: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{ast.NewIdent("errorTag")},
										Type:  ast.NewIdent("bool"),
									},
									{
										Names: []*ast.Ident{ast.NewIdent("throw")},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("dto"),
											Sel: ast.NewIdent("Throw"),
										},
									},
									{
										Names: []*ast.Ident{ast.NewIdent("out")},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("dto"),
											Sel: ast.NewIdent(v.info.ApiNameCapital + "Output"),
										},
									},
								},
							},
						},
					}
					tspec.Type.(*ast.InterfaceType).Methods.List = append(
						tspec.Type.(*ast.InterfaceType).Methods.List,
						newMethod,
					)
				}
			}

		}
	}
	return v
}

func getFileToString(path string) (file string, err error) {
	print(path)
	sourceFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer sourceFile.Close()

	buffer := make([]byte, 1024) // 1KB 的缓冲区
	for {
		bytesRead, readErr := sourceFile.Read(buffer)
		if readErr != nil && readErr != io.EOF {
			return "", readErr
		}

		if bytesRead == 0 {
			print(file)
			return file, nil
		}
		print(string(buffer[:bytesRead]))
		file = file + string(buffer[:bytesRead])
	}
}
