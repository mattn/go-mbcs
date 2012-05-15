package mbcs

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	LC_ALL      = 0
	LC_COLLATE  = 1
	LC_CTYPE    = 2
	LC_MONETARY = 3
	LC_NUMERIC  = 5
	LC_TIME     = 6
)

var msvcrt syscall.Handle
var setlocale uintptr
var wcstombs uintptr
var mbstowcs uintptr

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

func init() {
	var err error
	if msvcrt, err = syscall.LoadLibrary("msvcrt.dll"); err != nil {
		abort("LoadLibrary", err)
	}
	if setlocale, err = syscall.GetProcAddress(msvcrt, "setlocale"); err != nil {
		abort("GetProcAddress", err)
	}
	if wcstombs, err = syscall.GetProcAddress(msvcrt, "wcstombs"); err != nil {
		abort("GetProcAddress", err)
	}
	if mbstowcs, err = syscall.GetProcAddress(msvcrt, "mbstowcs"); err != nil {
		abort("GetProcAddress", err)
	}
}

func SetLocale(cat int, loc string) string {
	bs := []byte("")
	r, _, _ := syscall.Syscall(
		setlocale,
		2,
		0,
		uintptr(unsafe.Pointer(&bs[0])),
		0)

	bytes := ((*[1<<24]byte)(unsafe.Pointer(r)))[:]
	return string(bytes[0:clen(bytes)])
}

func WcsToMbs(wcs string) ([]byte, uint) {
	wcs32 := ([]rune)(wcs)
	var wcs16 []uint16 = make([]uint16, len(wcs32)+1)
	for n, c := range wcs32 {
		wcs16[n] = uint16(c)
	}
	r, _, _ := syscall.Syscall(
		wcstombs,
		0,
		uintptr(unsafe.Pointer(&wcs16[0])),
		0,
		0)
	l := uint16(r)
	var mbs []byte = make([]byte, l)
	r, _, _ = syscall.Syscall(
		uintptr(wcstombs),
		uintptr(unsafe.Pointer(&mbs[0])),
		uintptr(unsafe.Pointer(&wcs16[0])),
		uintptr(l),
		0)
	return mbs, uint(r)
}

func MbsToWcs(mbs []byte) (string, uint) {
	r, _, _ := syscall.Syscall(
		mbstowcs,
		0,
		uintptr(unsafe.Pointer(&mbs[0])),
		0,
		0)
	l := uint16(r)
	var wcs16 []uint16 = make([]uint16, l)
	r, _, _ = syscall.Syscall(
		mbstowcs,
		uintptr(unsafe.Pointer(&wcs16[0])),
		uintptr(unsafe.Pointer(&mbs[0])),
		uintptr(l),
		0)
	var wcs32 []rune = make([]rune, len(wcs16))
	for n, c := range wcs16 {
		wcs32[n] = rune(c)
	}
	return string(wcs32), uint(r)
}
