package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var runtime float64
	// fmt.Println("| day | time |")
	fmt.Println("|-----|-----:|")
	for _, f := range os.Args[1:] {
		fd, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}

		input := bufio.NewScanner(fd)
		times := make([]float64, 0, 2)
		for input.Scan() {
			line := input.Text()
			line = strings.Replace(line, "\\|", "", 1)
			if strings.Contains(line, "input.txt") {
				args := strings.Fields(strings.Split(line, "|")[2])
				f, _ := strconv.ParseFloat(args[0], 64)
				times = append(times, f)
			}
		}
		fd.Close()
		labels := strings.Split(f, "/")
		label := labels[len(labels)-2]
		fmt.Printf("| %s | %.1f |\n", label, times[1]-times[0])
		runtime += times[1] - times[0]
	}
	fmt.Printf("| total | %.1f |\n", runtime)
}
