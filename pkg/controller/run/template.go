package run

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func txtFuncMap() template.FuncMap {
	fncs := sprig.TxtFuncMap()
	delete(fncs, "env")
	delete(fncs, "expandenv")
	delete(fncs, "getHostByName")
	return fncs
}
