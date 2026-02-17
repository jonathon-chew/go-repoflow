package git_utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	TemporaryDirectory string = "./tmp"
)

func FindFilesInCurrentDirectory() (fileList []os.FileInfo) {

	// Look for all the file in the current directory, but not the sub folders!
	var files, issue = os.Open(".")
	if issue != nil {
		fmt.Printf("Error opening directory %s\n", issue)
		os.Exit(1)
	}

	// List the files in the folder
	fileList, err := files.Readdir(0)
	if err != nil {
		fmt.Printf("Error reading directory, %s\n", err)
		os.Exit(1)
	}

	return fileList
}

type FileInfo struct {
	Path     string
	IsDir    bool
	Name     string
	FullPath string
}

func FindAllFilesInCurrentDirectoryAndSubdirectories() []FileInfo {
	fileList := []FileInfo{}
	root, errGettingRoot := os.Getwd()
	if errGettingRoot != nil {
		log.Fatal(errGettingRoot)
		return fileList
	}

	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}

		fileList = append(fileList, FileInfo{
			Path:     path,
			IsDir:    d.IsDir(),
			Name:     info.Name(),
			FullPath: filepath.Join(root, path),
		})
		return nil
	})
	return fileList
}

func MakeDirectoryList(fileList []os.FileInfo) []string {

	// Initialise a list of the directories
	directoryList := []string{}

	// Loop through all the files and check if their a directory or not
	for _, i := range fileList {
		if i.IsDir() {
			directoryList = append(directoryList, i.Name())
		}
	}

	return directoryList
}

func NewDirectory() {
	_, ErrLookingForFile := os.Stat(TemporaryDirectory)
	if ErrLookingForFile == nil {
		return
	}

	ErrMakingDir := os.Mkdir(TemporaryDirectory, os.FileMode(0755))
	if ErrMakingDir != nil {
		log.Fatal(ErrMakingDir)
		return
	}
}

func FindGitRepos(root string) []string {
	var repos []string
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && d.Name() == ".git" {
			repos = append(repos, filepath.Dir(path))
			return filepath.SkipDir // stop traversing this subdir
		}
		return nil
	})
	return repos
}
