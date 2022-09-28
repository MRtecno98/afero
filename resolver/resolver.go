package resolver

import (
	"errors"
	"net/url"
	"runtime"
	"strings"

	"github.com/MRtecno98/afero"
	"github.com/MRtecno98/afero/sftpfs"
)

var protocols = map[string]func(*url.URL) (afero.Fs, error){
	"mem": func(*url.URL) (afero.Fs, error) {
		return afero.NewMemMapFs(), nil
	},

	"file": func(u *url.URL) (afero.Fs, error) {
		path := u.Path
		if runtime.GOOS == "windows" {
			path = strings.TrimLeft(path, "/")
			path = strings.TrimLeft(path, "\\")
		}

		return afero.NewBasePathFs(afero.NewOsFs(), path), nil
	},

	"sftp": sftpfs.Resolve,
	"ssh":  sftpfs.Resolve,
}

func init() {
	protocols[""] = protocols["file"] // No scheme: Default protocol
}

func OpenUrl(u string) (afero.Fs, error) {
	url, err := url.Parse(u)

	if err != nil {
		return nil, err
	}

	if proto, ok := protocols[url.Scheme]; ok {
		return proto(url)
	} else {
		return nil, errors.New("protocol not implemented")
	}
}
