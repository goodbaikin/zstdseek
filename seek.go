package zstdseek

/*
	#cgo CFLAGS: -I libzstd-seek -D _ZSTD_SEEK_DEBUG_=1
	#cgo LDFLAGS: -L build -l zstd-seek -l zstd
	#include <stdint.h>
	#include <zstd-seek.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

var (
	errCreateFailed = "failed to create new seeker"
	errSeekFailed   = "failed to seek"
)

type Seeker interface {
	Read([]byte) uint
	Seek(int, int) error
	Size() uint

	// for advanced users
	Table() Table
}

func CreateFromFile(file string, withoutTable bool) (Seeker, error) {
	cfile := C.CString(file)

	var cctx *C.ZSTDSeek_Context
	if withoutTable {
		cctx = C.ZSTDSeek_createFromFileWithoutJumpTable(cfile)
	} else {
		cctx = C.ZSTDSeek_createFromFile(cfile)
	}

	if cctx == nil {
		return nil, fmt.Errorf("%s: %s", errCreateFailed, file)
	}
	return &ctx{cctx: cctx}, nil
}

func CreateFromBytes(buff []byte, withoutTable bool) (Seeker, error) {
	cbuff := unsafe.Pointer(&buff[0])
	csize := C.size_t(len(buff))

	var cctx *C.ZSTDSeek_Context
	if withoutTable {
		cctx = C.ZSTDSeek_createWithoutJumpTable(cbuff, csize)
	} else {
		cctx = C.ZSTDSeek_create(cbuff, csize)
	}

	if cctx == nil {
		return nil, fmt.Errorf("%s", errCreateFailed)
	}
	return &ctx{cctx: cctx}, nil
}

func (c *ctx) Read(outBuff []byte) uint {
	cOutBuff := unsafe.Pointer(&outBuff[0])
	cOutBuffSize := C.size_t(len(outBuff))
	ret := C.ZSTDSeek_read(cOutBuff, cOutBuffSize, c.cctx)
	return uint(ret)
}

func (c *ctx) Seek(offset, origin int) error {
	cOffset := C.long(offset)
	cOrigin := C.int(origin)

	ret := C.ZSTDSeek_seek(c.cctx, cOffset, cOrigin)
	if ret != 0 {
		return fmt.Errorf("%s", errSeekFailed)
	}
	return nil
}

func (c *ctx) Size() uint {
	ret := C.ZSTDSeek_uncompressedFileSize(c.cctx)
	return uint(ret)
}

func (c *ctx) Table() Table {
	return c
}
