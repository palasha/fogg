package plugins

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// TypePluginFormat is the plugin format such as binary, zip, tar
type TypePluginFormat string

const (
	// TypePluginFormatTar is a tar archived plugin
	TypePluginFormatTar TypePluginFormat = "tar"
	// TypePluginFormatBin is a binary plugin
	TypePluginFormatBin TypePluginFormat = "bin"
)

// CustomPlugin is a custom plugin
type CustomPlugin struct {
	URL       string           `json:"url" validate:"required"`
	Format    TypePluginFormat `json:"format" validate:"required"`
	targetDir string
}

// Install installs the custom plugin
func (cp *CustomPlugin) Install(fs afero.Fs, pluginName string) error {
	if cp == nil {
		return errors.New("nil CustomPlugin")
	}
	if fs == nil {
		return errors.New("nil fs")
	}

	tmpPath, err := cp.fetch(pluginName)
	defer os.Remove(tmpPath) // clean up
	if err != nil {
		return err
	}
	return cp.process(fs, pluginName, tmpPath)
}

// SetTargetPath sets the target path for this plugin
func (cp *CustomPlugin) SetTargetPath(path string) {
	cp.targetDir = path
}

// fetch fetches the custom plugin at URL
func (cp *CustomPlugin) fetch(pluginName string) (string, error) {
	tmpFile, err := ioutil.TempFile("", fmt.Sprintf("%s-*.tmp", pluginName))
	if err != nil {
		return "", errors.Wrap(err, "could not create temporary directory")
	}
	resp, err := http.Get(cp.URL)
	if err != nil {
		return "", errors.Wrapf(err, "could not get %s", cp.URL)
	}
	defer resp.Body.Close()
	_, err = io.Copy(tmpFile, resp.Body)
	return tmpFile.Name(), errors.Wrap(err, "could not download file")
}

// process the custom plugin
func (cp *CustomPlugin) process(fs afero.Fs, pluginName string, path string) error {
	switch cp.Format {
	case TypePluginFormatTar:
		return cp.processTar(fs, path)
	case TypePluginFormatBin:
		return cp.processBin(fs, pluginName, path)
	default:
		return errors.Errorf("Unknown plugin format %s", cp.Format)
	}
}

func (cp *CustomPlugin) processBin(fs afero.Fs, name string, downloadPath string) error {
	err := fs.MkdirAll(cp.targetDir, 0755)
	if err != nil {
		return errors.Wrapf(err, "Could not create directory %s", cp.targetDir)
	}
	targetPath := path.Join(cp.targetDir, name)
	src, err := os.Open(downloadPath)
	if err != nil {
		return errors.Wrapf(err, "Could not open downloaded file at %s", downloadPath)
	}
	dst, err := fs.Create(targetPath)
	if err != nil {
		return errors.Wrapf(err, "Could not open target file at %s", targetPath)
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		return errors.Wrapf(err, "Could not move %s to %s", downloadPath, targetPath)
	}
	err = fs.Chmod(targetPath, os.FileMode(0755))
	return errors.Wrapf(err, "Error making %s executable", targetPath)
}

// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func (cp *CustomPlugin) processTar(fs afero.Fs, path string) error {
	err := fs.MkdirAll(cp.targetDir, 0755)
	if err != nil {
		return errors.Wrapf(err, "Could not create directory %s", cp.targetDir)
	}
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "could not read staged custom plugin")
	}
	defer f.Close()
	gzr, err := gzip.NewReader(f)
	if err != nil {
		return errors.Wrap(err, "could not create gzip reader")
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // no more files are found
		}
		if err != nil {
			return errors.Wrap(err, "Error reading tar")
		}
		if header == nil {
			return errors.New("Nil tar file header")
		}
		// the target location where the dir/file should be created
		target := filepath.Join(cp.targetDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir: // if its a dir and it doesn't exist create it
			err := fs.MkdirAll(target, 0755)
			if err != nil {
				return errors.Wrapf(err, "tar: could not create directory %s", target)
			}
		case tar.TypeReg: // if it is a file create it, preserving the file mode
			destFile, err := fs.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return errors.Wrapf(err, "tar: could not open destination file for %s", target)
			}
			_, err = io.Copy(destFile, tr)
			if err != nil {
				destFile.Close()
				return errors.Wrapf(err, "tar: could not copy file contents")
			}
			// Manually take care of closing file since defer will pile them up
			destFile.Close()
		default:
			log.Warnf("tar: unrecognized tar.Type %d", header.Typeflag)
		}
	}
	return nil
}