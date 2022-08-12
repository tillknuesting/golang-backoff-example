# golang-backoff-example

So when a client calls a server and fails, e.g., because it timeouts, we tend to retry the request.

This could lead to code like this that retries until it succeeds:


```
    res, err := SendRequest() 
    for err !=nil {
        res, err = SendRequest()
    }
   ```
   
The problem is that if every instance of the downstream service fails, the retry loop leads to potential thousand requests that can degrade the entire network and thus other services - this cascading failure is known as a retry storm. 

If we implement a retry mechanism, we should always include a backoff algorithm.

So if we want to retry, how long should we wait?


A naive solution would be this:



```
res, err := SendRequest() forerr!=nil{
         time.Sleep(2 * time.Second)
         res, err = SendRequest()
     }
```

A fixed-duration backoff delay might work fine if there are a tiny number of retrying instances, but it does not scale, so exponential backoff algorithms are used.


An element of randomness, called jitter, is included to avoid clustering requests.


```
        res, err := SendRequest()
        base, cap := time.Second, time.Minute
        for backoff := base; err != nil; backoff <<= 1 { if backoff > cap {
            backoff = cap
        }
        jitter := rand.Int63n(int64(backoff * 3))
        sleep := base + time.Duration(jitter)
        time.Sleep(sleep)
        res, err = SendRequest()

```

gRPC provides this out of the box https://github.com/grpc/grpc-go/tree/v1.48.0/examples/features/retry


    
