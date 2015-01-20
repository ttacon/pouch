package main

import (
	"fmt"

	"github.com/ttacon/pouch/pouch/defs"
)

// TODO(ttacon): make this smarter
// TODO(ttacon): refactor helpers out so it's not one giant glob
func diffAndReport(dbe, fe []*defs.StructInfo) {
	var (
		diffs []string
		dm    = make(map[string]*defs.StructInfo)
		fm    = make(map[string]*defs.StructInfo)
		seen  = make(map[string]struct{})
	)

	for _, s := range dbe {
		dm[s.Name] = s
	}

	for _, s := range fe {
		fm[s.Name] = s
	}

	for name, e := range dm {
		fent, ok := fm[name]
		seen[name] = struct{}{}
		if !ok {
			diffs = append(diffs, "no known struct: "+name)
			continue
		}
		dbFMap := fieldMapFromEntity(e)
		feFMap := fieldMapFromEntity(fent)
		for fieldName, field := range dbFMap {
			fField, ok := feFMap[fieldName]
			if !ok {
				diffs = append(diffs, "struct on file doesn't have field: "+fieldName)
				continue
			}
			if field.Column != fField.Column {
				diffs = append(diffs,
					fmt.Sprintf("columns for %s.%s differ: %s vs %q",
						name, fieldName, field.Column, fField.Column))
			}
			if field.IsPrimaryKey != fField.IsPrimaryKey {
				if field.IsPrimaryKey {
					diffs = append(diffs,
						fmt.Sprintf("column %s.%s is reported as a primary key by "+
							"the database but not by the struct definition",
							name, fieldName))
				} else {
					diffs = append(diffs,
						fmt.Sprintf("column %s.%s is not reported as a primary "+
							"key by the database but is by the struct definition",
							name, fieldName))
				}
			}
			if field.IsPointer != fField.IsPointer {
				if field.IsPointer {
					diffs = append(diffs,
						fmt.Sprintf("column %s.%s is reported as nullable by "+
							"the database but not by the struct definition",
							name, fieldName))
				} else {
					diffs = append(diffs,
						fmt.Sprintf("column %s.%s is reported as not nullable by "+
							"the database but nullable by the struct definition",
							name, fieldName))
				}
			}
			if field.Type != fField.Type {
				diffs = append(diffs,
					fmt.Sprintf("column %s.%s has a type conflict (db/code): %q vs %q",
						name, fieldName, field.Type, fField.Type))
			}
		}
	}

	// now just double check if any code defined structs are not in the db
	for name, _ := range fm {
		if _, ok := seen[name]; !ok {
			diffs = append(diffs,
				fmt.Sprintf("%s is defined by the code but not in the db", name))
		}
	}

	fmt.Printf(dbgenPrmpt+" %d differences to report\n", len(diffs))
	if len(diffs) == 0 {
		return
	}

	for i, diff := range diffs {
		fmt.Printf("[%d] %s\n", i+1, diff)
	}
}

func fieldMapFromEntity(s *defs.StructInfo) map[string]defs.FieldInfo {
	var m = make(map[string]defs.FieldInfo)
	for _, f := range s.Fields {
		m[f.Name] = f
	}
	return m
}

func diffContained(dbe, fe []*defs.StructInfo) []*defs.StructInfo {
	var (
		dm      = make(map[string]*defs.StructInfo)
		notSeen []*defs.StructInfo
	)

	for _, s := range dbe {
		dm[s.Name] = s
	}

	for _, s := range fe {
		if _, ok := dm[s.Name]; !ok {
			notSeen = append(notSeen, s)
		}
	}

	return notSeen
}
