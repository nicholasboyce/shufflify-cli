package main

import (
	"math/rand"
)

/*
If seed param is true, then a new seed based on the current time is created.
If false, then the same seed is used, generating the same sequence of pseudorandom numbers every time the program is run.
*/
func Shuffle(list []string) {
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
}
