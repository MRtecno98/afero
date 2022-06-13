package sftpfs

import (
	"net/url"
	"os"
	"os/exec"

	"github.com/pkg/sftp"
	"github.com/spf13/afero"
)

func Resolve(url *url.URL) (afero.Fs, error) {
	usr := url.User.Username()
	if usr != "" {
		usr += "@"
	}

	port := url.Port()
	if port != "" {
		port = "-p" + port
	}

	cmd := exec.Command("ssh", usr+url.Hostname(), port, "-s", "sftp")

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
