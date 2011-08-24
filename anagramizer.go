// anagramizer - An anagram solver in Go

// Copyright (c) 2011, Roberto Teixeira <robteix@robteix.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"http"
	"io"
	"os"
	"regexp"
	"strings"
)

var fileName *string = flag.String("f", "wordlist.txt", "Wordlist file to use")
var maxSize *int = flag.Int("max", 0, "Maximum word size (0 for no limit)")
var minSize *int = flag.Int("min", 1, "Minimum word size")
var sortResults *bool = flag.Bool("s", false, "Sort results by word size")
var quiet *bool = flag.Bool("q", false, "Don't show any message except for the solutions")
var count *int = flag.Int("c", 0, "Maximum number of results (or 0 for no limit)")
var reverse *bool = flag.Bool("r", false, "If true, -s will sort from larger to smaller size")
var subAnagrams *bool = flag.Bool("sub", false, "If true, allow sub-anagrams (not all letters required)")
var delimiter *string = flag.String("d", "\n", "Word separator/delimiter.")

// # of solutions
var solutions uint = 0

var wordList *WordSorter

func TestAnagram(word, dictword string, ch chan string) {
	if len(dictword) < *minSize {
		return
	}
	if len(dictword) > *maxSize && *maxSize > 0 {
		return
	}
	for _, char := range strings.ToLower(dictword) {
		if strings.Contains(word, string(char)) {
			word = strings.Replace(word, string(char), "", 1)
		} else {
			return // not a solution
		}
	}
	if *subAnagrams {
		solutions++
		ch <- dictword
	} else {
		if len(word) == 0 {
			solutions++
			ch <- dictword
		}
	}
}

func solutionsHandler(w http.ResponseWriter, r *http.Request) {

	hint := r.FormValue("hint")
        if hint == "" {
                w.WriteHeader(http.StatusInternalServerError)
                w.Header().Set("Content-Type", "text/plain;charset=UTF-8;")
                io.WriteString(w, "Required parameter 'hint' not received.\n")
                return
        }

	// Use a regexp to find all actual characters
	// we already know about
	realCharExp := regexp.MustCompile("[^*]")
	realChars := realCharExp.FindAllString(hint, -1)

	// Replace all '_' in the hint expression for
	// 'any character that's not currently known'
	newr_str := strings.Replace(hint, "*",
		fmt.Sprintf("[^%s]", strings.Join(realChars, "")), -1)
	finalExp := regexp.MustCompile(fmt.Sprintf("^%s$", newr_str))

	io.WriteString(w, fmt.Sprintf(`<html>
<head><title>Possible Solutions for %s</title></head>
<body><h1>Possible Solutions for %s</h1><ul>`, hint, hint));
	// Now go through the word list looking for matches
	for i := range wordList.Words() {
		if finalExp.MatchString(wordList.Word(i)) {
			io.WriteString(w, fmt.Sprintf("<li>%s</li>", wordList.Word(i)))
		}
	}
	io.WriteString(w, "</ul></body></html>");

}

func anagramHandler(w http.ResponseWriter, r *http.Request) {

	word := r.FormValue("word")
	if word == "" {
		w.WriteHeader(http.StatusInternalServerError)
                w.Header().Set("Content-Type", "text/plain;charset=UTF-8;")
                io.WriteString(w, "Required parameter 'word' not received.\n")
		return
	}

	ch := make(chan string, 100)
	go func() {
		for i := range wordList.Words() {
			TestAnagram(word, wordList.Word(i), ch)
		}
		close(ch)
	}()
	ws := new(WordSorter)
	//for i := uint(0); i < solutions; i++ {
	//	ws.Append(<-ch)
	//}

	for w := range ch {
		ws.Append(w)
	}

	if *sortResults {
		if *reverse {
			ws.SortReversed()
		} else {
			ws.Sort()
		}
	}
	io.WriteString(w, fmt.Sprintf("<html><head><title>Anagrams for %s</title></head><body><h1>Anagrams for %s</h1><ul>", word, word));
	for i := range ws.Words() {
		if *count > 0 && i >= *count {
			break
		}
		io.WriteString(w, fmt.Sprintf("<li>%s</li>", ws.Word(i)))
	}
	io.WriteString(w, "</ul></body></html>");

}

func main() {

	flag.Parse()

	f, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(f)
	s := new(Status)
	if !*quiet {
		s.Start("Identifying anagrams")
	}
	wordList = new(WordSorter)
	// convert string to byte
	sep := (*delimiter)[0]
	for {
		line, err := r.ReadSlice(sep)
		if err == os.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if line[len(line)-1] == sep {
                        line = line[:len(line)-1]
                }
		if line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		wordList.Append(string(line))
	}
	f.Close()
	if !*quiet {
		s.Done()
	}
	http.HandleFunc("/anagrams", anagramHandler)
	http.HandleFunc("/solutions", solutionsHandler)
	http.ListenAndServe(":8080", nil)
}
