package nxlib

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ziutek/telnet"
)

type CmdCallback interface {
	Exec(text string)
	Flush(text string)
}

type Telnet struct {
	processor     *telnet.Conn
	timeout       time.Duration
	currentPrompt string
	defaultPrompt string
	verbose       bool
	host          string
	user          string
	password      string
	port          int
}

type DefaultCallBack struct {
}

func (d *DefaultCallBack) Exec(text string) {
	fmt.Print(text)
}

func (d *DefaultCallBack) Flush(text string) {
	fmt.Print(text)
}

type SilentCallBack struct {
}

func (d *SilentCallBack) Exec(text string) {
}

func (d *SilentCallBack) Flush(text string) {
}

type StdErrCallBack struct {
}

func (d *StdErrCallBack) Exec(text string) {
	fmt.Fprintf(os.Stderr, text)
}

func (d *StdErrCallBack) Flush(text string) {
	fmt.Fprintf(os.Stderr, text)
}

type LoggerCallBack struct {
	Path string
	Dest *os.File
}

func NewLoggerCallBack(path string) *LoggerCallBack {
	var err error
	ret := &LoggerCallBack{Path: path}
	ret.Dest, err = os.Create(ret.Path)
	if err != nil {
		panic(err)
	}
	return ret
}

func (d *LoggerCallBack) Exec(text string) {
	fmt.Print(text)
	_, err := d.Dest.WriteString(text)
	if err != nil {
		panic(err)
	}
}

func (d *LoggerCallBack) Flush(text string) {
	fmt.Print(text)
	_, err := d.Dest.WriteString(text)
	if err != nil {
		panic(err)
	}
}

func (d *LoggerCallBack) Clear() {
	d.Dest.Truncate(0)
	d.Dest.Seek(0, 0)
}

func NewTelnetSpecificWithPort(host, user, password string, port int) *Telnet {
	t := &Telnet{host: host, user: user, password: password, currentPrompt: "# ", defaultPrompt: "# ", timeout: DEFAULT_TELNET_TIMEOUT, verbose: true}
	t.loginWith(port)
	return t
}

func NewTelnetSpecific(host, user, password string) *Telnet {
	t := &Telnet{host: host, user: user, password: password, currentPrompt: "# ", defaultPrompt: "# ", timeout: DEFAULT_TELNET_TIMEOUT, verbose: true}
	t.login()
	return t
}

func NewTelnet() *Telnet {
	return NewTelnetSpecific(DEFAULT_HOST, DEFAULT_USER, DEFAULT_PASSWORD)
}

func (t *Telnet) Close() {
	t.processor.Close()
}

func (t *Telnet) SetTimeout(to time.Duration) {
	t.timeout = to
}

func (t *Telnet) Cmd(cmd string) string {
	sendln(t.processor, t.timeout, cmd)
	return t.done(&DefaultCallBack{})
}

func (t *Telnet) CmdSilent(cmd string) string {
	sendln(t.processor, t.timeout, cmd)
	return t.done(&SilentCallBack{})
}

func (t *Telnet) CmdCallback(cmd string, callback CmdCallback) {
	sendln(t.processor, t.timeout, cmd)
	wait_specific_nonreturn(t.processor, t.timeout, t.verbose, t.currentPrompt, callback)
}

func (t *Telnet) Print(cmd string) {
	send(t.processor, t.timeout, cmd)
}

func (t *Telnet) UpdatePrompt(prompt string) {
	t.currentPrompt = prompt
}

func (t *Telnet) ResetPrompt() {
	t.currentPrompt = t.defaultPrompt
}

func (t *Telnet) done(c CmdCallback) string {
	return wait_specific(t.processor, t.timeout, t.verbose, t.currentPrompt, c)
}

// リアルタイムなprintしたいので自分で判定...
func wait_specific(c *telnet.Conn, timeout time.Duration, verbose bool, delim string, callback CmdCallback) string {
	//CheckErr(c.SetReadDeadline(time.Now().Add(timeout)))
	var all []byte
	var log []byte
	for {
		b, err := c.ReadByte()
		CheckErr(err)
		all = append(all, b)
		log = append(log, b)
		if len(all) < len(delim) {
			continue
		}
		last := all[len(all)-len(delim):]
		if bytes.Compare(last, []byte(delim)) == 0 {
			callback.Flush(string(log))
			//fmt.Printf("log[%v]\n", all)
			//fmt.Printf("last[%v] / delim: [%v]\n", string(last), delim)
			break
		}
		if b == '\n' {
			callback.Exec(string(log))
			//fmt.Print(string(log))
			log = log[:0]
		}
	}
	return string(all)
}

func wait_specific_nonreturn(c *telnet.Conn, timeout time.Duration, verbose bool, delim string, callback CmdCallback) {
	var current []byte
	for {
		b, err := c.ReadByte()
		CheckErr(err)
		current = append(current, b)
		if len(current) < len(delim) {
			continue
		}
		last := current[len(current)-len(delim):]
		if string(last) == delim {
			callback.Flush(string(current))
			break
		}
		if b == '\n' {
			callback.Exec(string(current))
			current = current[:0]
		}
	}
}

func send(t *telnet.Conn, timeout time.Duration, s string) {
	CheckErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s))
	copy(buf, s)
	_, err := t.Write(buf)
	CheckErr(err)
}

func sendln(t *telnet.Conn, timeout time.Duration, s string) {
	send(t, timeout, s+"\n")
}

func (t *Telnet) loginWith(port int) error {
	// connect to this socket
	conn, err := telnet.Dial("tcp", t.host+":"+strconv.Itoa(port))

	if err != nil {
		fmt.Printf("Some error %v", err)
		return err
	}

	t.processor = conn
	wait_specific(conn, t.timeout, t.verbose, "login: ", &StdErrCallBack{})
	sendln(conn, t.timeout, t.user)
	wait_specific(conn, t.timeout, t.verbose, "Password:", &StdErrCallBack{})
	sendln(conn, t.timeout, t.password)
	wait_specific(conn, t.timeout, t.verbose, t.currentPrompt, &StdErrCallBack{})
	return nil
}

func (t *Telnet) login() error {
	return t.loginWith(DEFAULT_TELNET_PORT)
}
