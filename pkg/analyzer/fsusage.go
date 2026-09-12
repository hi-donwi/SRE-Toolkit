//go:build unix

package analyzer

import "syscall"

// fsUsage is a filesystem occupancy snapshot, in the same terms df(1) reports:
// blocks that exist, blocks free to root, and blocks available to unprivileged
// users. The gap between free and available is the reserved superuser pool.
type fsUsage struct {
	BlockSize   uint64
	Blocks      uint64
	BlocksFree  uint64
	BlocksAvail uint64
	Inodes      uint64
	InodesFree  uint64
}

// statfsUsage reads live filesystem statistics for path.
func statfsUsage(path string) (fsUsage, error) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return fsUsage{}, err
	}
	return fsUsage{
		BlockSize:   uint64(st.Bsize),
		Blocks:      uint64(st.Blocks),
		BlocksFree:  uint64(st.Bfree),
		BlocksAvail: uint64(st.Bavail),
		Inodes:      uint64(st.Files),
		InodesFree:  uint64(st.Ffree),
	}, nil
}

// UsedPercent reports occupancy the way df(1) does: against the space actually
// usable by an unprivileged process, not against raw capacity. Counting the
// root-reserved pool as free is what makes a naive check report ~5% headroom on
// a filesystem where ordinary writes are already failing.
func (u fsUsage) UsedPercent() float64 {
	used := u.Blocks - u.BlocksFree
	usable := used + u.BlocksAvail
	if usable == 0 {
		return 0
	}
	return float64(used) / float64(usable) * 100.0
}

// InodePercent reports inode table occupancy.
func (u fsUsage) InodePercent() float64 {
	if u.Inodes == 0 {
		return 0
	}
	return float64(u.Inodes-u.InodesFree) / float64(u.Inodes) * 100.0
}

// TotalBytes reports raw filesystem capacity.
func (u fsUsage) TotalBytes() uint64 { return u.Blocks * u.BlockSize }

// UsedBytes reports consumed bytes.
func (u fsUsage) UsedBytes() uint64 { return (u.Blocks - u.BlocksFree) * u.BlockSize }

// AvailBytes reports bytes an unprivileged process may still write.
func (u fsUsage) AvailBytes() uint64 { return u.BlocksAvail * u.BlockSize }
