package main

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
)

func panicif(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	for _, arg := range os.Args[1:] {
		fi, err := os.Stat(arg)
		panicif(err)

		if fi.IsDir() {
			renameDir(arg)
		} else {
			renameFile(arg)
		}
	}
}

func renameDir(dirname string) {
	files, err := ioutil.ReadDir(dirname)
	panicif(err)

	for _, fi := range files {
		if fi.IsDir() {
			renameDir(path.Join(dirname, fi.Name()))
		} else {
			renameFile(path.Join(dirname, fi.Name()))
		}
	}
}

func renameFile(filename string) {
	fi, err := os.Stat(filename)
	panicif(err)

	ext := path.Ext(fi.Name())
	if !isSub(ext) {
		return
	}

	ok, s, e := getEp(fi.Name())
	if !ok {
		return
	}

	dir := path.Dir(filename)
	files, err := ioutil.ReadDir(dir)
	panicif(err)

	for _, fi2 := range files {
		ext2 := path.Ext(fi2.Name())
		if fi2.IsDir() || !isMov(ext2) {
			continue
		}

		ok, s2, e2 := getEp(fi2.Name())
		if !ok || s != s2 || e != e2 {
			continue
		}
		os.Rename(filename, path.Join(dir, fi2.Name()[0:len(fi2.Name())-len(ext2)]+ext))
		return
	}
}

var subExts = map[string]bool{
	".smi": true,
	".srt": true,
	".ass": true,
}

func isSub(ext string) bool {
	_, found := subExts[ext]
	return found
}

var movExts = map[string]bool{
	".mp4": true,
	".mkv": true,
	".avi": true,
	".flv": true,
}

func isMov(ext string) bool {
	_, found := movExts[ext]
	return found
}

var regs = []*regexp.Regexp{
	regexp.MustCompile("(?i)^(?:[\\S ]+)(?:s(\\d{1,2})e(\\d{1,2})).*$"),
	regexp.MustCompile("(?i)^(?:[\\S ]+)(?:(\\d{1,2})x(\\d{1,2})).*$"),
}

func getEp(name string) (ok bool, s int, e int) {
	var err error
	ok = false
	for _, r := range regs {
		m := r.FindStringSubmatch(name)
		if len(m) < 3 {
			continue
		}

		s, err = strconv.Atoi(m[1])
		if err != nil {
			continue
		}

		e, err = strconv.Atoi(m[2])
		if err != nil {
			continue
		}

		ok = true
		return
	}
	return
}
