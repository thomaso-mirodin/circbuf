package circbuf

import (
	"bytes"
	"io"
	"testing"
)

func TestBuffer_Impl(t *testing.T) {
	var _ io.Writer = &Buffer{}
	var _ io.Reader = &Buffer{}
}

func TestBuffer_ShortWrite(t *testing.T) {
	buf, err := NewBuffer(1024)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	inp := []byte("hello world")

	n, err := buf.Write(inp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if n != len(inp) {
		t.Fatalf("bad: %v", n)
	}

	if !bytes.Equal(buf.Bytes(), inp) {
		t.Fatalf("bad: %v", buf.Bytes())
	}
}

func TestBuffer_ShortRead(t *testing.T) {
	buf, err := NewBuffer(1024)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	inp := []byte("hello world")

	n, err := buf.Write(inp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if n != len(inp) {
		t.Fatalf("bad: %v", n)
	}

	out := make([]byte, len(inp)-2)
	buf.Read(out)

	expected := []byte("hello wor")
	if !bytes.Equal(out, expected) {
		t.Fatalf("bad: %v", buf.Bytes())
	}
}

func TestBuffer_FullWrite(t *testing.T) {
	inp := []byte("hello world")

	buf, err := NewBuffer(int64(len(inp)))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	n, err := buf.Write(inp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if n != len(inp) {
		t.Fatalf("bad: %v", n)
	}

	if !bytes.Equal(buf.Bytes(), inp) {
		t.Fatalf("bad: input=\"%v\" output=\"%v\"", inp, buf.Bytes())
	}
}

func TestBuffer_FullRead(t *testing.T) {
	inp := []byte("hello world")

	buf, err := NewBuffer(int64(len(inp)))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	n, err := buf.Write(inp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if n != len(inp) {
		t.Fatalf("bad: %v", n)
	}

	out := make([]byte, len(inp))
	buf.Read(out)

	if !bytes.Equal(out, inp) {
		t.Fatalf("bad: %v", out)
	}
}

func TestBuffer_LongWrite(t *testing.T) {
	inp := []byte("hello world")

	buf, err := NewBuffer(6)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	n, err := buf.Write(inp)
	if err == nil {
		t.Fatalf("err: %v", buf)
	}
	if int64(n) > buf.Capacity() {
		t.Fatalf("bad: %v", n)
	}

	expect := []byte("hello ")
	if !bytes.Equal(buf.Bytes(), expect) {
		t.Fatalf("bad: %s", buf.Bytes())
	}
}

func TestBuffer_LongRead(t *testing.T) {
	inp := []byte("hello world")

	buf, err := NewBuffer(6)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	n, err := buf.Write(inp)
	if err == nil {
		t.Fatalf("err: %v", buf)
	}

	out := make([]byte, len(inp))
	buf.Read(out)

	expect := []byte("hello ")
	if !bytes.Equal(expect, out[:n]) {
		t.Fatalf("bad: expected=\"%v\" got=\"%v\"", expect, out[:n])
	}
}

// func TestSimpleRead(t *testing.T) {
// 	buf, err := NewBuffer(16)

// 	buf.Write([]byte("hello world"))

// 	out := make([]byte, 6)
// 	_, err := buf.Read(out)

// 	if err != nil {

// 	}
// }

func TestReadBeforeWrite(t *testing.T) {
	buf, err := NewBuffer(8)

	out := make([]byte, 8)
	n, err := buf.Read(out)

	if n != 0 {
		t.Fatalf("err: Read %i bytes without any being written first", n)
	}

	if err != nil {
		t.Fatalf("err: Read should never return an error")
	}

}

func TestReadPastWritePointer(t *testing.T) {
	buf, _ := NewBuffer(16)

	length, _ := buf.Write([]byte("Hello World"))

	out := make([]byte, 16)
	n, err := buf.Read(out)

	if n > length {
		t.Fatal("err: Read past the write cursor")
	} else if n < length {
		t.Fatal("err: Didn't read the full length")
	}

	if err != nil {
		t.Fatal("err: buf.Read should never return an error")
	}
}

func TestBuffer_HugeWrite(t *testing.T) {
	inp := []byte("hello world")

	buf, err := NewBuffer(3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	n, err := buf.Write(inp)
	if err == nil {
		t.Fatalf("err: %v", err)
	}
	if int64(n) > buf.Capacity() {
		t.Fatalf("bad: %v", n)
	}

	expect := []byte("hel")
	if !bytes.Equal(buf.Bytes(), expect) {
		t.Fatalf("bad: %s", buf.Bytes())
	}
}

func TestBuffer_ManySmallWrites(t *testing.T) {
	inp := []byte("hello world")

	buf, err := NewBuffer(3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i, b := range inp {
		n, err := buf.Write([]byte{b})

		if int64(i) < buf.Capacity() {
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			if n != 1 {
				t.Fatalf("bad: %v", n)
			}
		} else {
			if err == nil {
				t.Fatal("err: Write should've failed")
			}

			if n != 0 {
				t.Fatal("bad: Write should've failed")
			}
		}
	}

	expect := []byte("hel")
	if !bytes.Equal(buf.Bytes(), expect) {
		t.Fatalf("bad: %v", buf.Bytes())
	}
}

func TestBuffer_MultiPart(t *testing.T) {
	inputs := [][]byte{
		[]byte("hello world\n"),
		[]byte("this is a test\n"),
		[]byte("my cool input\n"),
	}
	total := 0

	buf, err := NewBuffer(16)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i, b := range inputs {
		n, err := buf.Write(b)
		total += n

		if i == 0 {
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if n != len(b) {
				t.Fatalf("bad: %v", n)
			}
		}
	}

	expect := []byte("hello world\nthis")
	if !bytes.Equal(buf.Bytes(), expect) {
		t.Fatalf("bad: expected=\"%s\" got=\"%s\"", expect, buf.Bytes())
	}
}

func TestBytes(t *testing.T) {
	buf, err := NewBuffer(10)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	in := []byte{1, 2, 3}

	t.Logf("Writing []byte{%v} into the buffer", in)
	switch i, err := buf.Write(in); {
	case err != nil:
		t.Fatalf("Failed to write bytes into buffer, err{%v}, bytesWritten{%v}", err, i)
	case i != len(in):
		t.Fatalf("Failed to write all the bytes into the buffer, bytesWritten{%v} != len{%v}", i, len(in))
	}

	t.Logf("Bytes: %v", buf.Bytes())
	if !bytes.Equal(buf.Bytes(), in) {
		t.Errorf("buf.Bytes(){%v} != in{%v}", buf.Bytes(), in)
	}
}

func TestFree(t *testing.T) {
	buf, err := NewBuffer(10)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	t.Log("Initial stats", buf.Free(), buf.Bytes(), buf)

	t.Log("Bumping writeCursor by 2")
	buf.writeCursor += 2
	if buf.Free() != 8 {
		t.Fatalf("Free() returned '%v', it should've returned 8", buf.Free())
	}
	t.Log("Free bytes:", buf.Free())

	t.Log("Writing two bytes: [1, 2]")
	_, err = buf.Write([]byte{1, 2})
	if err != nil {
		t.Errorf("Unable to write out bytes: err{%v}", err)
	}
	t.Log("Free Bytes:", buf.Free())

	if !bytes.Equal(buf.Bytes(), []byte{0, 0, 1, 2}) {
		t.Fatalf("err: buf.Bytes(){%v} != []byte{[0, 0, 1, 2]}")
	}

	if buf.Free() != 6 {
		t.Fatal("err: buf.Free(){%v} != 6", buf.Free())
	}

	t.Log("Final stats", buf.Free(), buf.Bytes(), buf)
}

func TestProperMod(t *testing.T) {
	q := 13

	for n := -(q - 1); n < q; n++ {
		one := n % q
		two := (n + q) % q
		three := ((n % q) + q) % q

		if two != three {
			t.Errorf("err: (Two != Three) Actual: %v; Expected: %v\n", three, two)
		}

		if n < 0 {
			if one == two {
				t.Errorf("err: (One == Two) Actual: %v; Expected: %v\n", two, one)
			}
			if one == three {
				t.Errorf("err: (One == Three) Actual: %v; Expected: %v\n", three, one)
			}
		}
	}
}

func BenchmarkProperModOne(b *testing.B) {
	q := 13
	for n := 0; n < b.N; n++ {
		_ = n % q
	}
}

func BenchmarkProperModTwo(b *testing.B) {
	q := 13
	for n := 0; n < b.N; n++ {
		_ = (n + q) % q
	}
}

func BenchmarkProperModThree(b *testing.B) {
	q := 13
	for n := 0; n < b.N; n++ {
		_ = ((n % q) + q) % q
	}
}

func BenchmarkBuiltinCopy(b *testing.B) {
	in := make([]byte, 10240)
	out := make([]byte, 10240)

	for i := 0; i < b.N; i++ {
		copy(out, in)
	}
}

func BenchmarkForLoopCopy(b *testing.B) {
	in := make([]byte, 10240)
	out := make([]byte, 10240)

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(in); j++ {
			out[j] = in[j]
		}
	}
}
