package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/goodbaikin/zstdseek"
)

const file = "example.tar.zst"

func main() {
	seeker, err := zstdseek.CreateFromFile(file, false)
	if err != nil {
		printError(err)
	}

	buf := make([]byte, seeker.Size())
	fmt.Fprintf(os.Stdout, "uncompressed file size: %d\n", seeker.Size())

	read := seeker.Read(buf)
	fmt.Fprintf(os.Stdout, "read %d bytes\n\n", read)

	files, err := listFileInTar(seeker)
	if err != nil {
		printError(err)
	}

	for _, file := range files {
		fmt.Fprintf(os.Stdout, "%s\n", file)
	}
}

func stringFromBytes(bytes []byte) string {
	i := 0
	for ; i < len(bytes); i++ {
		if bytes[i] == 0 {
			break
		}
	}
	if i == len(bytes) {
		i--
	}

	return string(bytes[:i])
}

func listFileInTar(seeker zstdseek.Seeker) ([]string, error) {
	const headerSize = 512
	files := []string{}
	offset := 0

	for {
		if err := seeker.Seek(offset, os.SEEK_SET); err != nil {
			return nil, err
		}

		header := [headerSize]byte{}
		seeker.Read(header[:])
		if header[0] == 0 {
			break
		}
		files = append(files, stringFromBytes(header[:99]))

		sizeStr := string(header[124:135])
		fileSize, err := strconv.ParseInt(sizeStr, 8, 64)
		if err != nil {
			return nil, err
		}
		blockSize := (math.Ceil(float64(fileSize)/512.0) + 1) * 512
		offset += int(blockSize)
	}

	return files, nil
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
