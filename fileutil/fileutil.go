package fileutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CopyDir(dst, src string) error {
	walkFn := func(srcPath string, info os.FileInfo, err error) error {
		relPath, err := filepath.Rel(src, srcPath)
		if err != nil {
			return fmt.Errorf("failed getting relative path. %s", err)
		}
		dstPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		fmt.Fprintf(os.Stderr, "copying %q to %q\n", srcPath, dstPath)
		return CopyFile(dstPath, srcPath)
	}
	return filepath.Walk(src, walkFn)
}

func CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), filepath.Base(dst)+".")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	const perm = 0644
	if err := os.Chmod(tmp.Name(), perm); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), dst); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return nil
}

func GoFileWalk(src string) (filenames []string, err error) {
	walkFn := func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			filenames = append(filenames, path)
			return nil
		}
		return err
	}
	err = filepath.Walk(src, walkFn)
	return
}
