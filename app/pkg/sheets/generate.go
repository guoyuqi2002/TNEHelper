//go:generate go run generate.go

package main

import (
	_ "embed"
	"iter"
	"os"
	"slices"
	"strconv"
	"strings"
	"tne-helper/pkg/sheets/definition"
)

//go:embed template.go.txt
var template string

func main() {
	for _, def := range definition.Definitions {
		folder := "./generated/" + def.Name + "/"
		if err := os.MkdirAll(folder, 0755); err != nil {
			panic(err)
		}

		f, err := os.Create(folder + "generated.go")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.WriteString(gen(def)); nil != err {
			panic(err.Error())
		}
	}
}

func gen(def definition.Definition) string {

	content := template
	content = strings.ReplaceAll(content, "{{name}}", def.Name)
	content = strings.ReplaceAll(content, "{{sheetName}}", def.SheetName)
	content = strings.ReplaceAll(content, "{{headRowIndex}}", strconv.Itoa(def.HeadRowIndex))
	content = strings.ReplaceAll(content, "{{len(keys)}}", strconv.Itoa(len(def.Heads)))
	content = strings.ReplaceAll(content, "{{keys}}", strings.Join(slices.Collect(appendKey(def.Heads, "\t`", "`,")), "\n"))

	indices := []string{}
	legalKeys := duplicatedKeysRename(slices.Collect(makeLegalKey(def.Heads)))
	for idx, key := range legalKeys {
		indices = append(indices, "\tꝎ"+key+" Index"+" = "+strconv.Itoa(idx))
	}

	content = strings.ReplaceAll(content, "{{indices}}", strings.Join(indices, "\n"))

	return content
}

func makeLegalKey(list []string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range list {
			v = strings.ReplaceAll(v, " ", "_")
			v = strings.ReplaceAll(v, ",", "_")
			v = strings.ReplaceAll(v, "'", "_")
			v = strings.ReplaceAll(v, "(", "_")
			v = strings.ReplaceAll(v, ")", "_")
			v = strings.ReplaceAll(v, "/", "_")
			v = strings.ReplaceAll(v, "#", "_")
			v = strings.ReplaceAll(v, "-", "_")
			v = strings.ReplaceAll(v, "?", "_")
			if !yield(v) {
				return
			}
		}
	}
}

func appendKey(list []string, prefix, suffix string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range list {
			v = prefix + v + suffix
			if !yield(v) {
				return
			}
		}
	}
}

func duplicatedKeysRename(keys []string) []string {
	count := map[string]int{}
	ret := []string{}

	for _, key := range keys {
		count[key] += 1

		if count[key] == 1 {
			ret = append(ret, key)
		} else {
			ret = append(ret, key+strconv.Itoa(count[key]))
		}
	}

	return ret
}
