package sqlitevfs

import (
	"os"

	"github.com/MRtecno98/afero"
	vfs "github.com/psanford/sqlite3vfs"
)

var flagmap = map[vfs.OpenFlag]int{
	vfs.OpenReadOnly:      os.O_RDONLY,
	vfs.OpenReadWrite:     os.O_RDWR,
	vfs.OpenCreate:        os.O_CREATE,
	vfs.OpenDeleteOnClose: 0x00000000,
	vfs.OpenExclusive:     0x00000000,
	vfs.OpenAutoProxy:     0x00000000,
	vfs.OpenURI:           0x00000000,
	vfs.OpenMemory:        0x00000000,
	vfs.OpenMainDB:        0x00000000,
	vfs.OpenTempDB:        0x00000000,
	vfs.OpenTransientDB:   0x00000000,
	vfs.OpenMainJournal:   0x00000000,
	vfs.OpenTempJournal:   0x00000000,
	vfs.OpenSubJournal:    0x00000000,
	vfs.OpenSuperJournal:  0x00000000,
	vfs.OpenNoMutex:       0x00000000,
	vfs.OpenFullMutex:     0x00000000,
	vfs.OpenSharedCache:   0x00000000,
	vfs.OpenPrivateCache:  0x00000000,
	vfs.OpenWAL:           0x00000000,
	vfs.OpenNoFollow:      0x00000000,
}

type AferoVFS struct {
	Fs afero.Afero
}

type AferoVFSFile struct {
	afero.File

	VFS *AferoVFS
}

func RegisterVFS(name string, fs afero.Afero) {
	vfs.RegisterVFS(name, &AferoVFS{Fs: fs})
}

// Convert vfs.OpenFlag to os.OpenFile flags
func MapFlags(flags vfs.OpenFlag) (int, vfs.OpenFlag) {
	var osflags int
	var usedflags vfs.OpenFlag
	for f, o := range flagmap {
		if o != 0 && flags&f != 0 {
			osflags |= o
			usedflags |= f
		}
	}

	return osflags, usedflags
}

func (a *AferoVFS) Open(name string, flags vfs.OpenFlag) (vfs.File, vfs.OpenFlag, error) {
	osflags, usedflags := MapFlags(flags)

	f, err := a.Fs.OpenFile(name, osflags, 0666)
	if err != nil {
		return nil, 0, err
	}

	return &AferoVFSFile{File: f, VFS: a}, usedflags, nil
}

func (a *AferoVFS) Delete(name string, dirSync bool) error {
	return a.Fs.Remove(name)
}

func (a *AferoVFS) Access(name string, flags vfs.AccessFlag) (bool, error) {
	if flags == vfs.AccessExists {
		exists, err := a.Fs.Exists(name)
		if err != nil || !exists {
			return false, err
		}
	}

	info, err := a.Fs.Stat(name)

	if err != nil {
		if flags == vfs.AccessExists && os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	var res bool = true
	switch flags {
	case vfs.AccessReadWrite:
		res = res && info.Mode().Perm()&0600 != 0
	case vfs.AccessRead:
		res = res && info.Mode().Perm()&0400 != 0
	}

	return res, nil
}

func (a *AferoVFS) FullPathname(name string) string {
	return name
}

func (f *AferoVFSFile) FileSize() (int64, error) {
	info, err := f.File.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func (f *AferoVFSFile) Sync(flag vfs.SyncType) error {
	return f.File.Sync()
}

func (f *AferoVFSFile) Close() error {
	return f.File.Close()
}

func (f *AferoVFSFile) SectorSize() int64 {
	return 0
}

func (f *AferoVFSFile) DeviceCharacteristics() vfs.DeviceCharacteristic {
	return 0
}

// TODO: Implement locking functions

func (f *AferoVFSFile) Lock(elock vfs.LockType) error {
	return nil
}

func (f *AferoVFSFile) Unlock(elock vfs.LockType) error {
	return nil
}

func (f *AferoVFSFile) CheckReservedLock() (bool, error) {
	return false, nil
}
