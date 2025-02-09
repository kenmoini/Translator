package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// I get tired of typing this all the time
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func xl(fromLang string, toLang string, xlate string) string {
	// fix URLs because google translate changes [link](http://you.link) to
	// [link] (http://your.link) and it *also* will translate any path
	// components, thus breaking your URLs.
	reg := regexp.MustCompile(`]\([-a-zA-Z0-9@:%._\+~#=\/]{1,256}\)`)
	// get all the URLs with a single RegEx, keep them for later.
	var foundUrls [][]byte = reg.FindAll([]byte(xlate), -1)
	translated, err := translateTextWithModel(toLang, xlate, "nmt")
	checkError(err)
	// a bunch of regexs to fix other broken stuff
	reg = regexp.MustCompile(` (\*\*) ([A-za-z0-9]+) (\*\*)`) // fix bolds (**foo**)
	translated = string(reg.ReplaceAll([]byte(translated), []byte(" $1$2$3")))
	reg = regexp.MustCompile(`&quot;`) // fix escaped quotes
	translated = string(reg.ReplaceAll([]byte(translated), []byte("\"")))
	reg = regexp.MustCompile(`&gt;`) //fix >
	translated = string(reg.ReplaceAll([]byte(translated), []byte(">")))
	reg = regexp.MustCompile(`&lt;`) // fix <
	translated = string(reg.ReplaceAll([]byte(translated), []byte("<")))
	reg = regexp.MustCompile(`&#39;`) // fix '
	translated = string(reg.ReplaceAll([]byte(translated), []byte("'")))
	reg = regexp.MustCompile(` (\*) ([A-za-z0-9]+) (\*)`) // fix underline (*foo*)
	translated = string(reg.ReplaceAll([]byte(translated), []byte("$1$2$3")))
	reg = regexp.MustCompile(`({{)(<)[ ]{1,3}([vV]ideo)`) // fix video shortcodes
	translated = string(reg.ReplaceAll([]byte(translated), []byte("$1$2 video")))
	reg = regexp.MustCompile(`({{)(<)[ ]{1,3}([yY]outube)`) // fix youtube shortcodes
	translated = string(reg.ReplaceAll([]byte(translated), []byte("$1$2 youtube")))
	// Now it's time to go back and replace all the fucked up urls ...
	reg = regexp.MustCompile(`] \([-a-zA-Z0-9@:%._\+~#=\/ ]{1,256}\)`)
	for x := 0; x < len(foundUrls); x++ {
		// fmt.Println("FoundURL: ", string(foundUrls[x]))
		tmp := reg.FindIndex([]byte(translated))
		if tmp == nil {
			break
		}
		t := []byte(translated)
		translated = fmt.Sprintf("%s(%s%s", string(t[0:tmp[0]+1]), string(foundUrls[x][2:]), (string(t[tmp[1]:])))
	}
	return translated
}

// walk through the front matter, etc. and translate stuff
func doXlate(from string, lang string, readFile string, writeFile string) {
	file, err := os.Open(readFile)
	checkError(err)
	defer file.Close()
	xfile, err := os.Create(writeFile)
	checkError(err)
	defer xfile.Close()

	head := false
	code := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ln := scanner.Text()
		if strings.HasPrefix(ln, "{{") {
			xfile.WriteString(ln + "\n")
			continue
		}
		if strings.HasPrefix(ln, "```") { // deal with in-line code
			xfile.WriteString(ln + "\n")
			code = !code
			continue
		}
		if code { // I don't translate code!
			xfile.WriteString(ln + "\n")
			continue
		}
		if string(ln) == "---" { // start and end of front matter
			xfile.WriteString(ln + "\n")
			head = !head
		} else if !head {
			if strings.HasPrefix(ln, "!") { // translate the ALT-TEXT not the image path
				bar := strings.Split(ln, "]")
				desc := strings.Split(bar[0], "[")
				translated := xl(from, lang, desc[1])
				xfile.WriteString("![" + translated + "]" + bar[1] + "\n")
			} else { // blank lines and everything else
				if ln == "" { // handle blank lines.
					xfile.WriteString("\n")
				} else { // everything else
					translated := xl(from, lang, ln)
					xfile.WriteString(translated + "\n")
				}
			}
		} else { // handle header fields
			headString := strings.Split(ln, ":")
			if headString[0] == "title" { // title
				translated := xl(from, lang, headString[1])
				xfile.WriteString(headString[0] + ": " + translated + "\n")
			} else if headString[0] == "description" { // description
				translated := xl(from, lang, headString[1])
				xfile.WriteString(headString[0] + ": " + translated + "\n")
			} else { // all other header fields left as-is
				xfile.WriteString(ln + "\n")
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	xfile.Close()
	file.Close()
}

// isValueInList checks if a value is in the string slice and returns a boolean
func isValueInList(value string, list []string) bool { // Test Written
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// future work for automagically translating all files.
func getFile(from string, path string, lang string) {
	thisDir, err := os.ReadDir(path)
	checkError(err)
	for _, f := range thisDir {
		if f.IsDir() {
			if f.Name() == "images" {
				continue
			}
			//fmt.Println("going into ", path + "/" + f.Name())
			getFile(from, path+"/"+f.Name(), lang) // fucking hell, recursion!
		} else {
			if strings.Split(f.Name(), ".")[0] == "_index" || strings.Split(f.Name(), ".")[0] == "index" {
				fromFile := fmt.Sprintf("%s/%s.%s.md", path, strings.Split(f.Name(), ".")[0], from)
				toFile := fmt.Sprintf("%s/%s.%s.md", path, strings.Split(f.Name(), ".")[0], lang)
				// fmt.Println("From: ", fromFile)
				// fmt.Println(toFile)
				_, err := os.Stat(toFile)
				if !os.IsNotExist(err) {
					if !(strings.Split(f.Name(), ".")[0] == "_index") {
						addReadingTime(fromFile)
						addReadingTime(toFile)
					}
					// fmt.Printf("Already translated:\t %s/index.%s.md\n", path, lang)
					continue
				}
				addReadingTime(fromFile) // get the reading time first.
				// fmt.Printf("Found a file to translate:\t %s/%s\n", path, f.Name())
				fmt.Printf("Translating:\t %s\nto: \t\t%s\n", fromFile, toFile)
				doXlate(from, lang, fromFile, toFile)
				// }
				continue
			}
		}
	}
}

func main() {
	var parameters PathParameters

	// get command line parameters
	flag.StringVar(&parameters.path, "path", "", "Content source path")
	flag.BoolVar(&parameters.recursive, "recursive", false, "Recursive path search")
	flag.StringVar(&parameters.destination, "dest", "", "Target translate content directory")
	flag.StringVar(&parameters.fromLang, "from", "", "Original language code")
	flag.StringVar(&parameters.toLang, "to", "", "Target language code")
	flag.StringVar(&parameters.frontMatterTargets, "frontmatter", "bio,title,description", "Frontmatter fields to translate")
	flag.StringVar(&parameters.googleAuthJSON, "googleauth", "google-secret.json", "Google Auth JSON file path")
	flag.Parse()

	// Check to make sure we have all the parameters we need
	if parameters.path == "" || parameters.destination == "" || parameters.fromLang == "" || parameters.toLang == "" {
		// Exit and display usage information
		fmt.Println("ERROR: Missing required parameters!")
		flag.Usage()
		os.Exit(1)
	}

	// Set scoped variables
	sourcePath := strings.TrimSuffix(parameters.path, "/")
	targetPath := strings.TrimSuffix(parameters.destination, "/")

	// Set global variables
	FromLang = parameters.fromLang
	ToLang = parameters.toLang
	FrontMatterTargets = strings.Split(parameters.frontMatterTargets, ",")

	// Start general output
	fmt.Println("===== Input Parameters =====")
	fmt.Println("- From Language:       ", FromLang)
	fmt.Println("- To Language:         ", ToLang)
	fmt.Println("- From Path:           ", sourcePath)
	fmt.Println("- To Path:             ", targetPath)
	fmt.Println("- FrontMatter Targets: ", FrontMatterTargets)

	// Check to see if the source exists.
	fi, err := os.Stat(sourcePath)
	checkError(err)

	// Set what file mode the source is to determine if this is a one-file or directory
	sourceMode := fi.Mode()

	// Check to see if the target exists
	// If the target doesn't exist, create it based on what the source is
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		if sourceMode.IsDir() {
			// Create the target directory.
			fmt.Println("===== Creating target directory: ", targetPath)
			os.MkdirAll(targetPath, 0755)
		}
		if sourceMode.IsRegular() {
			// Pop the filename off the path and create the parent directory.
			fmt.Println("===== Creating target directory: ", filepath.Dir(targetPath))
			os.MkdirAll(filepath.Dir(targetPath), 0755)
		}
	}

	// Check if the source path is a directory or a file.
	if sourceMode.IsDir() {
		// do directory stuff
		if parameters.recursive {
			fmt.Println("===== Recursive directory mode =====")
		} else {
			fmt.Println("===== Directory mode =====")
		}

		// Get a list of files to translate.
		targetFiles := ReadDir(sourcePath, parameters.recursive)

		// Loop through the slice of files.
		for _, f := range targetFiles {
			toFile := targetPath + strings.Replace(f, sourcePath, "", 1)

			fmt.Println(" - From: ", strings.TrimSuffix(f, "/"))
			fmt.Println("     To: ", toFile)

			//doXlate(FromLang, ToLang, sourcePath, toFile)
		}
	}

	if sourceMode.IsRegular() {
		// we're just doing one file
		fmt.Println("===== Single file mode =====")

		// Get the filename from the source path.
		fileName := filepath.Base(sourcePath)

		// The target path should have already been created.
		// If the target path is a directory, then take the source filename and append it to the target path.
		// If the target path is not a directory and does not exist, then it must be a file that is being written to
		var toFile string

		// Check to see if the target path is a directory or a file.
		tFii, err := os.Stat(targetPath)
		// See if it exists
		if os.IsNotExist(err) {
			// The target path doesn't exist - must be a target filename
			toFile = targetPath
		} else {

			targetMode := tFii.Mode()

			if targetMode.IsDir() {
				toFile = strings.TrimSuffix(targetPath, "/") + "/" + fileName
			} else {
				toFile = targetPath
			}
		}

		fmt.Println(" - From: " + sourcePath)
		fmt.Println("     To: ", toFile)

		// Switch between possible file type parsers
		switch filepath.Ext(sourcePath) {
		case ".md":
			translatedContent, err := TranslateMarkdown(sourcePath)
			checkError(err)

			// Write the translated content to the target file.
			err = ioutil.WriteFile(toFile, translatedContent, 0644)
			checkError(err)

			fmt.Printf("%s", translatedContent)
		}
	}

}
