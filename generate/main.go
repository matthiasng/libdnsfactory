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
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/module"
	"golang.org/x/mod/modfile"
)

type tmplVars struct {
	Providers []providerInfo
}

type providerInfo struct {
	Name        string
	PackageName string
	Version     string
	Variables   []varInfo
	Readme      reamdeInfo
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

func includeModule(req *modfile.Require) bool {
	if strings.HasPrefix(req.Mod.Path, "github.com/libdns/libdns") {
		return false
	}

	if strings.HasPrefix(req.Mod.Path, "github.com/libdns/") || strings.HasPrefix(req.Mod.Path, "github.com/Alfschmalf/inwx") {
		return true
	}

	return false
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

	tmpl, err := template.New(tmplName).Funcs(funcMap).Parse(string(tmplFile))
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
	err = ioutil.WriteFile(targetFilename, buf.Bytes(), 0640)
	if err != nil {
		return err
	}

	fmt.Printf("update: %s -> %s\n", tmplFilename, targetFilename)
	return nil
}

func getProviders() ([]providerInfo, error) {
	goPath := os.Getenv("GOPATH")
	goModCache := filepath.Join(goPath, "pkg", "mod")

	data, err := ioutil.ReadFile("./go.mod")
	if err != nil {
		return nil, errors.Wrap(err, "loading modfile")
	}

	modFile, err := modfile.Parse("./go.mod", data, nil)
	if err != nil {
		return nil, errors.Wrap(err, "parsing modfile")
	}

	providers := []providerInfo{}

	for _, req := range modFile.Require {
		if !includeModule(req) {
			continue
		}

		encodedModName, err := module.EncodePath(req.Mod.Path)
		if err != nil {
			return nil, errors.Wrap(err, "encoding module path")
		}
		encodedModVersion, err := module.EncodeVersion(req.Mod.Version)
		if err != nil {
			return nil, errors.Wrap(err, "encoding module version")
		}

		dir := filepath.Join(goModCache, encodedModName+"@"+encodedModVersion)
		providerFilename := filepath.Join(dir, "provider.go")

		packageName, err := getPackageNameFromFile(providerFilename)
		if err != nil {
			return nil, errors.Wrap(err, "getting package name")
		}

		provider := providerInfo{
			Name:        filepath.Base(req.Mod.Path),
			PackageName: packageName,
			Version:     req.Mod.Version,
		}

		// get provider info from source
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

func getPackageNameFromFile(filename string) (string, error) {
	fset := token.NewFileSet()

	astFile, err := parser.ParseFile(fset, filename, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}

	if astFile.Name == nil {
		return "", fmt.Errorf("no package found")
	}

	return astFile.Name.Name, nil
}
