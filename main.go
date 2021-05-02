package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func clean(directory string) {
	dirRead, _ := os.Open(directory)
	dirFiles, _ := dirRead.Readdir(0)

	for index := range dirFiles {
		fileHere := dirFiles[index]
		nameHere := fileHere.Name()
		if strings.Contains(nameHere, ".pdf") && strings.Contains(nameHere, "page") {
			fullPath := directory + "/" + nameHere
			os.Remove(fullPath)
		}
	}
}

func unite(directory string, outfile string, pages []int) string {
	a := []string{}
	for _, p := range pages {
		a = append(a, fmt.Sprintf("%v/page%d.pdf", directory, p))
	}
	a = append(a, fmt.Sprintf("%v/%v", directory, outfile))
	cmd := exec.Command("pdfunite", a...)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(string(out))
	}
	return outfile
}

func actualNumOfPages(directory string) int {
	totals := 0
	dirRead, _ := os.Open(directory)
	dirFiles, _ := dirRead.Readdir(0)
	for index := range dirFiles {
		fileHere := dirFiles[index]
		nameHere := fileHere.Name()
		if strings.Contains(nameHere, "page") && strings.Contains(nameHere, ".pdf") {
			totals++
		}
	}
	return totals
}

type interval struct {
	Start int
	End   int
}

func NewIntervals(input string) ([]interval, error) {
	res := make([]interval, 0)
	for _, iv := range strings.Split(input, ",") {
		is := strings.Split(iv, "-")
		if len(is) == 2 {
			// case of interval
			from, _ := strconv.Atoi(is[0])
			end, _ := strconv.Atoi(is[1])
			res = append(res, interval{Start: from, End: end})
		} else {
			// single digit
			i, _ := strconv.Atoi(is[0])
			res = append(res, interval{Start: i, End: i})
		}
	}
	return res, nil
}

func isInInterval(page int, interval string) bool {
	if fmt.Sprintf("%d", page) == interval {
		return true
	}
	ints, _ := NewIntervals(interval)
	for _, v := range ints {
		matchLeft := v.Start != 0 && v.Start <= page
		matchRight := v.End != 0 && v.End >= page
		if matchLeft && matchRight {
			// at leas tone interval matches
			return true
		}

	}
	return false
}

var sInput = flag.String("i", "input.pdf", "name of the input file")

func main() {
	fmt.Println("PDF-RM")
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("Command argument must be a page or an interval, i.e. 4,5-10")
	}
	interval := flag.Args()[0]
	if vals, err := NewIntervals(interval); err != nil {
		log.Fatalf("interval parse error %s", err.Error())
	} else {
		fmt.Printf("Pages to be removed: %v\n", vals)
	}

	if _, err := os.Stat(*sInput); os.IsNotExist(err) {
		log.Fatalf("input file not found '%v'", *sInput)
	}

	pdfseparate, err := exec.LookPath("pdfseparate")
	if err != nil {
		log.Fatal("pdfseparate was not found")
	}
	fmt.Printf("pdfseparate is available at %s\n", pdfseparate)

	pdfunite, err := exec.LookPath("pdfunite")
	if err != nil {
		log.Fatal("pdfunite was not found")
	}
	fmt.Printf("pdfunite is available at %s\n", pdfunite)

	outDir := "."
	clean(outDir)

	cmd := exec.Command("pdfseparate", *sInput, "page%d.pdf")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	pageNum := actualNumOfPages(outDir)
	out := make([]int, 0, pageNum)
	for i := 1; i <= pageNum; i++ {
		if !isInInterval(i, interval) {
			out = append(out, i)
		}
	}

	fmt.Printf("input=%s %d -> %d pages\n", *sInput, pageNum, len(out))
	unite(outDir, *sInput, out)
	clean(outDir)
}
