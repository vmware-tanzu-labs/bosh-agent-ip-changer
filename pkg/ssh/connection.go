package ssh

import (
	"bufio"
	"context"
	boshhttp "github.com/cloudfoundry/bosh-utils/httpclient"
	proxy "github.com/cloudfoundry/socks5-proxy"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
}

func New() *Connection {
	return &Connection{}
}

func (c *Connection) Execute(address, user, password string, commands ...string) ([]byte, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	conn, err := newSSHClient(address, sshConfig)
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	defer func() { _ = session.Close() }()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4k
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4k
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return nil, err
	}

	in, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}

	out, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	var output []byte

	go func(in io.WriteCloser, out io.Reader, output *[]byte) {
		var (
			line string
			r    = bufio.NewReader(out)
		)
		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}

			*output = append(*output, b)

			if b == byte('\n') {
				line = ""
				continue
			}

			line += string(b)

			if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
				_, err = in.Write([]byte(password + "\n"))
				if err != nil {
					break
				}
			}
		}
	}(in, out, &output)

	cmd := strings.Join(commands, "; ")
	_, err = session.Output(cmd)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func newSSHClient(address string, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {
	dialer := net.Dialer{}
	dialContextFunc := dialer.DialContext

	socksProxy := proxy.NewSocks5Proxy(proxy.NewHostKey(), log.New(io.Discard, "", log.LstdFlags), 1*time.Minute)
	dialContextFunc = boshhttp.SOCKS5DialContextFuncFromEnvironment(&net.Dialer{}, socksProxy)

	conn, err := dialContextFunc(context.Background(), "tcp", address)
	if err != nil {
		return nil, err
	}

	c, ch, reqs, err := ssh.NewClientConn(conn, address, sshConfig)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(c, ch, reqs), nil
}
