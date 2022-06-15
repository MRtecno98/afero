package sftpfs

import (
	"net/url"
	"os"
	"os/exec"

	"github.com/MRtecno98/afero"
	"github.com/pkg/sftp"
)

func Resolve(url *url.URL) (afero.Fs, error) {
	var arg []string = make([]string, 0, 5)

	usr := url.User.Username()
	if usr != "" {
		usr += "@"
	}

	arg = append(arg, usr+url.Hostname())

	port := url.Port()
	if port != "" {
		arg = append(arg, "-p", port)
	}

	arg = append(arg, "-s", "sftp")

	// log.Print("ssh", arg, len(arg))

	cmd := exec.Command("ssh", arg...)

	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	client, err := sftp.NewClientPipe(stdout, stdin)
	if err != nil {
		return nil, err
	}

	return afero.NewBasePathFs(New(client), url.Path), nil
}
