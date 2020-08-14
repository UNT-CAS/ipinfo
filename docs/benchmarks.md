# Benchmarks

We cannot benchmark ipinfo.io, because it is not open source so we cannot test it in localhost. But, this is what [hey](https://github.com/rakyll/hey) tells us after doing 10000 requests to localhost:

- The average request completed in 15 milliseconds.
- 95% of the Requests were completed under 31 milliseconds.
- 10,000 Requests were completed with 200 request concurrency in under 1 second (.899 Seconds)

Now, what does this tell us?  It is fast.  SO fast, that it is negligable to your application.  If you need something FASTER, you should be paying for a service.

```
$ ./hey_linux_amd64 -n 10000 -c 200 http://localhost/

Summary:
  Total:        0.8998 secs
  Slowest:      0.0828 secs
  Fastest:      0.0002 secs
  Average:      0.0173 secs
  Requests/sec: 11113.9437

  Total data:   1860000 bytes
  Size/request: 186 bytes

Response time histogram:
  0.000 [1]     |
  0.008 [894]   |■■■■■■■
  0.017 [4770]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.025 [3226]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.033 [724]   |■■■■■■
  0.042 [171]   |■
  0.050 [20]    |
  0.058 [47]    |
  0.066 [90]    |■
  0.075 [51]    |
  0.083 [6]     |


Latency distribution:
  10% in 0.0089 secs
  25% in 0.0118 secs
  50% in 0.0157 secs
  75% in 0.0204 secs
  90% in 0.0255 secs
  95% in 0.0306 secs
  99% in 0.0642 secs

Details (average, fastest, slowest):
  DNS+dialup:    0.0007 secs, 0.0002 secs, 0.0828 secs
  DNS-lookup:    0.0002 secs, 0.0000 secs, 0.0374 secs
  req write:     0.0000 secs, 0.0000 secs, 0.0320 secs
  resp wait:     0.0158 secs, 0.0002 secs, 0.0519 secs
  resp read:     0.0006 secs, 0.0000 secs, 0.0269 secs

Status code distribution:
  [200] 10000 responses
```
