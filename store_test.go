package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestWriteStream(t *testing.T) {
	s := NewStore()

	key := "myfolders"
	data := []byte("statement of claims")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.readStream(key)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	b, _ := io.ReadAll(r)
	if string(b) != string(data) {
		t.Error(b)
	}

	fmt.Println(string(b))
}

func TestCreateTransformFunc(t *testing.T) {
	key := "mylawfirmfiles"

	res := TransformFunc(key)
	fmt.Println(res)
}

func TestDelete(t *testing.T) {
	s := NewStore()

	key := "myfolders"
	data := []byte("my law journals")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(t)
	}
}
