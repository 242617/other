package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	leadingDotRE    = regexp.MustCompile(`^\.(.*)$`)
	repeatesRE      = regexp.MustCompile(`([_ ]){2,}`)
	leadingDigitsRE = regexp.MustCompile(`^((\d{1,3})\.)(.*)$`)
	spacesRE        = regexp.MustCompile(`([_ ])`)
)

//	rename \
//			-re "haydn_-_(.{2})_(?P<name>.*),_part_(i*)\.flac" \
//			-re-to 'haydn_-_seasons_-_[$1]_${name}_(part_$3).flac'
func init() { log.SetFlags(log.Lshortfile) }
func main() {
	var (
		dryRunFlag       = flag.Bool("n", false, "Dry run")
		prefixString     = flag.String("prefix", "", "Filename prefix")
		renamePrefixFlag = flag.Bool("rename-prefix", true, "Fix prefix")
		removeFilesFlag  = flag.Bool("rm", true, "Remove files (e.g. with leading dot)")
		regexpString     = flag.String("re", "", `Regexp template (can use named groups), haydn_-_(.{2})_(?P<name>.*),_part_(i*)\.flac`)
		regexpToString   = flag.String("re-to", "", `String to replace to (can use $1), haydn_-_seasons_-_[$1]_${name}_(part_$3).flac`)
		checkExtention   = flag.Bool("check-extention", false, "Check extention if known (.mp3)")
	)
	flag.Parse()

	pwd, err := os.Getwd()
	die(err)
	rootPath := filepath.Clean(pwd)

	startBase := filepath.Base(rootPath)
	startDir := filepath.Dir(rootPath)
	root := os.DirFS(startDir)

	err = fs.WalkDir(root, startBase,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			oldBase := filepath.Base(path)

			newBase := oldBase
			if *renamePrefixFlag {
				*prefixString = rename(*prefixString)
			}
			newBase = *prefixString + newBase
			newBase = rename(newBase)
			if *regexpString != "" {
				newBase = regexp.MustCompile(*regexpString).
					ReplaceAllString(newBase, *regexpToString)
			}

			oldPath := fmt.Sprintf("%s/%s/%s", startDir, filepath.Dir(path), oldBase)
			newPath := fmt.Sprintf("%s/%s/%s", startDir, filepath.Dir(path), newBase)

			if *checkExtention {
				ext := filepath.Ext(path)
				switch ext {
				case ".mp2", ".mp3", ".flac", ".aif", ".ogg", ".wma", ".m4a":
				case ".jpg", ".png", ".txt", ".cue", ".nfo":
					// Just skip
					return nil
				default:
					fmt.Printf("[WARNING]\t%q\n\t%s: %q\n\n", oldPath, "unknown file extension", ext)
					return nil
				}
			}

			if *dryRunFlag {
				switch {
				case *removeFilesFlag && leadingDotRE.MatchString(oldBase):
					fmt.Printf("[DELETE]\t%q\n\n", oldPath)

				default:
					if oldPath == newPath {
						break
					}
					fmt.Printf("[RENAME]\t%q\n\tinto\t%q\n\n",
						oldPath, newPath)

				}
				return nil
			}

			switch {
			case *removeFilesFlag && leadingDotRE.MatchString(oldBase):
				return os.Remove(oldPath)
			default:
				return os.Rename(oldPath, newPath)
			}
		},
	)
	die(err)
}

type renamer = func(string) string

func rename(str string) string {
	for _, rename := range []renamer{
		strings.TrimSpace,
		strings.ToLower,
		func(str string) string { return repeatesRE.ReplaceAllString(str, "$1") },
		func(str string) string { return leadingDigitsRE.ReplaceAllString(str, `$2 $3`) },
		func(str string) string { return spacesRE.ReplaceAllString(str, `_`) },
	} {
		str = rename(str)
	}
	return str
}

func die(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: %s", file, line, err.Error())
		os.Exit(1)
	}
}
