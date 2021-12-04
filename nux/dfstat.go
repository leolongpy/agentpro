package nux

import (
	"agentpro/utils"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"syscall"
)

const nuxRootFs = "NUX_ROOTFS"

func init() {
	fmt.Printf("use %s = %s as file system root", nuxRootFs, Root())
}

func Root() string {
	root := os.Getenv(nuxRootFs)
	if !strings.HasPrefix(root, string(os.PathSeparator)) {
		return ""
	}
	root = strings.TrimSuffix(root, string(os.PathSeparator))
	if pathExists(root) {
		return root
	}
	return ""
}

func pathExists(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		return fi.IsDir()
	}
	return false
}

var FSFILE_PREFIX_IGNORE = []string{
	"/sys",
	"/net",
	"/misc",
	"/proc",
	"/lib",
}

func IgnoreFsFile(fs_file string) bool {
	for _, prefix := range FSFILE_PREFIX_IGNORE {
		if strings.HasPrefix(fs_file, prefix) {
			return true
		}
	}

	return false
}

var FSSPEC_IGNORE = map[string]struct{}{
	"none":      struct{}{},
	"nodev":     struct{}{},
	"proc":      struct{}{},
	"hugetlbfs": struct{}{},
	"mqueue":    struct{}{},
}

var FSTYPE_IGNORE = map[string]struct{}{
	"cgroup":     struct{}{},
	"debugfs":    struct{}{},
	"devpts":     struct{}{},
	"devtmpfs":   struct{}{},
	"iso9660":    struct{}{},
	"rpc_pipefs": struct{}{},
	"rootfs":     struct{}{},
	"overlay":    struct{}{},
	"tmpfs":      struct{}{},
	"squashfs":   struct{}{},
}

type DeviceUsage struct {
	FsSpec            string
	FsFile            string
	FsVfstype         string
	BlocksAll         uint64
	BlocksUsed        uint64
	BlocksFree        uint64
	BlocksUsedPercent float64
	BlocksFreePercent float64
	InodesAll         uint64
	InodesUsed        uint64
	InodesFree        uint64
	InodesUsedPercent float64
	InodesFreePercent float64
}

func (this *DeviceUsage) String() string {
	return fmt.Sprintf("<FsSpec:%s, FsFile:%s, FsVfstype:%s, BPFree:%f...>",
		this.FsSpec,
		this.FsFile,
		this.FsVfstype,
		this.BlocksFreePercent)
}

// return: [][$fsSpec, $fsFile, $fsVfstype]
func ListMountPoint() ([][3]string, error) {
	contents, err := ioutil.ReadFile(Root() + "/proc/mounts")
	if err != nil {
		return nil, err
	}

	ret := make([][3]string, 0)

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		line, err := utils.ReadLine(reader)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}

		fields := strings.Fields(string(line))

		fsSpec := fields[0]
		fsFile := fields[1]
		fsVfstype := fields[2]

		if _, exist := FSSPEC_IGNORE[fsSpec]; exist {
			continue
		}

		if _, exist := FSTYPE_IGNORE[fsVfstype]; exist {
			continue
		}

		if strings.HasPrefix(fsVfstype, "fuse") {
			continue
		}

		if IgnoreFsFile(fsFile) {
			continue
		}

		// keep /dev/xxx device with shorter fsFile (remove mount binds)
		if strings.HasPrefix(fsSpec, Root()+"/dev") {
			deviceFound := false
			for idx := range ret {
				if ret[idx][0] == fsSpec {
					deviceFound = true
					if len(fsFile) < len(ret[idx][1]) {
						ret[idx][1] = fsFile
					}
					break
				}
			}
			if !deviceFound {
				ret = append(ret, [3]string{fsSpec, fsFile, fsVfstype})
			}
		} else {
			ret = append(ret, [3]string{fsSpec, fsFile, fsVfstype})
		}
	}
	return ret, nil
}

func BuildDeviceUsage(_fsSpec, _fsFile, _fsVfstype string) (*DeviceUsage, error) {
	ret := &DeviceUsage{FsSpec: _fsSpec, FsFile: _fsFile, FsVfstype: _fsVfstype}

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(_fsFile, &fs)
	if err != nil {
		return nil, err
	}
	// blocks
	used := fs.Blocks - fs.Bfree
	// ret.BlocksAll = uint64(fs.Frsize) * fs.Blocks
	// ret.BlocksUsed = uint64(fs.Frsize) * used
	// ret.BlocksFree = uint64(fs.Frsize) * fs.Bavail
	if fs.Blocks == 0 {
		ret.BlocksUsedPercent = 0
		ret.BlocksFreePercent = 0
	} else {
		ret.BlocksUsedPercent = float64(used) * 100.0 / float64(used+fs.Bavail)
		ret.BlocksFreePercent = 100.0 - ret.BlocksUsedPercent
	}

	// inodes
	ret.InodesAll = fs.Files
	if fs.Ffree == math.MaxUint64 {
		ret.InodesFree = 0
		ret.InodesUsed = 0
	} else {
		ret.InodesFree = fs.Ffree
		ret.InodesUsed = fs.Files - fs.Ffree
	}
	if fs.Files == 0 {
		ret.InodesUsedPercent = 0
		ret.InodesFreePercent = 0
	} else {
		ret.InodesUsedPercent = float64(ret.InodesUsed) * 100.0 / float64(ret.InodesAll)
		ret.InodesFreePercent = 100.0 - ret.InodesUsedPercent
	}

	return ret, nil
}
