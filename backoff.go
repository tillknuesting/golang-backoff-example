package backofff

// Using withExponentialBackoff as the backoff function sends retries
// initially with a 1-second delay, but doubling after each attempt to
// a maximum delay of 1-minute. Each backoff time gets an additional
func withExponentialBackoffAndJitter() string {
	res, err := sendRequest()
	base, cap := time.Second, time.Minute

	for backoff := base; err != nil; backoff <<= 1 {
		if backoff > cap {
			backoff = cap
		}

		jitter := rand.Int63n(int64(backoff * 3))
		sleep := base + time.Duration(jitter)
		time.Sleep(sleep)
		res, err = sendRequest()
	}

	return res
}
