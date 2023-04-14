package smb

import (
	"errors"
	"net"
	"os"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/spf13/afero"
)

// Fs is an afero.Fs implementation that enables access to an SMB share using functions provided by the
// github.com/hirochachacha/go-smb2 library.
// Because of the nature of SMB, there are two important differences compared to other afero.Fs implementations:
//   - Chown is not supported and will always return an error.
//   - Some SMB-specific operations need to be performed before and after using the smb.Fs instance, breaking the afero.Fs abstraction.
//
// This leads to the following usage pattern:
//  1. Create a new smb.Fs instance using the NewSmbFs function.
//  2. Call the Mount method to mount a share.
//  3. Use smb.Fs as drop-in replacement for any other afero.Fs implementation.
//  4. Call the Umount method to unmount the share and the Logoff method to close the connection to the SMB server.
type Fs struct {
	conn    net.Conn
	session *smb2.Session
	share   *smb2.Share
}

func NewSmbFs(server, user, password string) (*Fs, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: password,
		},
	}

	session, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}

	return &Fs{conn: conn, session: session}, nil
}

func (s *Fs) errorIfNotMounted() error {
	if s.share == nil {
		return errors.New("need to mount a share first")
	}
	return nil
}

func (s *Fs) Name() string {
	return "smb.Fs"
}

func (s *Fs) Create(name string) (afero.File, error) {
	if err := s.errorIfNotMounted(); err != nil {
		return nil, err
	}
	return s.share.Create(name)
}

func (s *Fs) Mkdir(name string, perm os.FileMode) error {
	if err := s.errorIfNotMounted(); err != nil {
		return err
	}
	return s.share.Mkdir(name, perm)
}

func (s *Fs) MkdirAll(path string, perm os.FileMode) error {
	if err := s.errorIfNotMounted(); err != nil {
		return err
	}
	return s.share.MkdirAll(path, perm)
}

func (s *Fs) Open(name string) (afero.File, error) {
	if err := s.errorIfNotMounted(); err != nil {
		return nil, err
	}
	return s.share.Open(name)
}

func (s *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if err := s.errorIfNotMounted(); err != nil {
		return nil, err
	}
	return s.share.OpenFile(name, flag, perm)
}

func (s *Fs) Remove(name string) error {
	if err := s.errorIfNotMounted(); err != nil {
		return err
	}
	return s.share.Remove(name)
}

func (s *Fs) RemoveAll(path string) error {
	if err := s.errorIfNotMounted(); err != nil {
		return err
	}
	return s.share.RemoveAll(path)
}

func (s *Fs) Rename(oldname, newname string) error {
	if err := s.errorIfNotMounted(); err != nil {
		return err
	}
	return s.share.Rename(oldname, newname)
}

func (s *Fs) Stat(name string) (os.FileInfo, error) {
	if err := s.errorIfNotMounted(); err != nil {
		return nil, err
	}
	return s.share.Stat(name)
}

func (s *Fs) Chmod(name string, mode os.FileMode) error {
	if err := s.errorIfNotMounted(); err != nil {
		return err
	}
	return s.share.Chmod(name, mode)
}

func (s *Fs) Chown(_ string, _, _ int) error {
	return errors.New("chown is not supported")
}

func (s *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	if err := s.errorIfNotMounted(); err != nil {
		return err
	}
	return s.share.Chtimes(name, atime, mtime)
}

// SMB-specific methods

// Mount mounts the provided share. After calling this method, the smb.Fs instance can be used as an afero.Fs implementation.
func (s *Fs) Mount(sharename string) error {
	share, err := s.session.Mount(sharename)
	if err != nil {
		return err
	}

	s.share = share

	return nil
}

// Umount unmounts the share. After calling this method, the Mount method should be called again before using the smb.Fs instance.
func (s *Fs) Umount() error {
	if err := s.errorIfNotMounted(); err != nil {
		return err
	}

	if err := s.share.Umount(); err != nil {
		return err
	}

	s.share = nil

	return nil
}

// Logoff invalidates the session. After calling this method, the smb.Fs instance is no longer usable.
func (s *Fs) Logoff() error {
	if err := s.conn.Close(); err != nil {
		return err
	}
	return s.session.Logoff()
}
