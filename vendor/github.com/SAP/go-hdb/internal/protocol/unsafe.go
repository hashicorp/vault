/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protocol

import "unsafe"

func init() {
	//check if pointer could be stored in uint64
	var ptr uintptr
	var ui64 uint64

	if unsafe.Sizeof(ptr) > unsafe.Sizeof(ui64) {
		panic("pointer size exceeds uint64 size")
	}
}

// LobWriteDescrToPointer returns a pointer to a LobWriteDescr compatible to sql/driver/Value (int64).
func LobWriteDescrToPointer(w *LobWriteDescr) int64 {
	return int64(uintptr(unsafe.Pointer(w)))
}

func pointerToLobWriteDescr(ptr int64) *LobWriteDescr {
	return (*LobWriteDescr)(unsafe.Pointer(uintptr(ptr)))
}

func lobReadDescrToPointer(r *LobReadDescr) int64 {
	return int64(uintptr(unsafe.Pointer(r)))
}

// PointerToLobReadDescr returns the address of a LobReadDescr from an sql/driver/Value (int64) compatible pointer.
func PointerToLobReadDescr(ptr int64) *LobReadDescr {
	return (*LobReadDescr)(unsafe.Pointer(uintptr(ptr)))
}
