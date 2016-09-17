package WithGrace

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pschlump/demo-graceful/godebug"
)

//
// This is not an exact copy but the technique described is from:
//
//   http://blog.nella.org/zero-downtime-upgrades-of-tcp-servers-in-go/		-- Overview
//	 https://github.com/fvbock/endless										-- MIT Licensed
//	 http://grisha.org/blog/2014/06/03/graceful-restart-in-golang/			-- Nice Explanation of how/why it works
//	 https://github.com/fvbock/endless										-- MIT Licensed
//	 https://rcrowley.org/articles/golang-graceful-stop.html				-- some other licnese
//	 https://github.com/tylerb/graceful										-- MIT Licensed
//
// I use a modified version of the "endless" library in my own server, http://github.com/pschlump/Go-FTL
// It is also not clear to me how this is going to interact with HTTP2.0 and continuously open connections
// to the client.  I know that "graceful" has a bunch of speical code in it for shutdown with HTTP2.0.
//
//

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
		if godebug.DebugOn("grace-db1") {
			fmt.Printf("%sShutdown Signal Received, AT:%s%s\n", godebug.ColorYellow, godebug.LF(), godebug.ColorReset)
		}
		wg.shutdown = true
		wg.shutdownSignaled <- wg.Listener.Close()
	}()

	go func() {
		// Set up channel on which to send signal notifications.  We must use a buffered channel or risk missing the signal
		// if we're not ready to receive when the signal is sent.
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)

		s := <-c // Block until a signal is received.
		if godebug.DebugOn("grace-db1") {
			fmt.Printf("%sGot signal:%v%s\n", godebug.ColorRed, s, godebug.ColorReset)
		}
		wg.GracefulShutdownServer()
	}()

	return
}

func (wg *WithGrace) Accept() (newConn net.Conn, err error) {

	if wg.shutdown {
		if godebug.DebugOn("grace-db1") {
			fmt.Printf("%sShutdown Pending, request ignored, AT:%s%s\n", godebug.ColorRed, godebug.LF(), godebug.ColorReset)
		}
		return nil, ErrShutdownPending
	}

	t, err := wg.Listener.Accept() // call the wraped "Accept"
	if err != nil {
		return
	}

	if godebug.DebugOn("grace-db1") {
		fmt.Printf("%sNew Connection Returned, AT:%s%s\n", godebug.ColorGreen, godebug.LF(), godebug.ColorReset)
	}
	newConn = WithGraceConn{Conn: t, wg: wg}

	wg.startHandler()
	return
}

func (w WithGraceConn) Close() error {
	w.wg.finishHandler()
	if godebug.DebugOn("grace-db1") {
		fmt.Printf("Close Called, AT:%s\n", godebug.LF())
	}
	return w.Conn.Close()
}

func (wg *WithGrace) Close() error {
	if godebug.DebugOn("grace-db1") {
		fmt.Printf("%sOther Close Called ---- TOP , AT:%s%s\n", godebug.ColorCyan, godebug.LF(), godebug.ColorReset)
	}
	if wg.shutdown {
		return syscall.EINVAL
	}
	if godebug.DebugOn("grace-db1") {
		fmt.Printf("Other Close Called, AT:%s\n", godebug.LF())
	}
	wg.shutdownSignaled <- nil
	return <-wg.shutdownSignaled
}

func (wg *WithGrace) GracefulShutdownServer() {
	if godebug.DebugOn("grace-db1") {
		fmt.Printf("%sShutdown Called For, AT:%s%s\n", godebug.ColorYellow, godebug.LF(), godebug.ColorReset)
	}
	wg.shutdownSignaled <- ErrShutdownError
}

func (wg *WithGrace) startHandler() {
	if godebug.DebugOn("grace-db1") {
		fmt.Printf("ADD\n")
	}
	wg.nBefore++
	wg.waitForExit.Add(1)
}

func (wg *WithGrace) finishHandler() {
	if wg.shutdown {
		wg.nAfter++
	}
	if godebug.DebugOn("grace-db1") {
		fmt.Printf("SUB %d %d\n", wg.nBefore, wg.nAfter)
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
		// defer wg.waitForExit.Done()
		defer wg.finishHandler()
		err = server.Serve(wg)
	}()

	return
}

func (wg *WithGrace) WaitForTheEnd() {
	wg.waitForExit.Wait()
	if godebug.DebugOn("grace-db1") {
		fmt.Printf("Before=%d After=%d\n", wg.nBefore, wg.nAfter)
	}
}

/* vim: set noai ts=4 sw=4: */
