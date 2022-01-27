# random-stddev

A simple REST service writen in Go, supporting the following GET operations:
```
/random/stddev?requests={r}&length={l}
/random/mean?requests={r}&length={l}
```

which perform {r} concurrent requests to random.org API asking for {l} number of random integers in range 0 to 100.

For each {r} requests, the service calculates standard deviation of the drawn integers and standard deviation of integers drawn in all requests.
Results are returned in JSON.

# Example
```
GET /random/stddev?requests=2&length=5
```

Response:
```json
[
  {
    "stddev": 26.29524671875128,
    "data": [47,47,94,22,83]
  },
  {
    "stddev": 33.030894629119565,
    "data": [43,98,14,4,28]
  },
  {
    "stddev": 31.67964646267379,
    "data": [47,47,94,22,83,43,98,14,4,28]
  }
]
```
