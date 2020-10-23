package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/fatih/structtag"
	"golang.org/x/mod/modfile"
)

type tmplVars struct {
	Providers []providerInfo
}

type providerInfo struct {
	Name      string
	Version   string
	Variables []varInfo
	Readme    reamdeInfo
}

type varInfo struct {
	Name     string
	Doc      string
	Type     string
	Required bool
}

type reamdeInfo struct {
	Authentication string
}

func main() {
	providers, err := getProviders()
	if err != nil {
		panic(err)
	}

	// for _, info := range providers {
	// 	fmt.Println("-----------")
	// 	fmt.Println(info.Name)
	// 	fmt.Println("-----------")

	// 	fmt.Println("Vars:")
	// 	fmt.Println(info.Variables)
	// 	fmt.Println("Readme:")
	// 	fmt.Println(info.Readme.Authentication)

	// 	fmt.Println()
	// }

	for _, name := range []string{"factory.go", "docs.md"} {
		err = generateFile(name, providers)
		if err != nil {
			panic(err)
		}
	}
}

func generateFile(tmplName string, providers []providerInfo) error {
	tmplFilename := fmt.Sprintf("./tmpl/%s.tmpl", tmplName)
	tmplFile, err := ioutil.ReadFile(tmplFilename)
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"FirstToUpper": func(s string) string {
			return strings.ToUpper(s[0:1]) + s[1:]
		},
		"ToUpper": strings.ToUpper,
		"Replace": strings.ReplaceAll,
	}

	tmpl, err := template.New("test").Funcs(funcMap).Parse(string(tmplFile))
	if err != nil {
		return err
	}

	tmplVars := tmplVars{
		Providers: providers,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &tmplVars)
	if err != nil {
		return err
	}

	targetFilename := fmt.Sprintf("./%s", tmplName)
	err = ioutil.WriteFile(targetFilename, buf.Bytes(), 640)
	if err != nil {
		return err
	}

	fmt.Printf("update: %s -> %s\n", tmplFilename, targetFilename)

	return nil
}

func generateFactoryGo(providers []providerInfo) error {
	tmplFile, err := ioutil.ReadFile("./factory.go.tmpl")
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}

	tmpl, err := template.New("test").Funcs(funcMap).Parse(string(tmplFile))
	if err != nil {
		return err
	}

	tmplVars := tmplVars{
		Providers: providers,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &tmplVars)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./factory.go", buf.Bytes(), 640)
	if err != nil {
		return err
	}

	return nil
}

func getProviders() ([]providerInfo, error) {
	goPath := os.Getenv("GOPATH")
	goModCache := filepath.Join(goPath, "pkg", "mod")

	data, err := ioutil.ReadFile("./go.mod")
	if err != nil {
		return nil, err
	}

	modFile, err := modfile.Parse("./go.mod", data, nil)
	if err != nil {
		return nil, err
	}

	providers := []providerInfo{}

	for _, req := range modFile.Require {
		if !strings.HasPrefix(req.Mod.Path, "github.com/libdns/") || strings.HasPrefix(req.Mod.Path, "github.com/libdns/libdns") {
			continue
		}

		dir := filepath.Join(goModCache, req.Mod.String())

		provider := providerInfo{
			Name:    filepath.Base(req.Mod.Path),
			Version: req.Mod.Version,
		}

		// get provider info from source
		providerFilename := filepath.Join(dir, "provider.go")
		providerFile, err := ioutil.ReadFile(providerFilename)
		if err != nil {
			return nil, err
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, providerFilename, providerFile, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		for _, node := range f.Decls {
			switch node.(type) {
			case *ast.GenDecl:
				genDecl := node.(*ast.GenDecl)
				for _, spec := range genDecl.Specs {
					switch spec.(type) {
					case *ast.TypeSpec:
						typeSpec := spec.(*ast.TypeSpec)

						if typeSpec.Name.Name != "Provider" {
							break
						}

						switch typeSpec.Type.(type) {
						case *ast.StructType:
							structType := typeSpec.Type.(*ast.StructType)

							for _, field := range structType.Fields.List {
								if len(field.Names) == 0 {
									continue
								}

								name := field.Names[0].Name
								if !unicode.IsUpper(rune(name[0])) {
									continue
								}

								tags, err := structtag.Parse(strings.Trim(field.Tag.Value, "`"))
								if err != nil {
									panic(err)
								}

								jsonTags, err := tags.Get("json")
								if err != nil {
									panic(err)
								}

								required := true
								for _, val := range jsonTags.Options {
									if val == "omitempty" {
										required = false
										break
									}
								}

								provider.Variables = append(provider.Variables, varInfo{
									Name:     name,
									Doc:      strings.TrimSpace(field.Doc.Text()),
									Type:     string(providerFile[field.Type.Pos()-1 : field.Type.End()-1]),
									Required: required,
								})

							}
						}
					}
				}
			}
		}

		// get info from readme
		readmeFilename := filepath.Join(dir, "README.md")
		file, err := ioutil.ReadFile(readmeFilename)
		if err != nil {
			return nil, err
		}

		authSectionLines := []string{}

		lines := strings.Split(string(file), "\n")
		isInAuthSection := false
		for _, line := range lines {
			if !isInAuthSection {
				if strings.TrimSpace(line) == "## Authenticating" {
					isInAuthSection = true
					continue
				}
			} else {
				if strings.HasPrefix(strings.TrimSpace(line), "#") {
					break
				}
				authSectionLines = append(authSectionLines, line)
			}
		}

		provider.Readme.Authentication = strings.TrimSpace(strings.Join(authSectionLines, "\n"))

		providers = append(providers, provider)
	}

	return providers, nil
}
