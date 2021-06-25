package module

import (
	// "go/parser"
	// "go/printer"
	// 	"go/token"
	// "path/filepath"
	// "strings"
	//	"fmt"
	//	"github.com/fatih/structtag"

	"strings"
	"text/template"

	"alticeusa.com/maui/protoc-gen-firestore/firestore"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type mod struct {
	*pgs.ModuleBase
	pgsgo.Context
	tpl *template.Template
}

func New() pgs.Module {
	return &mod{
		ModuleBase: &pgs.ModuleBase{},
	}
}

func (m *mod) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())
}

func (mod) Name() string {
	return "gen-firestore"
}

func (m mod) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	//	buf := &bytes.Buffer{}

	module := m.Parameters().Str("module")
	tpl := template.New("firestore").Funcs(map[string]interface{}{
		"package":                 m.PackageName,
		"name":                    m.Name,
		"shouldGenerateFirestore": shouldGenerateFirestore,
	})
	m.tpl = template.Must(tpl.Parse(firestoreTpl))

	for _, f := range targets {
		m.Push(f.Name().String())
		defer m.Pop()

		// check if twe need to execute the template
		// we only want to generate the template if at least 1 message
		// has the option set for firestore, otherwise we end up with an
		// empty file with imports, which whill break in compilation
		skipThisFile := true
		for _, msg := range f.Messages() {
			if shouldGenerateFirestore(msg) {
				skipThisFile = false
				break
			}
		}

		if !skipThisFile {
			filename := m.Context.OutputPath(f).SetExt(".firestore.go").String()
			if module != "" {
				filename = strings.TrimPrefix(filename, module+"/")
			}
			m.AddGeneratorTemplateFile(filename, m.tpl, f)
		}
	}

	return m.Artifacts()
}

func shouldGenerateFirestore(m pgs.Message) bool {
	var tval bool
	ok, err := m.Extension(firestore.E_GenerateFirestore, &tval)
	if !ok || err != nil {
		return false
	}
	return tval
}

// This is an alternate way of doing this, using a visitor
// buf.Reset()

// v := newVisitor(buf, m)
// m.CheckErr(pgs.Walk(v, f), "unable to process firestore generation")
// out := buf.String()

// m.AddGeneratorFile(
// 	filename,
// 	out,
// )
