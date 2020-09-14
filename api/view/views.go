// Code generated by "esc -pkg view -o api/view/views.go -include .*\.html -prefix=api/view api/view"; DO NOT EDIT.

package view

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/groups.html": {
		name:    "groups.html",
		local:   "api/view/groups.html",
		size:    6191,
		modtime: 1599746508,
		compressed: `
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
`,
	},

	"/": {
		name:  "/",
		local: `api/view`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	"api/view": {
		_escData["/groups.html"],
	},
}
