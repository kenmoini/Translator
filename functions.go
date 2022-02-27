package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	ghodssyaml "github.com/ghodss/yaml"
)

// Read a directory and return a list of files
func ReadDir(dirname string, recursive bool) []string {
	// Open the directory
	f, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	// Read the directory
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	// Close the directory
	defer f.Close()

	matchedFiles := []string{}

	// Loop through the files and directories
	for _, v := range files {
		// If this is a directory, loop through it and add all the files if recursive
		if v.IsDir() && recursive {
			matchedFiles = append(matchedFiles, ReadDir(dirname+"/"+v.Name(), recursive)...)
		} else {
			if !v.IsDir() {
				matchedFiles = append(matchedFiles, dirname+"/"+v.Name())
			}
		}
	}

	return matchedFiles
}

func TranslateMarkdown(filename string) {

	// Check to see if this is a markdown file
	if filepath.Ext(filename) == ".md" {

		// Open the file to test
		file, err := os.Open(filename)
		checkError(err)

		defer file.Close()

		// Extract the frontmatter
		_, _, frontMatter, err := ExtractFrontmatter(file)
		checkError(err)

		//fmt.Printf("Resulting struct: %#v\n", frontMatter)
		//fmt.Printf("Resulting struct: %v\n", frontMatter)
		//fmt.Printf("Targeting: %v\n", FrontMatterTargets)

		// See if any of the targets exist in the extracted frontmatter
		for k, v := range frontMatter {
			if isValueInList(k, FrontMatterTargets) {
				fmt.Printf("%v: %v\n", k, v)
			}
		}

		// Read the file
		/*
			dat, err := os.ReadFile(filename)
			checkError(err)

			// If so, get front matter
			frontMatter, rest, err := GetFrontmatter(dat)
			checkError(err)

			fmt.Printf("%+v\n", frontMatter)
			fmt.Println(string(rest))
		*/

		// Translate title
		// Translate description
		// Translate the body
	} else {
		println("Not a markdown file!")
	}
}

// ExtractFrontmatter reads the frontmatter from a file and supports YAML, TOML, and JSON, returning the parsed data in JSON format
func ExtractFrontmatter(file *os.File) (string, int, map[string]interface{}, error) {

	// Create a new buffer to read the file
	scanner := bufio.NewScanner(file)

	var line int
	var frontMatterType string
	var frontMatterClosing string
	var frontMatterLastLine int
	var frontMatterData []byte
	var frontMatterStruct map[string]interface{}

	// Loop through the file
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		if line == 0 {
			switch scanner.Text() {
			case "---", "---yaml", "---yml":
				frontMatterType = "yaml"
				frontMatterClosing = "---"
			case "---toml":
				frontMatterType = "toml"
				frontMatterClosing = "---"
			case "+++":
				frontMatterType = "toml"
				frontMatterClosing = "+++"
			case ";;;":
				frontMatterType = "json"
				frontMatterClosing = ";;;"
			case "---json":
				frontMatterType = "json"
				frontMatterClosing = "---"
			case "{":
				frontMatterType = "json"
				frontMatterClosing = "}"
			default:
				frontMatterType = "unknown"
				frontMatterClosing = "unknown"
			}
		} else {
			// Look for the closing frontmatter string
			if scanner.Text() == frontMatterClosing {
				frontMatterLastLine = line
				break
			}
		}

		// Add the Frontmatter to the string slice
		frontMatterData = append(frontMatterData, scanner.Text()+"\n"...)

		// Increment the line counter
		line++
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
		return frontMatterType, frontMatterLastLine, frontMatterStruct, err
	}

	// General debug output
	//fmt.Println("===== Frontmatter format: ", frontMatterType)
	//fmt.Println("===== Frontmatter last line number: ", frontMatterLastLine)
	//fmt.Println("===== Frontmatter data:")
	//fmt.Printf("%s", frontMatterData)

	// Convert the FrontmatterData YAML to JSON
	jsonDoc, err := ghodssyaml.YAMLToJSON(frontMatterData)
	if err != nil {
		fmt.Printf("Error converting YAML to JSON: %s\n", err.Error())
		return frontMatterType, frontMatterLastLine, frontMatterStruct, err
	}

	err = json.Unmarshal(jsonDoc, &frontMatterStruct)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %s\n", err.Error())
		return frontMatterType, frontMatterLastLine, frontMatterStruct, err
	}

	//fmt.Printf("Resulting struct: %#v\n", frontMatterStruct)

	return frontMatterType, frontMatterLastLine, frontMatterStruct, nil
}
