package common

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

// mockConn implements net.Conn for testing
type mockConn struct {
	readBuf  *bytes.Buffer
	writeBuf *bytes.Buffer
	closed   bool
}

func newMockConn(data []byte) *mockConn {
	return &mockConn{
		readBuf:  bytes.NewBuffer(data),
		writeBuf: &bytes.Buffer{},
	}
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuf.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuf.Write(b)
}

func (m *mockConn) Close() error {
	m.closed = true
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9090}
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestBufConnStartStop(t *testing.T) {
	mock := newMockConn([]byte("test data"))
	bufConn := &BufConn{Conn: mock}

	if bufConn.isRun {
		t.Error("BufConn should start with isRun=false")
	}

	bufConn.Start()
	if !bufConn.isRun {
		t.Error("After Start(), isRun should be true")
	}

	bufConn.Stop()
	if bufConn.isRun {
		t.Error("After Stop(), isRun should be false")
	}
}

func TestBufConnBuffering(t *testing.T) {
	testData := []byte("Hello, World!")
	mock := newMockConn(testData)
	bufConn := &BufConn{Conn: mock}

	// Start buffering mode
	bufConn.Start()

	// Read data while buffering is active
	buf := make([]byte, 5)
	n, err := bufConn.Read(buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != 5 {
		t.Errorf("Expected to read 5 bytes, got %d", n)
	}
	if string(buf) != "Hello" {
		t.Errorf("Expected 'Hello', got %q", string(buf))
	}

	// Data should be buffered
	if len(bufConn.buf) != 5 {
		t.Errorf("Expected buffer length 5, got %d", len(bufConn.buf))
	}

	// Stop buffering
	bufConn.Stop()

	// Read from buffer
	buf2 := make([]byte, 3)
	n, err = bufConn.Read(buf2)
	if err != nil {
		t.Fatalf("Read from buffer failed: %v", err)
	}
	if n != 3 {
		t.Errorf("Expected to read 3 bytes from buffer, got %d", n)
	}
	if string(buf2) != "Hel" {
		t.Errorf("Expected 'Hel', got %q", string(buf2))
	}

	// Read remaining buffered data and new data
	buf3 := make([]byte, 20)
	n, err = bufConn.Read(buf3)
	if err != nil && err != io.EOF {
		t.Fatalf("Read failed: %v", err)
	}
	expected := "lo, World!"
	if string(buf3[:n]) != expected {
		t.Errorf("Expected %q, got %q", expected, string(buf3[:n]))
	}
}

func TestBufConnWrite(t *testing.T) {
	mock := newMockConn(nil)
	bufConn := &BufConn{Conn: mock}

	testData := []byte("test write")
	n, err := bufConn.Write(testData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, got %d", len(testData), n)
	}
	if mock.writeBuf.String() != string(testData) {
		t.Errorf("Expected %q, got %q", string(testData), mock.writeBuf.String())
	}
}

func TestBufConnClose(t *testing.T) {
	mock := newMockConn(nil)
	bufConn := &BufConn{Conn: mock}

	err := bufConn.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	if !mock.closed {
		t.Error("Underlying connection should be closed")
	}
}

func TestBufConnReadLargeBuffer(t *testing.T) {
	testData := []byte("1234567890")
	mock := newMockConn(testData)
	bufConn := &BufConn{Conn: mock}

	// Start buffering and read partial data
	bufConn.Start()
	buf1 := make([]byte, 5)
	bufConn.Read(buf1)
	bufConn.Stop()

	// Now read with a buffer larger than cached data
	buf2 := make([]byte, 20)
	n, err := bufConn.Read(buf2)
	if err != nil && err != io.EOF {
		t.Fatalf("Read failed: %v", err)
	}

	// Should read cached data (5 bytes) plus remaining data (5 bytes) = 10 total
	expected := "1234567890"
	if string(buf2[:n]) != expected {
		t.Errorf("Expected %q, got %q", expected, string(buf2[:n]))
	}
}
