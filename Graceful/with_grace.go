package WithGrace

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"time"
)

type WithGrace struct {
	*net.TCPListener          // standard libary listener embeded
	shutdownSignaled chan int // control channel
}

var ErrTypeConversion = errors.New("net.Listen returned an invalid type - can not convert")
var ErrShutdownPending = errors.New("Shutdown in progress - no new connections accepted")

func NewWithGraceListener(netName, laddr string) (rv *WithGrace, err error) {

	rv = &WithGrace{
		shutdownSignaled: make(chan int),
	}

	var ok bool
	var t net.Listener
	t, err = net.Listen(netName, laddr)
	if err != nil {
		return
	}
	rv.TCPListener, ok = t.(*net.TCPListener)
	if !ok {
		err = ErrTypeConversion
		return
	}

	return
}

func (wg *WithGrace) Accept() (net.Conn, error) {

	for {

		wg.SetDeadline(time.Now().Add(time.Duration(100) * time.Millisecond)) // 10 times a second, check for timeout/shutdown

		newConn, err := wg.TCPListener.Accept() // call the wraped "Accept"

		select {
		case <-wg.shutdownSignaled: // if closed
			return nil, ErrShutdownPending
		default: // see if this is just a regular timeout
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			// If this is a timeout, then continue to wait for new connections		-- xyzzy
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue // -- xyzzy - restructure this, get rid of "continue" -- maybee reverse logic
			}
		}

		return newConn, err
	}
}

func (wg *WithGrace) GracefulShutdownServer() {
	close(wg.shutdownSignaled)
}

func (wg *WithGrace) ListenAndServeGracefully() (err error) {
	var waitForExit sync.WaitGroup

	server := http.Server{}

	waitForExit.Add(1)
	go func() {
		defer waitForExit.Done()
		err = server.Serve(wg)
	}()

	waitForExit.Wait()
	return
}
