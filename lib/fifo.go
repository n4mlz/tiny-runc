package lib

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
)

func SetupPipes(pipePaths ...string) error {
	for _, path := range pipePaths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("error removing existing pipe: %v", err)
		}
		if err := syscall.Mkfifo(path, 0666); err != nil {
			return fmt.Errorf("error creating pipe: %v", err)
		}
	}
	return nil
}

func CleanupPipes(pipePaths ...string) {
	for _, path := range pipePaths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: failed to remove pipe %s: %v\n", path, err)
		}
	}
}

type PipeHandler struct {
	ReadPipe  *os.File
	WritePipe *os.File
}

func OpenPipesA(readPath, writePath string) (*PipeHandler, error) {
	writePipe, err := os.OpenFile(writePath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("error opening write pipe: %v", err)
	}

	readPipe, err := os.OpenFile(readPath, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("error opening read pipe: %v", err)
	}

	return &PipeHandler{
		ReadPipe:  readPipe,
		WritePipe: writePipe,
	}, nil
}

func OpenPipesB(readPath, writePath string) (*PipeHandler, error) {
	readPipe, err := os.OpenFile(readPath, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("error opening read pipe: %v", err)
	}

	writePipe, err := os.OpenFile(writePath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("error opening write pipe: %v", err)
	}

	return &PipeHandler{
		ReadPipe:  readPipe,
		WritePipe: writePipe,
	}, nil
}

func (ph *PipeHandler) Close() {
	ph.ReadPipe.Close()
	ph.WritePipe.Close()
}

func (ph *PipeHandler) SendMessage(message string) error {
	_, err := ph.WritePipe.Write([]byte(message + "\n"))
	return err
}

func (ph *PipeHandler) ReceiveMessage() (string, error) {
	reader := bufio.NewReader(ph.ReadPipe)
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return message[:len(message)-1], nil
}
