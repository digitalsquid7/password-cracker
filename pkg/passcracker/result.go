package passcracker

import "time"

type Result struct {
	Hash        string
	Success     bool
	Password    string
	TimeElapsed time.Duration
}
