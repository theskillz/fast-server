http://localhost:8080/ entrypoint
http://localhost:8080/stats stats

```bash
ab -n 40000 -c 1000 "127.0.0.1:8080/"

This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 4000 requests
Completed 8000 requests
Completed 12000 requests
Completed 16000 requests
Completed 20000 requests
Completed 24000 requests
Completed 28000 requests
Completed 32000 requests
Completed 36000 requests
Completed 40000 requests
Finished 40000 requests


Server Software:        
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /
Document Length:        0 bytes

Concurrency Level:      1000
Time taken for tests:   5.303 seconds
Complete requests:      40000
Failed requests:        0
Total transferred:      3000000 bytes
HTML transferred:       0 bytes
Requests per second:    7542.74 [#/sec] (mean)
Time per request:       132.578 [ms] (mean)
Time per request:       0.133 [ms] (mean, across all concurrent requests)
Transfer rate:          552.45 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   12  14.5      7      98
Processing:    30  118  38.6    114     419
Waiting:       13  109  35.6    107     279
Total:         42  130  37.7    126     430

Percentage of the requests served within a certain time (ms)
  50%    126
  66%    140
  75%    150
  80%    156
  90%    175
  95%    196
  98%    226
  99%    246
 100%    430 (longest request)


```