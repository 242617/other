package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/242617/other/mdprc/document"
)

func init() { log.SetFlags(log.Lshortfile) }
func main() {
	source, target := flag.String("source", "", "Source path (directory)"), flag.String("target", "", "Target path (directory)")
	flag.Parse()

	if *source == "" {
		log.Fatal(errors.New("empty source"))
	}
	if *target == "" {
		log.Fatal(errors.New("empty target"))
	}

	src, trg := filepath.Clean(*source), filepath.Clean(*target)

	for _, path := range []string{src, trg} {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(errors.Wrap(err, "os open"))
		}
		defer f.Close()
		if info, err := f.Stat(); err != nil {
			log.Fatal(errors.Wrap(err, "f stat"))
		} else if !info.IsDir() {
			log.Fatal(fmt.Errorf("path %q is not directory, as needed", path))
		}
	}

	if err := filepath.WalkDir(src, walk(src, trg)); err != nil {
		log.Fatal(errors.Wrap(err, "filepath walk dir"))
	}
}

func walk(source, target string) fs.WalkDirFunc {
	return func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		targetPath, err := filepath.Rel(source, path)
		if err != nil {
			return errors.Wrap(err, "filepath rel")
		}
		targetPath = filepath.Join(target, targetPath)

		info, err := dirEntry.Info()
		if err != nil {
			return errors.Wrap(err, "dir entry info")
		}

		if dirEntry.Type().IsDir() {
			if err := os.MkdirAll(targetPath, info.Mode().Perm()); err != nil {
				return errors.Wrap(err, "mkdir all")
			}
			return nil
		}

		if err := place(path, source, targetPath, info.Mode().Perm(), filepath.Ext(dirEntry.Name()) == ".md"); err != nil {
			return errors.Wrap(err, "place")
		}

		return nil
	}
}

func place(from, source, to string, perm fs.FileMode, process bool) error {
	if !process {
		b, err := os.ReadFile(from)
		if err != nil {
			return errors.Wrap(err, "os read file")
		}
		if err := os.WriteFile(to, b, perm); err != nil {
			return errors.Wrap(err, "os write file")
		}
		return nil
	}

	document, err := document.Parse(source, from)
	if err != nil {
		return errors.Wrap(err, "parse")
	}

	if document.SkipPlace() {
		return nil
	}

	if err := os.WriteFile(to, document.Bytes(), perm); err != nil {
		return errors.Wrap(err, "os write file")
	}
	return nil
}
