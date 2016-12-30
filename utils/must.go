package utils

import "log"

func MustFunc(f func() error) {
	if err := f(); err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
}
