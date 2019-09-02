package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
)

type PacFile struct {
	Name   [0x20]byte
	Size   uint32
	Offset uint32
}

func (p PacFile) String() string {
	return fmt.Sprintf("%s\tStarting: 0x%.8X\tSize: 0x%.8X\n", p.Name[:bytes.IndexByte(p.Name[:], byte(0))], p.Offset, p.Size)
}

func (p PacFile) Extract(file *os.File) {
	_, err := file.Seek(int64(p.Offset), 0)

	if err != nil {
		log.Fatal("Error while seeking into the file", err)
	}

	data := readNextByte(file, int(p.Size))
	path := file.Name()[:strings.Index(file.Name(), ".pac")] + "/"
	filename := string(p.Name[:bytes.IndexByte(p.Name[:], byte(0))])

	os.Mkdir(path, os.ModeDir)

	output, err := os.OpenFile(path+filename, os.O_CREATE, 0755)

	if err != nil {
		log.Fatal("Error while opening or creating the outputfile", err)
	}

	defer output.Close()

	output.Write(data)
}

func main() {
	path := os.Args[1]
	f, err := os.Open(path)
	files := make([]PacFile, 0)

	if err != nil {
		log.Fatal("Error while opening the file", err)
	}
	defer f.Close()

	sig := readNextByte(f, 4)

	if string(sig) != "PAC\x20" {
		log.Fatal("Unsupported file provided")
	}

	f.Seek(0x08, 0)

	nbFile := binary.LittleEndian.Uint32(readNextByte(f, 4))
	f.Seek(0x804, 0)

	i := uint32(0)
	for i < nbFile {
		pac := PacFile{}
		data := readNextByte(f, 0x28)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.LittleEndian, &pac)

		if err != nil {
			log.Fatal("binary.read failed", err)
		}

		files = append(files, pac)
		fmt.Print(pac)
		i++
	}

	for i := uint32(0); i < nbFile; i++ {
		files[i].Extract(f)
	}

}

func readNextByte(file *os.File, number int) []byte {

	bytes := make([]byte, number)

	_, err := file.Read(bytes)

	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
