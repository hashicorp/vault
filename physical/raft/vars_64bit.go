// +build !386,!arm

package raft

const initialMmapSize = 100 * 1024 * 1024 * 1024 // 100GB
