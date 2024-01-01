package main

import "os"

// if pathExists, delete path

func logout(path string) {
	if pathValid(path) {
		os.RemoveAll(path)
	}
}
