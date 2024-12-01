package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strconv"
	"text/template"
)

const (
	SRC = "aoc%d.go"
	OUT = "aoc%d_header.go"

	AUTH = "Erik Adelbert"
	MAIL = "erik_AT_adelbert_DOT_fr"
	LINK = "https://github.com/erik-adelbert/aoc"
)

const HEADER = `// aoc{{.Day}}.go --
// advent of code {{.Year}} day {{.Day}}
//
// https://adventofcode.com/{{.Year}}/day/{{.Day}}
// {{.Link}}
//
// (ɔ) {{.Author}} - {{.Mail}}
// -------------------------------------------
// {{.Year}}-12-{{.Day}}: initial commit

`

type header struct {
	Author, Mail, Link string
	Year, Day          int
}

var setup header

func init() {
	setup.getpwd()
	setup.getenv()
}

func main() {
	var err error

	src := fmt.Sprintf(SRC, setup.Day)

	var info fs.FileInfo
	if info, err = os.Stat(src); err != nil {
		fatal(err)
	}

	var data []byte
	if data, err = os.ReadFile(src); err != nil {
		fatal(err)
	}

	if bytes.HasPrefix(data, []byte("// "+src)) {
		warnf("aborting: %s has header", src)
		os.Exit(0)
	}

	var tmpl *template.Template
	if tmpl, err = template.New("header").Parse(HEADER); err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(HEADER)+len(data)))

	if err = tmpl.Execute(buf, setup); err != nil { // write header
		panic(err)
	}

	if _, err = buf.Write(data); err != nil {
		panic(err)
	}

	out := fmt.Sprintf(OUT, setup.Day)
	data = append(bytes.TrimSpace(buf.Bytes()), '\n')
	if err = os.WriteFile(out, data, info.Mode()); err != nil {
		fatal(err)
	}

	if err = os.Rename(out, src); err != nil {
		panic(err)
	}
}

func (h *header) getenv() {
	const (
		MAXAUTH = 127
		MAXMAIL = 254
		MAXLINK = 254

		ENVAUTH = "AUTHOR"
		ENVMAIL = "MAIL"
		ENVLINK = "LINK"
	)

	valids := []struct {
		s   *string
		env string
		max int
		def string
	}{
		{&(h.Author), ENVAUTH, MAXAUTH, AUTH},
		{&(h.Mail), ENVMAIL, MAXMAIL, MAIL},
		{&(h.Link), ENVMAIL, MAXLINK, LINK},
	}

	for _, field := range valids {
		*field.s = os.Getenv(field.env)
		if len(*field.s) == 0 || len(*field.s) >= field.max {
			*field.s = field.def
		}
	}

	return
}

func (h *header) getpwd() {
	var err error

	var pwd string
	if pwd, err = os.Getwd(); err != nil {
		fatal(err)
	}
	pwd = path.Clean(pwd)

	if h.Year, err = strconv.Atoi(path.Base(path.Dir(pwd))); err != nil {
		fatal(err)
	}

	if h.Day, err = strconv.Atoi(path.Base(pwd)); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	log.Fatal(err.Error())
}

func warnf(format string, v ...any) {
	log.Printf(format, v...)
}
