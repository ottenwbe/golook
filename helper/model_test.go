package helper

import (
	"encoding/json"
	"io"
	"log"
	"testing"
	"time"
)

const (
	testFileName = "test.txt"
	testSysName  = "test"
)

func TestDecodeFile(t *testing.T) {
	c := make(chan bool)
	pipeReader, pipeWriter := io.Pipe()
	go PipeWriteFile(pipeWriter)
	go PipeReadFile(pipeReader, c)
	result := <-c

	if !result {
		t.Error("Decoding failed")
	}
}

func TestDecodeSystem(t *testing.T) {
	c := make(chan bool)
	pipeReader, pipeWriter := io.Pipe()
	go PipeWriteSystem(pipeWriter)
	go PipeReadSystem(pipeReader, c)
	result := <-c

	if !result {
		t.Error("Decoding failed")
	}
}

func TestDecodeFiles(t *testing.T) {
	c := make(chan bool)
	pipeReader, pipeWriter := io.Pipe()
	go PipeWriteFiles(pipeWriter)
	go PipeReadFiles(pipeReader, c)
	result := <-c

	if !result {
		t.Error("Decoding failed")
	}
}

func newTestSystem(sysName string) *System {
	s := &System{
		Name:  sysName,
		OS:    "Linux",
		IP:    "localhost",
		UUID:  "uuid",
		Files: nil,
	}
	return s
}

func newTestFile(fileName string) File {
	f := File{}
	f.Name = fileName
	f.Accessed = time.Now()
	f.Created = time.Now()
	f.Modified = time.Now()
	return f
}

func PipeWriteSystem(pipeWriter *io.PipeWriter) {
	b, _ := json.Marshal(newTestSystem(testSysName))
	pipeWriter.Write(b)
}

func PipeReadSystem(pipeReader *io.PipeReader, c chan bool) {
	if s, err := DecodeSystem(pipeReader); err == nil && s.Name == testSysName {
		c <- true
	} else {
		log.Printf("Error expected nil, got %s", err)
		log.Printf("File, expected %s got %s", testSysName, s.Name)
	}
	c <- false
}

func PipeWriteFiles(pipeWriter *io.PipeWriter) {
	b, _ := json.Marshal([]File{newTestFile(testFileName), newTestFile(testFileName)})
	pipeWriter.Write(b)
}

func PipeReadFiles(pipeReader *io.PipeReader, c chan bool) {
	if f, err := DecodeFiles(pipeReader); err == nil && f[0].Name == testFileName {
		c <- true
	} else {
		log.Printf("Error expected nil, got %s", err)
		log.Printf("File, expected %s got %s", testFileName, f[0].Name)
	}
	c <- false
}

func PipeWriteFile(pipeWriter *io.PipeWriter) {
	b, _ := json.Marshal(newTestFile(testFileName))
	pipeWriter.Write(b)
}

func PipeReadFile(pipeReader *io.PipeReader, c chan bool) {
	if f, err := DecodeFile(pipeReader); err == nil && f.Name == testFileName {
		c <- true
	} else {
		log.Printf("Error expected nil, got %s", err)
		log.Printf("File, expected %s got %s", testFileName, f.Name)
	}
	c <- false
}
