package resolver

import (
	"errors"
	"net/url"

	"github.com/MRtecno98/afero/sftpfs"
	"github.com/spf13/afero"
)

var protocols = map[string]func(*url.URL) (afero.Fs, error){
	"mem": func(*url.URL) (afero.Fs, error) {
		return afero.NewMemMapFs(), nil
	},

	"sftp": sftpfs.Resolve,
	"ssh":  sftpfs.Resolve,
}

func OpenUrl(u string) (afero.Fs, error) {
	urlparse, err := url.Parse(u)

	if err != nil {
		return nil, err
	}

	if proto, ok := protocols[urlparse.Scheme]; ok {
		return proto(urlparse)
	} else {
		return nil, errors.New("protocol not implemented")
	}
}
