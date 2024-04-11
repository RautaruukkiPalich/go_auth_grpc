package auth

import (
	"encoding/base64"
	"math/rand"
)

func generatePassword(str string) string {

	const passLength = 15
	var shuffled string

	for len(shuffled) < 30 {
		shuffled += str
	}
	data := []byte(shuffled)

	rand.Shuffle(len(shuffled), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	return base64.StdEncoding.EncodeToString(data)[:passLength]
}