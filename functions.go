package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	ghodssyaml "github.com/ghodss/yaml"
)

// ReadDir reads a directory and return a list of files, with an optional recursive flag
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

func TranslateString(input string) string {
	return input
}

func ProcessContent(content []byte) []byte {
	return content
}

// TranslateMarkdown translates the markdown file to the target language
func TranslateMarkdown(filename string) ([]byte, error) {

	var fullContent []byte
	var translatedFrontMatter []byte

	// Extract the frontmatter
	frontMatterMeta, lastLine, frontMatter, frontMatterJSON, err := ExtractFrontmatter(filename)
	checkError(err)

	//fmt.Printf("Resulting struct: %#v\n", frontMatter)
	//fmt.Printf("Resulting struct: %v\n", frontMatter)
	//fmt.Printf("Targeting: %v\n", FrontMatterTargets)

	// See if any of the targets exist in the extracted frontmatter and translate it
	for k, v := range frontMatter {
		if isValueInList(k, FrontMatterTargets) {
			// fmt.Printf("%v: %v\n", k, v)
			frontMatter[k] = TranslateString(v.(string))
		}
	}

	// Convert FrontMatter back into YAML
	// TODO: Handle other frontmatter types
	switch frontMatterMeta.Type {
	case "yaml":
		y, err := ghodssyaml.JSONToYAML(frontMatterJSON)

		if err != nil {
			fmt.Printf("err: %v\n", err)
			return []byte{}, err
		}

		//fmt.Println(string(y))
		var buf []byte
		buf = append(buf, frontMatterMeta.StartingLine+"\n"...)
		buf = append(buf, y...)
		buf = append(buf, frontMatterMeta.EndingLine+"\n"...)
		translatedFrontMatter = buf
	}

	// Extract the content with the last line number
	content, err := ExtractContent(filename, lastLine)
	checkError(err)

	//fmt.Printf("%s", content)

	// Glue the frontmatter and content together
	fullContent = append(fullContent, translatedFrontMatter...)
	fullContent = append(fullContent, content...)

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
	return fullContent, nil
}

// ExtractContent skips past the last line of the FrontMatter and returns the content from the file
func ExtractContent(filename string, lastLine int) ([]byte, error) {

	// Open the file to test
	file, err := os.Open(filename)
	checkError(err)

	defer file.Close()

	// Create a new buffer to read the file
	scanner := bufio.NewScanner(file)

	var line int
	var contentData []byte

	// Loop through the file
	for scanner.Scan() {
		//fmt.Printf("%v\n", scanner.Text())
		if line > lastLine {
			// Add the content to the byte slice
			contentData = append(contentData, scanner.Text()+"\n"...)
		}
		// Increment the line counter
		line++
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
		return []byte{}, err
	}

	return contentData, nil
}

// ExtractFrontmatter reads the frontmatter from a file and supports YAML, TOML, and JSON, returning the parsed data in JSON format
func ExtractFrontmatter(filename string) (FrontMatterMeta, int, map[string]interface{}, []byte, error) {

	// Open the file to test
	file, err := os.Open(filename)
	checkError(err)

	defer file.Close()

	// Create a new buffer to read the file
	scanner := bufio.NewScanner(file)

	var line int
	// var frontMatterType string
	var frontMatterClosing string
	var frontMatterLastLine int
	var frontMatterData []byte
	var frontMatterStruct map[string]interface{}
	var frontMatterMeta FrontMatterMeta

	// Loop through the file
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		if line == 0 {
			switch scanner.Text() {
			case "---":
				//frontMatterType = "yaml"
				frontMatterClosing = "---"
				frontMatterMeta.Type = "yaml"
				frontMatterMeta.StartingLine = "---"
				frontMatterMeta.EndingLine = "---"
			case "---yaml":
				// frontMatterType = "yaml"
				frontMatterClosing = "---"
				frontMatterMeta.Type = "yaml"
				frontMatterMeta.StartingLine = "---yaml"
				frontMatterMeta.EndingLine = "---"
			case "---yml":
				// frontMatterType = "yaml"
				frontMatterClosing = "---"
				frontMatterMeta.Type = "yaml"
				frontMatterMeta.StartingLine = "---yml"
				frontMatterMeta.EndingLine = "---"
			case "---toml":
				// frontMatterType = "toml"
				frontMatterClosing = "---"
				frontMatterMeta.Type = "toml"
				frontMatterMeta.StartingLine = "---toml"
				frontMatterMeta.EndingLine = "---"
			case "+++":
				// frontMatterType = "toml"
				frontMatterClosing = "+++"
				frontMatterMeta.Type = "toml"
				frontMatterMeta.StartingLine = "+++"
				frontMatterMeta.EndingLine = "+++"
			case ";;;":
				// frontMatterType = "json"
				frontMatterClosing = ";;;"
				frontMatterMeta.Type = "json"
				frontMatterMeta.StartingLine = ";;;"
				frontMatterMeta.EndingLine = ";;;"
			case "---json":
				// frontMatterType = "json"
				frontMatterClosing = "---"
				frontMatterMeta.Type = "json"
				frontMatterMeta.StartingLine = "---json"
				frontMatterMeta.EndingLine = "---"
			case "{":
				// frontMatterType = "json"
				frontMatterClosing = "}"
				frontMatterMeta.Type = "json"
				frontMatterMeta.StartingLine = "{"
				frontMatterMeta.EndingLine = "}"
			default:
				// frontMatterType = "unknown"
				frontMatterClosing = "unknown"
				frontMatterMeta.Type = "unknown"
				frontMatterMeta.StartingLine = "unknown"
				frontMatterMeta.EndingLine = "unknown"
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
		return frontMatterMeta, frontMatterLastLine, frontMatterStruct, []byte{}, err
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
		return frontMatterMeta, frontMatterLastLine, frontMatterStruct, []byte{}, err
	}

	err = json.Unmarshal(jsonDoc, &frontMatterStruct)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %s\n", err.Error())
		return frontMatterMeta, frontMatterLastLine, frontMatterStruct, []byte{}, err
	}

	//fmt.Printf("Resulting struct: %#v\n", frontMatterStruct)

	return frontMatterMeta, frontMatterLastLine, frontMatterStruct, jsonDoc, nil
}
