package passcracker

import (
	"crypto/sha256"
	"fmt"
	"password-cracker/pkg/wordhasher"
	"strings"
	"sync"
	"time"
)

type BruteForce struct {
	searchSpace    string
	workerTotal    int
	searchDuration time.Duration
}

func (d *BruteForce) Crack(passwordHash string) *Result {
	now := time.Now()
	done := make(chan interface{})
	passwordStream := d.createWorkers(done, passwordHash)

	select {
	case password, ok := <-passwordStream:
		close(done)
		return &Result{
			Hash:        passwordHash,
			Success:     ok,
			Password:    password,
			TimeElapsed: time.Since(now),
		}
	}
}

func (d *BruteForce) createWorkers(done chan interface{}, passwordHash string) <-chan string {
	passwordStream := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(d.workerTotal)
	windowCalc := windowIndexCalculator{len(d.searchSpace), len(d.searchSpace) / d.workerTotal}

	for i := 0; i < d.workerTotal; i++ {
		go func(i int) {
			defer wg.Done()
			l, r := windowCalc.calculate(i)
			w := worker{searchSpace: d.searchSpace, workerTotal: d.workerTotal, searchDuration: d.searchDuration}
			w.work(done, passwordStream, passwordHash, l, r)
		}(i)
	}

	go func() {
		wg.Wait()
		close(passwordStream)
	}()

	return passwordStream
}

func NewBruteForce(searchSpace string, workers int, searchDuration time.Duration) *BruteForce {
	return &BruteForce{
		searchSpace:    searchSpace,
		workerTotal:    workers,
		searchDuration: searchDuration,
	}
}

type worker struct {
	searchSpace    string
	workerTotal    int
	searchDuration time.Duration
}

func (d *worker) work(done <-chan interface{}, passwordStream chan<- string, passwordHash string, l int, r int) {
	hasher := wordhasher.New(sha256.New())
	baseWord := ""
	timer := time.After(d.searchDuration)

	checkCombinations := func(l int, r int) string {
		for _, char := range d.searchSpace[l:r] {
			password := fmt.Sprintf("%s%c", baseWord, char)
			wordHash := hasher.Hash(password)
			if wordHash == passwordHash {
				return password
			}
		}
		return ""
	}

	password := checkCombinations(l, r)

	for {
		select {
		case <-done:
			return
		case <-timer:
			return
		default:
		}

		if password != "" {
			passwordStream <- password
			return
		}

		baseWord = d.updateBaseWord(baseWord, l, r)
		password = checkCombinations(0, len(d.searchSpace))
	}
}

func (d *worker) updateBaseWord(baseWord string, l int, r int) string {
	switch len(baseWord) {
	case 0:
		return string(d.searchSpace[l])
	case 1:
		currIndex := strings.IndexByte(d.searchSpace, baseWord[0])
		if currIndex < r-1 {
			return string(d.searchSpace[currIndex+1])
		}
	default:
		for i := len(baseWord) - 1; i > 0; i-- {
			currIndex := strings.IndexByte(d.searchSpace, baseWord[i])
			if currIndex < len(d.searchSpace)-1 {
				return baseWord[:i] + string(d.searchSpace[currIndex+1]) + strings.Repeat(string(d.searchSpace[0]), len(baseWord[i+1:]))
			}
		}
		currIndex := strings.IndexByte(d.searchSpace, baseWord[0])
		if currIndex < r-1 {
			return string(d.searchSpace[currIndex+1]) + strings.Repeat(string(d.searchSpace[0]), len(baseWord[1:]))
		}
	}

	return string(d.searchSpace[l]) + strings.Repeat(string(d.searchSpace[0]), len(baseWord))
}
