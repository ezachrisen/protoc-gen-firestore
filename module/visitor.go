package module

// This file has an alternate way of traversing the nodes in the proto file

// import (
// 	"fmt"
// 	"io"

// 	"alticeusa.com/maui/protoc-gen-firestore/firestore"
// 	pgs "github.com/lyft/protoc-gen-star"
// )

// // visitor visits each node of a .proto file's structure
// type visitor struct {
// 	pgs.Visitor
// 	pgs.DebuggerCommon
// 	w io.Writer
// }

// func newVisitor(w io.Writer, d pgs.DebuggerCommon) pgs.Visitor {
// 	return visitor{
// 		w:              w,
// 		Visitor:        pgs.NilVisitor(),
// 		DebuggerCommon: d,
// 	}
// }

// func (v visitor) VisitPackage(p pgs.Package) (pgs.Visitor, error) {
// 	v.Debugf("Package %s", p.ProtoName())
// 	return v, nil
// }

// func (v visitor) VisitMessage(m pgs.Message) (pgs.Visitor, error) {
// 	v.Debugf("  Message: %s", m.Name())
// 	var tval bool
// 	ok, err := m.Extension(firestore.E_GenerateFirestore, &tval)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if ok {
// 		v.Debugf("    - Option: %t:%t", ok, tval)
// 	}

// 	fmt.Fprintf(v.w, "Writing option for %s", m.Name())

// 	return v, nil
// }

// func (v visitor) VisitField(f pgs.Field) (pgs.Visitor, error) {
// 	v.Debugf("  + Field: %s", f.Name())
// 	return v, nil
// }

// func (v visitor) VisitFile(f pgs.File) (pgs.Visitor, error) {
// 	v.Debugf("File: %s", f.InputPath())
// 	return v, nil
// }

// func (v visitor) VisitService(s pgs.Service) (pgs.Visitor, error) {
// 	v.Debugf("  Sevice: %s", s.Name())
// 	return v, nil
// }

// // VisitMethod logs out ServiceName#MethodName for m.
// func (v visitor) VisitMethod(m pgs.Method) (pgs.Visitor, error) {
// 	// m.Service().Name()
// 	v.Debugf("     Method: %s", m.Name())
// 	return nil, nil
// }
