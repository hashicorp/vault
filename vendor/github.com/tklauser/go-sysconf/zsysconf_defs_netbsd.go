// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs sysconf_defs_netbsd.go

//go:build netbsd
// +build netbsd

package sysconf

const (
	SC_AIO_LISTIO_MAX               = 0x33
	SC_AIO_MAX                      = 0x34
	SC_ARG_MAX                      = 0x1
	SC_ATEXIT_MAX                   = 0x28
	SC_BC_BASE_MAX                  = 0x9
	SC_BC_DIM_MAX                   = 0xa
	SC_BC_SCALE_MAX                 = 0xb
	SC_BC_STRING_MAX                = 0xc
	SC_CHILD_MAX                    = 0x2
	SC_CLK_TCK                      = 0x27
	SC_COLL_WEIGHTS_MAX             = 0xd
	SC_EXPR_NEST_MAX                = 0xe
	SC_HOST_NAME_MAX                = 0x45
	SC_IOV_MAX                      = 0x20
	SC_LINE_MAX                     = 0xf
	SC_LOGIN_NAME_MAX               = 0x25
	SC_MQ_OPEN_MAX                  = 0x36
	SC_MQ_PRIO_MAX                  = 0x37
	SC_NGROUPS_MAX                  = 0x4
	SC_OPEN_MAX                     = 0x5
	SC_PAGE_SIZE                    = 0x1c
	SC_PAGESIZE                     = 0x1c
	SC_THREAD_DESTRUCTOR_ITERATIONS = 0x39
	SC_THREAD_KEYS_MAX              = 0x3a
	SC_THREAD_STACK_MIN             = 0x3b
	SC_THREAD_THREADS_MAX           = 0x3c
	SC_RE_DUP_MAX                   = 0x10
	SC_STREAM_MAX                   = 0x1a
	SC_SYMLOOP_MAX                  = 0x49
	SC_TTY_NAME_MAX                 = 0x44
	SC_TZNAME_MAX                   = 0x1b

	SC_ASYNCHRONOUS_IO = 0x32
	SC_BARRIERS        = 0x2b
	SC_FSYNC           = 0x1d
	SC_JOB_CONTROL     = 0x6
	SC_MAPPED_FILES    = 0x21
	SC_SEMAPHORES      = 0x2a
	SC_SHELL           = 0x48
	SC_THREADS         = 0x29
	SC_TIMERS          = 0x2c
	SC_VERSION         = 0x8

	SC_2_VERSION   = 0x11
	SC_2_C_DEV     = 0x13
	SC_2_FORT_DEV  = 0x15
	SC_2_FORT_RUN  = 0x16
	SC_2_LOCALEDEF = 0x17
	SC_2_SW_DEV    = 0x18
	SC_2_UPE       = 0x19

	SC_PHYS_PAGES       = 0x79
	SC_MONOTONIC_CLOCK  = 0x26
	SC_NPROCESSORS_CONF = 0x3e9
	SC_NPROCESSORS_ONLN = 0x3ea
)

const (
	_MAXHOSTNAMELEN = 0x100
	_MAXLOGNAME     = 0x10
	_MAXSYMLINKS    = 0x20

	_POSIX_ARG_MAX                      = 0x1000
	_POSIX_CHILD_MAX                    = 0x19
	_POSIX_SHELL                        = 0x1
	_POSIX_THREAD_DESTRUCTOR_ITERATIONS = 0x4
	_POSIX_THREAD_KEYS_MAX              = 0x100
	_POSIX_VERSION                      = 0x30db0

	_POSIX2_VERSION = 0x30db0

	_FOPEN_MAX  = 0x14
	_NAME_MAX   = 0x1ff
	_RE_DUP_MAX = 0xff

	_BC_BASE_MAX      = 0x7fffffff
	_BC_DIM_MAX       = 0xffff
	_BC_SCALE_MAX     = 0x7fffffff
	_BC_STRING_MAX    = 0x7fffffff
	_COLL_WEIGHTS_MAX = 0x2
	_EXPR_NEST_MAX    = 0x20
	_LINE_MAX         = 0x800

	_PATH_DEV      = "/dev/"
	_PATH_ZONEINFO = "/usr/share/zoneinfo"
)

const _PC_NAME_MAX = 0x4
