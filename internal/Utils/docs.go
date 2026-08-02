package Utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	CHANGELOG = "changelog.md"
	README    = "readme.md"
	LICSENCE  = "license"
)

func Docs() {
	filepath.WalkDir(".", func(p string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		fileName := strings.ToLower(d.Name())

		if fileName == CHANGELOG {
			fmt.Println("Found CHANGELOG")
			return nil
		}

		if fileName == README {
			fmt.Println("Found README")
			return nil
		}

		if fileName == LICSENCE {
			fmt.Println("Found LICSENCE")
			return nil
		}

		return nil
	})
}
