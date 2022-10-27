package zstdseek

/*
	#cgo CFLAGS: -I libzstd-seek
	#cgo LDFLAGS: -L build -l zstd-seek -l zstd
	#include <stdint.h>
	#include <zstd-seek.h>
*/
import "C"
import "fmt"

type Table interface {
	Add(uint, uint) error
	Initialize() error
	InitializeUpUntilPos(uint) error
	IsInitialized() bool
}

type ctx struct {
	cctx *C.ZSTDSeek_Context
}

var (
	errGetFailed        = "failed to get table"
	errInitializeFailed = "failed to initialize"
)

func (c *ctx) Add(compressedPos uint, uncompressedPos uint) error {
	jt := C.ZSTDSeek_getJumpTableOfContext(c.cctx)
	if jt == nil {
		return fmt.Errorf(errGetFailed)
	}

	ccompressedPos := C.size_t(compressedPos)
	cuncompressedPos := C.size_t(uncompressedPos)
	C.ZSTDSeek_addJumpTableRecord(jt, ccompressedPos, cuncompressedPos)
	return nil
}

func (c *ctx) Initialize() error {
	ret := C.ZSTDSeek_initializeJumpTable(c.cctx)
	if ret != 0 {
		return fmt.Errorf(errInitializeFailed)
	}
	return nil
}

func (c *ctx) InitializeUpUntilPos(upUntilPos uint) error {
	size_t := C.size_t(upUntilPos)
	ret := C.ZSTDSeek_initializeJumpTableUpUntilPos(c.cctx, size_t)

	if ret != 0 {
		return fmt.Errorf(errInitializeFailed)
	}
	return nil
}

func (c *ctx) IsInitialized() bool {
	ret := C.ZSTDSeek_jumpTableIsInitialized(c.cctx)
	return ret != 0
}
