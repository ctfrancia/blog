# Go's io.Reader

The `io.Reader` interface has a single method, `Read`, which reads data into a byte slice.
The `Read` method returns the number of bytes read and an error, if any.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```
