package main

import (
	"fmt"
	"strings"
	"text/template"
)

type Enum struct {
	Name       string
	withValues bool
	consts     []string
	values     []int64
}

func NewEnum(name string) *Enum {
	return &Enum{Name: name}
}

func (e *Enum) First() string {
	if e.withValues {
		return ""
	}
	return fmt.Sprintf("%s %s = iota", e.consts[0], e.Name)
}

func (e *Enum) Others() string {
	if e.withValues {
		var lines string
		for i, c := range e.consts {
			lines += fmt.Sprintf("%s %s = %d\n", c, e.Name, e.values[i])
		}
		return lines
	}
	return strings.Join(e.consts[1:], "\n")
}

type TplData struct {
	Enums  []Enum
	Errors []string
}

func (e *TplData) Stringers() string {
	names := make([]string, len(e.Enums))
	for i := range e.Enums {
		names[i] = e.Enums[i].Name
	}
	return strings.Join(names, ",")
}

const TplStr = `
//go:generate stringer -type={{.Stringers}}
package toxenums

import "fmt"

{{range .Enums}}
type {{.Name}} int
const(
	{{.First}}
	{{.Others}}
)
{{end}}

{{range .Errors}}
func (e {{.}}) Error() string { return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e) } {{end}}
`

var tpl = template.Must(template.New("toxcore_consts").Parse(TplStr))
