package ircclient

import (
	"net"
	"os"
	"log"
	"bufio"
	"strings"
)

type IRCConn struct {
	conn *net.TCPConn
	bio  *bufio.ReadWriter
	tmgr *ThrottleIrcu
	done chan bool

	Err    chan os.Error
	Output chan string
	Input  chan string
}

func NewIRCConn() *IRCConn {
	return &IRCConn{done: make(chan bool), Output: make(chan string, 50), Input: make(chan string, 50), tmgr: new(ThrottleIrcu), Err: make(chan os.Error, 5)}
}

func (ic *IRCConn) Connect(hostport string) os.Error {
	ic.conn.SetTimeout(1)
	if len(hostport) == 0 {
		return os.NewError("empty server addr, not connecting")
	}
	if ic.conn != nil {
		log.Printf("warning: already connected to " + ic.conn.RemoteAddr().String())
	}
	c, err := net.Dial("tcp", hostport)
	if err != nil {
		return err
	}
	ic.conn, _ = c.(*net.TCPConn)
	ic.bio = bufio.NewReadWriter(bufio.NewReader(ic.conn), bufio.NewWriter(ic.conn))

	go func() {
		for {
			// TODO: err
			s, err := ic.bio.ReadString('\n')
			if err != nil {
				select {
				case d := <-ic.done:
					ic.done <- d
					return
				default:
					ic.Err <- os.NewError("ircmessage: receive: " + err.String())
					ic.Quit()
					return
				}
			}
			s = strings.Trim(s, "\r\n")
			ic.Input <- s
			log.Println("<< " + s)
		}
	}()
	go func() {
		for {
			select {
			case s := <-ic.Output:
				s = s + "\r\n"
				ic.tmgr.WaitSend(s)
				log.Print(">> " + s)
				if _, err = ic.bio.WriteString(s); err != nil {
					ic.Err <- os.NewError("ircmessage: send: " + err.String())
					log.Println("Send failed: " + err.String())
					ic.Quit()
					return
				}
				ic.bio.Flush()
			case d := <-ic.done:
				ic.done <- d
				if d {
					return
				}
			}
		}
	}()

	return nil
}

func (ic *IRCConn) Flush() {
	// TODO implement
	// Should block until all data to server has been sent
	// and bufio flushed
}

func (ic *IRCConn) Quit() {
	ic.conn.Close()
	//close(ic.Output)
	//_, x, y, _ := runtime.Caller(1)
	//log.Println(x)
	//log.Println(y)
	//log.Println("Closing channel")
	close(ic.Input)
	ic.done <- true
}
