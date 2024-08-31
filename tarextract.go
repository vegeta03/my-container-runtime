package main

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func extractTar(tarFile, destDir string) error {
    file, err := os.Open(tarFile)
    if err != nil {
        return err
    }
    defer file.Close()

    tr := tar.NewReader(file)

    for {
        header, err := tr.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        target := filepath.Join(destDir, header.Name)

        switch header.Typeflag {
        case tar.TypeDir:
            if err := os.MkdirAll(target, 0755); err != nil {
                return err
            }
        case tar.TypeReg:
            outFile, err := os.Create(target)
            if err != nil {
                return err
            }
            if _, err := io.Copy(outFile, tr); err != nil {
                outFile.Close()
                return err
            }
            outFile.Close()
        case tar.TypeSymlink:
            if err := os.Symlink(header.Linkname, target); err != nil {
                return err
            }
        }
    }
    return nil
}