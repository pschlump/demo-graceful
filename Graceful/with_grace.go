package WithGrace

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"syscall"
	"time"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

type WithGrace struct {
	net.Listener                    // standard libary listener embeded
	shutdownSignaled chan error     // control channel
	waitForExit      sync.WaitGroup //
	shutdown         bool           //
	nBefore          int            // track some stats so it is easier to debug
	nAfter           int            //
	laddr            string         //
}

type WithGraceConn struct {
	net.Conn            //
	wg       *WithGrace //
}

// Let compile staticly check that the types/interfaces are complete
var _ net.Listener = (*WithGrace)(nil)
var _ net.Conn = (*WithGraceConn)(nil)

// var ErrTypeConversion = errors.New("net.Listen returned an invalid type - can not convert")
var ErrShutdownPending = errors.New("Shutdown in progress - no new connections accepted")
var ErrShutdownError = errors.New("Shutdown in started - no new connections accepted")

func NewWithGraceListener(netName, laddr string) (wg *WithGrace, err error) {

	tt, err := net.Listen(netName, laddr)
	if err != nil {
		return
	}
	wg = &WithGrace{
		Listener:         tt,
		shutdownSignaled: make(chan error),
		laddr:            laddr,
	}

	go func() {
		_ = <-wg.shutdownSignaled
		fmt.Printf("%sShutdown Signal Received, AT:%s%s\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset)
		wg.shutdown = true
		wg.shutdownSignaled <- wg.Listener.Close()
	}()

	return
}

func (wg *WithGrace) Accept() (newConn net.Conn, err error) {

	if wg.shutdown {
		fmt.Printf("%sShutdown Pending, request ignored, AT:%s%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
		return nil, ErrShutdownPending
	}

	t, err := wg.Listener.Accept() // call the wraped "Accept"
	if err != nil {
		return
	}

	fmt.Printf("%sNew Connection Returned, AT:%s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
	newConn = WithGraceConn{Conn: t, wg: wg}

	wg.startHandler()
	return
}

func (w WithGraceConn) Close() error {
	w.wg.finishHandler()
	return w.Conn.Close()
}

func (wg *WithGrace) Close() error {
	if wg.shutdown {
		return syscall.EINVAL
	}
	wg.shutdownSignaled <- nil
	return <-wg.shutdownSignaled
}

func (wg *WithGrace) GracefulShutdownServer() {
	fmt.Printf("%sShutdown Called For !!!!!!!!!!!!!!!!!!!, AT:%s%s\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset)
	wg.shutdownSignaled <- ErrShutdownError
}

func (wg *WithGrace) startHandler() {
	fmt.Printf("ADD\n")
	wg.waitForExit.Add(1)
}

func (wg *WithGrace) finishHandler() {
	fmt.Printf("SUB\n")
	if wg.shutdown {
		wg.nAfter++
	} else {
		wg.nBefore++
	}
	wg.waitForExit.Done()
}

func (wg *WithGrace) ListenAndServeGracefully() (err error) {
	server := http.Server{
		Addr:           wg.laddr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	wg.waitForExit.Add(1)
	go func() {
		defer wg.waitForExit.Done()
		err = server.Serve(wg)
	}()

	return
}

func (wg *WithGrace) WaitForTheEnd() {
	wg.waitForExit.Wait()
	fmt.Printf("Before=%d After=%d\n", wg.nBefore, wg.nAfter)
}
