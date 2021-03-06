circbuf
=======

Don't use this, go use https://github.com/glycerine/rbuf instead.

---
This repository provides the `circbuf` package. This provides a `Buffer` object
that implements a non-overwriting circular (or ring) buffer. It has a fixed size. The non-wrapping nature means that reads will block if all the available data has been read and that writes will block when there is no more data in the buffer is unread.

The buffer implements the `io.Writer` and `io.Reader` interfaces.

It is not safe for use in a shared concurrent situation

Documentation
=============

Full documentation can be found on [Godoc](http://godoc.org/github.com/thomaso-mirodin/circbuf)

Usage
=====

The `circbuf` package is very easy to use:

```go
buf, _ := NewBuffer(6)
buf.Write([]byte("hello world"))

if string(buf.Bytes()) != "hello " {
    panic("should only have the first 6 bytes!")
}

```

