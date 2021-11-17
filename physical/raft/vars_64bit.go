// +build amd64 arm64 s390x

package raft

const initialMmapSize = 100 * 1024 * 1024 * 1024 // 100GB
