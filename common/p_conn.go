package common

import (
	"net"
	"sync"
	"time"
)

// CloseWriter interface for connections that support half-close
type CloseWriter interface {
	CloseWrite() error
}

// BufConn wraps a net.Conn with buffering capabilities
// It allows reading data that was previously buffered during protocol detection
type BufConn struct {
	Conn         net.Conn  // Underlying network connection
	buf          []byte    // Internal buffer for storing read data
	isRun        bool      // Flag indicating if buffering is active
	bufDataIndex int       // Current position in buffer (unused)
	mu           sync.Mutex // Mutex for thread-safe operations
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (buf *BufConn) Close() error {
	return buf.Conn.Close()
}

// LocalAddr returns the local network address, if known.
func (buf *BufConn) LocalAddr() net.Addr {
	return buf.Conn.LocalAddr()
}

// RemoteAddr returns the remote network address, if known.
func (buf *BufConn) RemoteAddr() net.Addr {
	return buf.Conn.RemoteAddr()
}

func (buf *BufConn) SetDeadline(t time.Time) error {
	return buf.Conn.SetDeadline(t)
}

func (buf *BufConn) SetReadDeadline(t time.Time) error {
	return buf.Conn.SetReadDeadline(t)
}

func (buf *BufConn) SetWriteDeadline(t time.Time) error {
	return buf.Conn.SetWriteDeadline(t)
}

// Write writes data to the underlying connection
func (buf *BufConn) Write(b []byte) (int, error) {
	return buf.Conn.Write(b)
}

// Read reads data from the connection, using the buffer if available
func (buf *BufConn) Read(b []byte) (int, error) {

	if buf.isRun {
		read, err := buf.Conn.Read(b)
		if err != nil {
			return 0, err
		}
		buf.buf = append(buf.buf, b[:read]...)
		return read, err
	} else {
		dstLen := len(b)
		// Check the size of data to be read
		cacheLen := len(buf.buf)
		if cacheLen > 0 {
			// If there is cached data
			if dstLen > len(buf.buf) {
				i := copy(b, buf.buf)
				buf.buf = buf.buf[:0]
				read, err := buf.Conn.Read(b[i:])
				if err != nil {
					return 0, err
				}
				return i + read, nil
			} else {
				i := copy(b, buf.buf[:dstLen])
				buf.buf = buf.buf[dstLen:]
				return i, nil
			}
		} else {
			return buf.Conn.Read(b)
		}
	}
}

// CloseWrite closes the write side of the connection if supported
func (buf *BufConn) CloseWrite() error {
	if v, ok := buf.Conn.(CloseWriter); ok {
		return v.CloseWrite()
	}
	return nil
}

// Start enables buffering mode for protocol detection
func (buf *BufConn) Start() {
	buf.mu.Lock()
	defer buf.mu.Unlock()
	buf.isRun = true
}

// Stop disables buffering mode after protocol detection
func (buf *BufConn) Stop() {
	buf.mu.Lock()
	defer buf.mu.Unlock()
	buf.isRun = false
}
