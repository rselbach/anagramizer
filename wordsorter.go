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
	"sort"
	"sync"
	)

type WordSorter struct {
	lock sync.Mutex
	reverse bool
	words []string
	count int
}

func (w *WordSorter) Init() {
	w.lock.Lock()
	w.reverse = false
	w.count = 0
	w.words = make([]string, 0)
	w.lock.Unlock()
}

func (w *WordSorter) Len() int {
	return w.count
}

func (w *WordSorter) Swap(i, j int) {
	w.lock.Lock()
	tmp := w.words[i]
	w.words[i] = w.words[j]
	w.words[j] = tmp
	w.lock.Unlock()
}

func (w *WordSorter) Less(i, j int) bool {
	if w.reverse {
		return len(w.words[i]) > len(w.words[j])
	}
	return len(w.words[i]) < len(w.words[j])
}

func (w *WordSorter) Append(data string) {
	w.lock.Lock()
	if w.count+1 > cap(w.words) { // reallocate
		// Allocate double what's needed, for future growth.
		newSlice := make([]string, (w.count+1)*2)
		// The copy function is predeclared and works for any w.words type.
		copy(newSlice, w.words)
		w.words = newSlice
	}
	w.words = w.words[0 : w.count+1]
	w.words[w.count] = data
	w.count++
	w.lock.Unlock()
}

func (w *WordSorter) Words() []string {
	return w.words
}

func (w *WordSorter) Word(i int) string {
	return w.words[i]
}

func (w *WordSorter) Sort() {
	w.lock.Lock()
	w.reverse = false
	w.lock.Unlock()
	sort.Sort(w)
}

func (w *WordSorter) SortReversed() {
	w.lock.Lock()
	w.reverse = true
	w.lock.Unlock()
	sort.Sort(w)
}
