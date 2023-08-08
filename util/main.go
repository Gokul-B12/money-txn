package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//RandomInt to generate random interger within a range

func RandomInt(min, max int64) int64 {
	n := min + rand.Int63n(max-min+1)

	return n
}

//Randomstring to generate random string from "const alphabet set"

func RandomString(n int) string {
	var sb strings.Builder

	len := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len)]

		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(5)
}

func RandomBalance() int64 {
	return RandomInt(1, 1000)
}

func RandomCurrency() string {
	s := []string{"INR", "USD", "EUR"}

	n := len(s)

	c := s[rand.Intn(n)]

	return c
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(5))
}
