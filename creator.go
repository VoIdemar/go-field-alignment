package main

import "go/ast"

// ============= Creator Item
func createItemInfoName(data interface{}) string {
	switch elem := data.(type) {
	case *ast.TypeSpec:
		return elem.Name.Name
	case *ast.Field:
		if len(elem.Names) == 0 {
			return ""
		}
		return elem.Names[0].Name
	}
	return ""
}

func createItemInfoPath(name, parentName string) string {
	if parentName != "" {
		return parentName + "/" + name
	}
	return name
}
func createItemInfo(data interface{}, parentData *ItemInfo, mapper map[string]*ItemInfo) *ItemInfo {
	switch Elem := data.(type) {
	case *ast.TypeSpec:
		newItem := &ItemInfo{
			Name:        createItemInfoName(Elem),
			Root:        Elem,
			StructType:  Elem.Type,
			IsStructure: true,
			StringType:  getTypeString(Elem.Type),
		}
		if newItem.Name == "" {
			return nil
		}
		newItem.Path = createItemInfoPath(newItem.Name, "")
		mapper[newItem.Path] = newItem
		if elem, ok := Elem.Type.(*ast.StructType); ok {
			for _, field := range elem.Fields.List {
				newField := createItemInfo(field, newItem, mapper)
				if newField != nil {
					newItem.NestedFields = append(newItem.NestedFields, newField)
				}
			}
		}
		return newItem
	case *ast.Field:
		newItem := &ItemInfo{
			Name:       createItemInfoName(Elem),
			RootField:  Elem,
			StructType: Elem.Type,
			StringType: getTypeString(Elem.Type),
		}
		if newItem.Name == "" {
			return nil
		}
		newItem.Path = createItemInfoPath(newItem.Name, parentData.Path)
		mapper[newItem.Path] = newItem
		if ident, ok := Elem.Type.(*ast.Ident); ok {
			newItem.IsStructure = ident.Obj != nil
		}
		if elem, ok := Elem.Type.(*ast.StructType); ok {
			newItem.IsStructure = true
			for _, field := range elem.Fields.List {
				newField := createItemInfo(field, newItem, mapper)
				if newField != nil {
					newItem.NestedFields = append(newItem.NestedFields, newField)
				}
			}
		}
		return newItem
	}
	return nil
}
