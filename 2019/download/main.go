package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"strconv"
)

const (
	OUT = "input.txt"
	URL = "https://adventofcode.com/%d/day/%d/input"
)

var client http.Client

func main() {
	if _, err := os.Stat(OUT); !errors.Is(err, os.ErrNotExist) {
		log.Printf("input file %s already exists", OUT)
		os.Exit(0)
	}

	session := os.Getenv("SESSION")
	if len(session) == 0 {
		log.Println("$SESSION not set")
		os.Exit(0)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error finding out pwd: %s", err.Error())
	}
	pwd = path.Clean(pwd)

	day, err := strconv.ParseInt(path.Base(pwd), 10, 0)
	if err != nil {
		log.Fatalf("error computing day from pwd: %s", err.Error())
	}

	year, err := strconv.ParseInt(path.Base(path.Dir(pwd)), 10, 0)
	if err != nil {
		log.Fatalf("error computing year from pwd: %s", err.Error())
	}

	url := fmt.Sprintf(URL, year, day)
	log.Printf("downloading from %s in session %s...", url, session[:6])

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("error creating cookie jar: %s\n", err.Error())
	}
	client = http.Client{
		Jar: jar,
	}

	if err := downloadInputFile(OUT, url, session); err != nil {
		log.Fatalf("error downloading input file: %s", err.Error())
	}
}

func downloadInputFile(filepath string, url string, session string) error {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  session,
		MaxAge: 300,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error building GET request: %w", err)
	}

	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error while copying data: %w", err)
	}

	return nil
}
