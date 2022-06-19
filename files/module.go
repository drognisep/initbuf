package files

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrNotAModule   = errors.New("not a go module")
	errFoundModName = errors.New("found module, stop iterating")
	moduleRegex     = regexp.MustCompile(`^\s*module\s+(\S+)$`)
)

// getCurrentGoModule attempts to find the go module name from a go.mod in the current or parent directory.
func getCurrentGoModuleCwd() (string, error) {
	return getCurrentGoModule(".")
}

// getCurrentGoModule attempts to find the go module name from a go.mod in the base or base parent directory.
func getCurrentGoModule(base string) (string, error) {
	cwd, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}

	volumeName := filepath.VolumeName(cwd)
	volumeNameSlash := volumeName + string(filepath.Separator)
	for cwd != volumeName && cwd != volumeNameSlash {
		var module string
		err := filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
			cwd := cwd
			if err != nil {
				return err
			}
			if d.IsDir() && path != cwd {
				return filepath.SkipDir
			}
			if filepath.Base(path) == "go.mod" {
				var err error
				module, err = parseModule(path)
				if err != nil {
					// Continue iterating
					return nil
				}
				return errFoundModName
			}
			return nil
		})
		if errors.Is(err, errFoundModName) {
			return module, nil
		}
		cwd = filepath.Dir(cwd)
	}
	return "", ErrNotAModule
}

func parseModule(modPath string) (string, error) {
	modFile, err := os.Open(modPath)
	if err != nil {
		return "", err
	}
	defer modFile.Close()
	scanner := bufio.NewScanner(modFile)

	for scanner.Scan() {
		line := scanner.Text()
		if moduleRegex.MatchString(line) {
			groups := moduleRegex.FindStringSubmatch(line)
			if len(groups[1]) > 0 {
				return groups[1], nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", ErrNotAModule
}

func ModuleRoot() (string, error) {
	return moduleRoot(".")
}

func moduleRoot(base string) (string, error) {
	cwd, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}

	volumeName := filepath.VolumeName(cwd)
	volumeNameSlash := volumeName + string(filepath.Separator)

	for cwd != volumeName && cwd != volumeNameSlash {
		var goModDir string
		err := filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
			cwd := cwd
			if err != nil {
				return err
			}
			if d.IsDir() && path != cwd {
				return filepath.SkipDir
			}
			if filepath.Base(path) == "go.mod" {
				goModDir = cwd
				return errFoundModName
			}
			return nil
		})
		if errors.Is(err, errFoundModName) {
			return goModDir, nil
		}
		cwd = filepath.Dir(cwd)
	}
	return "", ErrNotAModule
}
