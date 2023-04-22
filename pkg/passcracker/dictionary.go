package passcracker

import (
	"crypto/sha256"
	"password-cracker/pkg/wordhasher"
	"sync"
	"time"
)

type Dictionary struct {
	words       []string
	workerTotal int
}

func (d *Dictionary) Crack(passwordHash string) *Result {
	now := time.Now()
	done := make(chan interface{})
	passwordStream := d.createWorkers(done, passwordHash)

	select {
	case password, ok := <-passwordStream:
		if ok {
			close(done)
		}
		return &Result{
			Hash:        passwordHash,
			Success:     ok,
			Password:    password,
			TimeElapsed: time.Since(now),
		}
	}
}

func (d *Dictionary) createWorkers(done <-chan interface{}, passwordHash string) <-chan string {
	wg := sync.WaitGroup{}
	passwordStream := make(chan string)
	wg.Add(d.workerTotal)

	for i := 0; i < d.workerTotal; i++ {
		go d.createWorker(done, i, &wg, passwordStream, passwordHash)
	}

	go func() {
		wg.Wait()
		close(passwordStream)
	}()
	return passwordStream
}

func (d *Dictionary) createWorker(done <-chan interface{}, i int, wg *sync.WaitGroup, passwordStream chan<- string, passwordHash string) {
	defer wg.Done()
	windowCalc := windowIndexCalculator{totalSize: len(d.words), windowSize: len(d.words) / d.workerTotal}
	l, r := windowCalc.calculate(i)
	hasher := wordhasher.New(sha256.New())
	for _, word := range d.words[l:r] {
		select {
		case <-done:
			return
		default:
			wordHash := hasher.Hash(word)
			if wordHash == passwordHash {
				passwordStream <- word
			}
		}
	}
}

func NewDictionary(words []string, workers int) *Dictionary {
	return &Dictionary{
		words:       words,
		workerTotal: workers,
	}
}
