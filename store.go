package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type PathProps struct {
	PathName string
	Filename string
}

func (p PathProps) FirstPath() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathProps) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

func TransformFunc(key string) PathProps {
	hash := sha1.Sum([]byte(key))
	hs := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hs) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hs[from:to]
	}

	return PathProps{
		PathName: strings.Join(paths, "/"),
		Filename: hs,
	}
}

type Store struct {
	Root          string
	TransformFunc func(string) PathProps
}

func NewStore() *Store {
	return &Store{
		Root:          "netfiles",
		TransformFunc: TransformFunc,
	}
}

func (s *Store) Has(key string) bool {
	pathProps := s.TransformFunc(key)
	_, err := os.Stat(fmt.Sprintf("%s/%s", s.Root, pathProps.FullPath()))
	return errors.Is(err, os.ErrNotExist)
}

func (s *Store) Delete(key string) error {
	pathProps := s.TransformFunc(key)
	defer func() {
		log.Println("deleted")
	}()
	// filepath with root
	return os.RemoveAll(fmt.Sprintf("%s/%s", s.Root, pathProps.FirstPath()))
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathProps := s.TransformFunc(key)
	return os.Open(fmt.Sprintf("%s/%s", s.Root, pathProps.FullPath()))
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathProps := s.TransformFunc(key)
	root := fmt.Sprintf("%s/%s", s.Root, pathProps.PathName)
	if err := os.MkdirAll(root, os.ModePerm); err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s", s.Root, pathProps.FullPath())
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	wr, err := io.Copy(file, r)
	if err != nil {
		return err
	}

	log.Printf("written %d bytes to disk: %s", wr, filePath)
	return nil
}
