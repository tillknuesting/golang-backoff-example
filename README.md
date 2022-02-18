# golang-backoff-example

So when a client calls a server and fails e.g. because it timeouted we tend to retry the request, right?

This could lead to code like this that retries until it succeeds:


```
    res, err := SendRequest() 
    for err !=nil {
        res, err = SendRequest()
    }
   ```
   
The problem is that if every instance of the downstream service failed the retry loop leads to potential thousand requests that can degrade the entire network and thus other services - this cascading failure is known as retry storm. 

If we implement a retry mechanism we should always include a backoff algorithm.

So if we want to retry, how long should we wait?


A naive solution would this:



```
res, err := SendRequest() forerr!=nil{
         time.Sleep(2 * time.Second)
         res, err = SendRequest()
     }
```

A fixed-duration backoff delay might work fine if you have a very small number of retrying instances, but it doesn't scale and therefore exponential backoff algorithms are used.


To avoid clustering of requests, an element of randomness, called jitter, is included.


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




    
