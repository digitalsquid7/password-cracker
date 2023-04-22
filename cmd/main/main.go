package main

import (
	"fmt"
	"password-cracker/pkg/filereader"
	"password-cracker/pkg/passcracker"
	"strings"
	"time"
)

type Cracker interface {
	Crack(hash string) *passcracker.Result
}

type CrackerMethod struct {
	name     string
	cracker  Cracker
	hashFile string
}

func main() {
	mostUsedPasswords, err := filereader.ReadLines("dictionary.txt")
	if err != nil {
		fmt.Println(err)
	}

	methods := []CrackerMethod{
		{
			name:     "Dictionary",
			cracker:  passcracker.NewDictionary(mostUsedPasswords, 4),
			hashFile: "hashes_dictionary.txt",
		},
		{
			name:     "Brute Force",
			cracker:  passcracker.NewBruteForce("abcdefghijklmnopqrstuvwxyz", 4, 5*time.Second),
			hashFile: "hashes_brute_force.txt",
		},
	}

	fmt.Printf("%-10s | %-11s | %-64s | %-10s | %-15s\n", "Success", "Method", "Hash", "Password", "Time Elapsed")
	fmt.Println(strings.Repeat("=", 119))

	for _, method := range methods {
		hashes, err := filereader.ReadLines(method.hashFile)

		if err != nil {
			fmt.Println(err)
		}

		for _, hash := range hashes {
			result := method.cracker.Crack(hash)
			fmt.Printf("%-10t | %-11s | %-64s | %-10s | %-15s\n", result.Success, method.name, hash, result.Password, result.TimeElapsed)
		}
	}
}
