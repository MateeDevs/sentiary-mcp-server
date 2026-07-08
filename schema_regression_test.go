package main

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/MateeDevs/sentiary-mcp-server/internal/sentiary"
	"github.com/google/jsonschema-go/jsonschema"
)

// toolSchemaTypes lists every input and output type used by the tools
// registered in registerTools. The MCP go-sdk derives each tool's
// input/output JSON Schema from these types via jsonschema.ForType, so this is
// the authoritative set to guard.
var toolSchemaTypes = []reflect.Type{
	// Inputs
	reflect.TypeFor[sentiary.ListStringsInput](),
	reflect.TypeFor[sentiary.SearchStringInput](),
	reflect.TypeFor[sentiary.GetStringInput](),
	reflect.TypeFor[sentiary.GetStringByKeyInput](),
	reflect.TypeFor[sentiary.AddStringInput](),
	reflect.TypeFor[sentiary.EditStringInput](),
	reflect.TypeFor[sentiary.RemoveStringInput](),
	reflect.TypeFor[sentiary.SetTranslationInput](),
	reflect.TypeFor[sentiary.RemoveTranslationInput](),
	reflect.TypeFor[sentiary.ProjectInput](),
	reflect.TypeFor[sentiary.ExportInput](),
	reflect.TypeFor[sentiary.ImportInput](),
	// Outputs
	reflect.TypeFor[sentiary.Paging](),
	reflect.TypeFor[sentiary.Term](),
	reflect.TypeFor[sentiary.DeleteOutput](),
	reflect.TypeFor[sentiary.ProjectBatchInfo](),
	reflect.TypeFor[sentiary.ExportOutput](),
	reflect.TypeFor[sentiary.ImportResult](),
}

// TestToolSchemasHaveNoBooleanProperties guards against regressing to a schema
// where a `properties` entry serializes to a JSON boolean (e.g. an `any` field
// emitting `true`). Claude Code's strict validator rejects boolean property
// subschemas and drops the entire tool list, so a single offending field takes
// down every mcp__sentiary__* tool.
func TestToolSchemasHaveNoBooleanProperties(t *testing.T) {
	for _, rt := range toolSchemaTypes {
		schema, err := jsonschema.ForType(rt, &jsonschema.ForOptions{})
		if err != nil {
			t.Fatalf("ForType(%s): %v", rt.Name(), err)
		}
		raw, err := json.Marshal(schema)
		if err != nil {
			t.Fatalf("marshal schema for %s: %v", rt.Name(), err)
		}
		var doc any
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatalf("unmarshal schema for %s: %v", rt.Name(), err)
		}
		if paths := booleanPropertyPaths(doc, rt.Name()); len(paths) > 0 {
			t.Errorf("type %s emits boolean-valued property subschema(s): %v\nschema: %s", rt.Name(), paths, raw)
		}
	}
}

// booleanPropertyPaths walks a JSON Schema document and returns the paths of any
// entry inside a `properties` map whose value is a JSON boolean. It intentionally
// ignores boolean `additionalProperties`, which is valid and accepted.
func booleanPropertyPaths(node any, path string) []string {
	obj, ok := node.(map[string]any)
	if !ok {
		return nil
	}
	var hits []string
	if props, ok := obj["properties"].(map[string]any); ok {
		for key, value := range props {
			childPath := path + ".properties." + key
			if _, isBool := value.(bool); isBool {
				hits = append(hits, childPath)
				continue
			}
			hits = append(hits, booleanPropertyPaths(value, childPath)...)
		}
	}
	if items, ok := obj["items"].(map[string]any); ok {
		hits = append(hits, booleanPropertyPaths(items, path+".items")...)
	}
	if ap, ok := obj["additionalProperties"].(map[string]any); ok {
		hits = append(hits, booleanPropertyPaths(ap, path+".additionalProperties")...)
	}
	for _, combiner := range []string{"allOf", "anyOf", "oneOf"} {
		if list, ok := obj[combiner].([]any); ok {
			for i, sub := range list {
				hits = append(hits, booleanPropertyPaths(sub, path+"."+combiner+"["+strconv.Itoa(i)+"]")...)
			}
		}
	}
	return hits
}
