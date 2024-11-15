package system // import "github.com/ory/dockertest/docker/pkg/system"

import "syscall"

// fromStatT converts a syscall.Stat_t type to a system.Stat_t type
func fromStatT(s *syscall.Stat_t) (*StatT, error) {
	return &StatT{size: s.Size,
		mode: s.Mode,
		uid:  s.Uid,
		gid:  s.Gid,
		rdev: s.Rdev,
		mtim: s.Mtim}, nil
}

// FromStatT converts a syscall.Stat_t type to a system.Stat_t type
// This is exposed on Linux as pkg/archive/changes uses it.
func FromStatT(s *syscall.Stat_t) (*StatT, error) {
	return fromStatT(s)
}
