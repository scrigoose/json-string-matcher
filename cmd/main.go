package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Match struct {
	Value     string
	SourceKey string
	TargetKey string
}

var skipNums bool

func main() {
	sourceFile, targetFile := parseArgs()

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

	sourceMap := make(map[string][]string)
	targetMap := make(map[string][]string)

	flattenJson(sourceJson, "", sourceMap)
	flattenJson(targetJson, "", targetMap)

	matches := findMatches(sourceMap, targetMap)
	printMatches(matches)

}

func parseArgs() (string, string) {
	flag.BoolVar(&skipNums, "skipnums", false, "Skip comparing numeric values")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go [-skipnums] <source.json> <target.json>")
		os.Exit(1)
	}
	return args[0], args[1]
}

func flattenJson(data interface{}, prefix string, flatMap map[string][]string) {
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
	case reflect.Bool:
		return
	default:
		valStr := fmt.Sprintf("%v", data)
		flatMap[valStr] = append(flatMap[valStr], prefix)
	}
}

func findMatches(source, target map[string][]string) []Match {
	var matches []Match
	for val, sourceKeys := range source {
		if skipNums && isNumeric(val) {
			continue
		}
		if targetKeys, exists := target[val]; exists {
			for _, sKey := range sourceKeys {
				for _, tKey := range targetKeys {
					matches = append(matches, Match{Value: val, SourceKey: sKey, TargetKey: tKey})
				}
			}
		}
	}
	return matches
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func printMatches(matches []Match) {
	for _, match := range matches {
		fmt.Printf("---\nValue: %s\nSource Key: %s\nTarget Key: %s\n", match.Value, match.SourceKey, match.TargetKey)
	}
	fmt.Println("---")
}
