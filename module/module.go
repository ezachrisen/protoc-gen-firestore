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
	firestoreTpl *template.Template
	serviceTpl   *template.Template
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

	module := m.Parameters().Str("module")

	tpl := template.New("firestore").Funcs(map[string]interface{}{
		"package":                 m.PackageName,
		"name":                    m.Name,
		"shouldGenerateFirestore": shouldGenerateFirestore,
	})
	m.firestoreTpl = template.Must(tpl.Parse(firestoreTpl))

	tpl2 := template.New("service").Funcs(map[string]interface{}{
		"package":               m.PackageName,
		"name":                  m.Name,
		"getVerbEntity":         getVerbEntity,
		"isEntityMethod":        isEntityMethod,
		"shouldGenerateService": shouldGenerateService,
	})
	m.serviceTpl = template.Must(tpl2.Parse(serviceTpl))

	for _, f := range targets {
		m.Push(f.Name().String())
		defer m.Pop()

		filename := m.Context.OutputPath(f).SetExt(".firestore.go").String()
		if module != "" {
			filename = strings.TrimPrefix(filename, module+"/")
		}

		// FIRESTORE CLIENT GENERATION

		// check if twe need to execute the firestore client template
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
			m.AddGeneratorTemplateFile(filename, m.firestoreTpl, f)
		}

		// SERVICE GENERATION
		filename = m.Context.OutputPath(f).SetExt(".service.go").String()
		if module != "" {
			filename = strings.TrimPrefix(filename, module+"/")
		}

		m.AddGeneratorTemplateFile(filename, m.serviceTpl, f)

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

func shouldGenerateService(m pgs.Service) bool {
	var tval bool
	ok, err := m.Extension(firestore.E_GenerateService, &tval)
	if !ok || err != nil {
		return false
	}
	return tval
}

type verbEntity struct {
	Verb   string
	Entity string
}

func isEntityMethod(m pgs.Method) bool {
	ve := getVerbEntity(m)
	if ve.Verb == "" {
		return false
	}
	return true
}

func getVerbEntity(m pgs.Method) verbEntity {
	ve := verbEntity{}
	n := m.Name().String()

	if strings.HasPrefix(n, "Create") {
		ve.Verb = "Create"
		ve.Entity = strings.TrimPrefix(n, "Create")
		return ve
	}

	if strings.HasPrefix(n, "Get") {
		ve.Verb = "Get"
		ve.Entity = strings.TrimPrefix(n, "Get")
		return ve
	}

	if strings.HasPrefix(n, "Delete") {
		ve.Verb = "Delete"
		ve.Entity = strings.TrimPrefix(n, "Delete")
		return ve
	}

	if strings.HasPrefix(n, "Update") {
		ve.Verb = "Update"
		ve.Entity = strings.TrimPrefix(n, "Update")
		return ve
	}

	if strings.HasPrefix(n, "List") {
		ve.Verb = "List"
		ve.Entity = strings.TrimPrefix(n, "List")
		// list methods should be of the form "ListBooks", so
		// to get the entity ("Book") we must remove the s
		ve.Entity = strings.TrimSuffix(ve.Entity, "s")
		return ve
	}

	return ve

}
