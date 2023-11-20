package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type Match struct {
	Value     string
	SourceKey string
	TargetKey string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <source.json> <target.json>")
		os.Exit(1)
	}

	sourceFile, targetFile := os.Args[1], os.Args[2]

	sourceData, err := os.ReadFile(sourceFile)
	if err != nil {
		panic(err)
	}

	targetData, err := os.ReadFile(targetFile)
	if err != nil {
		panic(err)
	}

	var sourceJson, targetJson interface{}
	json.Unmarshal(sourceData, &sourceJson)
	json.Unmarshal(targetData, &targetJson)

	sourceMap := make(map[string]string)
	targetMap := make(map[string]string)

	flattenJson(sourceJson, "", sourceMap)
	flattenJson(targetJson, "", targetMap)

	matches := findMatches(sourceMap, targetMap)
	for _, match := range matches {
		fmt.Printf("---\nValue: %s\nSource Key: %s\nTarget Key: %s\n", match.Value, match.SourceKey, match.TargetKey)
	}
}

func flattenJson(data interface{}, prefix string, flatMap map[string]string) {
	if reflect.TypeOf(data) == nil {
		return
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		for key, value := range data.(map[string]interface{}) {
			flattenJson(value, prefix+key+".", flatMap)
		}
	case reflect.Slice:
		for i, value := range data.([]interface{}) {
			flattenJson(value, fmt.Sprintf("%s%d.", prefix, i), flatMap)
		}
	default:
		flatMap[prefix] = fmt.Sprintf("%v", data)
	}
}

func findMatches(source, target map[string]string) []Match {
	var matches []Match
	for sourceKey, sourceVal := range source {
		for targetKey, targetVal := range target {
			if sourceVal == targetVal {
				matches = append(matches, Match{Value: sourceVal, SourceKey: sourceKey, TargetKey: targetKey})
			}
		}
	}
	return matches
}
