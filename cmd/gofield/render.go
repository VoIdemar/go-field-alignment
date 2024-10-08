package main

import (
	"fmt"
	"strings"
)

// ============= Render

// renderStructure generates a string representation of an Structure structure.
// It handles both top-level structures and nested fields, including their
// documentation, tags, and comments. The function recursively processes
// nested structures to create a complete representation.
//
// The function performs the following tasks:
// - Checks if the element is a valid custom type or a structure
// - Generates the struct definition with its name (for top-level structures)
// - Iterates through all nested fields, rendering each one
// - Includes field documentation, tags, and comments
// - Handles both root-level and nested comments
//
// Parameters:
// - elem: Pointer to an Structure structure to be rendered
//
// Returns:
// - A string containing the rendered structure
func renderStructure(elem *Structure) string {
	data := ""

	isValidCustomNameType := isValidCustomTypeName(elem.StringType)
	if !elem.IsStructure || isValidCustomNameType {
		return elem.StringType
	}

	if elem.Root != nil {
		// Don't add "type " here, as current structure may be inside a "type" block
		data += fmt.Sprintf("%s struct{", elem.Name)
	} else {
		data += fmt.Sprintf("struct {")
	}
	if len(elem.NestedFields) > 0 {
		data += fmt.Sprintf("\n")
	}

	for idx, field := range elem.NestedFields {
		// Doc
		if field.RootField != nil && field.RootField.Doc != nil && len(field.RootField.Doc.List) > 0 {
			for _, comment := range field.RootField.Doc.List {
				data += fmt.Sprintln(comment.Text)
			}
		}
		if strings.HasPrefix(field.Name, "!") {
			field.Name = ""
		}
		data += fmt.Sprintf("%s %s ", field.Name, renderStructure(field))
		// Tag
		if field.RootField != nil {
			// Tags
			if field.RootField.Tag != nil && len(field.RootField.Tag.Value) > 0 {
				data += fmt.Sprintf(" %s", field.RootField.Tag.Value)
			}
			// Comment
			if field.RootField.Comment != nil && len(field.RootField.Comment.List) > 0 {
				for _, comment := range field.RootField.Comment.List {
					data += fmt.Sprintf(" %s", comment.Text)
				}
			}
		}
		if idx != len(elem.NestedFields) {
			data += fmt.Sprintf("\n")
		}
	}

	data += fmt.Sprintf("}")
	// Comments
	if elem.RootField != nil {
		if elem.RootField.Comment != nil && len(elem.RootField.Comment.List) > 0 {
			for _, comment := range elem.RootField.Comment.List {
				data += fmt.Sprintf("%s", comment.Text)
			}
		}
	} else if elem.Root != nil {
		if elem.Root.Comment != nil && len(elem.Root.Comment.List) > 0 {
			for _, comment := range elem.Root.Comment.List {
				data += fmt.Sprintf("%s", comment.Text)
			}
		}
	}
	return data
}

func renderTextStructures(structures []*Structure) {
	for _, structure := range structures {
		// Don't format code here - "renderStructure" generates a replacement for a part of target Go file,
		// not a valid piece of Go code per-se.
		//
		// Code will be formatted afterwards.
		structure.MetaData.Data = []byte(renderStructure(structure))
	}
}
