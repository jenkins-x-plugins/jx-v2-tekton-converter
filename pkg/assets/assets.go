// Code generated by go-bindata.
// sources:
// resources/git/git-clone.yaml
// DO NOT EDIT!

package assets

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _resourcesGitGitCloneYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x57\xdd\x6e\xdb\x46\x13\xbd\xd7\x53\x9c\x30\x06\x2c\x7d\x88\xa8\x2f\xb7\x2a\xdc\x20\x88\x05\x38\xa8\x5b\x1b\x96\x9c\xa6\x48\x0c\x63\x45\x0e\xc5\xad\x96\xbb\xec\xee\xd0\xb2\x53\x17\xe8\xd3\xf4\xc1\xfa\x24\xc5\x2e\x29\x8a\xb2\x94\x28\x3f\x17\x81\xc2\x9d\x3d\x67\xfe\xce\xcc\x46\x94\xf2\x1d\x59\x27\x8d\x1e\x83\x69\xc9\x46\xc7\x29\xdd\x8d\xee\x5e\x0a\x55\xe6\xe2\x65\x6f\x29\x75\x3a\xc6\x4c\xb8\x65\xaf\x20\x16\xa9\x60\x31\xee\x01\x5a\x14\x34\xc6\x42\xf2\x30\x51\x46\x53\xcf\x95\x94\xf8\xef\x2b\x63\x97\xae\x14\x09\x39\xff\xaf\x61\x63\x67\x2a\x2e\x2b\xee\x01\x40\x4a\x2e\xb1\xb2\xe4\x40\x38\xcb\xc9\x63\xc0\x52\x69\xb0\x92\x4a\x61\x4e\x08\x80\x29\x8c\x66\x03\xce\x09\x77\x46\x55\x05\x61\x2e\x92\xa5\xd4\x0b\x70\x2e\xdd\x86\xa6\x07\x94\xc2\x8a\x62\x8b\xee\x6a\x72\x79\x71\x7b\x7d\x75\xbe\x4b\xe8\xc9\x2a\xab\xc0\xa6\xa6\x09\x16\xfc\x50\xd2\x18\x8e\xad\xd4\x8b\x0e\xca\xe5\xf5\xf9\xf9\x6d\xf8\x6b\x7a\xf6\x7a\x3f\x94\xa5\x3b\xe9\x53\x17\xf0\x72\x4a\x96\xa6\x62\xf4\xe7\x56\xe8\x24\x7f\x01\x16\x8b\x17\x70\xb9\x78\x01\x4b\xd9\xbf\x7f\xff\x33\xd8\xc7\xe6\x51\x33\x51\x29\x1e\xa3\x10\x8e\xc9\x76\x3c\x70\xd5\xbc\x30\x69\xa5\xc8\xed\xd2\xa7\x94\x49\x4d\x0e\x32\x0b\x49\xb2\xe4\x4c\x65\x13\x82\xcb\x4d\xa5\x52\x48\x2d\x59\x0a\x25\x3f\x11\x84\x4e\x91\x11\x27\x79\x30\x7c\x82\xf9\x79\x6f\x22\xb6\x15\x45\x1d\x6f\x52\x2a\x39\xdf\x75\xa4\x24\x9b\x19\x5b\x38\x08\x1f\xab\x52\x66\x55\xa7\x16\xab\x9c\x2c\xc1\x68\xf5\x10\x88\x0b\xe3\x7c\xc2\x12\xd2\x8c\xc4\x14\x85\xe4\xbe\x1b\x84\xa2\x07\x4c\xf8\xd2\x07\x37\x29\x3d\xe4\xd9\xcb\xae\x5b\xce\xa9\x77\x64\x65\xf6\xf0\xc5\x1c\xe5\xcc\x65\xdc\x9a\xae\xb3\x34\x27\x38\x62\x5f\x3d\x1f\x2c\x8c\x45\x26\x94\x23\x48\xed\x5d\x6e\xfc\x5a\x28\x33\x17\x2a\xd4\x3b\x31\x3a\x93\x8b\x6f\x4d\x9c\xab\xe6\xa9\xb4\x94\xb0\xb1\x7b\x9c\xec\x9e\x42\x6a\x27\x53\x0a\xf9\x8a\x6a\xd1\x44\x9b\x5e\x6f\xbb\x36\x9c\xaf\x75\xd3\x78\x29\x35\x9b\x43\x8e\x39\x9b\x6c\x17\x54\x11\xd3\xe4\x5e\x3a\xde\xd8\x76\x3c\x4b\x14\x09\xed\xb5\x1b\xf8\x12\xa3\x99\x34\x3b\x98\x75\xc7\x95\xe6\xd8\xf9\x1b\x2c\xb5\xf0\x37\xb0\x89\xa3\x2f\xb3\xb5\x5f\x0c\xa1\x2c\x89\xf4\x01\xe4\x99\xdc\x00\x73\xca\x8c\x25\xb0\x7d\x08\x7a\xee\x06\x15\x06\x01\xfb\xc6\x39\x14\x4b\x28\x54\x37\x1a\x5f\xe1\x4b\x6b\xee\xf7\xa4\xd8\xa7\xea\x6c\x36\xbb\x44\xe9\xcf\xe1\xc8\xde\x91\x45\x66\x2c\xb4\xd1\xc3\xe9\xf4\x1c\x96\xfe\xa8\xc8\xf1\x41\x4d\x3c\x25\x74\x07\x18\xa7\xbb\x94\xdf\x4b\xa7\xcd\x17\xb8\xb4\x69\x88\x86\x30\x25\x87\x9a\x99\xac\xfe\xe4\x73\xec\x5d\x19\xd5\xfe\x7c\x0b\xb5\x25\x57\x29\xde\x1a\xad\xb5\x72\xf7\x4f\xf2\xd2\x52\x22\x1d\x35\x36\x98\x9e\xbd\x06\xe7\x82\xb1\x12\x6e\xad\x6c\xcc\x1f\xea\x01\x1e\x16\x0a\xe0\x98\xca\x6d\xfc\x76\x2a\xcb\x42\x2c\xfc\x8a\x49\x6c\x2c\xcd\xa8\x5e\x4d\x43\x4b\x8a\x84\x23\x37\x5a\x48\xce\xab\x79\x9c\x98\xa2\x39\x4a\xd2\x51\x29\x4b\x52\x52\xd3\x28\x29\x52\x6f\x30\xf4\x43\x70\xac\x04\x93\xab\x3d\xae\xfd\x1d\xe3\xb1\xe9\xcd\x37\x67\x93\x37\x3f\x5d\x5c\xcf\x6e\x4f\xdf\x5e\x9d\x44\xa3\x56\x68\xa3\xa3\xbe\xd4\x65\xc5\x2e\xae\x77\x4b\xdc\xd5\xe8\x20\x6a\x6e\x3f\xdf\xbe\x7e\xd4\xdf\xec\xbe\xb8\xd6\x6e\x5c\x0a\xce\x07\x87\xd0\x1a\xb8\xa0\xb6\x54\xda\xfe\x00\x7f\x36\x9f\x80\xe7\x38\x0d\x22\x85\xd0\x8d\x7e\x7c\x39\xf7\x49\xb1\x23\x3e\x99\x79\xd5\xd5\x6a\x8b\x37\x50\x1d\xd0\x5f\x09\xa9\xd1\xc7\x8c\xdf\x2b\xc7\x88\x6c\x81\xa1\xcd\x70\xd4\x0d\x28\xc2\x9c\x12\x51\x39\xda\xfe\x8c\x42\x2e\x72\xf6\x83\x33\x1a\x45\x1d\x44\x63\x6b\x4f\x8c\x09\x9d\x27\x50\x98\x4a\x33\xa5\xcd\xf6\xde\xb8\x21\x33\x7c\xf8\x80\x61\x8a\xe8\x09\xdd\xcd\x0d\x7e\xf0\x18\xba\x35\xed\x84\xef\x75\x9a\xcb\x34\x25\x8d\x4c\x2a\x72\x61\xa9\xad\x23\x96\xcd\x2e\xab\xff\x34\xc1\x3c\x81\x1f\xfd\x6f\x1f\xec\x5e\x2c\x38\x16\x36\xe4\x79\x25\x39\x47\x8c\x79\xe5\xb3\x99\xa8\x2a\xf5\x1f\xe3\xf8\x30\x59\xfc\xe1\x59\x7c\xf3\xfd\x8c\x31\x4a\x55\xb9\x50\x73\xe3\x87\x21\x92\x5c\x58\x91\xd4\x8f\x83\x43\xd4\xf1\xab\x0d\x71\x26\x9b\x9f\x7f\xad\xbb\xac\x4e\x7f\xf4\xb4\x25\xb7\x57\xc1\x20\xc2\xc9\x49\xb3\xc7\xf6\xd4\x65\xdd\xa9\xbd\x96\xa4\xf9\xe5\x85\x86\xe1\xa7\x5d\xf8\x76\x36\x0f\x22\x3c\x3e\x82\xee\x4b\x63\xeb\x01\x79\x7b\x79\x75\xf1\xfe\xb7\x93\xcf\x5f\xf8\x1a\x68\xb7\x1f\x7b\xfa\x05\x70\xf7\x75\xe8\xcd\xd0\xdd\x82\xfe\xe5\xe2\x33\xb8\x6b\xe3\x75\x3a\x46\x4b\x33\x14\x65\xd9\x8e\x22\x7c\x6c\x53\x38\xf4\x8f\xd0\x1d\xb6\xf5\xbb\x75\x10\x75\x4d\xdb\x47\xe6\x8e\xfd\xd6\x0b\x75\xfb\x92\x9f\x3c\x3b\x12\xeb\x9c\x6f\x9e\x42\x3b\xa8\xed\xd1\x36\xe2\xe6\xe1\xb8\xe7\x4a\x7b\xb6\x7d\x27\x3c\x19\xf7\x75\x5b\xc9\x79\x3b\x47\x93\x9d\x59\xd0\x1c\x5c\x4d\xa6\xd7\xe7\x33\x1f\x9b\x9f\xae\xcd\x7b\x7b\x58\x0a\xeb\x08\x67\x93\xd7\xa7\x78\x04\x5b\x3f\x4a\x8e\x3f\xea\xe3\x16\x6e\xf2\xfe\xed\xec\xf6\xcd\xc5\xe9\xe4\x24\x3a\x7a\x15\x75\xda\x1e\xd1\x51\x7b\x16\xe1\xd9\x09\xfe\x8f\x9b\x75\xfd\xbb\xed\x4d\xf7\x92\xb1\x31\xed\x3d\x91\xd2\x73\xfc\x2c\x96\xfe\x1d\x6d\x09\xab\xf5\x1c\x15\x69\x0a\x01\xb6\x42\x2a\x2f\x63\x4d\x2b\xbf\x88\xd0\xfc\x0f\xa6\xde\xa3\xcf\x1a\x00\x4a\x72\x83\xa1\x2f\xe7\x26\xc0\x08\x3f\xe2\xa8\xdf\xec\xdb\xb8\x5e\xa0\xf5\xf6\xe8\xfd\x17\x00\x00\xff\xff\xab\x51\x20\xa4\xa1\x0d\x00\x00")

func resourcesGitGitCloneYamlBytes() ([]byte, error) {
	return bindataRead(
		_resourcesGitGitCloneYaml,
		"resources/git/git-clone.yaml",
	)
}

func resourcesGitGitCloneYaml() (*asset, error) {
	bytes, err := resourcesGitGitCloneYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/git/git-clone.yaml", size: 3489, mode: os.FileMode(420), modTime: time.Unix(1601047620, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"resources/git/git-clone.yaml": resourcesGitGitCloneYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"resources": &bintree{nil, map[string]*bintree{
		"git": &bintree{nil, map[string]*bintree{
			"git-clone.yaml": &bintree{resourcesGitGitCloneYaml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
