package main

import (
	"fmt"
	"gopkg.in/src-d/go-git-fixtures.v3"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BytesToString converts an array of bytes that need to be converted to a string
// data is the string array of bytes that need to be converted to a string
func BytesToString(data []byte) string {
	return string(data[:])
}

func loadGitModulesFile(searchDir string) {
	path := fixtures.ByTag("submodule").One().Worktree().Root()
	fmt.Println(path)
}

// visitGitDirectoryClosure given a searchDir as a string search for all .git
// folders in the current directory. Obtain their corresponding remote URL and
// verify if these exist in the root gitmodules folder. If the git subproject
// URL doesn't exist add this accordingly.
// Note: In the version I will assume subproject dont have submodules themselves!
func visitGitDirectoryClosure(searchDir string) func(string, os.FileInfo, error) error {
	return func(fp string, fi os.FileInfo, err error) error {
		//fmt.Println("fp: ",fp)
		//fmt.Println("searchDir: ",searchDir)
		if err != nil {
			log.Fatal("Error 1: ", err) // can't walk here,
			return nil                  // but continue walking elsewhere
		}
		if fi.IsDir() {
			file := filepath.Base(fp)
			//fmt.Println("file: ",file)
			gitRootDir := searchDir + "/.git"
			if file == ".git" && fp != gitRootDir {
				cmd := exec.Command("git", "config", "--get", "remote.origin.url")
				cmd.Dir = fp
				byteArray, err := cmd.Output()

				if err != nil {
					log.Fatal("Error 2: ", err) // can't walk here,
					return nil                  // but continue walking elsewhere
				}
				gitSubModuleURL := strings.Trim(BytesToString(byteArray[:]), "\n")
				fmt.Printf("gitSubModulePath: %s gitSubModuleURL: %s\n", gitSubModuleURL, fp)

				//Utilize the go
				loadGitModulesFile(searchDir)

				//Verify gitSubModuleURL does exit in the root gitModules file.

				/*
				//Check to see if it exists in .gitmodules
								cmd = exec.Command("grep", "-c", gitSubModuleURL, ".gitmodules")
								cmd.Dir = searchDir
								byteArray, err = cmd.Output()

								if err != nil {
									fmt.Printf("git submodule add %existsInGitModulesFile %existsInGitModulesFile\n", gitSubModuleURL, fp)
									log.Fatal("Error 3: ", err) // can't walk here,
									return nil                // but continue walking elsewhere
								}
								existsInGitModulesFile := strings.Trim(BytesToString(byteArray[:]), "\n")
								fmt.Println("existsInGitModulesFile: ", existsInGitModulesFile)

								if existsInGitModulesFile == "0" {
									fmt.Printf("git submodule add %existsInGitModulesFile %existsInGitModulesFile\n", gitSubModuleURL, fp)
								}*/
				return nil
			}
		}
		return nil // not a directory.  ignore.
	}
}

func main() {
	searchDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error 4: ", err)
	}
	visitGitDirectory := visitGitDirectoryClosure(searchDir)
	filepath.Walk(searchDir, visitGitDirectory)
}

//find . -type d -name .git -exec sh -c 'cd {} && git config --get remote.origin.url && cd - >/dev/null' \; |
// awk '{printf("grep \"%s\" .gitmodules || echo %s \n",$0,$0);}' | sh
