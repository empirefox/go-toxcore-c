// ./tools/tools -f /usr/local/include/tox/tox.h > toxenums/toxcore_consts.go
// ./tools/tools -f /usr/local/include/tox/toxav.h -e TOXAV_ERR_ > toxenums/toxav_consts.go
// ./tools/tools -f /usr/local/include/tox/toxencryptsave.h > toxenums/toxencryptsave_consts.go
// go generate ./toxenums
package main

import (
	"bytes"
	"flag"
	"go/format"
	"log"
	"os"
	"strings"

	"github.com/go-clang/v3.8/clang"
)

var (
	filename  = flag.String("f", "/usr/include/tox/tox.h", "file name to import consts")
	errPrefix = flag.String("e", "TOX_ERR_", "gen error if has this prefix")
)

func main() {
	flag.Parse()

	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	cmdArgs := []string{
		"-std=c99",
		"-I/usr/lib/clang/3.8/include",
	}
	tu := idx.ParseTranslationUnit(*filename, cmdArgs, nil, 0)
	defer tu.Dispose()

	for _, d := range tu.Diagnostics() {
		log.Println("PROBLEM:", d.Spelling())
	}

	var enumDecls = make(map[uint32]bool)
	var enumConstDecls = make(map[uint32]bool)

	var data TplData
	var currentEnum *Enum
	tuc := tu.TranslationUnitCursor()
	ok := tuc.Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
		switch cursor.Kind() {
		case clang.Cursor_EnumDecl:
			if _, ok := enumDecls[cursor.HashCursor()]; ok {
				break
			}
			enumDecls[cursor.HashCursor()] = true

			if currentEnum != nil {
				data.Enums = append(data.Enums, *currentEnum)
			}
			// log.Println(cursor.BriefCommentText())
			// log.Println(cursor.Type().Kind().String(), cursor.Type().Spelling())
			currentEnum = NewEnum(cursor.Spelling())
			if strings.HasPrefix(currentEnum.Name, *errPrefix) {
				data.Errors = append(data.Errors, currentEnum.Name)
			}

		case clang.Cursor_EnumConstantDecl:
			if _, ok := enumConstDecls[cursor.HashCursor()]; ok {
				break
			}
			enumConstDecls[cursor.HashCursor()] = true

			// log.Println(cursor.BriefCommentText())
			// log.Println(cursor.Kind().String(), cursor.Type().Kind().String(), cursor.Type().Spelling(), cursor.Spelling())
			currentEnum.consts = append(currentEnum.consts, cursor.Spelling())
			currentEnum.values = append(currentEnum.values, cursor.EnumConstantDeclValue())

		default:
			// log.Println(cursor.Kind().String(), cursor.Type().Kind().String(), cursor.Type().Spelling(), cursor.Spelling())
		}

		return clang.ChildVisit_Recurse
	})

	if !ok {
		log.Fatalln("clang Visit failed")
	}

	if currentEnum == nil {
		log.Fatalln("no enums found")
	}
	data.Enums = append(data.Enums, *currentEnum)

	for i := range data.Enums {
		data.Enums[i].withValues = withValues(data.Enums[i].values)
	}

	var buf bytes.Buffer
	err := tpl.Execute(&buf, &data)
	if err != nil {
		log.Fatalln(err)
	}

	source, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalln(err)
	}

	os.Stdout.Write(source)
}

func withValues(values []int64) bool {
	for i, v := range values {
		if int64(i) != v {
			return true
		}
	}
	return false
}
