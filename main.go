package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	cheatExtension  = ".cht"
	fixedFileSuffix = "_fixed"
	keyDesc         = "desc"
	keyCode         = "code"
	keyEnable       = "enable"
)

type (
	Cheat struct {
		Desc   string
		Code   string
		Enable bool
	}
)

func main() {
	var dirName, fileName, newFile string

	flag.StringVar(&fileName, "in", "", "Name of the file to fix.")
	flag.StringVar(&dirName, "dir", "", "Directory to parse for files to fix.")
	flag.StringVar(&newFile, "out", "", "Name of the file to output too.")

	flag.Parse()

	if fileName == "" && dirName == "" {
		panic("You must provide a file name or directory to fix.")
	}

	if fileName != "" {
		dirName = ""
		pathName := getPath(fileName)
		fileName = pathName + strings.TrimSuffix(fileName, cheatExtension)

		readFile(fileName, newFile)
	}

	if dirName != "" {
		parseDirectory(dirName)
	}
}

func parseDirectory(dir string) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		panic(err)
	}

	for _, f := range files {
		fileName := f.Name()
		fileName = dir + string(os.PathSeparator) + fileName

		info, err := os.Stat(fileName)

		if err != nil {
			panic(err)
		}

		if info.IsDir() {
			parseDirectory(fileName)
		} else {
			ext := filepath.Ext(fileName)

			if ext == cheatExtension {
				readFile(fileName, "")
			}
		}
	}
}

func readFile(fileName, newFile string) {
	if newFile == "" {
		newFile = strings.TrimSuffix(fileName, filepath.Ext(fileName))
		newFile = newFile + fixedFileSuffix
	} else {
		newFile = strings.TrimSuffix(newFile, filepath.Ext(newFile))
	}

	newFile = newFile + cheatExtension

	if err := parseAndFix(fileName, newFile); err != nil {
		panic(err)
	}

	fmt.Sprintf("Cheat file %s fixed and saved to %s.\n", fileName, newFile)
}

func getPath(fileName string) string {
	dirName := filepath.Dir(fileName)
	pathName, err := filepath.Abs(dirName)

	if err != nil {
		panic(err)
	}

	pathName += string(os.PathSeparator)

	return pathName
}

func parseAndFix(in, out string) error {
	inFile, err := os.Open(in)

	if err != nil {
		return err
	}

	defer inFile.Close()

	inScan := bufio.NewScanner(inFile)

	line := 0
	cheatGroup := make([]string, 0)
	cheats := make([]Cheat, 0)

	for inScan.Scan() {
		if line > 1 {
			lineText := inScan.Text()

			if lineText != "" {
				cheatGroup = append(cheatGroup, inScan.Text())

				if len(cheatGroup) > 2 {
					cheat := newCheatFromGroup(cheatGroup)
					cheats = append(cheats, parseCheatGroup(cheat)...)
					cheatGroup = make([]string, 0)
				}
			}
		}

		line++
	}

	return outputToFile(out, cheats)
}

func outputToFile(out string, cheats []Cheat) error {
	numOfCheats := len(cheats)
	cheatLine := fmt.Sprintf("cheats = %d\n\n", numOfCheats)
	output := []byte(cheatLine)
	cheatIdx := 0
	outputFmt := `cheat%[1]d_` + keyDesc + ` = "%[2]s"
cheat%[1]d_` + keyCode + ` = "%[3]s"
cheat%[1]d_` + keyEnable + ` = %[4]t `

	for i, c := range cheats {
		line := fmt.Sprintf(outputFmt, cheatIdx, c.Desc, c.Code, c.Enable)

		if i < numOfCheats-1 {
			line += "\n\n"
		}

		lineBytes := []byte(line)
		output = append(output, lineBytes...)
		cheatIdx++
	}

	err := ioutil.WriteFile(out, output, 0777)

	return err
}

func newCheatFromGroup(group []string) Cheat {
	c := Cheat{}

	for _, s := range group {
		lineData := strings.Split(s, " = ")
		lineData[1] = strings.Trim(lineData[1], `"`)

		if strings.Contains(strings.ToLower(lineData[0]), keyDesc) {
			c.Desc = lineData[1]
		}

		if strings.Contains(strings.ToLower(lineData[0]), keyCode) {
			c.Code = strings.Replace(lineData[1], ":", "", -1)
		}

		if strings.Contains(strings.ToLower(lineData[0]), keyEnable) {
			c.Enable = lineData[1] == "true"
		}
	}

	return c
}

func parseCheatGroup(cheat Cheat) []Cheat {
	cheats := make([]Cheat, 0)
	codes := strings.Split(cheat.Code, "+")

	for _, code := range codes {
		newCheat := Cheat{
			Desc:   cheat.Desc,
			Code:   code,
			Enable: cheat.Enable,
		}

		cheats = append(cheats, newCheat)
	}

	return cheats
}
