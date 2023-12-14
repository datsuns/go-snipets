/*
  Ftp.Get()で以下の様なエラーが出たら？
    An operation on a socket could not be performed because the system lacked
    sufficient buffer space or because a queue was full.

  対策：レジストリをいじる
    https://support.microsoft.com/ja-jp/help/196271/when-you-try-to-connect-from-tcp-ports-greater-than-5000-you-receive-the-error-wsaenobufs-10055

  利用するポート番号が良くない(5000を超えるもの)みたいなのでそこを直せたらいいが・・・
*/

package nxlib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jlaffaye/ftp"
)

type Ftp struct {
	host string
	conn *ftp.ServerConn
}

func NewFtp(host, user, password string) *Ftp {
	return NewFtpWith(host, user, password, DEFAULT_FTP_PORT)
}

func NewFtpWith(host, user, password string, port int) *Ftp {
	f := &Ftp{host: host}
	url := f.host + ":" + strconv.Itoa(port)
	fmt.Println("ftp connect to " + url)
	conn, err := ftp.Connect(url)
	CheckErr(err)
	f.conn = conn
	err = f.conn.Login(user, password)
	CheckErr(err)
	return f
}

func (f *Ftp) Close() {
	f.conn.Quit()
}

func (f *Ftp) List() []string {
	log, err := f.conn.NameList("/root")
	CheckErr(err)
	return log
}

func (f *Ftp) Pwd() string {
	pwd, err := f.conn.CurrentDir()
	CheckErr(err)
	return pwd
}

func (f *Ftp) cd_to_filepath(path string) (string, string, string) {
	dir := filepath.ToSlash(filepath.Dir(path))
	file := filepath.Base(path)
	pwd := f.Pwd()
	err := f.conn.ChangeDir(dir)
	CheckErr(err)
	return pwd, dir, file
}

// jlaffaye/ftp's get(Retr) only support current (or relative) path
func (f *Ftp) Get(path string) {
	var err error
	downloaded := 0

	pwd, _, file := f.cd_to_filepath(path)
	defer f.conn.ChangeDir(pwd)

	r, err := f.conn.Retr(file)
	CheckErr(err)
	defer r.Close()
	dest, err := os.Create(file)
	CheckErr(err)
	defer dest.Close()
	buf := make([]byte, 1024*1024*16)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		} else if n == 0 {
			break
		} else {
			_, err := dest.Write(buf[:n])
			if err != nil {
				panic(err)
			}
			downloaded += n
		}
	}
	fmt.Printf("%d bytes downloaded\n", downloaded)
}

func (f *Ftp) Put(src, dest string) {
	fp, err := os.Open(src)
	CheckErr(err)
	defer fp.Close()
	err = f.conn.Stor(dest, fp)
	CheckErr(err)
}

func (f *Ftp) Delete(path string) {
	err := f.conn.Delete(path)
	CheckErr(err)
}
