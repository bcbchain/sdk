package fs

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tendermint/tmlibs/log"
)

// PathExists returns whether the given file or directory exists or not
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// MakeDir make dir with Permission 0777
func MakeDir(dir string) (bool, error) {
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UnTarGz takes a destination path and a reader; a tar reader loops over the tar.gz file
// creating the file structure at 'dst' along the way, and writing any files
// nolint gocyclo // 這個方法鬼子寫的，抄來的，不改了
func UnTarGz(dst string, r io.Reader, l log.Logger) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		if l != nil {
			l.Error("unTar", "err", err.Error())
		} else {
			fmt.Println("unTar err : " + err.Error())
		}
		return err
	}
	defer func() {
		if e := gzr.Close(); e != nil {
			if l != nil {
				l.Warn("UnTarGz close Reader Error", "err", err)
			} else {
				fmt.Println("UnTarGz close reader error:", err.Error())
			}
		}
	}()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

			// return any other error
		case err != nil:
			return err

			// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0750); err != nil {
					if l != nil {
						l.Error("unTar MkdirAll", "err", err.Error())
					} else {
						fmt.Println("unTar MkdirAll err : " + err.Error())
					}
					return err
				}
			}

			// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err = io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			if err = f.Close(); err != nil {
				if l != nil {
					l.Warn("file can't be closed", "f", target)
				} else {
					fmt.Println("file can't be closed: " + target)
				}
			}
		}
	}
}
