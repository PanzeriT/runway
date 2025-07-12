package generator

// TODO: generator needs a defaul strct that can be passed as first argument for all pages

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/panzerit/runway/internal/generator/asset"
	"github.com/panzerit/runway/internal/generator/template/code"
	"github.com/panzerit/runway/internal/generator/template/html"
	"github.com/spf13/viper"
)

type Generator struct {
	mainDir   string
	outputDir string
	schemaDir string

	debugMode bool

	template *template.Template
}

type OptionFunc func(*Generator) *Generator

func (g *Generator) debug(a ...any) {
	if g.debugMode {
		fmt.Println(a...)
	}
}

func (g *Generator) loadDefaults() {
	if g.mainDir == "" {
		mainDir, err := os.Getwd()
		CheckError(err, 1, "Error getting current directory")
		g.mainDir = mainDir
	}

	if g.outputDir == "" {
		g.outputDir = g.mainDir
	} else {
		g.outputDir = makePathAbsolute(g.outputDir, g.mainDir)
	}

	if g.schemaDir == "" {
		g.schemaDir = filepath.Join(g.outputDir, "schemas")
	}
}

func New(opts ...OptionFunc) *Generator {
	g := &Generator{}

	for _, opt := range opts {
		opt(g)
	}

	g.loadDefaults()
	g.loadTemplates()

	g.debug(fmt.Sprintf("Generator initialized:\n  mainDir:   %s\n  outputDir: %s\n  schemaDir: %s",
		g.mainDir, g.outputDir, g.schemaDir))

	return g
}

func WithOutputDir(name string) OptionFunc {
	return func(g *Generator) *Generator {
		g.outputDir = name
		return g
	}
}

func WithDebugMode(on bool) OptionFunc {
	return func(g *Generator) *Generator {
		g.debugMode = on
		return g
	}
}

func (g Generator) Run() {
	g.generateMainGo()
	g.generateStaticFiles()
	g.copyHTMLCode()
	g.copyAssets()
	// g.processSchemas()  TODO: activate again
}

func (g *Generator) loadTemplates() {
	template, err := template.New("generator").ParseFS(code.FS, "*.tmpl")
	CheckError(err, Success)

	g.template = template
}

func (g Generator) generateMainGo() {
	err := g.writeFile("main.go", "main.go", map[string]string{
		"ModuleName": viper.GetString("module_name"),
	})
	CheckError(err, ErrWritingFile)
}

func (g *Generator) Init(moduleName string) *Generator {
	viper.Set("module_name", moduleName)
	viper.WriteConfig()

	g.writeFile("go.mod", "go.mod", map[string]string{
		"ModuleName": moduleName,
	})

	return g
}

func (g *Generator) SaveConfig() {
}

type StaticFile struct {
	filePath string
	data     any
}

func (g Generator) generateStaticFiles() {
	// TODO: find files automatically
	staticFiles := []StaticFile{
		{filePath: "./internal/config/config.go"},
		{filePath: "./internal/config/loader.go"},

		{filePath: "./internal/server/admin/admin.go"},
		{filePath: "./internal/server/admin/auth.go"},
		{filePath: "./internal/server/admin/data.go"},
		{filePath: "./internal/server/admin/dashboard.go"},
		{filePath: "./internal/server/admin/middleware.go"},
	}

	// TODO: this has to use the file system
	for _, s := range staticFiles {
		parts := strings.Split(s.filePath, "/")
		g.writeFile(s.filePath, parts[len(parts)-1], map[string]string{
			"ModuleName": viper.GetString("module_name"),
		})
	}
}

func (g Generator) writeFile(relPath, name string, data any) error {
	filePath := filepath.Join(g.outputDir, relPath)
	os.MkdirAll(filepath.Dir(filePath), 0755)

	f, err := os.Create(filePath)
	CheckError(err, ErrCreatingFile, filePath)
	defer f.Close()

	err = g.template.ExecuteTemplate(f, name, data)
	if err != nil {
		return err
	}

	g.debug("wrote file", filePath)

	return nil
}

func (g Generator) copyHTMLCode() {
	dst := filepath.Join(g.outputDir, "internal/server/admin/template")

	err := os.RemoveAll(dst)
	CheckError(err, ErrDeletingFile, dst)

	err = fs.WalkDir(html.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			dstPath := getTemplateFileName(path)
			err := copyFSFile(html.FS, path, filepath.Join(dst, dstPath))
			CheckError(err, ErrCopyingFile, path)
		} else {
			os.MkdirAll(filepath.Join(dst, path), 0755)
		}
		return nil
	})
	CheckError(err, ErrCreatingFile)
}

// copyAssets removes previous assets and copies new ones from the embedded asset
// filesystem.
func (g Generator) copyAssets() {
	dst := filepath.Join(g.outputDir, "internal/server/admin/asset")

	err := os.RemoveAll(dst)
	CheckError(err, ErrDeletingFile, dst)

	err = os.CopyFS(dst, asset.FS)
	CheckError(err, ErrCreatingFile, dst)
}

func (g Generator) processSchemas() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || filepath.Ext(path) != ".go" {
			return nil
		}
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				fmt.Printf("Struct: %s\n", typeSpec.Name.Name)
				if genDecl.Doc != nil {
					fmt.Printf("  Doc: %s\n", genDecl.Doc.Text())
				}
				for _, field := range structType.Fields.List {
					for _, name := range field.Names {
						fmt.Printf("  Field: %s\n", name.Name)
					}
					if field.Tag != nil {
						fmt.Printf("    Tag: %s\n", field.Tag.Value)
					}
					if field.Doc != nil {
						fmt.Printf("    Doc: %s\n", field.Doc.Text())
					}

				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Tidy up...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Run()

	fmt.Println("Done!")
}

func getTemplateFileName(fileName string) string {
return strings.ReplaceAll(fileName, "_templ", "")
}
