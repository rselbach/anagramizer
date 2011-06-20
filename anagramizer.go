/*
 * anagramizer - An anagram solver in Go
 * 
 * Copyright (c) 2011, Roberto Teixeira <robteix@robteix.com>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *   * Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *
 *   * Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *
 *   * Neither the name of the copyright holder nor the names of its
 *     contributors may be used to endorse or promote products derived from
 *     this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

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

// # of solutions
var solutions uint = 0

func TestAnagram(word, dictword string, ch chan string) {
	if len(dictword) < *minSize {
		return
	}
	if len(dictword) > *maxSize && *maxSize > 0 {
		return
	}
	for _, char := range dictword {
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
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [letters]\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	word := flag.Arg(0)

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
	for {
		line, _, err := r.ReadLine()
		if err == os.EOF {
			break
		}
		if err != nil {
			panic(err)
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
