package btrfs 

/*
#include <stdlib.h>
#include <dirent.h>
#include <btrfs/ioctl.h>
*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"
)

func Free(p *C.char) {
	C.free(unsafe.Pointer(p))
}

func OpenDir(dirPath string) (*C.DIR, error) {
	path := C.CString(dirPath)
	defer Free(path)

	dir := C.opendir(path)
	if dir == nil {
		return nil, fmt.Errorf("Error opening directory")
	}
	return dir, nil
}

func CloseDir(dir *C.DIR) {
	if dir != nil {
		C.closedir(dir)
	}
}

func GetDirFd(dir *C.DIR) uintptr {
	return uintptr(C.dirfd(dir))
}

func SubvolumeCreate(dst string, name string) error {
	dstdir, err := OpenDir(dst)
	if err != nil {
		return err
	}
	defer CloseDir(dstdir)
	
	var args C.struct_btrfs_ioctl_vol_args
	for i, c := range []byte(name) {
		args.name[i] = C.char(c)
	}
	
	_, _, errorno := syscall.Syscall(syscall.SYS_IOCTL, GetDirFd(dstdir), C.BTRFS_IOC_SUBVOL_CREATE, uintptr(unsafe.Pointer(&args)))

	if errorno != 0 {
		return fmt.Errorf("Error creating subvolume: %v", errorno.Error())
	}

	return nil
	
}

func SubvolumeSnapshot(src string, dst string, name string, readonly bool) error {
	srcDir, err := OpenDir(src)
	if err != nil {
		return err
	}
	defer CloseDir(srcDir)

	dstDir, err2 := OpenDir(dst)
	if err2 != nil {
		return err2
	}
	defer CloseDir(dstDir)

	var args C.struct_btrfs_ioctl_vol_args_v2
	for i, c := range []byte(name) {
		args.name[i] = C.char(c)
	}
	args.fd = C.__s64(GetDirFd(srcDir))
	
	if(readonly) {
		args.flags |= C.BTRFS_SUBVOL_RDONLY
	}
	
	_, _, errorno := syscall.Syscall(syscall.SYS_IOCTL, GetDirFd(dstDir), C.BTRFS_IOC_SNAP_CREATE_V2, uintptr(unsafe.Pointer(&args)))

	if errorno != 0 {
		fmt.Println(errorno.Error())
		return fmt.Errorf("Error creating snapshot: %v", errorno.Error())
	}
	
	Sync(dstDir)
	
	return nil
}

func Sync(dir *C.DIR) error {
	_, _, errorno := syscall.Syscall(syscall.SYS_IOCTL, GetDirFd(dir), C.BTRFS_IOC_START_SYNC, uintptr(unsafe.Pointer(nil)))
	if errorno != 0 {
		return fmt.Errorf("Error sync: %v", errorno.Error())
	} 
	
	_,_, errorno2 := syscall.Syscall(syscall.SYS_IOCTL, GetDirFd(dir), C.BTRFS_IOC_WAIT_SYNC, uintptr(unsafe.Pointer(nil)))
	if errorno2 != 0 {
		return fmt.Errorf("Error wait sync %v", errorno.Error())
	}
	return nil;
}

func SubvolumeDelete(dstdir string, name string) error {
	dir, err := OpenDir(dstdir)
	if err != nil {
		return err
	}	
	defer CloseDir(dir)

	var args C.struct_btrfs_ioctl_vol_args
	for i, c := range []byte(name) {
		args.name[i] = C.char(c)
	}
	
	_,_, errorno := syscall.Syscall(syscall.SYS_IOCTL, GetDirFd(dir), C.BTRFS_IOC_SNAP_DESTROY, uintptr(unsafe.Pointer(&args)))
	
	if errorno != 0 {
		return fmt.Errorf("Error deleting subvolume: %v", errorno.Error())
	}

	return nil
}

