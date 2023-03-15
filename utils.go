package sharedlib

import "math/rand"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomString(length int) string {
	s := make([]rune, length)
	for i := 0; i < length; i++ {
		if i == 0 {
			s[i] = letters[rand.Intn(len(letters)-10)]
			continue
		}
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}
