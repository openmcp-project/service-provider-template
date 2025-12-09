package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// exec from project root
const templatesDir = "cmd/template/files"

type TemplateData struct {
	Group       string
	Version     string
	Kind        string
	Package     string
	Module      string
	WithExample bool
}

func main() {
	group := flag.String("group", "foo", "GVK group prefix (will always be suffixed with services.openmcp.cloud)")
	version := flag.String("version", "v1alpha1", "GVK version")
	kind := flag.String("kind", "FooService", "GVK kind")
	includeSample := flag.Bool("sample", false, "Generate with sample implementation")
	module := flag.String("module", "github.com/openmcp-project/service-provider-template", "Go module")
	flag.Parse()
	data := TemplateData{
		Group:       *group,
		Version:     *version,
		Kind:        *kind,
		Package:     strings.ToLower(*group),
		Module:      *module,
		WithExample: *includeSample,
	}
	apiDir := filepath.Join("api", "v1alpha1")
	controllerDir := filepath.Join("internal", "controller")
	e2eDir := filepath.Join("test", "e2e")
	// api
	execTemplate("types.go.tmpl",
		filepath.Join(apiDir, strings.ToLower(fmt.Sprintf("%s_types.go", *kind))),
		data)
	execTemplate("groupversion_info.go.tmpl",
		filepath.Join(apiDir, "groupversion_info.go"),
		data)
	// controller
	execTemplate("controller.go.tmpl",
		filepath.Join(controllerDir, strings.ToLower(fmt.Sprintf("%s_controller.go", *kind))),
		data)
	// e2e tests
	execTemplate("test.go.tmpl",
		filepath.Join(e2eDir, "serviceprovider_test.go"),
		data)
	// call go mod edit followed by gopls to rename module
	if err := exec.Command("go", "mod", "edit", "-module", *module).Run(); err != nil {
		log.Fatalf("go mod edit failed: %v", err)
	}
	// fix import statements
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not determine current working directory: %v", err)
	}
	log.Println(rootDir)
	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".go") {
			return replaceImports(path, "github.com/openmcp-project/service-provider-template", *module)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("rename imports failed: %v", err)
	}
	// clean up repo
	// remove /cmd/template
	fmt.Printf("Generated service-provider for %s/%s' in %s\n", data.Group, data.Kind, *module)
}

func execTemplate(templateName, outPath string, data TemplateData) {
	tplPath := filepath.Join(templatesDir, templateName)
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		log.Fatalf("failed parsing template %s: %v", templateName, err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("failed creating file %s: %v", outPath, err)
	}
	log.Default().Println(outPath)
	defer close(outPath, f)
	if err := tpl.Execute(f, data); err != nil {
		log.Fatalf("failed executing template %s: %v", templateName, err)
	}
}

func replaceImports(filename, oldRepo, newRepo string) error {
	input, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer close(filename, input)
	var result []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, oldRepo, newRepo, 1)
		result = append(result, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	// keep empty line
	result = append(result, "")
	return os.WriteFile(filename, []byte(strings.Join(result, "\n")), 0644)
}

func close(filename string, c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatalf("please reset/checkout %s and try again: %v", filename, err)
	}
}
