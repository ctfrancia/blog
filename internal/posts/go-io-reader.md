+++
title = "Go's io.Reader"
description = "The `io.Reader` interface in Go"
date = "2024-11-14"

[author]
  name = "Christian Francia"
  email = "ctfrancia@gmail.com"
  github = "ctfrancia"
+++

The `io.Reader` interface has a single method, `Read`, which reads data into a byte slice.
The `Read` method returns the number of bytes read and an error, if any.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```
