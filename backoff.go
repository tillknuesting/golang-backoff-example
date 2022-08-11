package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	s, err := retryWithExponentialBackoffAndJitter(2)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(s)
}

var count = 0

func sendRequest() (string, error) {
	count++

	if count <= 50 {
		return "", errors.New("error from send request")
	}

	return "success", nil
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}

	return v
}

// Using withExponentialBackoff as the backoff function sends retries
// initially with a 1-second delay, but doubling after each attempt to
// a maximum delay of 1-minute. The jitter is a randomization factor.
// Randomness is generated cryptographic secure and nondeterministic.
func retryWithExponentialBackoffAndJitter(retries int) (string, error) {
	res, err := sendRequest()

	base, caP := time.Second, time.Minute

	var src cryptoSource
	rnd := rand.New(src)

	var retryCounter int

	for backoff := base; err != nil; backoff <<= 1 {
		if retryCounter >= retries {
			return res, fmt.Errorf("maximum retries reached: %w", err)
		}

		if backoff > caP {
			backoff = caP
		}

		jitter := rnd.Int63n(int64(backoff * 3))
		sleep := base + time.Duration(jitter)
		log.Println("Error while sending request, retrying in", sleep)
		time.Sleep(sleep)

		res, err = sendRequest()

		retryCounter++
	}

	return res, nil
}
