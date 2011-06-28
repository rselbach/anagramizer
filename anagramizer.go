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
	"os"
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

func main() {

	flag.Parse()

	if flag.NArg() != 1 || len(*delimiter) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [letters]\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	word := strings.ToLower(flag.Arg(0))

	f, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}

	ch := make(chan string)
	r := bufio.NewReader(f)
	s := new(Status)
	if !*quiet {
		s.Start("Identifying anagrams")
	}
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
		go TestAnagram(word, string(line), ch)
	}
	f.Close()
	if !*quiet {
		s.Done()
	}

	if !*quiet {
		s.Start("Compiling results")
	}
	ws := new(WordSorter)
	for i := uint(0); i < solutions; i++ {
		ws.Append(<-ch)
	}
	if !*quiet {
		s.Done()
	}
	if *sortResults {
		if !*quiet {
			s.Start("Sorting results")
		}
		if *reverse {
			ws.SortReversed()
		} else {
			ws.Sort()
		}
		if !*quiet {
			s.Done()
		}
	}
	for i := range ws.Words() {
		if *count > 0 && i >= *count {
			break
		}
		fmt.Printf("%s\n", ws.Word(i))
	}

}
