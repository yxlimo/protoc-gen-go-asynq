package gen

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

var _ pgs.Module = &Module{}

type Module struct {
	base *pgs.ModuleBase
	ctx  pgsgo.Context
	tpl  *template.Template
}

// Module returns an initialized JSONifyPlugin
func NewModule() *Module { return &Module{base: &pgs.ModuleBase{}} }

func (f *Module) InitContext(c pgs.BuildContext) {
	f.base.InitContext(c)
	f.ctx = pgsgo.InitContext(c.Parameters())
	tpl := template.New("fednode").Funcs(map[string]interface{}{
		"package": f.ctx.PackageName,
		"name":    f.ctx.Name,
	})
	f.tpl = template.Must(tpl.Parse(fednodeTpl))
}

func (f *Module) Name() string {
	return "fednode"
}

func (f *Module) Execute(targets map[string]pgs.File, _ map[string]pgs.Package) []pgs.Artifact {
	for _, t := range targets {
		f.generate(t)
	}
	return f.base.Artifacts()
}

func (f *Module) generate(file pgs.File) {
	if len(file.Services()) == 0 {
		return
	}
	name := f.ctx.OutputPath(file).SetExt(".fednode.go")
	f.base.AddGeneratorTemplateFile(name.String(), f.tpl, file)
}
