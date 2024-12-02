package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Copy the file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions
func CopyDir(src string, dst string) error {
	// Get properties of the source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Copy files
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetCommandPath() (err error, exeDir string) {
	exePath, err := os.Executable()
	if err != nil {
		return err, ""
	}
	exeDir = filepath.Dir(exePath)
	return
}

func CreateBackupDir() (err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get project struct error: %s", err.Error())
	}
	err, dir := GetCommandPath()
	if err != nil {
		return fmt.Errorf("get command path error: %s", err.Error())
	}
	cwdModel := filepath.Join(cwd, "model")
	dir1 := filepath.Join(dir, "test1")
	dir2 := filepath.Join(dir, "test2")
	dirModel1 := filepath.Join(dir1, "model")
	dirModel2 := filepath.Join(dir2, "model")
	if !fileExists(filepath.Join(cwd, "api")) ||
		!fileExists(filepath.Join(cwd, "service")) ||
		!fileExists(filepath.Join(cwd, "model")) ||
		!fileExists(filepath.Join(cwdModel, "dto")) {
		return fmt.Errorf("project struct path error: %s", err.Error())
	}
	err = CopyDir(filepath.Join(cwd, "api"), filepath.Join(dir1, "api"))
	if err != nil {
		return fmt.Errorf("copy %s path error: %s", filepath.Join(cwd, "api"), err.Error())
	}
	err = CopyDir(filepath.Join(cwd, "service"), filepath.Join(dir1, "service"))
	if err != nil {
		return fmt.Errorf("copy %s path error: %s", filepath.Join(cwd, "service"), err.Error())
	}
	err = CopyDir(filepath.Join(cwdModel, "dto"), filepath.Join(dirModel1, "dto"))
	if err != nil {
		return fmt.Errorf("copy %s path error: %s", filepath.Join(cwdModel, "dto"), err.Error())
	}

	err = CopyDir(filepath.Join(cwd, "api"), filepath.Join(dir2, "api"))
	if err != nil {
		return fmt.Errorf("copy %s path error: %s", filepath.Join(cwd, "api"), err.Error())
	}
	err = CopyDir(filepath.Join(cwd, "service"), filepath.Join(dir2, "service"))
	if err != nil {
		return fmt.Errorf("copy %s path error: %s", filepath.Join(cwd, "service"), err.Error())
	}
	err = CopyDir(filepath.Join(cwdModel, "dto"), filepath.Join(dirModel2, "dto"))
	if err != nil {
		return fmt.Errorf("copy %s path error: %s", filepath.Join(cwdModel, "dto"), err.Error())
	}
	return nil
}

func CopyAllFileToProject() (err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get project struct error: %s", err.Error())
	}
	err, dir := GetCommandPath()
	if err != nil {
		return fmt.Errorf("get command path error: %s", err.Error())
	}
	cwdModel := filepath.Join(cwd, "model")
	dir2 := filepath.Join(dir, "test2")
	dirModel2 := filepath.Join(dir2, "model")
	err = CopyDir(filepath.Join(dir2, "api"), filepath.Join(cwd, "api"))
	if err != nil {
		return fmt.Errorf("copy %s path error: %s", filepath.Join(dir2, "api"), err.Error())
	}
	err = CopyDir(filepath.Join(dir2, "service"), filepath.Join(cwd, "service"))
	if err != nil {
		return fmt.Errorf("copy %s path error: %s", filepath.Join(dir2, "service"), err.Error())
	}
	err = CopyDir(filepath.Join(dirModel2, "dto"), filepath.Join(cwdModel, "dto"))
	if err != nil {
		return fmt.Errorf("copy %s path error: %s", filepath.Join(dirModel2, "dto"), err.Error())
	}
	return nil
}

func RemoveAllCopyFile() (err error) {
	err, dir := GetCommandPath()
	if err != nil {
		return fmt.Errorf("get command path error: %s", err.Error())
	}
	dir1 := filepath.Join(dir, "test1")
	dir2 := filepath.Join(dir, "test2")
	err = os.RemoveAll(dir1)
	if err != nil {
		return fmt.Errorf("remove test1 path error: %s", err.Error())
	}
	err = os.RemoveAll(dir2)
	if err != nil {
		return fmt.Errorf("remove test2 path error: %s", err.Error())
	}
	return nil
}
