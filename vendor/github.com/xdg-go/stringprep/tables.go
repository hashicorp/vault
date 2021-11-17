// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package stringprep

var tableA1 = Set{
	RuneRange{0x0221, 0x0221},
	RuneRange{0x0234, 0x024F},
	RuneRange{0x02AE, 0x02AF},
	RuneRange{0x02EF, 0x02FF},
	RuneRange{0x0350, 0x035F},
	RuneRange{0x0370, 0x0373},
	RuneRange{0x0376, 0x0379},
	RuneRange{0x037B, 0x037D},
	RuneRange{0x037F, 0x0383},
	RuneRange{0x038B, 0x038B},
	RuneRange{0x038D, 0x038D},
	RuneRange{0x03A2, 0x03A2},
	RuneRange{0x03CF, 0x03CF},
	RuneRange{0x03F7, 0x03FF},
	RuneRange{0x0487, 0x0487},
	RuneRange{0x04CF, 0x04CF},
	RuneRange{0x04F6, 0x04F7},
	RuneRange{0x04FA, 0x04FF},
	RuneRange{0x0510, 0x0530},
	RuneRange{0x0557, 0x0558},
	RuneRange{0x0560, 0x0560},
	RuneRange{0x0588, 0x0588},
	RuneRange{0x058B, 0x0590},
	RuneRange{0x05A2, 0x05A2},
	RuneRange{0x05BA, 0x05BA},
	RuneRange{0x05C5, 0x05CF},
	RuneRange{0x05EB, 0x05EF},
	RuneRange{0x05F5, 0x060B},
	RuneRange{0x060D, 0x061A},
	RuneRange{0x061C, 0x061E},
	RuneRange{0x0620, 0x0620},
	RuneRange{0x063B, 0x063F},
	RuneRange{0x0656, 0x065F},
	RuneRange{0x06EE, 0x06EF},
	RuneRange{0x06FF, 0x06FF},
	RuneRange{0x070E, 0x070E},
	RuneRange{0x072D, 0x072F},
	RuneRange{0x074B, 0x077F},
	RuneRange{0x07B2, 0x0900},
	RuneRange{0x0904, 0x0904},
	RuneRange{0x093A, 0x093B},
	RuneRange{0x094E, 0x094F},
	RuneRange{0x0955, 0x0957},
	RuneRange{0x0971, 0x0980},
	RuneRange{0x0984, 0x0984},
	RuneRange{0x098D, 0x098E},
	RuneRange{0x0991, 0x0992},
	RuneRange{0x09A9, 0x09A9},
	RuneRange{0x09B1, 0x09B1},
	RuneRange{0x09B3, 0x09B5},
	RuneRange{0x09BA, 0x09BB},
	RuneRange{0x09BD, 0x09BD},
	RuneRange{0x09C5, 0x09C6},
	RuneRange{0x09C9, 0x09CA},
	RuneRange{0x09CE, 0x09D6},
	RuneRange{0x09D8, 0x09DB},
	RuneRange{0x09DE, 0x09DE},
	RuneRange{0x09E4, 0x09E5},
	RuneRange{0x09FB, 0x0A01},
	RuneRange{0x0A03, 0x0A04},
	RuneRange{0x0A0B, 0x0A0E},
	RuneRange{0x0A11, 0x0A12},
	RuneRange{0x0A29, 0x0A29},
	RuneRange{0x0A31, 0x0A31},
	RuneRange{0x0A34, 0x0A34},
	RuneRange{0x0A37, 0x0A37},
	RuneRange{0x0A3A, 0x0A3B},
	RuneRange{0x0A3D, 0x0A3D},
	RuneRange{0x0A43, 0x0A46},
	RuneRange{0x0A49, 0x0A4A},
	RuneRange{0x0A4E, 0x0A58},
	RuneRange{0x0A5D, 0x0A5D},
	RuneRange{0x0A5F, 0x0A65},
	RuneRange{0x0A75, 0x0A80},
	RuneRange{0x0A84, 0x0A84},
	RuneRange{0x0A8C, 0x0A8C},
	RuneRange{0x0A8E, 0x0A8E},
	RuneRange{0x0A92, 0x0A92},
	RuneRange{0x0AA9, 0x0AA9},
	RuneRange{0x0AB1, 0x0AB1},
	RuneRange{0x0AB4, 0x0AB4},
	RuneRange{0x0ABA, 0x0ABB},
	RuneRange{0x0AC6, 0x0AC6},
	RuneRange{0x0ACA, 0x0ACA},
	RuneRange{0x0ACE, 0x0ACF},
	RuneRange{0x0AD1, 0x0ADF},
	RuneRange{0x0AE1, 0x0AE5},
	RuneRange{0x0AF0, 0x0B00},
	RuneRange{0x0B04, 0x0B04},
	RuneRange{0x0B0D, 0x0B0E},
	RuneRange{0x0B11, 0x0B12},
	RuneRange{0x0B29, 0x0B29},
	RuneRange{0x0B31, 0x0B31},
	RuneRange{0x0B34, 0x0B35},
	RuneRange{0x0B3A, 0x0B3B},
	RuneRange{0x0B44, 0x0B46},
	RuneRange{0x0B49, 0x0B4A},
	RuneRange{0x0B4E, 0x0B55},
	RuneRange{0x0B58, 0x0B5B},
	RuneRange{0x0B5E, 0x0B5E},
	RuneRange{0x0B62, 0x0B65},
	RuneRange{0x0B71, 0x0B81},
	RuneRange{0x0B84, 0x0B84},
	RuneRange{0x0B8B, 0x0B8D},
	RuneRange{0x0B91, 0x0B91},
	RuneRange{0x0B96, 0x0B98},
	RuneRange{0x0B9B, 0x0B9B},
	RuneRange{0x0B9D, 0x0B9D},
	RuneRange{0x0BA0, 0x0BA2},
	RuneRange{0x0BA5, 0x0BA7},
	RuneRange{0x0BAB, 0x0BAD},
	RuneRange{0x0BB6, 0x0BB6},
	RuneRange{0x0BBA, 0x0BBD},
	RuneRange{0x0BC3, 0x0BC5},
	RuneRange{0x0BC9, 0x0BC9},
	RuneRange{0x0BCE, 0x0BD6},
	RuneRange{0x0BD8, 0x0BE6},
	RuneRange{0x0BF3, 0x0C00},
	RuneRange{0x0C04, 0x0C04},
	RuneRange{0x0C0D, 0x0C0D},
	RuneRange{0x0C11, 0x0C11},
	RuneRange{0x0C29, 0x0C29},
	RuneRange{0x0C34, 0x0C34},
	RuneRange{0x0C3A, 0x0C3D},
	RuneRange{0x0C45, 0x0C45},
	RuneRange{0x0C49, 0x0C49},
	RuneRange{0x0C4E, 0x0C54},
	RuneRange{0x0C57, 0x0C5F},
	RuneRange{0x0C62, 0x0C65},
	RuneRange{0x0C70, 0x0C81},
	RuneRange{0x0C84, 0x0C84},
	RuneRange{0x0C8D, 0x0C8D},
	RuneRange{0x0C91, 0x0C91},
	RuneRange{0x0CA9, 0x0CA9},
	RuneRange{0x0CB4, 0x0CB4},
	RuneRange{0x0CBA, 0x0CBD},
	RuneRange{0x0CC5, 0x0CC5},
	RuneRange{0x0CC9, 0x0CC9},
	RuneRange{0x0CCE, 0x0CD4},
	RuneRange{0x0CD7, 0x0CDD},
	RuneRange{0x0CDF, 0x0CDF},
	RuneRange{0x0CE2, 0x0CE5},
	RuneRange{0x0CF0, 0x0D01},
	RuneRange{0x0D04, 0x0D04},
	RuneRange{0x0D0D, 0x0D0D},
	RuneRange{0x0D11, 0x0D11},
	RuneRange{0x0D29, 0x0D29},
	RuneRange{0x0D3A, 0x0D3D},
	RuneRange{0x0D44, 0x0D45},
	RuneRange{0x0D49, 0x0D49},
	RuneRange{0x0D4E, 0x0D56},
	RuneRange{0x0D58, 0x0D5F},
	RuneRange{0x0D62, 0x0D65},
	RuneRange{0x0D70, 0x0D81},
	RuneRange{0x0D84, 0x0D84},
	RuneRange{0x0D97, 0x0D99},
	RuneRange{0x0DB2, 0x0DB2},
	RuneRange{0x0DBC, 0x0DBC},
	RuneRange{0x0DBE, 0x0DBF},
	RuneRange{0x0DC7, 0x0DC9},
	RuneRange{0x0DCB, 0x0DCE},
	RuneRange{0x0DD5, 0x0DD5},
	RuneRange{0x0DD7, 0x0DD7},
	RuneRange{0x0DE0, 0x0DF1},
	RuneRange{0x0DF5, 0x0E00},
	RuneRange{0x0E3B, 0x0E3E},
	RuneRange{0x0E5C, 0x0E80},
	RuneRange{0x0E83, 0x0E83},
	RuneRange{0x0E85, 0x0E86},
	RuneRange{0x0E89, 0x0E89},
	RuneRange{0x0E8B, 0x0E8C},
	RuneRange{0x0E8E, 0x0E93},
	RuneRange{0x0E98, 0x0E98},
	RuneRange{0x0EA0, 0x0EA0},
	RuneRange{0x0EA4, 0x0EA4},
	RuneRange{0x0EA6, 0x0EA6},
	RuneRange{0x0EA8, 0x0EA9},
	RuneRange{0x0EAC, 0x0EAC},
	RuneRange{0x0EBA, 0x0EBA},
	RuneRange{0x0EBE, 0x0EBF},
	RuneRange{0x0EC5, 0x0EC5},
	RuneRange{0x0EC7, 0x0EC7},
	RuneRange{0x0ECE, 0x0ECF},
	RuneRange{0x0EDA, 0x0EDB},
	RuneRange{0x0EDE, 0x0EFF},
	RuneRange{0x0F48, 0x0F48},
	RuneRange{0x0F6B, 0x0F70},
	RuneRange{0x0F8C, 0x0F8F},
	RuneRange{0x0F98, 0x0F98},
	RuneRange{0x0FBD, 0x0FBD},
	RuneRange{0x0FCD, 0x0FCE},
	RuneRange{0x0FD0, 0x0FFF},
	RuneRange{0x1022, 0x1022},
	RuneRange{0x1028, 0x1028},
	RuneRange{0x102B, 0x102B},
	RuneRange{0x1033, 0x1035},
	RuneRange{0x103A, 0x103F},
	RuneRange{0x105A, 0x109F},
	RuneRange{0x10C6, 0x10CF},
	RuneRange{0x10F9, 0x10FA},
	RuneRange{0x10FC, 0x10FF},
	RuneRange{0x115A, 0x115E},
	RuneRange{0x11A3, 0x11A7},
	RuneRange{0x11FA, 0x11FF},
	RuneRange{0x1207, 0x1207},
	RuneRange{0x1247, 0x1247},
	RuneRange{0x1249, 0x1249},
	RuneRange{0x124E, 0x124F},
	RuneRange{0x1257, 0x1257},
	RuneRange{0x1259, 0x1259},
	RuneRange{0x125E, 0x125F},
	RuneRange{0x1287, 0x1287},
	RuneRange{0x1289, 0x1289},
	RuneRange{0x128E, 0x128F},
	RuneRange{0x12AF, 0x12AF},
	RuneRange{0x12B1, 0x12B1},
	RuneRange{0x12B6, 0x12B7},
	RuneRange{0x12BF, 0x12BF},
	RuneRange{0x12C1, 0x12C1},
	RuneRange{0x12C6, 0x12C7},
	RuneRange{0x12CF, 0x12CF},
	RuneRange{0x12D7, 0x12D7},
	RuneRange{0x12EF, 0x12EF},
	RuneRange{0x130F, 0x130F},
	RuneRange{0x1311, 0x1311},
	RuneRange{0x1316, 0x1317},
	RuneRange{0x131F, 0x131F},
	RuneRange{0x1347, 0x1347},
	RuneRange{0x135B, 0x1360},
	RuneRange{0x137D, 0x139F},
	RuneRange{0x13F5, 0x1400},
	RuneRange{0x1677, 0x167F},
	RuneRange{0x169D, 0x169F},
	RuneRange{0x16F1, 0x16FF},
	RuneRange{0x170D, 0x170D},
	RuneRange{0x1715, 0x171F},
	RuneRange{0x1737, 0x173F},
	RuneRange{0x1754, 0x175F},
	RuneRange{0x176D, 0x176D},
	RuneRange{0x1771, 0x1771},
	RuneRange{0x1774, 0x177F},
	RuneRange{0x17DD, 0x17DF},
	RuneRange{0x17EA, 0x17FF},
	RuneRange{0x180F, 0x180F},
	RuneRange{0x181A, 0x181F},
	RuneRange{0x1878, 0x187F},
	RuneRange{0x18AA, 0x1DFF},
	RuneRange{0x1E9C, 0x1E9F},
	RuneRange{0x1EFA, 0x1EFF},
	RuneRange{0x1F16, 0x1F17},
	RuneRange{0x1F1E, 0x1F1F},
	RuneRange{0x1F46, 0x1F47},
	RuneRange{0x1F4E, 0x1F4F},
	RuneRange{0x1F58, 0x1F58},
	RuneRange{0x1F5A, 0x1F5A},
	RuneRange{0x1F5C, 0x1F5C},
	RuneRange{0x1F5E, 0x1F5E},
	RuneRange{0x1F7E, 0x1F7F},
	RuneRange{0x1FB5, 0x1FB5},
	RuneRange{0x1FC5, 0x1FC5},
	RuneRange{0x1FD4, 0x1FD5},
	RuneRange{0x1FDC, 0x1FDC},
	RuneRange{0x1FF0, 0x1FF1},
	RuneRange{0x1FF5, 0x1FF5},
	RuneRange{0x1FFF, 0x1FFF},
	RuneRange{0x2053, 0x2056},
	RuneRange{0x2058, 0x205E},
	RuneRange{0x2064, 0x2069},
	RuneRange{0x2072, 0x2073},
	RuneRange{0x208F, 0x209F},
	RuneRange{0x20B2, 0x20CF},
	RuneRange{0x20EB, 0x20FF},
	RuneRange{0x213B, 0x213C},
	RuneRange{0x214C, 0x2152},
	RuneRange{0x2184, 0x218F},
	RuneRange{0x23CF, 0x23FF},
	RuneRange{0x2427, 0x243F},
	RuneRange{0x244B, 0x245F},
	RuneRange{0x24FF, 0x24FF},
	RuneRange{0x2614, 0x2615},
	RuneRange{0x2618, 0x2618},
	RuneRange{0x267E, 0x267F},
	RuneRange{0x268A, 0x2700},
	RuneRange{0x2705, 0x2705},
	RuneRange{0x270A, 0x270B},
	RuneRange{0x2728, 0x2728},
	RuneRange{0x274C, 0x274C},
	RuneRange{0x274E, 0x274E},
	RuneRange{0x2753, 0x2755},
	RuneRange{0x2757, 0x2757},
	RuneRange{0x275F, 0x2760},
	RuneRange{0x2795, 0x2797},
	RuneRange{0x27B0, 0x27B0},
	RuneRange{0x27BF, 0x27CF},
	RuneRange{0x27EC, 0x27EF},
	RuneRange{0x2B00, 0x2E7F},
	RuneRange{0x2E9A, 0x2E9A},
	RuneRange{0x2EF4, 0x2EFF},
	RuneRange{0x2FD6, 0x2FEF},
	RuneRange{0x2FFC, 0x2FFF},
	RuneRange{0x3040, 0x3040},
	RuneRange{0x3097, 0x3098},
	RuneRange{0x3100, 0x3104},
	RuneRange{0x312D, 0x3130},
	RuneRange{0x318F, 0x318F},
	RuneRange{0x31B8, 0x31EF},
	RuneRange{0x321D, 0x321F},
	RuneRange{0x3244, 0x3250},
	RuneRange{0x327C, 0x327E},
	RuneRange{0x32CC, 0x32CF},
	RuneRange{0x32FF, 0x32FF},
	RuneRange{0x3377, 0x337A},
	RuneRange{0x33DE, 0x33DF},
	RuneRange{0x33FF, 0x33FF},
	RuneRange{0x4DB6, 0x4DFF},
	RuneRange{0x9FA6, 0x9FFF},
	RuneRange{0xA48D, 0xA48F},
	RuneRange{0xA4C7, 0xABFF},
	RuneRange{0xD7A4, 0xD7FF},
	RuneRange{0xFA2E, 0xFA2F},
	RuneRange{0xFA6B, 0xFAFF},
	RuneRange{0xFB07, 0xFB12},
	RuneRange{0xFB18, 0xFB1C},
	RuneRange{0xFB37, 0xFB37},
	RuneRange{0xFB3D, 0xFB3D},
	RuneRange{0xFB3F, 0xFB3F},
	RuneRange{0xFB42, 0xFB42},
	RuneRange{0xFB45, 0xFB45},
	RuneRange{0xFBB2, 0xFBD2},
	RuneRange{0xFD40, 0xFD4F},
	RuneRange{0xFD90, 0xFD91},
	RuneRange{0xFDC8, 0xFDCF},
	RuneRange{0xFDFD, 0xFDFF},
	RuneRange{0xFE10, 0xFE1F},
	RuneRange{0xFE24, 0xFE2F},
	RuneRange{0xFE47, 0xFE48},
	RuneRange{0xFE53, 0xFE53},
	RuneRange{0xFE67, 0xFE67},
	RuneRange{0xFE6C, 0xFE6F},
	RuneRange{0xFE75, 0xFE75},
	RuneRange{0xFEFD, 0xFEFE},
	RuneRange{0xFF00, 0xFF00},
	RuneRange{0xFFBF, 0xFFC1},
	RuneRange{0xFFC8, 0xFFC9},
	RuneRange{0xFFD0, 0xFFD1},
	RuneRange{0xFFD8, 0xFFD9},
	RuneRange{0xFFDD, 0xFFDF},
	RuneRange{0xFFE7, 0xFFE7},
	RuneRange{0xFFEF, 0xFFF8},
	RuneRange{0x10000, 0x102FF},
	RuneRange{0x1031F, 0x1031F},
	RuneRange{0x10324, 0x1032F},
	RuneRange{0x1034B, 0x103FF},
	RuneRange{0x10426, 0x10427},
	RuneRange{0x1044E, 0x1CFFF},
	RuneRange{0x1D0F6, 0x1D0FF},
	RuneRange{0x1D127, 0x1D129},
	RuneRange{0x1D1DE, 0x1D3FF},
	RuneRange{0x1D455, 0x1D455},
	RuneRange{0x1D49D, 0x1D49D},
	RuneRange{0x1D4A0, 0x1D4A1},
	RuneRange{0x1D4A3, 0x1D4A4},
	RuneRange{0x1D4A7, 0x1D4A8},
	RuneRange{0x1D4AD, 0x1D4AD},
	RuneRange{0x1D4BA, 0x1D4BA},
	RuneRange{0x1D4BC, 0x1D4BC},
	RuneRange{0x1D4C1, 0x1D4C1},
	RuneRange{0x1D4C4, 0x1D4C4},
	RuneRange{0x1D506, 0x1D506},
	RuneRange{0x1D50B, 0x1D50C},
	RuneRange{0x1D515, 0x1D515},
	RuneRange{0x1D51D, 0x1D51D},
	RuneRange{0x1D53A, 0x1D53A},
	RuneRange{0x1D53F, 0x1D53F},
	RuneRange{0x1D545, 0x1D545},
	RuneRange{0x1D547, 0x1D549},
	RuneRange{0x1D551, 0x1D551},
	RuneRange{0x1D6A4, 0x1D6A7},
	RuneRange{0x1D7CA, 0x1D7CD},
	RuneRange{0x1D800, 0x1FFFD},
	RuneRange{0x2A6D7, 0x2F7FF},
	RuneRange{0x2FA1E, 0x2FFFD},
	RuneRange{0x30000, 0x3FFFD},
	RuneRange{0x40000, 0x4FFFD},
	RuneRange{0x50000, 0x5FFFD},
	RuneRange{0x60000, 0x6FFFD},
	RuneRange{0x70000, 0x7FFFD},
	RuneRange{0x80000, 0x8FFFD},
	RuneRange{0x90000, 0x9FFFD},
	RuneRange{0xA0000, 0xAFFFD},
	RuneRange{0xB0000, 0xBFFFD},
	RuneRange{0xC0000, 0xCFFFD},
	RuneRange{0xD0000, 0xDFFFD},
	RuneRange{0xE0000, 0xE0000},
	RuneRange{0xE0002, 0xE001F},
	RuneRange{0xE0080, 0xEFFFD},
}

// TableA1 represents RFC-3454 Table A.1.
var TableA1 Set = tableA1

var tableB1 = Mapping{
	0x00AD: []rune{}, // Map to nothing
	0x034F: []rune{}, // Map to nothing
	0x180B: []rune{}, // Map to nothing
	0x180C: []rune{}, // Map to nothing
	0x180D: []rune{}, // Map to nothing
	0x200B: []rune{}, // Map to nothing
	0x200C: []rune{}, // Map to nothing
	0x200D: []rune{}, // Map to nothing
	0x2060: []rune{}, // Map to nothing
	0xFE00: []rune{}, // Map to nothing
	0xFE01: []rune{}, // Map to nothing
	0xFE02: []rune{}, // Map to nothing
	0xFE03: []rune{}, // Map to nothing
	0xFE04: []rune{}, // Map to nothing
	0xFE05: []rune{}, // Map to nothing
	0xFE06: []rune{}, // Map to nothing
	0xFE07: []rune{}, // Map to nothing
	0xFE08: []rune{}, // Map to nothing
	0xFE09: []rune{}, // Map to nothing
	0xFE0A: []rune{}, // Map to nothing
	0xFE0B: []rune{}, // Map to nothing
	0xFE0C: []rune{}, // Map to nothing
	0xFE0D: []rune{}, // Map to nothing
	0xFE0E: []rune{}, // Map to nothing
	0xFE0F: []rune{}, // Map to nothing
	0xFEFF: []rune{}, // Map to nothing
}

// TableB1 represents RFC-3454 Table B.1.
var TableB1 Mapping = tableB1

var tableB2 = Mapping{
	0x0041:  []rune{0x0061},                         // Case map
	0x0042:  []rune{0x0062},                         // Case map
	0x0043:  []rune{0x0063},                         // Case map
	0x0044:  []rune{0x0064},                         // Case map
	0x0045:  []rune{0x0065},                         // Case map
	0x0046:  []rune{0x0066},                         // Case map
	0x0047:  []rune{0x0067},                         // Case map
	0x0048:  []rune{0x0068},                         // Case map
	0x0049:  []rune{0x0069},                         // Case map
	0x004A:  []rune{0x006A},                         // Case map
	0x004B:  []rune{0x006B},                         // Case map
	0x004C:  []rune{0x006C},                         // Case map
	0x004D:  []rune{0x006D},                         // Case map
	0x004E:  []rune{0x006E},                         // Case map
	0x004F:  []rune{0x006F},                         // Case map
	0x0050:  []rune{0x0070},                         // Case map
	0x0051:  []rune{0x0071},                         // Case map
	0x0052:  []rune{0x0072},                         // Case map
	0x0053:  []rune{0x0073},                         // Case map
	0x0054:  []rune{0x0074},                         // Case map
	0x0055:  []rune{0x0075},                         // Case map
	0x0056:  []rune{0x0076},                         // Case map
	0x0057:  []rune{0x0077},                         // Case map
	0x0058:  []rune{0x0078},                         // Case map
	0x0059:  []rune{0x0079},                         // Case map
	0x005A:  []rune{0x007A},                         // Case map
	0x00B5:  []rune{0x03BC},                         // Case map
	0x00C0:  []rune{0x00E0},                         // Case map
	0x00C1:  []rune{0x00E1},                         // Case map
	0x00C2:  []rune{0x00E2},                         // Case map
	0x00C3:  []rune{0x00E3},                         // Case map
	0x00C4:  []rune{0x00E4},                         // Case map
	0x00C5:  []rune{0x00E5},                         // Case map
	0x00C6:  []rune{0x00E6},                         // Case map
	0x00C7:  []rune{0x00E7},                         // Case map
	0x00C8:  []rune{0x00E8},                         // Case map
	0x00C9:  []rune{0x00E9},                         // Case map
	0x00CA:  []rune{0x00EA},                         // Case map
	0x00CB:  []rune{0x00EB},                         // Case map
	0x00CC:  []rune{0x00EC},                         // Case map
	0x00CD:  []rune{0x00ED},                         // Case map
	0x00CE:  []rune{0x00EE},                         // Case map
	0x00CF:  []rune{0x00EF},                         // Case map
	0x00D0:  []rune{0x00F0},                         // Case map
	0x00D1:  []rune{0x00F1},                         // Case map
	0x00D2:  []rune{0x00F2},                         // Case map
	0x00D3:  []rune{0x00F3},                         // Case map
	0x00D4:  []rune{0x00F4},                         // Case map
	0x00D5:  []rune{0x00F5},                         // Case map
	0x00D6:  []rune{0x00F6},                         // Case map
	0x00D8:  []rune{0x00F8},                         // Case map
	0x00D9:  []rune{0x00F9},                         // Case map
	0x00DA:  []rune{0x00FA},                         // Case map
	0x00DB:  []rune{0x00FB},                         // Case map
	0x00DC:  []rune{0x00FC},                         // Case map
	0x00DD:  []rune{0x00FD},                         // Case map
	0x00DE:  []rune{0x00FE},                         // Case map
	0x00DF:  []rune{0x0073, 0x0073},                 // Case map
	0x0100:  []rune{0x0101},                         // Case map
	0x0102:  []rune{0x0103},                         // Case map
	0x0104:  []rune{0x0105},                         // Case map
	0x0106:  []rune{0x0107},                         // Case map
	0x0108:  []rune{0x0109},                         // Case map
	0x010A:  []rune{0x010B},                         // Case map
	0x010C:  []rune{0x010D},                         // Case map
	0x010E:  []rune{0x010F},                         // Case map
	0x0110:  []rune{0x0111},                         // Case map
	0x0112:  []rune{0x0113},                         // Case map
	0x0114:  []rune{0x0115},                         // Case map
	0x0116:  []rune{0x0117},                         // Case map
	0x0118:  []rune{0x0119},                         // Case map
	0x011A:  []rune{0x011B},                         // Case map
	0x011C:  []rune{0x011D},                         // Case map
	0x011E:  []rune{0x011F},                         // Case map
	0x0120:  []rune{0x0121},                         // Case map
	0x0122:  []rune{0x0123},                         // Case map
	0x0124:  []rune{0x0125},                         // Case map
	0x0126:  []rune{0x0127},                         // Case map
	0x0128:  []rune{0x0129},                         // Case map
	0x012A:  []rune{0x012B},                         // Case map
	0x012C:  []rune{0x012D},                         // Case map
	0x012E:  []rune{0x012F},                         // Case map
	0x0130:  []rune{0x0069, 0x0307},                 // Case map
	0x0132:  []rune{0x0133},                         // Case map
	0x0134:  []rune{0x0135},                         // Case map
	0x0136:  []rune{0x0137},                         // Case map
	0x0139:  []rune{0x013A},                         // Case map
	0x013B:  []rune{0x013C},                         // Case map
	0x013D:  []rune{0x013E},                         // Case map
	0x013F:  []rune{0x0140},                         // Case map
	0x0141:  []rune{0x0142},                         // Case map
	0x0143:  []rune{0x0144},                         // Case map
	0x0145:  []rune{0x0146},                         // Case map
	0x0147:  []rune{0x0148},                         // Case map
	0x0149:  []rune{0x02BC, 0x006E},                 // Case map
	0x014A:  []rune{0x014B},                         // Case map
	0x014C:  []rune{0x014D},                         // Case map
	0x014E:  []rune{0x014F},                         // Case map
	0x0150:  []rune{0x0151},                         // Case map
	0x0152:  []rune{0x0153},                         // Case map
	0x0154:  []rune{0x0155},                         // Case map
	0x0156:  []rune{0x0157},                         // Case map
	0x0158:  []rune{0x0159},                         // Case map
	0x015A:  []rune{0x015B},                         // Case map
	0x015C:  []rune{0x015D},                         // Case map
	0x015E:  []rune{0x015F},                         // Case map
	0x0160:  []rune{0x0161},                         // Case map
	0x0162:  []rune{0x0163},                         // Case map
	0x0164:  []rune{0x0165},                         // Case map
	0x0166:  []rune{0x0167},                         // Case map
	0x0168:  []rune{0x0169},                         // Case map
	0x016A:  []rune{0x016B},                         // Case map
	0x016C:  []rune{0x016D},                         // Case map
	0x016E:  []rune{0x016F},                         // Case map
	0x0170:  []rune{0x0171},                         // Case map
	0x0172:  []rune{0x0173},                         // Case map
	0x0174:  []rune{0x0175},                         // Case map
	0x0176:  []rune{0x0177},                         // Case map
	0x0178:  []rune{0x00FF},                         // Case map
	0x0179:  []rune{0x017A},                         // Case map
	0x017B:  []rune{0x017C},                         // Case map
	0x017D:  []rune{0x017E},                         // Case map
	0x017F:  []rune{0x0073},                         // Case map
	0x0181:  []rune{0x0253},                         // Case map
	0x0182:  []rune{0x0183},                         // Case map
	0x0184:  []rune{0x0185},                         // Case map
	0x0186:  []rune{0x0254},                         // Case map
	0x0187:  []rune{0x0188},                         // Case map
	0x0189:  []rune{0x0256},                         // Case map
	0x018A:  []rune{0x0257},                         // Case map
	0x018B:  []rune{0x018C},                         // Case map
	0x018E:  []rune{0x01DD},                         // Case map
	0x018F:  []rune{0x0259},                         // Case map
	0x0190:  []rune{0x025B},                         // Case map
	0x0191:  []rune{0x0192},                         // Case map
	0x0193:  []rune{0x0260},                         // Case map
	0x0194:  []rune{0x0263},                         // Case map
	0x0196:  []rune{0x0269},                         // Case map
	0x0197:  []rune{0x0268},                         // Case map
	0x0198:  []rune{0x0199},                         // Case map
	0x019C:  []rune{0x026F},                         // Case map
	0x019D:  []rune{0x0272},                         // Case map
	0x019F:  []rune{0x0275},                         // Case map
	0x01A0:  []rune{0x01A1},                         // Case map
	0x01A2:  []rune{0x01A3},                         // Case map
	0x01A4:  []rune{0x01A5},                         // Case map
	0x01A6:  []rune{0x0280},                         // Case map
	0x01A7:  []rune{0x01A8},                         // Case map
	0x01A9:  []rune{0x0283},                         // Case map
	0x01AC:  []rune{0x01AD},                         // Case map
	0x01AE:  []rune{0x0288},                         // Case map
	0x01AF:  []rune{0x01B0},                         // Case map
	0x01B1:  []rune{0x028A},                         // Case map
	0x01B2:  []rune{0x028B},                         // Case map
	0x01B3:  []rune{0x01B4},                         // Case map
	0x01B5:  []rune{0x01B6},                         // Case map
	0x01B7:  []rune{0x0292},                         // Case map
	0x01B8:  []rune{0x01B9},                         // Case map
	0x01BC:  []rune{0x01BD},                         // Case map
	0x01C4:  []rune{0x01C6},                         // Case map
	0x01C5:  []rune{0x01C6},                         // Case map
	0x01C7:  []rune{0x01C9},                         // Case map
	0x01C8:  []rune{0x01C9},                         // Case map
	0x01CA:  []rune{0x01CC},                         // Case map
	0x01CB:  []rune{0x01CC},                         // Case map
	0x01CD:  []rune{0x01CE},                         // Case map
	0x01CF:  []rune{0x01D0},                         // Case map
	0x01D1:  []rune{0x01D2},                         // Case map
	0x01D3:  []rune{0x01D4},                         // Case map
	0x01D5:  []rune{0x01D6},                         // Case map
	0x01D7:  []rune{0x01D8},                         // Case map
	0x01D9:  []rune{0x01DA},                         // Case map
	0x01DB:  []rune{0x01DC},                         // Case map
	0x01DE:  []rune{0x01DF},                         // Case map
	0x01E0:  []rune{0x01E1},                         // Case map
	0x01E2:  []rune{0x01E3},                         // Case map
	0x01E4:  []rune{0x01E5},                         // Case map
	0x01E6:  []rune{0x01E7},                         // Case map
	0x01E8:  []rune{0x01E9},                         // Case map
	0x01EA:  []rune{0x01EB},                         // Case map
	0x01EC:  []rune{0x01ED},                         // Case map
	0x01EE:  []rune{0x01EF},                         // Case map
	0x01F0:  []rune{0x006A, 0x030C},                 // Case map
	0x01F1:  []rune{0x01F3},                         // Case map
	0x01F2:  []rune{0x01F3},                         // Case map
	0x01F4:  []rune{0x01F5},                         // Case map
	0x01F6:  []rune{0x0195},                         // Case map
	0x01F7:  []rune{0x01BF},                         // Case map
	0x01F8:  []rune{0x01F9},                         // Case map
	0x01FA:  []rune{0x01FB},                         // Case map
	0x01FC:  []rune{0x01FD},                         // Case map
	0x01FE:  []rune{0x01FF},                         // Case map
	0x0200:  []rune{0x0201},                         // Case map
	0x0202:  []rune{0x0203},                         // Case map
	0x0204:  []rune{0x0205},                         // Case map
	0x0206:  []rune{0x0207},                         // Case map
	0x0208:  []rune{0x0209},                         // Case map
	0x020A:  []rune{0x020B},                         // Case map
	0x020C:  []rune{0x020D},                         // Case map
	0x020E:  []rune{0x020F},                         // Case map
	0x0210:  []rune{0x0211},                         // Case map
	0x0212:  []rune{0x0213},                         // Case map
	0x0214:  []rune{0x0215},                         // Case map
	0x0216:  []rune{0x0217},                         // Case map
	0x0218:  []rune{0x0219},                         // Case map
	0x021A:  []rune{0x021B},                         // Case map
	0x021C:  []rune{0x021D},                         // Case map
	0x021E:  []rune{0x021F},                         // Case map
	0x0220:  []rune{0x019E},                         // Case map
	0x0222:  []rune{0x0223},                         // Case map
	0x0224:  []rune{0x0225},                         // Case map
	0x0226:  []rune{0x0227},                         // Case map
	0x0228:  []rune{0x0229},                         // Case map
	0x022A:  []rune{0x022B},                         // Case map
	0x022C:  []rune{0x022D},                         // Case map
	0x022E:  []rune{0x022F},                         // Case map
	0x0230:  []rune{0x0231},                         // Case map
	0x0232:  []rune{0x0233},                         // Case map
	0x0345:  []rune{0x03B9},                         // Case map
	0x037A:  []rune{0x0020, 0x03B9},                 // Additional folding
	0x0386:  []rune{0x03AC},                         // Case map
	0x0388:  []rune{0x03AD},                         // Case map
	0x0389:  []rune{0x03AE},                         // Case map
	0x038A:  []rune{0x03AF},                         // Case map
	0x038C:  []rune{0x03CC},                         // Case map
	0x038E:  []rune{0x03CD},                         // Case map
	0x038F:  []rune{0x03CE},                         // Case map
	0x0390:  []rune{0x03B9, 0x0308, 0x0301},         // Case map
	0x0391:  []rune{0x03B1},                         // Case map
	0x0392:  []rune{0x03B2},                         // Case map
	0x0393:  []rune{0x03B3},                         // Case map
	0x0394:  []rune{0x03B4},                         // Case map
	0x0395:  []rune{0x03B5},                         // Case map
	0x0396:  []rune{0x03B6},                         // Case map
	0x0397:  []rune{0x03B7},                         // Case map
	0x0398:  []rune{0x03B8},                         // Case map
	0x0399:  []rune{0x03B9},                         // Case map
	0x039A:  []rune{0x03BA},                         // Case map
	0x039B:  []rune{0x03BB},                         // Case map
	0x039C:  []rune{0x03BC},                         // Case map
	0x039D:  []rune{0x03BD},                         // Case map
	0x039E:  []rune{0x03BE},                         // Case map
	0x039F:  []rune{0x03BF},                         // Case map
	0x03A0:  []rune{0x03C0},                         // Case map
	0x03A1:  []rune{0x03C1},                         // Case map
	0x03A3:  []rune{0x03C3},                         // Case map
	0x03A4:  []rune{0x03C4},                         // Case map
	0x03A5:  []rune{0x03C5},                         // Case map
	0x03A6:  []rune{0x03C6},                         // Case map
	0x03A7:  []rune{0x03C7},                         // Case map
	0x03A8:  []rune{0x03C8},                         // Case map
	0x03A9:  []rune{0x03C9},                         // Case map
	0x03AA:  []rune{0x03CA},                         // Case map
	0x03AB:  []rune{0x03CB},                         // Case map
	0x03B0:  []rune{0x03C5, 0x0308, 0x0301},         // Case map
	0x03C2:  []rune{0x03C3},                         // Case map
	0x03D0:  []rune{0x03B2},                         // Case map
	0x03D1:  []rune{0x03B8},                         // Case map
	0x03D2:  []rune{0x03C5},                         // Additional folding
	0x03D3:  []rune{0x03CD},                         // Additional folding
	0x03D4:  []rune{0x03CB},                         // Additional folding
	0x03D5:  []rune{0x03C6},                         // Case map
	0x03D6:  []rune{0x03C0},                         // Case map
	0x03D8:  []rune{0x03D9},                         // Case map
	0x03DA:  []rune{0x03DB},                         // Case map
	0x03DC:  []rune{0x03DD},                         // Case map
	0x03DE:  []rune{0x03DF},                         // Case map
	0x03E0:  []rune{0x03E1},                         // Case map
	0x03E2:  []rune{0x03E3},                         // Case map
	0x03E4:  []rune{0x03E5},                         // Case map
	0x03E6:  []rune{0x03E7},                         // Case map
	0x03E8:  []rune{0x03E9},                         // Case map
	0x03EA:  []rune{0x03EB},                         // Case map
	0x03EC:  []rune{0x03ED},                         // Case map
	0x03EE:  []rune{0x03EF},                         // Case map
	0x03F0:  []rune{0x03BA},                         // Case map
	0x03F1:  []rune{0x03C1},                         // Case map
	0x03F2:  []rune{0x03C3},                         // Case map
	0x03F4:  []rune{0x03B8},                         // Case map
	0x03F5:  []rune{0x03B5},                         // Case map
	0x0400:  []rune{0x0450},                         // Case map
	0x0401:  []rune{0x0451},                         // Case map
	0x0402:  []rune{0x0452},                         // Case map
	0x0403:  []rune{0x0453},                         // Case map
	0x0404:  []rune{0x0454},                         // Case map
	0x0405:  []rune{0x0455},                         // Case map
	0x0406:  []rune{0x0456},                         // Case map
	0x0407:  []rune{0x0457},                         // Case map
	0x0408:  []rune{0x0458},                         // Case map
	0x0409:  []rune{0x0459},                         // Case map
	0x040A:  []rune{0x045A},                         // Case map
	0x040B:  []rune{0x045B},                         // Case map
	0x040C:  []rune{0x045C},                         // Case map
	0x040D:  []rune{0x045D},                         // Case map
	0x040E:  []rune{0x045E},                         // Case map
	0x040F:  []rune{0x045F},                         // Case map
	0x0410:  []rune{0x0430},                         // Case map
	0x0411:  []rune{0x0431},                         // Case map
	0x0412:  []rune{0x0432},                         // Case map
	0x0413:  []rune{0x0433},                         // Case map
	0x0414:  []rune{0x0434},                         // Case map
	0x0415:  []rune{0x0435},                         // Case map
	0x0416:  []rune{0x0436},                         // Case map
	0x0417:  []rune{0x0437},                         // Case map
	0x0418:  []rune{0x0438},                         // Case map
	0x0419:  []rune{0x0439},                         // Case map
	0x041A:  []rune{0x043A},                         // Case map
	0x041B:  []rune{0x043B},                         // Case map
	0x041C:  []rune{0x043C},                         // Case map
	0x041D:  []rune{0x043D},                         // Case map
	0x041E:  []rune{0x043E},                         // Case map
	0x041F:  []rune{0x043F},                         // Case map
	0x0420:  []rune{0x0440},                         // Case map
	0x0421:  []rune{0x0441},                         // Case map
	0x0422:  []rune{0x0442},                         // Case map
	0x0423:  []rune{0x0443},                         // Case map
	0x0424:  []rune{0x0444},                         // Case map
	0x0425:  []rune{0x0445},                         // Case map
	0x0426:  []rune{0x0446},                         // Case map
	0x0427:  []rune{0x0447},                         // Case map
	0x0428:  []rune{0x0448},                         // Case map
	0x0429:  []rune{0x0449},                         // Case map
	0x042A:  []rune{0x044A},                         // Case map
	0x042B:  []rune{0x044B},                         // Case map
	0x042C:  []rune{0x044C},                         // Case map
	0x042D:  []rune{0x044D},                         // Case map
	0x042E:  []rune{0x044E},                         // Case map
	0x042F:  []rune{0x044F},                         // Case map
	0x0460:  []rune{0x0461},                         // Case map
	0x0462:  []rune{0x0463},                         // Case map
	0x0464:  []rune{0x0465},                         // Case map
	0x0466:  []rune{0x0467},                         // Case map
	0x0468:  []rune{0x0469},                         // Case map
	0x046A:  []rune{0x046B},                         // Case map
	0x046C:  []rune{0x046D},                         // Case map
	0x046E:  []rune{0x046F},                         // Case map
	0x0470:  []rune{0x0471},                         // Case map
	0x0472:  []rune{0x0473},                         // Case map
	0x0474:  []rune{0x0475},                         // Case map
	0x0476:  []rune{0x0477},                         // Case map
	0x0478:  []rune{0x0479},                         // Case map
	0x047A:  []rune{0x047B},                         // Case map
	0x047C:  []rune{0x047D},                         // Case map
	0x047E:  []rune{0x047F},                         // Case map
	0x0480:  []rune{0x0481},                         // Case map
	0x048A:  []rune{0x048B},                         // Case map
	0x048C:  []rune{0x048D},                         // Case map
	0x048E:  []rune{0x048F},                         // Case map
	0x0490:  []rune{0x0491},                         // Case map
	0x0492:  []rune{0x0493},                         // Case map
	0x0494:  []rune{0x0495},                         // Case map
	0x0496:  []rune{0x0497},                         // Case map
	0x0498:  []rune{0x0499},                         // Case map
	0x049A:  []rune{0x049B},                         // Case map
	0x049C:  []rune{0x049D},                         // Case map
	0x049E:  []rune{0x049F},                         // Case map
	0x04A0:  []rune{0x04A1},                         // Case map
	0x04A2:  []rune{0x04A3},                         // Case map
	0x04A4:  []rune{0x04A5},                         // Case map
	0x04A6:  []rune{0x04A7},                         // Case map
	0x04A8:  []rune{0x04A9},                         // Case map
	0x04AA:  []rune{0x04AB},                         // Case map
	0x04AC:  []rune{0x04AD},                         // Case map
	0x04AE:  []rune{0x04AF},                         // Case map
	0x04B0:  []rune{0x04B1},                         // Case map
	0x04B2:  []rune{0x04B3},                         // Case map
	0x04B4:  []rune{0x04B5},                         // Case map
	0x04B6:  []rune{0x04B7},                         // Case map
	0x04B8:  []rune{0x04B9},                         // Case map
	0x04BA:  []rune{0x04BB},                         // Case map
	0x04BC:  []rune{0x04BD},                         // Case map
	0x04BE:  []rune{0x04BF},                         // Case map
	0x04C1:  []rune{0x04C2},                         // Case map
	0x04C3:  []rune{0x04C4},                         // Case map
	0x04C5:  []rune{0x04C6},                         // Case map
	0x04C7:  []rune{0x04C8},                         // Case map
	0x04C9:  []rune{0x04CA},                         // Case map
	0x04CB:  []rune{0x04CC},                         // Case map
	0x04CD:  []rune{0x04CE},                         // Case map
	0x04D0:  []rune{0x04D1},                         // Case map
	0x04D2:  []rune{0x04D3},                         // Case map
	0x04D4:  []rune{0x04D5},                         // Case map
	0x04D6:  []rune{0x04D7},                         // Case map
	0x04D8:  []rune{0x04D9},                         // Case map
	0x04DA:  []rune{0x04DB},                         // Case map
	0x04DC:  []rune{0x04DD},                         // Case map
	0x04DE:  []rune{0x04DF},                         // Case map
	0x04E0:  []rune{0x04E1},                         // Case map
	0x04E2:  []rune{0x04E3},                         // Case map
	0x04E4:  []rune{0x04E5},                         // Case map
	0x04E6:  []rune{0x04E7},                         // Case map
	0x04E8:  []rune{0x04E9},                         // Case map
	0x04EA:  []rune{0x04EB},                         // Case map
	0x04EC:  []rune{0x04ED},                         // Case map
	0x04EE:  []rune{0x04EF},                         // Case map
	0x04F0:  []rune{0x04F1},                         // Case map
	0x04F2:  []rune{0x04F3},                         // Case map
	0x04F4:  []rune{0x04F5},                         // Case map
	0x04F8:  []rune{0x04F9},                         // Case map
	0x0500:  []rune{0x0501},                         // Case map
	0x0502:  []rune{0x0503},                         // Case map
	0x0504:  []rune{0x0505},                         // Case map
	0x0506:  []rune{0x0507},                         // Case map
	0x0508:  []rune{0x0509},                         // Case map
	0x050A:  []rune{0x050B},                         // Case map
	0x050C:  []rune{0x050D},                         // Case map
	0x050E:  []rune{0x050F},                         // Case map
	0x0531:  []rune{0x0561},                         // Case map
	0x0532:  []rune{0x0562},                         // Case map
	0x0533:  []rune{0x0563},                         // Case map
	0x0534:  []rune{0x0564},                         // Case map
	0x0535:  []rune{0x0565},                         // Case map
	0x0536:  []rune{0x0566},                         // Case map
	0x0537:  []rune{0x0567},                         // Case map
	0x0538:  []rune{0x0568},                         // Case map
	0x0539:  []rune{0x0569},                         // Case map
	0x053A:  []rune{0x056A},                         // Case map
	0x053B:  []rune{0x056B},                         // Case map
	0x053C:  []rune{0x056C},                         // Case map
	0x053D:  []rune{0x056D},                         // Case map
	0x053E:  []rune{0x056E},                         // Case map
	0x053F:  []rune{0x056F},                         // Case map
	0x0540:  []rune{0x0570},                         // Case map
	0x0541:  []rune{0x0571},                         // Case map
	0x0542:  []rune{0x0572},                         // Case map
	0x0543:  []rune{0x0573},                         // Case map
	0x0544:  []rune{0x0574},                         // Case map
	0x0545:  []rune{0x0575},                         // Case map
	0x0546:  []rune{0x0576},                         // Case map
	0x0547:  []rune{0x0577},                         // Case map
	0x0548:  []rune{0x0578},                         // Case map
	0x0549:  []rune{0x0579},                         // Case map
	0x054A:  []rune{0x057A},                         // Case map
	0x054B:  []rune{0x057B},                         // Case map
	0x054C:  []rune{0x057C},                         // Case map
	0x054D:  []rune{0x057D},                         // Case map
	0x054E:  []rune{0x057E},                         // Case map
	0x054F:  []rune{0x057F},                         // Case map
	0x0550:  []rune{0x0580},                         // Case map
	0x0551:  []rune{0x0581},                         // Case map
	0x0552:  []rune{0x0582},                         // Case map
	0x0553:  []rune{0x0583},                         // Case map
	0x0554:  []rune{0x0584},                         // Case map
	0x0555:  []rune{0x0585},                         // Case map
	0x0556:  []rune{0x0586},                         // Case map
	0x0587:  []rune{0x0565, 0x0582},                 // Case map
	0x1E00:  []rune{0x1E01},                         // Case map
	0x1E02:  []rune{0x1E03},                         // Case map
	0x1E04:  []rune{0x1E05},                         // Case map
	0x1E06:  []rune{0x1E07},                         // Case map
	0x1E08:  []rune{0x1E09},                         // Case map
	0x1E0A:  []rune{0x1E0B},                         // Case map
	0x1E0C:  []rune{0x1E0D},                         // Case map
	0x1E0E:  []rune{0x1E0F},                         // Case map
	0x1E10:  []rune{0x1E11},                         // Case map
	0x1E12:  []rune{0x1E13},                         // Case map
	0x1E14:  []rune{0x1E15},                         // Case map
	0x1E16:  []rune{0x1E17},                         // Case map
	0x1E18:  []rune{0x1E19},                         // Case map
	0x1E1A:  []rune{0x1E1B},                         // Case map
	0x1E1C:  []rune{0x1E1D},                         // Case map
	0x1E1E:  []rune{0x1E1F},                         // Case map
	0x1E20:  []rune{0x1E21},                         // Case map
	0x1E22:  []rune{0x1E23},                         // Case map
	0x1E24:  []rune{0x1E25},                         // Case map
	0x1E26:  []rune{0x1E27},                         // Case map
	0x1E28:  []rune{0x1E29},                         // Case map
	0x1E2A:  []rune{0x1E2B},                         // Case map
	0x1E2C:  []rune{0x1E2D},                         // Case map
	0x1E2E:  []rune{0x1E2F},                         // Case map
	0x1E30:  []rune{0x1E31},                         // Case map
	0x1E32:  []rune{0x1E33},                         // Case map
	0x1E34:  []rune{0x1E35},                         // Case map
	0x1E36:  []rune{0x1E37},                         // Case map
	0x1E38:  []rune{0x1E39},                         // Case map
	0x1E3A:  []rune{0x1E3B},                         // Case map
	0x1E3C:  []rune{0x1E3D},                         // Case map
	0x1E3E:  []rune{0x1E3F},                         // Case map
	0x1E40:  []rune{0x1E41},                         // Case map
	0x1E42:  []rune{0x1E43},                         // Case map
	0x1E44:  []rune{0x1E45},                         // Case map
	0x1E46:  []rune{0x1E47},                         // Case map
	0x1E48:  []rune{0x1E49},                         // Case map
	0x1E4A:  []rune{0x1E4B},                         // Case map
	0x1E4C:  []rune{0x1E4D},                         // Case map
	0x1E4E:  []rune{0x1E4F},                         // Case map
	0x1E50:  []rune{0x1E51},                         // Case map
	0x1E52:  []rune{0x1E53},                         // Case map
	0x1E54:  []rune{0x1E55},                         // Case map
	0x1E56:  []rune{0x1E57},                         // Case map
	0x1E58:  []rune{0x1E59},                         // Case map
	0x1E5A:  []rune{0x1E5B},                         // Case map
	0x1E5C:  []rune{0x1E5D},                         // Case map
	0x1E5E:  []rune{0x1E5F},                         // Case map
	0x1E60:  []rune{0x1E61},                         // Case map
	0x1E62:  []rune{0x1E63},                         // Case map
	0x1E64:  []rune{0x1E65},                         // Case map
	0x1E66:  []rune{0x1E67},                         // Case map
	0x1E68:  []rune{0x1E69},                         // Case map
	0x1E6A:  []rune{0x1E6B},                         // Case map
	0x1E6C:  []rune{0x1E6D},                         // Case map
	0x1E6E:  []rune{0x1E6F},                         // Case map
	0x1E70:  []rune{0x1E71},                         // Case map
	0x1E72:  []rune{0x1E73},                         // Case map
	0x1E74:  []rune{0x1E75},                         // Case map
	0x1E76:  []rune{0x1E77},                         // Case map
	0x1E78:  []rune{0x1E79},                         // Case map
	0x1E7A:  []rune{0x1E7B},                         // Case map
	0x1E7C:  []rune{0x1E7D},                         // Case map
	0x1E7E:  []rune{0x1E7F},                         // Case map
	0x1E80:  []rune{0x1E81},                         // Case map
	0x1E82:  []rune{0x1E83},                         // Case map
	0x1E84:  []rune{0x1E85},                         // Case map
	0x1E86:  []rune{0x1E87},                         // Case map
	0x1E88:  []rune{0x1E89},                         // Case map
	0x1E8A:  []rune{0x1E8B},                         // Case map
	0x1E8C:  []rune{0x1E8D},                         // Case map
	0x1E8E:  []rune{0x1E8F},                         // Case map
	0x1E90:  []rune{0x1E91},                         // Case map
	0x1E92:  []rune{0x1E93},                         // Case map
	0x1E94:  []rune{0x1E95},                         // Case map
	0x1E96:  []rune{0x0068, 0x0331},                 // Case map
	0x1E97:  []rune{0x0074, 0x0308},                 // Case map
	0x1E98:  []rune{0x0077, 0x030A},                 // Case map
	0x1E99:  []rune{0x0079, 0x030A},                 // Case map
	0x1E9A:  []rune{0x0061, 0x02BE},                 // Case map
	0x1E9B:  []rune{0x1E61},                         // Case map
	0x1EA0:  []rune{0x1EA1},                         // Case map
	0x1EA2:  []rune{0x1EA3},                         // Case map
	0x1EA4:  []rune{0x1EA5},                         // Case map
	0x1EA6:  []rune{0x1EA7},                         // Case map
	0x1EA8:  []rune{0x1EA9},                         // Case map
	0x1EAA:  []rune{0x1EAB},                         // Case map
	0x1EAC:  []rune{0x1EAD},                         // Case map
	0x1EAE:  []rune{0x1EAF},                         // Case map
	0x1EB0:  []rune{0x1EB1},                         // Case map
	0x1EB2:  []rune{0x1EB3},                         // Case map
	0x1EB4:  []rune{0x1EB5},                         // Case map
	0x1EB6:  []rune{0x1EB7},                         // Case map
	0x1EB8:  []rune{0x1EB9},                         // Case map
	0x1EBA:  []rune{0x1EBB},                         // Case map
	0x1EBC:  []rune{0x1EBD},                         // Case map
	0x1EBE:  []rune{0x1EBF},                         // Case map
	0x1EC0:  []rune{0x1EC1},                         // Case map
	0x1EC2:  []rune{0x1EC3},                         // Case map
	0x1EC4:  []rune{0x1EC5},                         // Case map
	0x1EC6:  []rune{0x1EC7},                         // Case map
	0x1EC8:  []rune{0x1EC9},                         // Case map
	0x1ECA:  []rune{0x1ECB},                         // Case map
	0x1ECC:  []rune{0x1ECD},                         // Case map
	0x1ECE:  []rune{0x1ECF},                         // Case map
	0x1ED0:  []rune{0x1ED1},                         // Case map
	0x1ED2:  []rune{0x1ED3},                         // Case map
	0x1ED4:  []rune{0x1ED5},                         // Case map
	0x1ED6:  []rune{0x1ED7},                         // Case map
	0x1ED8:  []rune{0x1ED9},                         // Case map
	0x1EDA:  []rune{0x1EDB},                         // Case map
	0x1EDC:  []rune{0x1EDD},                         // Case map
	0x1EDE:  []rune{0x1EDF},                         // Case map
	0x1EE0:  []rune{0x1EE1},                         // Case map
	0x1EE2:  []rune{0x1EE3},                         // Case map
	0x1EE4:  []rune{0x1EE5},                         // Case map
	0x1EE6:  []rune{0x1EE7},                         // Case map
	0x1EE8:  []rune{0x1EE9},                         // Case map
	0x1EEA:  []rune{0x1EEB},                         // Case map
	0x1EEC:  []rune{0x1EED},                         // Case map
	0x1EEE:  []rune{0x1EEF},                         // Case map
	0x1EF0:  []rune{0x1EF1},                         // Case map
	0x1EF2:  []rune{0x1EF3},                         // Case map
	0x1EF4:  []rune{0x1EF5},                         // Case map
	0x1EF6:  []rune{0x1EF7},                         // Case map
	0x1EF8:  []rune{0x1EF9},                         // Case map
	0x1F08:  []rune{0x1F00},                         // Case map
	0x1F09:  []rune{0x1F01},                         // Case map
	0x1F0A:  []rune{0x1F02},                         // Case map
	0x1F0B:  []rune{0x1F03},                         // Case map
	0x1F0C:  []rune{0x1F04},                         // Case map
	0x1F0D:  []rune{0x1F05},                         // Case map
	0x1F0E:  []rune{0x1F06},                         // Case map
	0x1F0F:  []rune{0x1F07},                         // Case map
	0x1F18:  []rune{0x1F10},                         // Case map
	0x1F19:  []rune{0x1F11},                         // Case map
	0x1F1A:  []rune{0x1F12},                         // Case map
	0x1F1B:  []rune{0x1F13},                         // Case map
	0x1F1C:  []rune{0x1F14},                         // Case map
	0x1F1D:  []rune{0x1F15},                         // Case map
	0x1F28:  []rune{0x1F20},                         // Case map
	0x1F29:  []rune{0x1F21},                         // Case map
	0x1F2A:  []rune{0x1F22},                         // Case map
	0x1F2B:  []rune{0x1F23},                         // Case map
	0x1F2C:  []rune{0x1F24},                         // Case map
	0x1F2D:  []rune{0x1F25},                         // Case map
	0x1F2E:  []rune{0x1F26},                         // Case map
	0x1F2F:  []rune{0x1F27},                         // Case map
	0x1F38:  []rune{0x1F30},                         // Case map
	0x1F39:  []rune{0x1F31},                         // Case map
	0x1F3A:  []rune{0x1F32},                         // Case map
	0x1F3B:  []rune{0x1F33},                         // Case map
	0x1F3C:  []rune{0x1F34},                         // Case map
	0x1F3D:  []rune{0x1F35},                         // Case map
	0x1F3E:  []rune{0x1F36},                         // Case map
	0x1F3F:  []rune{0x1F37},                         // Case map
	0x1F48:  []rune{0x1F40},                         // Case map
	0x1F49:  []rune{0x1F41},                         // Case map
	0x1F4A:  []rune{0x1F42},                         // Case map
	0x1F4B:  []rune{0x1F43},                         // Case map
	0x1F4C:  []rune{0x1F44},                         // Case map
	0x1F4D:  []rune{0x1F45},                         // Case map
	0x1F50:  []rune{0x03C5, 0x0313},                 // Case map
	0x1F52:  []rune{0x03C5, 0x0313, 0x0300},         // Case map
	0x1F54:  []rune{0x03C5, 0x0313, 0x0301},         // Case map
	0x1F56:  []rune{0x03C5, 0x0313, 0x0342},         // Case map
	0x1F59:  []rune{0x1F51},                         // Case map
	0x1F5B:  []rune{0x1F53},                         // Case map
	0x1F5D:  []rune{0x1F55},                         // Case map
	0x1F5F:  []rune{0x1F57},                         // Case map
	0x1F68:  []rune{0x1F60},                         // Case map
	0x1F69:  []rune{0x1F61},                         // Case map
	0x1F6A:  []rune{0x1F62},                         // Case map
	0x1F6B:  []rune{0x1F63},                         // Case map
	0x1F6C:  []rune{0x1F64},                         // Case map
	0x1F6D:  []rune{0x1F65},                         // Case map
	0x1F6E:  []rune{0x1F66},                         // Case map
	0x1F6F:  []rune{0x1F67},                         // Case map
	0x1F80:  []rune{0x1F00, 0x03B9},                 // Case map
	0x1F81:  []rune{0x1F01, 0x03B9},                 // Case map
	0x1F82:  []rune{0x1F02, 0x03B9},                 // Case map
	0x1F83:  []rune{0x1F03, 0x03B9},                 // Case map
	0x1F84:  []rune{0x1F04, 0x03B9},                 // Case map
	0x1F85:  []rune{0x1F05, 0x03B9},                 // Case map
	0x1F86:  []rune{0x1F06, 0x03B9},                 // Case map
	0x1F87:  []rune{0x1F07, 0x03B9},                 // Case map
	0x1F88:  []rune{0x1F00, 0x03B9},                 // Case map
	0x1F89:  []rune{0x1F01, 0x03B9},                 // Case map
	0x1F8A:  []rune{0x1F02, 0x03B9},                 // Case map
	0x1F8B:  []rune{0x1F03, 0x03B9},                 // Case map
	0x1F8C:  []rune{0x1F04, 0x03B9},                 // Case map
	0x1F8D:  []rune{0x1F05, 0x03B9},                 // Case map
	0x1F8E:  []rune{0x1F06, 0x03B9},                 // Case map
	0x1F8F:  []rune{0x1F07, 0x03B9},                 // Case map
	0x1F90:  []rune{0x1F20, 0x03B9},                 // Case map
	0x1F91:  []rune{0x1F21, 0x03B9},                 // Case map
	0x1F92:  []rune{0x1F22, 0x03B9},                 // Case map
	0x1F93:  []rune{0x1F23, 0x03B9},                 // Case map
	0x1F94:  []rune{0x1F24, 0x03B9},                 // Case map
	0x1F95:  []rune{0x1F25, 0x03B9},                 // Case map
	0x1F96:  []rune{0x1F26, 0x03B9},                 // Case map
	0x1F97:  []rune{0x1F27, 0x03B9},                 // Case map
	0x1F98:  []rune{0x1F20, 0x03B9},                 // Case map
	0x1F99:  []rune{0x1F21, 0x03B9},                 // Case map
	0x1F9A:  []rune{0x1F22, 0x03B9},                 // Case map
	0x1F9B:  []rune{0x1F23, 0x03B9},                 // Case map
	0x1F9C:  []rune{0x1F24, 0x03B9},                 // Case map
	0x1F9D:  []rune{0x1F25, 0x03B9},                 // Case map
	0x1F9E:  []rune{0x1F26, 0x03B9},                 // Case map
	0x1F9F:  []rune{0x1F27, 0x03B9},                 // Case map
	0x1FA0:  []rune{0x1F60, 0x03B9},                 // Case map
	0x1FA1:  []rune{0x1F61, 0x03B9},                 // Case map
	0x1FA2:  []rune{0x1F62, 0x03B9},                 // Case map
	0x1FA3:  []rune{0x1F63, 0x03B9},                 // Case map
	0x1FA4:  []rune{0x1F64, 0x03B9},                 // Case map
	0x1FA5:  []rune{0x1F65, 0x03B9},                 // Case map
	0x1FA6:  []rune{0x1F66, 0x03B9},                 // Case map
	0x1FA7:  []rune{0x1F67, 0x03B9},                 // Case map
	0x1FA8:  []rune{0x1F60, 0x03B9},                 // Case map
	0x1FA9:  []rune{0x1F61, 0x03B9},                 // Case map
	0x1FAA:  []rune{0x1F62, 0x03B9},                 // Case map
	0x1FAB:  []rune{0x1F63, 0x03B9},                 // Case map
	0x1FAC:  []rune{0x1F64, 0x03B9},                 // Case map
	0x1FAD:  []rune{0x1F65, 0x03B9},                 // Case map
	0x1FAE:  []rune{0x1F66, 0x03B9},                 // Case map
	0x1FAF:  []rune{0x1F67, 0x03B9},                 // Case map
	0x1FB2:  []rune{0x1F70, 0x03B9},                 // Case map
	0x1FB3:  []rune{0x03B1, 0x03B9},                 // Case map
	0x1FB4:  []rune{0x03AC, 0x03B9},                 // Case map
	0x1FB6:  []rune{0x03B1, 0x0342},                 // Case map
	0x1FB7:  []rune{0x03B1, 0x0342, 0x03B9},         // Case map
	0x1FB8:  []rune{0x1FB0},                         // Case map
	0x1FB9:  []rune{0x1FB1},                         // Case map
	0x1FBA:  []rune{0x1F70},                         // Case map
	0x1FBB:  []rune{0x1F71},                         // Case map
	0x1FBC:  []rune{0x03B1, 0x03B9},                 // Case map
	0x1FBE:  []rune{0x03B9},                         // Case map
	0x1FC2:  []rune{0x1F74, 0x03B9},                 // Case map
	0x1FC3:  []rune{0x03B7, 0x03B9},                 // Case map
	0x1FC4:  []rune{0x03AE, 0x03B9},                 // Case map
	0x1FC6:  []rune{0x03B7, 0x0342},                 // Case map
	0x1FC7:  []rune{0x03B7, 0x0342, 0x03B9},         // Case map
	0x1FC8:  []rune{0x1F72},                         // Case map
	0x1FC9:  []rune{0x1F73},                         // Case map
	0x1FCA:  []rune{0x1F74},                         // Case map
	0x1FCB:  []rune{0x1F75},                         // Case map
	0x1FCC:  []rune{0x03B7, 0x03B9},                 // Case map
	0x1FD2:  []rune{0x03B9, 0x0308, 0x0300},         // Case map
	0x1FD3:  []rune{0x03B9, 0x0308, 0x0301},         // Case map
	0x1FD6:  []rune{0x03B9, 0x0342},                 // Case map
	0x1FD7:  []rune{0x03B9, 0x0308, 0x0342},         // Case map
	0x1FD8:  []rune{0x1FD0},                         // Case map
	0x1FD9:  []rune{0x1FD1},                         // Case map
	0x1FDA:  []rune{0x1F76},                         // Case map
	0x1FDB:  []rune{0x1F77},                         // Case map
	0x1FE2:  []rune{0x03C5, 0x0308, 0x0300},         // Case map
	0x1FE3:  []rune{0x03C5, 0x0308, 0x0301},         // Case map
	0x1FE4:  []rune{0x03C1, 0x0313},                 // Case map
	0x1FE6:  []rune{0x03C5, 0x0342},                 // Case map
	0x1FE7:  []rune{0x03C5, 0x0308, 0x0342},         // Case map
	0x1FE8:  []rune{0x1FE0},                         // Case map
	0x1FE9:  []rune{0x1FE1},                         // Case map
	0x1FEA:  []rune{0x1F7A},                         // Case map
	0x1FEB:  []rune{0x1F7B},                         // Case map
	0x1FEC:  []rune{0x1FE5},                         // Case map
	0x1FF2:  []rune{0x1F7C, 0x03B9},                 // Case map
	0x1FF3:  []rune{0x03C9, 0x03B9},                 // Case map
	0x1FF4:  []rune{0x03CE, 0x03B9},                 // Case map
	0x1FF6:  []rune{0x03C9, 0x0342},                 // Case map
	0x1FF7:  []rune{0x03C9, 0x0342, 0x03B9},         // Case map
	0x1FF8:  []rune{0x1F78},                         // Case map
	0x1FF9:  []rune{0x1F79},                         // Case map
	0x1FFA:  []rune{0x1F7C},                         // Case map
	0x1FFB:  []rune{0x1F7D},                         // Case map
	0x1FFC:  []rune{0x03C9, 0x03B9},                 // Case map
	0x20A8:  []rune{0x0072, 0x0073},                 // Additional folding
	0x2102:  []rune{0x0063},                         // Additional folding
	0x2103:  []rune{0x00B0, 0x0063},                 // Additional folding
	0x2107:  []rune{0x025B},                         // Additional folding
	0x2109:  []rune{0x00B0, 0x0066},                 // Additional folding
	0x210B:  []rune{0x0068},                         // Additional folding
	0x210C:  []rune{0x0068},                         // Additional folding
	0x210D:  []rune{0x0068},                         // Additional folding
	0x2110:  []rune{0x0069},                         // Additional folding
	0x2111:  []rune{0x0069},                         // Additional folding
	0x2112:  []rune{0x006C},                         // Additional folding
	0x2115:  []rune{0x006E},                         // Additional folding
	0x2116:  []rune{0x006E, 0x006F},                 // Additional folding
	0x2119:  []rune{0x0070},                         // Additional folding
	0x211A:  []rune{0x0071},                         // Additional folding
	0x211B:  []rune{0x0072},                         // Additional folding
	0x211C:  []rune{0x0072},                         // Additional folding
	0x211D:  []rune{0x0072},                         // Additional folding
	0x2120:  []rune{0x0073, 0x006D},                 // Additional folding
	0x2121:  []rune{0x0074, 0x0065, 0x006C},         // Additional folding
	0x2122:  []rune{0x0074, 0x006D},                 // Additional folding
	0x2124:  []rune{0x007A},                         // Additional folding
	0x2126:  []rune{0x03C9},                         // Case map
	0x2128:  []rune{0x007A},                         // Additional folding
	0x212A:  []rune{0x006B},                         // Case map
	0x212B:  []rune{0x00E5},                         // Case map
	0x212C:  []rune{0x0062},                         // Additional folding
	0x212D:  []rune{0x0063},                         // Additional folding
	0x2130:  []rune{0x0065},                         // Additional folding
	0x2131:  []rune{0x0066},                         // Additional folding
	0x2133:  []rune{0x006D},                         // Additional folding
	0x213E:  []rune{0x03B3},                         // Additional folding
	0x213F:  []rune{0x03C0},                         // Additional folding
	0x2145:  []rune{0x0064},                         // Additional folding
	0x2160:  []rune{0x2170},                         // Case map
	0x2161:  []rune{0x2171},                         // Case map
	0x2162:  []rune{0x2172},                         // Case map
	0x2163:  []rune{0x2173},                         // Case map
	0x2164:  []rune{0x2174},                         // Case map
	0x2165:  []rune{0x2175},                         // Case map
	0x2166:  []rune{0x2176},                         // Case map
	0x2167:  []rune{0x2177},                         // Case map
	0x2168:  []rune{0x2178},                         // Case map
	0x2169:  []rune{0x2179},                         // Case map
	0x216A:  []rune{0x217A},                         // Case map
	0x216B:  []rune{0x217B},                         // Case map
	0x216C:  []rune{0x217C},                         // Case map
	0x216D:  []rune{0x217D},                         // Case map
	0x216E:  []rune{0x217E},                         // Case map
	0x216F:  []rune{0x217F},                         // Case map
	0x24B6:  []rune{0x24D0},                         // Case map
	0x24B7:  []rune{0x24D1},                         // Case map
	0x24B8:  []rune{0x24D2},                         // Case map
	0x24B9:  []rune{0x24D3},                         // Case map
	0x24BA:  []rune{0x24D4},                         // Case map
	0x24BB:  []rune{0x24D5},                         // Case map
	0x24BC:  []rune{0x24D6},                         // Case map
	0x24BD:  []rune{0x24D7},                         // Case map
	0x24BE:  []rune{0x24D8},                         // Case map
	0x24BF:  []rune{0x24D9},                         // Case map
	0x24C0:  []rune{0x24DA},                         // Case map
	0x24C1:  []rune{0x24DB},                         // Case map
	0x24C2:  []rune{0x24DC},                         // Case map
	0x24C3:  []rune{0x24DD},                         // Case map
	0x24C4:  []rune{0x24DE},                         // Case map
	0x24C5:  []rune{0x24DF},                         // Case map
	0x24C6:  []rune{0x24E0},                         // Case map
	0x24C7:  []rune{0x24E1},                         // Case map
	0x24C8:  []rune{0x24E2},                         // Case map
	0x24C9:  []rune{0x24E3},                         // Case map
	0x24CA:  []rune{0x24E4},                         // Case map
	0x24CB:  []rune{0x24E5},                         // Case map
	0x24CC:  []rune{0x24E6},                         // Case map
	0x24CD:  []rune{0x24E7},                         // Case map
	0x24CE:  []rune{0x24E8},                         // Case map
	0x24CF:  []rune{0x24E9},                         // Case map
	0x3371:  []rune{0x0068, 0x0070, 0x0061},         // Additional folding
	0x3373:  []rune{0x0061, 0x0075},                 // Additional folding
	0x3375:  []rune{0x006F, 0x0076},                 // Additional folding
	0x3380:  []rune{0x0070, 0x0061},                 // Additional folding
	0x3381:  []rune{0x006E, 0x0061},                 // Additional folding
	0x3382:  []rune{0x03BC, 0x0061},                 // Additional folding
	0x3383:  []rune{0x006D, 0x0061},                 // Additional folding
	0x3384:  []rune{0x006B, 0x0061},                 // Additional folding
	0x3385:  []rune{0x006B, 0x0062},                 // Additional folding
	0x3386:  []rune{0x006D, 0x0062},                 // Additional folding
	0x3387:  []rune{0x0067, 0x0062},                 // Additional folding
	0x338A:  []rune{0x0070, 0x0066},                 // Additional folding
	0x338B:  []rune{0x006E, 0x0066},                 // Additional folding
	0x338C:  []rune{0x03BC, 0x0066},                 // Additional folding
	0x3390:  []rune{0x0068, 0x007A},                 // Additional folding
	0x3391:  []rune{0x006B, 0x0068, 0x007A},         // Additional folding
	0x3392:  []rune{0x006D, 0x0068, 0x007A},         // Additional folding
	0x3393:  []rune{0x0067, 0x0068, 0x007A},         // Additional folding
	0x3394:  []rune{0x0074, 0x0068, 0x007A},         // Additional folding
	0x33A9:  []rune{0x0070, 0x0061},                 // Additional folding
	0x33AA:  []rune{0x006B, 0x0070, 0x0061},         // Additional folding
	0x33AB:  []rune{0x006D, 0x0070, 0x0061},         // Additional folding
	0x33AC:  []rune{0x0067, 0x0070, 0x0061},         // Additional folding
	0x33B4:  []rune{0x0070, 0x0076},                 // Additional folding
	0x33B5:  []rune{0x006E, 0x0076},                 // Additional folding
	0x33B6:  []rune{0x03BC, 0x0076},                 // Additional folding
	0x33B7:  []rune{0x006D, 0x0076},                 // Additional folding
	0x33B8:  []rune{0x006B, 0x0076},                 // Additional folding
	0x33B9:  []rune{0x006D, 0x0076},                 // Additional folding
	0x33BA:  []rune{0x0070, 0x0077},                 // Additional folding
	0x33BB:  []rune{0x006E, 0x0077},                 // Additional folding
	0x33BC:  []rune{0x03BC, 0x0077},                 // Additional folding
	0x33BD:  []rune{0x006D, 0x0077},                 // Additional folding
	0x33BE:  []rune{0x006B, 0x0077},                 // Additional folding
	0x33BF:  []rune{0x006D, 0x0077},                 // Additional folding
	0x33C0:  []rune{0x006B, 0x03C9},                 // Additional folding
	0x33C1:  []rune{0x006D, 0x03C9},                 // Additional folding
	0x33C3:  []rune{0x0062, 0x0071},                 // Additional folding
	0x33C6:  []rune{0x0063, 0x2215, 0x006B, 0x0067}, // Additional folding
	0x33C7:  []rune{0x0063, 0x006F, 0x002E},         // Additional folding
	0x33C8:  []rune{0x0064, 0x0062},                 // Additional folding
	0x33C9:  []rune{0x0067, 0x0079},                 // Additional folding
	0x33CB:  []rune{0x0068, 0x0070},                 // Additional folding
	0x33CD:  []rune{0x006B, 0x006B},                 // Additional folding
	0x33CE:  []rune{0x006B, 0x006D},                 // Additional folding
	0x33D7:  []rune{0x0070, 0x0068},                 // Additional folding
	0x33D9:  []rune{0x0070, 0x0070, 0x006D},         // Additional folding
	0x33DA:  []rune{0x0070, 0x0072},                 // Additional folding
	0x33DC:  []rune{0x0073, 0x0076},                 // Additional folding
	0x33DD:  []rune{0x0077, 0x0062},                 // Additional folding
	0xFB00:  []rune{0x0066, 0x0066},                 // Case map
	0xFB01:  []rune{0x0066, 0x0069},                 // Case map
	0xFB02:  []rune{0x0066, 0x006C},                 // Case map
	0xFB03:  []rune{0x0066, 0x0066, 0x0069},         // Case map
	0xFB04:  []rune{0x0066, 0x0066, 0x006C},         // Case map
	0xFB05:  []rune{0x0073, 0x0074},                 // Case map
	0xFB06:  []rune{0x0073, 0x0074},                 // Case map
	0xFB13:  []rune{0x0574, 0x0576},                 // Case map
	0xFB14:  []rune{0x0574, 0x0565},                 // Case map
	0xFB15:  []rune{0x0574, 0x056B},                 // Case map
	0xFB16:  []rune{0x057E, 0x0576},                 // Case map
	0xFB17:  []rune{0x0574, 0x056D},                 // Case map
	0xFF21:  []rune{0xFF41},                         // Case map
	0xFF22:  []rune{0xFF42},                         // Case map
	0xFF23:  []rune{0xFF43},                         // Case map
	0xFF24:  []rune{0xFF44},                         // Case map
	0xFF25:  []rune{0xFF45},                         // Case map
	0xFF26:  []rune{0xFF46},                         // Case map
	0xFF27:  []rune{0xFF47},                         // Case map
	0xFF28:  []rune{0xFF48},                         // Case map
	0xFF29:  []rune{0xFF49},                         // Case map
	0xFF2A:  []rune{0xFF4A},                         // Case map
	0xFF2B:  []rune{0xFF4B},                         // Case map
	0xFF2C:  []rune{0xFF4C},                         // Case map
	0xFF2D:  []rune{0xFF4D},                         // Case map
	0xFF2E:  []rune{0xFF4E},                         // Case map
	0xFF2F:  []rune{0xFF4F},                         // Case map
	0xFF30:  []rune{0xFF50},                         // Case map
	0xFF31:  []rune{0xFF51},                         // Case map
	0xFF32:  []rune{0xFF52},                         // Case map
	0xFF33:  []rune{0xFF53},                         // Case map
	0xFF34:  []rune{0xFF54},                         // Case map
	0xFF35:  []rune{0xFF55},                         // Case map
	0xFF36:  []rune{0xFF56},                         // Case map
	0xFF37:  []rune{0xFF57},                         // Case map
	0xFF38:  []rune{0xFF58},                         // Case map
	0xFF39:  []rune{0xFF59},                         // Case map
	0xFF3A:  []rune{0xFF5A},                         // Case map
	0x10400: []rune{0x10428},                        // Case map
	0x10401: []rune{0x10429},                        // Case map
	0x10402: []rune{0x1042A},                        // Case map
	0x10403: []rune{0x1042B},                        // Case map
	0x10404: []rune{0x1042C},                        // Case map
	0x10405: []rune{0x1042D},                        // Case map
	0x10406: []rune{0x1042E},                        // Case map
	0x10407: []rune{0x1042F},                        // Case map
	0x10408: []rune{0x10430},                        // Case map
	0x10409: []rune{0x10431},                        // Case map
	0x1040A: []rune{0x10432},                        // Case map
	0x1040B: []rune{0x10433},                        // Case map
	0x1040C: []rune{0x10434},                        // Case map
	0x1040D: []rune{0x10435},                        // Case map
	0x1040E: []rune{0x10436},                        // Case map
	0x1040F: []rune{0x10437},                        // Case map
	0x10410: []rune{0x10438},                        // Case map
	0x10411: []rune{0x10439},                        // Case map
	0x10412: []rune{0x1043A},                        // Case map
	0x10413: []rune{0x1043B},                        // Case map
	0x10414: []rune{0x1043C},                        // Case map
	0x10415: []rune{0x1043D},                        // Case map
	0x10416: []rune{0x1043E},                        // Case map
	0x10417: []rune{0x1043F},                        // Case map
	0x10418: []rune{0x10440},                        // Case map
	0x10419: []rune{0x10441},                        // Case map
	0x1041A: []rune{0x10442},                        // Case map
	0x1041B: []rune{0x10443},                        // Case map
	0x1041C: []rune{0x10444},                        // Case map
	0x1041D: []rune{0x10445},                        // Case map
	0x1041E: []rune{0x10446},                        // Case map
	0x1041F: []rune{0x10447},                        // Case map
	0x10420: []rune{0x10448},                        // Case map
	0x10421: []rune{0x10449},                        // Case map
	0x10422: []rune{0x1044A},                        // Case map
	0x10423: []rune{0x1044B},                        // Case map
	0x10424: []rune{0x1044C},                        // Case map
	0x10425: []rune{0x1044D},                        // Case map
	0x1D400: []rune{0x0061},                         // Additional folding
	0x1D401: []rune{0x0062},                         // Additional folding
	0x1D402: []rune{0x0063},                         // Additional folding
	0x1D403: []rune{0x0064},                         // Additional folding
	0x1D404: []rune{0x0065},                         // Additional folding
	0x1D405: []rune{0x0066},                         // Additional folding
	0x1D406: []rune{0x0067},                         // Additional folding
	0x1D407: []rune{0x0068},                         // Additional folding
	0x1D408: []rune{0x0069},                         // Additional folding
	0x1D409: []rune{0x006A},                         // Additional folding
	0x1D40A: []rune{0x006B},                         // Additional folding
	0x1D40B: []rune{0x006C},                         // Additional folding
	0x1D40C: []rune{0x006D},                         // Additional folding
	0x1D40D: []rune{0x006E},                         // Additional folding
	0x1D40E: []rune{0x006F},                         // Additional folding
	0x1D40F: []rune{0x0070},                         // Additional folding
	0x1D410: []rune{0x0071},                         // Additional folding
	0x1D411: []rune{0x0072},                         // Additional folding
	0x1D412: []rune{0x0073},                         // Additional folding
	0x1D413: []rune{0x0074},                         // Additional folding
	0x1D414: []rune{0x0075},                         // Additional folding
	0x1D415: []rune{0x0076},                         // Additional folding
	0x1D416: []rune{0x0077},                         // Additional folding
	0x1D417: []rune{0x0078},                         // Additional folding
	0x1D418: []rune{0x0079},                         // Additional folding
	0x1D419: []rune{0x007A},                         // Additional folding
	0x1D434: []rune{0x0061},                         // Additional folding
	0x1D435: []rune{0x0062},                         // Additional folding
	0x1D436: []rune{0x0063},                         // Additional folding
	0x1D437: []rune{0x0064},                         // Additional folding
	0x1D438: []rune{0x0065},                         // Additional folding
	0x1D439: []rune{0x0066},                         // Additional folding
	0x1D43A: []rune{0x0067},                         // Additional folding
	0x1D43B: []rune{0x0068},                         // Additional folding
	0x1D43C: []rune{0x0069},                         // Additional folding
	0x1D43D: []rune{0x006A},                         // Additional folding
	0x1D43E: []rune{0x006B},                         // Additional folding
	0x1D43F: []rune{0x006C},                         // Additional folding
	0x1D440: []rune{0x006D},                         // Additional folding
	0x1D441: []rune{0x006E},                         // Additional folding
	0x1D442: []rune{0x006F},                         // Additional folding
	0x1D443: []rune{0x0070},                         // Additional folding
	0x1D444: []rune{0x0071},                         // Additional folding
	0x1D445: []rune{0x0072},                         // Additional folding
	0x1D446: []rune{0x0073},                         // Additional folding
	0x1D447: []rune{0x0074},                         // Additional folding
	0x1D448: []rune{0x0075},                         // Additional folding
	0x1D449: []rune{0x0076},                         // Additional folding
	0x1D44A: []rune{0x0077},                         // Additional folding
	0x1D44B: []rune{0x0078},                         // Additional folding
	0x1D44C: []rune{0x0079},                         // Additional folding
	0x1D44D: []rune{0x007A},                         // Additional folding
	0x1D468: []rune{0x0061},                         // Additional folding
	0x1D469: []rune{0x0062},                         // Additional folding
	0x1D46A: []rune{0x0063},                         // Additional folding
	0x1D46B: []rune{0x0064},                         // Additional folding
	0x1D46C: []rune{0x0065},                         // Additional folding
	0x1D46D: []rune{0x0066},                         // Additional folding
	0x1D46E: []rune{0x0067},                         // Additional folding
	0x1D46F: []rune{0x0068},                         // Additional folding
	0x1D470: []rune{0x0069},                         // Additional folding
	0x1D471: []rune{0x006A},                         // Additional folding
	0x1D472: []rune{0x006B},                         // Additional folding
	0x1D473: []rune{0x006C},                         // Additional folding
	0x1D474: []rune{0x006D},                         // Additional folding
	0x1D475: []rune{0x006E},                         // Additional folding
	0x1D476: []rune{0x006F},                         // Additional folding
	0x1D477: []rune{0x0070},                         // Additional folding
	0x1D478: []rune{0x0071},                         // Additional folding
	0x1D479: []rune{0x0072},                         // Additional folding
	0x1D47A: []rune{0x0073},                         // Additional folding
	0x1D47B: []rune{0x0074},                         // Additional folding
	0x1D47C: []rune{0x0075},                         // Additional folding
	0x1D47D: []rune{0x0076},                         // Additional folding
	0x1D47E: []rune{0x0077},                         // Additional folding
	0x1D47F: []rune{0x0078},                         // Additional folding
	0x1D480: []rune{0x0079},                         // Additional folding
	0x1D481: []rune{0x007A},                         // Additional folding
	0x1D49C: []rune{0x0061},                         // Additional folding
	0x1D49E: []rune{0x0063},                         // Additional folding
	0x1D49F: []rune{0x0064},                         // Additional folding
	0x1D4A2: []rune{0x0067},                         // Additional folding
	0x1D4A5: []rune{0x006A},                         // Additional folding
	0x1D4A6: []rune{0x006B},                         // Additional folding
	0x1D4A9: []rune{0x006E},                         // Additional folding
	0x1D4AA: []rune{0x006F},                         // Additional folding
	0x1D4AB: []rune{0x0070},                         // Additional folding
	0x1D4AC: []rune{0x0071},                         // Additional folding
	0x1D4AE: []rune{0x0073},                         // Additional folding
	0x1D4AF: []rune{0x0074},                         // Additional folding
	0x1D4B0: []rune{0x0075},                         // Additional folding
	0x1D4B1: []rune{0x0076},                         // Additional folding
	0x1D4B2: []rune{0x0077},                         // Additional folding
	0x1D4B3: []rune{0x0078},                         // Additional folding
	0x1D4B4: []rune{0x0079},                         // Additional folding
	0x1D4B5: []rune{0x007A},                         // Additional folding
	0x1D4D0: []rune{0x0061},                         // Additional folding
	0x1D4D1: []rune{0x0062},                         // Additional folding
	0x1D4D2: []rune{0x0063},                         // Additional folding
	0x1D4D3: []rune{0x0064},                         // Additional folding
	0x1D4D4: []rune{0x0065},                         // Additional folding
	0x1D4D5: []rune{0x0066},                         // Additional folding
	0x1D4D6: []rune{0x0067},                         // Additional folding
	0x1D4D7: []rune{0x0068},                         // Additional folding
	0x1D4D8: []rune{0x0069},                         // Additional folding
	0x1D4D9: []rune{0x006A},                         // Additional folding
	0x1D4DA: []rune{0x006B},                         // Additional folding
	0x1D4DB: []rune{0x006C},                         // Additional folding
	0x1D4DC: []rune{0x006D},                         // Additional folding
	0x1D4DD: []rune{0x006E},                         // Additional folding
	0x1D4DE: []rune{0x006F},                         // Additional folding
	0x1D4DF: []rune{0x0070},                         // Additional folding
	0x1D4E0: []rune{0x0071},                         // Additional folding
	0x1D4E1: []rune{0x0072},                         // Additional folding
	0x1D4E2: []rune{0x0073},                         // Additional folding
	0x1D4E3: []rune{0x0074},                         // Additional folding
	0x1D4E4: []rune{0x0075},                         // Additional folding
	0x1D4E5: []rune{0x0076},                         // Additional folding
	0x1D4E6: []rune{0x0077},                         // Additional folding
	0x1D4E7: []rune{0x0078},                         // Additional folding
	0x1D4E8: []rune{0x0079},                         // Additional folding
	0x1D4E9: []rune{0x007A},                         // Additional folding
	0x1D504: []rune{0x0061},                         // Additional folding
	0x1D505: []rune{0x0062},                         // Additional folding
	0x1D507: []rune{0x0064},                         // Additional folding
	0x1D508: []rune{0x0065},                         // Additional folding
	0x1D509: []rune{0x0066},                         // Additional folding
	0x1D50A: []rune{0x0067},                         // Additional folding
	0x1D50D: []rune{0x006A},                         // Additional folding
	0x1D50E: []rune{0x006B},                         // Additional folding
	0x1D50F: []rune{0x006C},                         // Additional folding
	0x1D510: []rune{0x006D},                         // Additional folding
	0x1D511: []rune{0x006E},                         // Additional folding
	0x1D512: []rune{0x006F},                         // Additional folding
	0x1D513: []rune{0x0070},                         // Additional folding
	0x1D514: []rune{0x0071},                         // Additional folding
	0x1D516: []rune{0x0073},                         // Additional folding
	0x1D517: []rune{0x0074},                         // Additional folding
	0x1D518: []rune{0x0075},                         // Additional folding
	0x1D519: []rune{0x0076},                         // Additional folding
	0x1D51A: []rune{0x0077},                         // Additional folding
	0x1D51B: []rune{0x0078},                         // Additional folding
	0x1D51C: []rune{0x0079},                         // Additional folding
	0x1D538: []rune{0x0061},                         // Additional folding
	0x1D539: []rune{0x0062},                         // Additional folding
	0x1D53B: []rune{0x0064},                         // Additional folding
	0x1D53C: []rune{0x0065},                         // Additional folding
	0x1D53D: []rune{0x0066},                         // Additional folding
	0x1D53E: []rune{0x0067},                         // Additional folding
	0x1D540: []rune{0x0069},                         // Additional folding
	0x1D541: []rune{0x006A},                         // Additional folding
	0x1D542: []rune{0x006B},                         // Additional folding
	0x1D543: []rune{0x006C},                         // Additional folding
	0x1D544: []rune{0x006D},                         // Additional folding
	0x1D546: []rune{0x006F},                         // Additional folding
	0x1D54A: []rune{0x0073},                         // Additional folding
	0x1D54B: []rune{0x0074},                         // Additional folding
	0x1D54C: []rune{0x0075},                         // Additional folding
	0x1D54D: []rune{0x0076},                         // Additional folding
	0x1D54E: []rune{0x0077},                         // Additional folding
	0x1D54F: []rune{0x0078},                         // Additional folding
	0x1D550: []rune{0x0079},                         // Additional folding
	0x1D56C: []rune{0x0061},                         // Additional folding
	0x1D56D: []rune{0x0062},                         // Additional folding
	0x1D56E: []rune{0x0063},                         // Additional folding
	0x1D56F: []rune{0x0064},                         // Additional folding
	0x1D570: []rune{0x0065},                         // Additional folding
	0x1D571: []rune{0x0066},                         // Additional folding
	0x1D572: []rune{0x0067},                         // Additional folding
	0x1D573: []rune{0x0068},                         // Additional folding
	0x1D574: []rune{0x0069},                         // Additional folding
	0x1D575: []rune{0x006A},                         // Additional folding
	0x1D576: []rune{0x006B},                         // Additional folding
	0x1D577: []rune{0x006C},                         // Additional folding
	0x1D578: []rune{0x006D},                         // Additional folding
	0x1D579: []rune{0x006E},                         // Additional folding
	0x1D57A: []rune{0x006F},                         // Additional folding
	0x1D57B: []rune{0x0070},                         // Additional folding
	0x1D57C: []rune{0x0071},                         // Additional folding
	0x1D57D: []rune{0x0072},                         // Additional folding
	0x1D57E: []rune{0x0073},                         // Additional folding
	0x1D57F: []rune{0x0074},                         // Additional folding
	0x1D580: []rune{0x0075},                         // Additional folding
	0x1D581: []rune{0x0076},                         // Additional folding
	0x1D582: []rune{0x0077},                         // Additional folding
	0x1D583: []rune{0x0078},                         // Additional folding
	0x1D584: []rune{0x0079},                         // Additional folding
	0x1D585: []rune{0x007A},                         // Additional folding
	0x1D5A0: []rune{0x0061},                         // Additional folding
	0x1D5A1: []rune{0x0062},                         // Additional folding
	0x1D5A2: []rune{0x0063},                         // Additional folding
	0x1D5A3: []rune{0x0064},                         // Additional folding
	0x1D5A4: []rune{0x0065},                         // Additional folding
	0x1D5A5: []rune{0x0066},                         // Additional folding
	0x1D5A6: []rune{0x0067},                         // Additional folding
	0x1D5A7: []rune{0x0068},                         // Additional folding
	0x1D5A8: []rune{0x0069},                         // Additional folding
	0x1D5A9: []rune{0x006A},                         // Additional folding
	0x1D5AA: []rune{0x006B},                         // Additional folding
	0x1D5AB: []rune{0x006C},                         // Additional folding
	0x1D5AC: []rune{0x006D},                         // Additional folding
	0x1D5AD: []rune{0x006E},                         // Additional folding
	0x1D5AE: []rune{0x006F},                         // Additional folding
	0x1D5AF: []rune{0x0070},                         // Additional folding
	0x1D5B0: []rune{0x0071},                         // Additional folding
	0x1D5B1: []rune{0x0072},                         // Additional folding
	0x1D5B2: []rune{0x0073},                         // Additional folding
	0x1D5B3: []rune{0x0074},                         // Additional folding
	0x1D5B4: []rune{0x0075},                         // Additional folding
	0x1D5B5: []rune{0x0076},                         // Additional folding
	0x1D5B6: []rune{0x0077},                         // Additional folding
	0x1D5B7: []rune{0x0078},                         // Additional folding
	0x1D5B8: []rune{0x0079},                         // Additional folding
	0x1D5B9: []rune{0x007A},                         // Additional folding
	0x1D5D4: []rune{0x0061},                         // Additional folding
	0x1D5D5: []rune{0x0062},                         // Additional folding
	0x1D5D6: []rune{0x0063},                         // Additional folding
	0x1D5D7: []rune{0x0064},                         // Additional folding
	0x1D5D8: []rune{0x0065},                         // Additional folding
	0x1D5D9: []rune{0x0066},                         // Additional folding
	0x1D5DA: []rune{0x0067},                         // Additional folding
	0x1D5DB: []rune{0x0068},                         // Additional folding
	0x1D5DC: []rune{0x0069},                         // Additional folding
	0x1D5DD: []rune{0x006A},                         // Additional folding
	0x1D5DE: []rune{0x006B},                         // Additional folding
	0x1D5DF: []rune{0x006C},                         // Additional folding
	0x1D5E0: []rune{0x006D},                         // Additional folding
	0x1D5E1: []rune{0x006E},                         // Additional folding
	0x1D5E2: []rune{0x006F},                         // Additional folding
	0x1D5E3: []rune{0x0070},                         // Additional folding
	0x1D5E4: []rune{0x0071},                         // Additional folding
	0x1D5E5: []rune{0x0072},                         // Additional folding
	0x1D5E6: []rune{0x0073},                         // Additional folding
	0x1D5E7: []rune{0x0074},                         // Additional folding
	0x1D5E8: []rune{0x0075},                         // Additional folding
	0x1D5E9: []rune{0x0076},                         // Additional folding
	0x1D5EA: []rune{0x0077},                         // Additional folding
	0x1D5EB: []rune{0x0078},                         // Additional folding
	0x1D5EC: []rune{0x0079},                         // Additional folding
	0x1D5ED: []rune{0x007A},                         // Additional folding
	0x1D608: []rune{0x0061},                         // Additional folding
	0x1D609: []rune{0x0062},                         // Additional folding
	0x1D60A: []rune{0x0063},                         // Additional folding
	0x1D60B: []rune{0x0064},                         // Additional folding
	0x1D60C: []rune{0x0065},                         // Additional folding
	0x1D60D: []rune{0x0066},                         // Additional folding
	0x1D60E: []rune{0x0067},                         // Additional folding
	0x1D60F: []rune{0x0068},                         // Additional folding
	0x1D610: []rune{0x0069},                         // Additional folding
	0x1D611: []rune{0x006A},                         // Additional folding
	0x1D612: []rune{0x006B},                         // Additional folding
	0x1D613: []rune{0x006C},                         // Additional folding
	0x1D614: []rune{0x006D},                         // Additional folding
	0x1D615: []rune{0x006E},                         // Additional folding
	0x1D616: []rune{0x006F},                         // Additional folding
	0x1D617: []rune{0x0070},                         // Additional folding
	0x1D618: []rune{0x0071},                         // Additional folding
	0x1D619: []rune{0x0072},                         // Additional folding
	0x1D61A: []rune{0x0073},                         // Additional folding
	0x1D61B: []rune{0x0074},                         // Additional folding
	0x1D61C: []rune{0x0075},                         // Additional folding
	0x1D61D: []rune{0x0076},                         // Additional folding
	0x1D61E: []rune{0x0077},                         // Additional folding
	0x1D61F: []rune{0x0078},                         // Additional folding
	0x1D620: []rune{0x0079},                         // Additional folding
	0x1D621: []rune{0x007A},                         // Additional folding
	0x1D63C: []rune{0x0061},                         // Additional folding
	0x1D63D: []rune{0x0062},                         // Additional folding
	0x1D63E: []rune{0x0063},                         // Additional folding
	0x1D63F: []rune{0x0064},                         // Additional folding
	0x1D640: []rune{0x0065},                         // Additional folding
	0x1D641: []rune{0x0066},                         // Additional folding
	0x1D642: []rune{0x0067},                         // Additional folding
	0x1D643: []rune{0x0068},                         // Additional folding
	0x1D644: []rune{0x0069},                         // Additional folding
	0x1D645: []rune{0x006A},                         // Additional folding
	0x1D646: []rune{0x006B},                         // Additional folding
	0x1D647: []rune{0x006C},                         // Additional folding
	0x1D648: []rune{0x006D},                         // Additional folding
	0x1D649: []rune{0x006E},                         // Additional folding
	0x1D64A: []rune{0x006F},                         // Additional folding
	0x1D64B: []rune{0x0070},                         // Additional folding
	0x1D64C: []rune{0x0071},                         // Additional folding
	0x1D64D: []rune{0x0072},                         // Additional folding
	0x1D64E: []rune{0x0073},                         // Additional folding
	0x1D64F: []rune{0x0074},                         // Additional folding
	0x1D650: []rune{0x0075},                         // Additional folding
	0x1D651: []rune{0x0076},                         // Additional folding
	0x1D652: []rune{0x0077},                         // Additional folding
	0x1D653: []rune{0x0078},                         // Additional folding
	0x1D654: []rune{0x0079},                         // Additional folding
	0x1D655: []rune{0x007A},                         // Additional folding
	0x1D670: []rune{0x0061},                         // Additional folding
	0x1D671: []rune{0x0062},                         // Additional folding
	0x1D672: []rune{0x0063},                         // Additional folding
	0x1D673: []rune{0x0064},                         // Additional folding
	0x1D674: []rune{0x0065},                         // Additional folding
	0x1D675: []rune{0x0066},                         // Additional folding
	0x1D676: []rune{0x0067},                         // Additional folding
	0x1D677: []rune{0x0068},                         // Additional folding
	0x1D678: []rune{0x0069},                         // Additional folding
	0x1D679: []rune{0x006A},                         // Additional folding
	0x1D67A: []rune{0x006B},                         // Additional folding
	0x1D67B: []rune{0x006C},                         // Additional folding
	0x1D67C: []rune{0x006D},                         // Additional folding
	0x1D67D: []rune{0x006E},                         // Additional folding
	0x1D67E: []rune{0x006F},                         // Additional folding
	0x1D67F: []rune{0x0070},                         // Additional folding
	0x1D680: []rune{0x0071},                         // Additional folding
	0x1D681: []rune{0x0072},                         // Additional folding
	0x1D682: []rune{0x0073},                         // Additional folding
	0x1D683: []rune{0x0074},                         // Additional folding
	0x1D684: []rune{0x0075},                         // Additional folding
	0x1D685: []rune{0x0076},                         // Additional folding
	0x1D686: []rune{0x0077},                         // Additional folding
	0x1D687: []rune{0x0078},                         // Additional folding
	0x1D688: []rune{0x0079},                         // Additional folding
	0x1D689: []rune{0x007A},                         // Additional folding
	0x1D6A8: []rune{0x03B1},                         // Additional folding
	0x1D6A9: []rune{0x03B2},                         // Additional folding
	0x1D6AA: []rune{0x03B3},                         // Additional folding
	0x1D6AB: []rune{0x03B4},                         // Additional folding
	0x1D6AC: []rune{0x03B5},                         // Additional folding
	0x1D6AD: []rune{0x03B6},                         // Additional folding
	0x1D6AE: []rune{0x03B7},                         // Additional folding
	0x1D6AF: []rune{0x03B8},                         // Additional folding
	0x1D6B0: []rune{0x03B9},                         // Additional folding
	0x1D6B1: []rune{0x03BA},                         // Additional folding
	0x1D6B2: []rune{0x03BB},                         // Additional folding
	0x1D6B3: []rune{0x03BC},                         // Additional folding
	0x1D6B4: []rune{0x03BD},                         // Additional folding
	0x1D6B5: []rune{0x03BE},                         // Additional folding
	0x1D6B6: []rune{0x03BF},                         // Additional folding
	0x1D6B7: []rune{0x03C0},                         // Additional folding
	0x1D6B8: []rune{0x03C1},                         // Additional folding
	0x1D6B9: []rune{0x03B8},                         // Additional folding
	0x1D6BA: []rune{0x03C3},                         // Additional folding
	0x1D6BB: []rune{0x03C4},                         // Additional folding
	0x1D6BC: []rune{0x03C5},                         // Additional folding
	0x1D6BD: []rune{0x03C6},                         // Additional folding
	0x1D6BE: []rune{0x03C7},                         // Additional folding
	0x1D6BF: []rune{0x03C8},                         // Additional folding
	0x1D6C0: []rune{0x03C9},                         // Additional folding
	0x1D6D3: []rune{0x03C3},                         // Additional folding
	0x1D6E2: []rune{0x03B1},                         // Additional folding
	0x1D6E3: []rune{0x03B2},                         // Additional folding
	0x1D6E4: []rune{0x03B3},                         // Additional folding
	0x1D6E5: []rune{0x03B4},                         // Additional folding
	0x1D6E6: []rune{0x03B5},                         // Additional folding
	0x1D6E7: []rune{0x03B6},                         // Additional folding
	0x1D6E8: []rune{0x03B7},                         // Additional folding
	0x1D6E9: []rune{0x03B8},                         // Additional folding
	0x1D6EA: []rune{0x03B9},                         // Additional folding
	0x1D6EB: []rune{0x03BA},                         // Additional folding
	0x1D6EC: []rune{0x03BB},                         // Additional folding
	0x1D6ED: []rune{0x03BC},                         // Additional folding
	0x1D6EE: []rune{0x03BD},                         // Additional folding
	0x1D6EF: []rune{0x03BE},                         // Additional folding
	0x1D6F0: []rune{0x03BF},                         // Additional folding
	0x1D6F1: []rune{0x03C0},                         // Additional folding
	0x1D6F2: []rune{0x03C1},                         // Additional folding
	0x1D6F3: []rune{0x03B8},                         // Additional folding
	0x1D6F4: []rune{0x03C3},                         // Additional folding
	0x1D6F5: []rune{0x03C4},                         // Additional folding
	0x1D6F6: []rune{0x03C5},                         // Additional folding
	0x1D6F7: []rune{0x03C6},                         // Additional folding
	0x1D6F8: []rune{0x03C7},                         // Additional folding
	0x1D6F9: []rune{0x03C8},                         // Additional folding
	0x1D6FA: []rune{0x03C9},                         // Additional folding
	0x1D70D: []rune{0x03C3},                         // Additional folding
	0x1D71C: []rune{0x03B1},                         // Additional folding
	0x1D71D: []rune{0x03B2},                         // Additional folding
	0x1D71E: []rune{0x03B3},                         // Additional folding
	0x1D71F: []rune{0x03B4},                         // Additional folding
	0x1D720: []rune{0x03B5},                         // Additional folding
	0x1D721: []rune{0x03B6},                         // Additional folding
	0x1D722: []rune{0x03B7},                         // Additional folding
	0x1D723: []rune{0x03B8},                         // Additional folding
	0x1D724: []rune{0x03B9},                         // Additional folding
	0x1D725: []rune{0x03BA},                         // Additional folding
	0x1D726: []rune{0x03BB},                         // Additional folding
	0x1D727: []rune{0x03BC},                         // Additional folding
	0x1D728: []rune{0x03BD},                         // Additional folding
	0x1D729: []rune{0x03BE},                         // Additional folding
	0x1D72A: []rune{0x03BF},                         // Additional folding
	0x1D72B: []rune{0x03C0},                         // Additional folding
	0x1D72C: []rune{0x03C1},                         // Additional folding
	0x1D72D: []rune{0x03B8},                         // Additional folding
	0x1D72E: []rune{0x03C3},                         // Additional folding
	0x1D72F: []rune{0x03C4},                         // Additional folding
	0x1D730: []rune{0x03C5},                         // Additional folding
	0x1D731: []rune{0x03C6},                         // Additional folding
	0x1D732: []rune{0x03C7},                         // Additional folding
	0x1D733: []rune{0x03C8},                         // Additional folding
	0x1D734: []rune{0x03C9},                         // Additional folding
	0x1D747: []rune{0x03C3},                         // Additional folding
	0x1D756: []rune{0x03B1},                         // Additional folding
	0x1D757: []rune{0x03B2},                         // Additional folding
	0x1D758: []rune{0x03B3},                         // Additional folding
	0x1D759: []rune{0x03B4},                         // Additional folding
	0x1D75A: []rune{0x03B5},                         // Additional folding
	0x1D75B: []rune{0x03B6},                         // Additional folding
	0x1D75C: []rune{0x03B7},                         // Additional folding
	0x1D75D: []rune{0x03B8},                         // Additional folding
	0x1D75E: []rune{0x03B9},                         // Additional folding
	0x1D75F: []rune{0x03BA},                         // Additional folding
	0x1D760: []rune{0x03BB},                         // Additional folding
	0x1D761: []rune{0x03BC},                         // Additional folding
	0x1D762: []rune{0x03BD},                         // Additional folding
	0x1D763: []rune{0x03BE},                         // Additional folding
	0x1D764: []rune{0x03BF},                         // Additional folding
	0x1D765: []rune{0x03C0},                         // Additional folding
	0x1D766: []rune{0x03C1},                         // Additional folding
	0x1D767: []rune{0x03B8},                         // Additional folding
	0x1D768: []rune{0x03C3},                         // Additional folding
	0x1D769: []rune{0x03C4},                         // Additional folding
	0x1D76A: []rune{0x03C5},                         // Additional folding
	0x1D76B: []rune{0x03C6},                         // Additional folding
	0x1D76C: []rune{0x03C7},                         // Additional folding
	0x1D76D: []rune{0x03C8},                         // Additional folding
	0x1D76E: []rune{0x03C9},                         // Additional folding
	0x1D781: []rune{0x03C3},                         // Additional folding
	0x1D790: []rune{0x03B1},                         // Additional folding
	0x1D791: []rune{0x03B2},                         // Additional folding
	0x1D792: []rune{0x03B3},                         // Additional folding
	0x1D793: []rune{0x03B4},                         // Additional folding
	0x1D794: []rune{0x03B5},                         // Additional folding
	0x1D795: []rune{0x03B6},                         // Additional folding
	0x1D796: []rune{0x03B7},                         // Additional folding
	0x1D797: []rune{0x03B8},                         // Additional folding
	0x1D798: []rune{0x03B9},                         // Additional folding
	0x1D799: []rune{0x03BA},                         // Additional folding
	0x1D79A: []rune{0x03BB},                         // Additional folding
	0x1D79B: []rune{0x03BC},                         // Additional folding
	0x1D79C: []rune{0x03BD},                         // Additional folding
	0x1D79D: []rune{0x03BE},                         // Additional folding
	0x1D79E: []rune{0x03BF},                         // Additional folding
	0x1D79F: []rune{0x03C0},                         // Additional folding
	0x1D7A0: []rune{0x03C1},                         // Additional folding
	0x1D7A1: []rune{0x03B8},                         // Additional folding
	0x1D7A2: []rune{0x03C3},                         // Additional folding
	0x1D7A3: []rune{0x03C4},                         // Additional folding
	0x1D7A4: []rune{0x03C5},                         // Additional folding
	0x1D7A5: []rune{0x03C6},                         // Additional folding
	0x1D7A6: []rune{0x03C7},                         // Additional folding
	0x1D7A7: []rune{0x03C8},                         // Additional folding
	0x1D7A8: []rune{0x03C9},                         // Additional folding
	0x1D7BB: []rune{0x03C3},                         // Additional folding
}

// TableB2 represents RFC-3454 Table B.2.
var TableB2 Mapping = tableB2

var tableB3 = Mapping{
	0x0041:  []rune{0x0061},                 // Case map
	0x0042:  []rune{0x0062},                 // Case map
	0x0043:  []rune{0x0063},                 // Case map
	0x0044:  []rune{0x0064},                 // Case map
	0x0045:  []rune{0x0065},                 // Case map
	0x0046:  []rune{0x0066},                 // Case map
	0x0047:  []rune{0x0067},                 // Case map
	0x0048:  []rune{0x0068},                 // Case map
	0x0049:  []rune{0x0069},                 // Case map
	0x004A:  []rune{0x006A},                 // Case map
	0x004B:  []rune{0x006B},                 // Case map
	0x004C:  []rune{0x006C},                 // Case map
	0x004D:  []rune{0x006D},                 // Case map
	0x004E:  []rune{0x006E},                 // Case map
	0x004F:  []rune{0x006F},                 // Case map
	0x0050:  []rune{0x0070},                 // Case map
	0x0051:  []rune{0x0071},                 // Case map
	0x0052:  []rune{0x0072},                 // Case map
	0x0053:  []rune{0x0073},                 // Case map
	0x0054:  []rune{0x0074},                 // Case map
	0x0055:  []rune{0x0075},                 // Case map
	0x0056:  []rune{0x0076},                 // Case map
	0x0057:  []rune{0x0077},                 // Case map
	0x0058:  []rune{0x0078},                 // Case map
	0x0059:  []rune{0x0079},                 // Case map
	0x005A:  []rune{0x007A},                 // Case map
	0x00B5:  []rune{0x03BC},                 // Case map
	0x00C0:  []rune{0x00E0},                 // Case map
	0x00C1:  []rune{0x00E1},                 // Case map
	0x00C2:  []rune{0x00E2},                 // Case map
	0x00C3:  []rune{0x00E3},                 // Case map
	0x00C4:  []rune{0x00E4},                 // Case map
	0x00C5:  []rune{0x00E5},                 // Case map
	0x00C6:  []rune{0x00E6},                 // Case map
	0x00C7:  []rune{0x00E7},                 // Case map
	0x00C8:  []rune{0x00E8},                 // Case map
	0x00C9:  []rune{0x00E9},                 // Case map
	0x00CA:  []rune{0x00EA},                 // Case map
	0x00CB:  []rune{0x00EB},                 // Case map
	0x00CC:  []rune{0x00EC},                 // Case map
	0x00CD:  []rune{0x00ED},                 // Case map
	0x00CE:  []rune{0x00EE},                 // Case map
	0x00CF:  []rune{0x00EF},                 // Case map
	0x00D0:  []rune{0x00F0},                 // Case map
	0x00D1:  []rune{0x00F1},                 // Case map
	0x00D2:  []rune{0x00F2},                 // Case map
	0x00D3:  []rune{0x00F3},                 // Case map
	0x00D4:  []rune{0x00F4},                 // Case map
	0x00D5:  []rune{0x00F5},                 // Case map
	0x00D6:  []rune{0x00F6},                 // Case map
	0x00D8:  []rune{0x00F8},                 // Case map
	0x00D9:  []rune{0x00F9},                 // Case map
	0x00DA:  []rune{0x00FA},                 // Case map
	0x00DB:  []rune{0x00FB},                 // Case map
	0x00DC:  []rune{0x00FC},                 // Case map
	0x00DD:  []rune{0x00FD},                 // Case map
	0x00DE:  []rune{0x00FE},                 // Case map
	0x00DF:  []rune{0x0073, 0x0073},         // Case map
	0x0100:  []rune{0x0101},                 // Case map
	0x0102:  []rune{0x0103},                 // Case map
	0x0104:  []rune{0x0105},                 // Case map
	0x0106:  []rune{0x0107},                 // Case map
	0x0108:  []rune{0x0109},                 // Case map
	0x010A:  []rune{0x010B},                 // Case map
	0x010C:  []rune{0x010D},                 // Case map
	0x010E:  []rune{0x010F},                 // Case map
	0x0110:  []rune{0x0111},                 // Case map
	0x0112:  []rune{0x0113},                 // Case map
	0x0114:  []rune{0x0115},                 // Case map
	0x0116:  []rune{0x0117},                 // Case map
	0x0118:  []rune{0x0119},                 // Case map
	0x011A:  []rune{0x011B},                 // Case map
	0x011C:  []rune{0x011D},                 // Case map
	0x011E:  []rune{0x011F},                 // Case map
	0x0120:  []rune{0x0121},                 // Case map
	0x0122:  []rune{0x0123},                 // Case map
	0x0124:  []rune{0x0125},                 // Case map
	0x0126:  []rune{0x0127},                 // Case map
	0x0128:  []rune{0x0129},                 // Case map
	0x012A:  []rune{0x012B},                 // Case map
	0x012C:  []rune{0x012D},                 // Case map
	0x012E:  []rune{0x012F},                 // Case map
	0x0130:  []rune{0x0069, 0x0307},         // Case map
	0x0132:  []rune{0x0133},                 // Case map
	0x0134:  []rune{0x0135},                 // Case map
	0x0136:  []rune{0x0137},                 // Case map
	0x0139:  []rune{0x013A},                 // Case map
	0x013B:  []rune{0x013C},                 // Case map
	0x013D:  []rune{0x013E},                 // Case map
	0x013F:  []rune{0x0140},                 // Case map
	0x0141:  []rune{0x0142},                 // Case map
	0x0143:  []rune{0x0144},                 // Case map
	0x0145:  []rune{0x0146},                 // Case map
	0x0147:  []rune{0x0148},                 // Case map
	0x0149:  []rune{0x02BC, 0x006E},         // Case map
	0x014A:  []rune{0x014B},                 // Case map
	0x014C:  []rune{0x014D},                 // Case map
	0x014E:  []rune{0x014F},                 // Case map
	0x0150:  []rune{0x0151},                 // Case map
	0x0152:  []rune{0x0153},                 // Case map
	0x0154:  []rune{0x0155},                 // Case map
	0x0156:  []rune{0x0157},                 // Case map
	0x0158:  []rune{0x0159},                 // Case map
	0x015A:  []rune{0x015B},                 // Case map
	0x015C:  []rune{0x015D},                 // Case map
	0x015E:  []rune{0x015F},                 // Case map
	0x0160:  []rune{0x0161},                 // Case map
	0x0162:  []rune{0x0163},                 // Case map
	0x0164:  []rune{0x0165},                 // Case map
	0x0166:  []rune{0x0167},                 // Case map
	0x0168:  []rune{0x0169},                 // Case map
	0x016A:  []rune{0x016B},                 // Case map
	0x016C:  []rune{0x016D},                 // Case map
	0x016E:  []rune{0x016F},                 // Case map
	0x0170:  []rune{0x0171},                 // Case map
	0x0172:  []rune{0x0173},                 // Case map
	0x0174:  []rune{0x0175},                 // Case map
	0x0176:  []rune{0x0177},                 // Case map
	0x0178:  []rune{0x00FF},                 // Case map
	0x0179:  []rune{0x017A},                 // Case map
	0x017B:  []rune{0x017C},                 // Case map
	0x017D:  []rune{0x017E},                 // Case map
	0x017F:  []rune{0x0073},                 // Case map
	0x0181:  []rune{0x0253},                 // Case map
	0x0182:  []rune{0x0183},                 // Case map
	0x0184:  []rune{0x0185},                 // Case map
	0x0186:  []rune{0x0254},                 // Case map
	0x0187:  []rune{0x0188},                 // Case map
	0x0189:  []rune{0x0256},                 // Case map
	0x018A:  []rune{0x0257},                 // Case map
	0x018B:  []rune{0x018C},                 // Case map
	0x018E:  []rune{0x01DD},                 // Case map
	0x018F:  []rune{0x0259},                 // Case map
	0x0190:  []rune{0x025B},                 // Case map
	0x0191:  []rune{0x0192},                 // Case map
	0x0193:  []rune{0x0260},                 // Case map
	0x0194:  []rune{0x0263},                 // Case map
	0x0196:  []rune{0x0269},                 // Case map
	0x0197:  []rune{0x0268},                 // Case map
	0x0198:  []rune{0x0199},                 // Case map
	0x019C:  []rune{0x026F},                 // Case map
	0x019D:  []rune{0x0272},                 // Case map
	0x019F:  []rune{0x0275},                 // Case map
	0x01A0:  []rune{0x01A1},                 // Case map
	0x01A2:  []rune{0x01A3},                 // Case map
	0x01A4:  []rune{0x01A5},                 // Case map
	0x01A6:  []rune{0x0280},                 // Case map
	0x01A7:  []rune{0x01A8},                 // Case map
	0x01A9:  []rune{0x0283},                 // Case map
	0x01AC:  []rune{0x01AD},                 // Case map
	0x01AE:  []rune{0x0288},                 // Case map
	0x01AF:  []rune{0x01B0},                 // Case map
	0x01B1:  []rune{0x028A},                 // Case map
	0x01B2:  []rune{0x028B},                 // Case map
	0x01B3:  []rune{0x01B4},                 // Case map
	0x01B5:  []rune{0x01B6},                 // Case map
	0x01B7:  []rune{0x0292},                 // Case map
	0x01B8:  []rune{0x01B9},                 // Case map
	0x01BC:  []rune{0x01BD},                 // Case map
	0x01C4:  []rune{0x01C6},                 // Case map
	0x01C5:  []rune{0x01C6},                 // Case map
	0x01C7:  []rune{0x01C9},                 // Case map
	0x01C8:  []rune{0x01C9},                 // Case map
	0x01CA:  []rune{0x01CC},                 // Case map
	0x01CB:  []rune{0x01CC},                 // Case map
	0x01CD:  []rune{0x01CE},                 // Case map
	0x01CF:  []rune{0x01D0},                 // Case map
	0x01D1:  []rune{0x01D2},                 // Case map
	0x01D3:  []rune{0x01D4},                 // Case map
	0x01D5:  []rune{0x01D6},                 // Case map
	0x01D7:  []rune{0x01D8},                 // Case map
	0x01D9:  []rune{0x01DA},                 // Case map
	0x01DB:  []rune{0x01DC},                 // Case map
	0x01DE:  []rune{0x01DF},                 // Case map
	0x01E0:  []rune{0x01E1},                 // Case map
	0x01E2:  []rune{0x01E3},                 // Case map
	0x01E4:  []rune{0x01E5},                 // Case map
	0x01E6:  []rune{0x01E7},                 // Case map
	0x01E8:  []rune{0x01E9},                 // Case map
	0x01EA:  []rune{0x01EB},                 // Case map
	0x01EC:  []rune{0x01ED},                 // Case map
	0x01EE:  []rune{0x01EF},                 // Case map
	0x01F0:  []rune{0x006A, 0x030C},         // Case map
	0x01F1:  []rune{0x01F3},                 // Case map
	0x01F2:  []rune{0x01F3},                 // Case map
	0x01F4:  []rune{0x01F5},                 // Case map
	0x01F6:  []rune{0x0195},                 // Case map
	0x01F7:  []rune{0x01BF},                 // Case map
	0x01F8:  []rune{0x01F9},                 // Case map
	0x01FA:  []rune{0x01FB},                 // Case map
	0x01FC:  []rune{0x01FD},                 // Case map
	0x01FE:  []rune{0x01FF},                 // Case map
	0x0200:  []rune{0x0201},                 // Case map
	0x0202:  []rune{0x0203},                 // Case map
	0x0204:  []rune{0x0205},                 // Case map
	0x0206:  []rune{0x0207},                 // Case map
	0x0208:  []rune{0x0209},                 // Case map
	0x020A:  []rune{0x020B},                 // Case map
	0x020C:  []rune{0x020D},                 // Case map
	0x020E:  []rune{0x020F},                 // Case map
	0x0210:  []rune{0x0211},                 // Case map
	0x0212:  []rune{0x0213},                 // Case map
	0x0214:  []rune{0x0215},                 // Case map
	0x0216:  []rune{0x0217},                 // Case map
	0x0218:  []rune{0x0219},                 // Case map
	0x021A:  []rune{0x021B},                 // Case map
	0x021C:  []rune{0x021D},                 // Case map
	0x021E:  []rune{0x021F},                 // Case map
	0x0220:  []rune{0x019E},                 // Case map
	0x0222:  []rune{0x0223},                 // Case map
	0x0224:  []rune{0x0225},                 // Case map
	0x0226:  []rune{0x0227},                 // Case map
	0x0228:  []rune{0x0229},                 // Case map
	0x022A:  []rune{0x022B},                 // Case map
	0x022C:  []rune{0x022D},                 // Case map
	0x022E:  []rune{0x022F},                 // Case map
	0x0230:  []rune{0x0231},                 // Case map
	0x0232:  []rune{0x0233},                 // Case map
	0x0345:  []rune{0x03B9},                 // Case map
	0x0386:  []rune{0x03AC},                 // Case map
	0x0388:  []rune{0x03AD},                 // Case map
	0x0389:  []rune{0x03AE},                 // Case map
	0x038A:  []rune{0x03AF},                 // Case map
	0x038C:  []rune{0x03CC},                 // Case map
	0x038E:  []rune{0x03CD},                 // Case map
	0x038F:  []rune{0x03CE},                 // Case map
	0x0390:  []rune{0x03B9, 0x0308, 0x0301}, // Case map
	0x0391:  []rune{0x03B1},                 // Case map
	0x0392:  []rune{0x03B2},                 // Case map
	0x0393:  []rune{0x03B3},                 // Case map
	0x0394:  []rune{0x03B4},                 // Case map
	0x0395:  []rune{0x03B5},                 // Case map
	0x0396:  []rune{0x03B6},                 // Case map
	0x0397:  []rune{0x03B7},                 // Case map
	0x0398:  []rune{0x03B8},                 // Case map
	0x0399:  []rune{0x03B9},                 // Case map
	0x039A:  []rune{0x03BA},                 // Case map
	0x039B:  []rune{0x03BB},                 // Case map
	0x039C:  []rune{0x03BC},                 // Case map
	0x039D:  []rune{0x03BD},                 // Case map
	0x039E:  []rune{0x03BE},                 // Case map
	0x039F:  []rune{0x03BF},                 // Case map
	0x03A0:  []rune{0x03C0},                 // Case map
	0x03A1:  []rune{0x03C1},                 // Case map
	0x03A3:  []rune{0x03C3},                 // Case map
	0x03A4:  []rune{0x03C4},                 // Case map
	0x03A5:  []rune{0x03C5},                 // Case map
	0x03A6:  []rune{0x03C6},                 // Case map
	0x03A7:  []rune{0x03C7},                 // Case map
	0x03A8:  []rune{0x03C8},                 // Case map
	0x03A9:  []rune{0x03C9},                 // Case map
	0x03AA:  []rune{0x03CA},                 // Case map
	0x03AB:  []rune{0x03CB},                 // Case map
	0x03B0:  []rune{0x03C5, 0x0308, 0x0301}, // Case map
	0x03C2:  []rune{0x03C3},                 // Case map
	0x03D0:  []rune{0x03B2},                 // Case map
	0x03D1:  []rune{0x03B8},                 // Case map
	0x03D5:  []rune{0x03C6},                 // Case map
	0x03D6:  []rune{0x03C0},                 // Case map
	0x03D8:  []rune{0x03D9},                 // Case map
	0x03DA:  []rune{0x03DB},                 // Case map
	0x03DC:  []rune{0x03DD},                 // Case map
	0x03DE:  []rune{0x03DF},                 // Case map
	0x03E0:  []rune{0x03E1},                 // Case map
	0x03E2:  []rune{0x03E3},                 // Case map
	0x03E4:  []rune{0x03E5},                 // Case map
	0x03E6:  []rune{0x03E7},                 // Case map
	0x03E8:  []rune{0x03E9},                 // Case map
	0x03EA:  []rune{0x03EB},                 // Case map
	0x03EC:  []rune{0x03ED},                 // Case map
	0x03EE:  []rune{0x03EF},                 // Case map
	0x03F0:  []rune{0x03BA},                 // Case map
	0x03F1:  []rune{0x03C1},                 // Case map
	0x03F2:  []rune{0x03C3},                 // Case map
	0x03F4:  []rune{0x03B8},                 // Case map
	0x03F5:  []rune{0x03B5},                 // Case map
	0x0400:  []rune{0x0450},                 // Case map
	0x0401:  []rune{0x0451},                 // Case map
	0x0402:  []rune{0x0452},                 // Case map
	0x0403:  []rune{0x0453},                 // Case map
	0x0404:  []rune{0x0454},                 // Case map
	0x0405:  []rune{0x0455},                 // Case map
	0x0406:  []rune{0x0456},                 // Case map
	0x0407:  []rune{0x0457},                 // Case map
	0x0408:  []rune{0x0458},                 // Case map
	0x0409:  []rune{0x0459},                 // Case map
	0x040A:  []rune{0x045A},                 // Case map
	0x040B:  []rune{0x045B},                 // Case map
	0x040C:  []rune{0x045C},                 // Case map
	0x040D:  []rune{0x045D},                 // Case map
	0x040E:  []rune{0x045E},                 // Case map
	0x040F:  []rune{0x045F},                 // Case map
	0x0410:  []rune{0x0430},                 // Case map
	0x0411:  []rune{0x0431},                 // Case map
	0x0412:  []rune{0x0432},                 // Case map
	0x0413:  []rune{0x0433},                 // Case map
	0x0414:  []rune{0x0434},                 // Case map
	0x0415:  []rune{0x0435},                 // Case map
	0x0416:  []rune{0x0436},                 // Case map
	0x0417:  []rune{0x0437},                 // Case map
	0x0418:  []rune{0x0438},                 // Case map
	0x0419:  []rune{0x0439},                 // Case map
	0x041A:  []rune{0x043A},                 // Case map
	0x041B:  []rune{0x043B},                 // Case map
	0x041C:  []rune{0x043C},                 // Case map
	0x041D:  []rune{0x043D},                 // Case map
	0x041E:  []rune{0x043E},                 // Case map
	0x041F:  []rune{0x043F},                 // Case map
	0x0420:  []rune{0x0440},                 // Case map
	0x0421:  []rune{0x0441},                 // Case map
	0x0422:  []rune{0x0442},                 // Case map
	0x0423:  []rune{0x0443},                 // Case map
	0x0424:  []rune{0x0444},                 // Case map
	0x0425:  []rune{0x0445},                 // Case map
	0x0426:  []rune{0x0446},                 // Case map
	0x0427:  []rune{0x0447},                 // Case map
	0x0428:  []rune{0x0448},                 // Case map
	0x0429:  []rune{0x0449},                 // Case map
	0x042A:  []rune{0x044A},                 // Case map
	0x042B:  []rune{0x044B},                 // Case map
	0x042C:  []rune{0x044C},                 // Case map
	0x042D:  []rune{0x044D},                 // Case map
	0x042E:  []rune{0x044E},                 // Case map
	0x042F:  []rune{0x044F},                 // Case map
	0x0460:  []rune{0x0461},                 // Case map
	0x0462:  []rune{0x0463},                 // Case map
	0x0464:  []rune{0x0465},                 // Case map
	0x0466:  []rune{0x0467},                 // Case map
	0x0468:  []rune{0x0469},                 // Case map
	0x046A:  []rune{0x046B},                 // Case map
	0x046C:  []rune{0x046D},                 // Case map
	0x046E:  []rune{0x046F},                 // Case map
	0x0470:  []rune{0x0471},                 // Case map
	0x0472:  []rune{0x0473},                 // Case map
	0x0474:  []rune{0x0475},                 // Case map
	0x0476:  []rune{0x0477},                 // Case map
	0x0478:  []rune{0x0479},                 // Case map
	0x047A:  []rune{0x047B},                 // Case map
	0x047C:  []rune{0x047D},                 // Case map
	0x047E:  []rune{0x047F},                 // Case map
	0x0480:  []rune{0x0481},                 // Case map
	0x048A:  []rune{0x048B},                 // Case map
	0x048C:  []rune{0x048D},                 // Case map
	0x048E:  []rune{0x048F},                 // Case map
	0x0490:  []rune{0x0491},                 // Case map
	0x0492:  []rune{0x0493},                 // Case map
	0x0494:  []rune{0x0495},                 // Case map
	0x0496:  []rune{0x0497},                 // Case map
	0x0498:  []rune{0x0499},                 // Case map
	0x049A:  []rune{0x049B},                 // Case map
	0x049C:  []rune{0x049D},                 // Case map
	0x049E:  []rune{0x049F},                 // Case map
	0x04A0:  []rune{0x04A1},                 // Case map
	0x04A2:  []rune{0x04A3},                 // Case map
	0x04A4:  []rune{0x04A5},                 // Case map
	0x04A6:  []rune{0x04A7},                 // Case map
	0x04A8:  []rune{0x04A9},                 // Case map
	0x04AA:  []rune{0x04AB},                 // Case map
	0x04AC:  []rune{0x04AD},                 // Case map
	0x04AE:  []rune{0x04AF},                 // Case map
	0x04B0:  []rune{0x04B1},                 // Case map
	0x04B2:  []rune{0x04B3},                 // Case map
	0x04B4:  []rune{0x04B5},                 // Case map
	0x04B6:  []rune{0x04B7},                 // Case map
	0x04B8:  []rune{0x04B9},                 // Case map
	0x04BA:  []rune{0x04BB},                 // Case map
	0x04BC:  []rune{0x04BD},                 // Case map
	0x04BE:  []rune{0x04BF},                 // Case map
	0x04C1:  []rune{0x04C2},                 // Case map
	0x04C3:  []rune{0x04C4},                 // Case map
	0x04C5:  []rune{0x04C6},                 // Case map
	0x04C7:  []rune{0x04C8},                 // Case map
	0x04C9:  []rune{0x04CA},                 // Case map
	0x04CB:  []rune{0x04CC},                 // Case map
	0x04CD:  []rune{0x04CE},                 // Case map
	0x04D0:  []rune{0x04D1},                 // Case map
	0x04D2:  []rune{0x04D3},                 // Case map
	0x04D4:  []rune{0x04D5},                 // Case map
	0x04D6:  []rune{0x04D7},                 // Case map
	0x04D8:  []rune{0x04D9},                 // Case map
	0x04DA:  []rune{0x04DB},                 // Case map
	0x04DC:  []rune{0x04DD},                 // Case map
	0x04DE:  []rune{0x04DF},                 // Case map
	0x04E0:  []rune{0x04E1},                 // Case map
	0x04E2:  []rune{0x04E3},                 // Case map
	0x04E4:  []rune{0x04E5},                 // Case map
	0x04E6:  []rune{0x04E7},                 // Case map
	0x04E8:  []rune{0x04E9},                 // Case map
	0x04EA:  []rune{0x04EB},                 // Case map
	0x04EC:  []rune{0x04ED},                 // Case map
	0x04EE:  []rune{0x04EF},                 // Case map
	0x04F0:  []rune{0x04F1},                 // Case map
	0x04F2:  []rune{0x04F3},                 // Case map
	0x04F4:  []rune{0x04F5},                 // Case map
	0x04F8:  []rune{0x04F9},                 // Case map
	0x0500:  []rune{0x0501},                 // Case map
	0x0502:  []rune{0x0503},                 // Case map
	0x0504:  []rune{0x0505},                 // Case map
	0x0506:  []rune{0x0507},                 // Case map
	0x0508:  []rune{0x0509},                 // Case map
	0x050A:  []rune{0x050B},                 // Case map
	0x050C:  []rune{0x050D},                 // Case map
	0x050E:  []rune{0x050F},                 // Case map
	0x0531:  []rune{0x0561},                 // Case map
	0x0532:  []rune{0x0562},                 // Case map
	0x0533:  []rune{0x0563},                 // Case map
	0x0534:  []rune{0x0564},                 // Case map
	0x0535:  []rune{0x0565},                 // Case map
	0x0536:  []rune{0x0566},                 // Case map
	0x0537:  []rune{0x0567},                 // Case map
	0x0538:  []rune{0x0568},                 // Case map
	0x0539:  []rune{0x0569},                 // Case map
	0x053A:  []rune{0x056A},                 // Case map
	0x053B:  []rune{0x056B},                 // Case map
	0x053C:  []rune{0x056C},                 // Case map
	0x053D:  []rune{0x056D},                 // Case map
	0x053E:  []rune{0x056E},                 // Case map
	0x053F:  []rune{0x056F},                 // Case map
	0x0540:  []rune{0x0570},                 // Case map
	0x0541:  []rune{0x0571},                 // Case map
	0x0542:  []rune{0x0572},                 // Case map
	0x0543:  []rune{0x0573},                 // Case map
	0x0544:  []rune{0x0574},                 // Case map
	0x0545:  []rune{0x0575},                 // Case map
	0x0546:  []rune{0x0576},                 // Case map
	0x0547:  []rune{0x0577},                 // Case map
	0x0548:  []rune{0x0578},                 // Case map
	0x0549:  []rune{0x0579},                 // Case map
	0x054A:  []rune{0x057A},                 // Case map
	0x054B:  []rune{0x057B},                 // Case map
	0x054C:  []rune{0x057C},                 // Case map
	0x054D:  []rune{0x057D},                 // Case map
	0x054E:  []rune{0x057E},                 // Case map
	0x054F:  []rune{0x057F},                 // Case map
	0x0550:  []rune{0x0580},                 // Case map
	0x0551:  []rune{0x0581},                 // Case map
	0x0552:  []rune{0x0582},                 // Case map
	0x0553:  []rune{0x0583},                 // Case map
	0x0554:  []rune{0x0584},                 // Case map
	0x0555:  []rune{0x0585},                 // Case map
	0x0556:  []rune{0x0586},                 // Case map
	0x0587:  []rune{0x0565, 0x0582},         // Case map
	0x1E00:  []rune{0x1E01},                 // Case map
	0x1E02:  []rune{0x1E03},                 // Case map
	0x1E04:  []rune{0x1E05},                 // Case map
	0x1E06:  []rune{0x1E07},                 // Case map
	0x1E08:  []rune{0x1E09},                 // Case map
	0x1E0A:  []rune{0x1E0B},                 // Case map
	0x1E0C:  []rune{0x1E0D},                 // Case map
	0x1E0E:  []rune{0x1E0F},                 // Case map
	0x1E10:  []rune{0x1E11},                 // Case map
	0x1E12:  []rune{0x1E13},                 // Case map
	0x1E14:  []rune{0x1E15},                 // Case map
	0x1E16:  []rune{0x1E17},                 // Case map
	0x1E18:  []rune{0x1E19},                 // Case map
	0x1E1A:  []rune{0x1E1B},                 // Case map
	0x1E1C:  []rune{0x1E1D},                 // Case map
	0x1E1E:  []rune{0x1E1F},                 // Case map
	0x1E20:  []rune{0x1E21},                 // Case map
	0x1E22:  []rune{0x1E23},                 // Case map
	0x1E24:  []rune{0x1E25},                 // Case map
	0x1E26:  []rune{0x1E27},                 // Case map
	0x1E28:  []rune{0x1E29},                 // Case map
	0x1E2A:  []rune{0x1E2B},                 // Case map
	0x1E2C:  []rune{0x1E2D},                 // Case map
	0x1E2E:  []rune{0x1E2F},                 // Case map
	0x1E30:  []rune{0x1E31},                 // Case map
	0x1E32:  []rune{0x1E33},                 // Case map
	0x1E34:  []rune{0x1E35},                 // Case map
	0x1E36:  []rune{0x1E37},                 // Case map
	0x1E38:  []rune{0x1E39},                 // Case map
	0x1E3A:  []rune{0x1E3B},                 // Case map
	0x1E3C:  []rune{0x1E3D},                 // Case map
	0x1E3E:  []rune{0x1E3F},                 // Case map
	0x1E40:  []rune{0x1E41},                 // Case map
	0x1E42:  []rune{0x1E43},                 // Case map
	0x1E44:  []rune{0x1E45},                 // Case map
	0x1E46:  []rune{0x1E47},                 // Case map
	0x1E48:  []rune{0x1E49},                 // Case map
	0x1E4A:  []rune{0x1E4B},                 // Case map
	0x1E4C:  []rune{0x1E4D},                 // Case map
	0x1E4E:  []rune{0x1E4F},                 // Case map
	0x1E50:  []rune{0x1E51},                 // Case map
	0x1E52:  []rune{0x1E53},                 // Case map
	0x1E54:  []rune{0x1E55},                 // Case map
	0x1E56:  []rune{0x1E57},                 // Case map
	0x1E58:  []rune{0x1E59},                 // Case map
	0x1E5A:  []rune{0x1E5B},                 // Case map
	0x1E5C:  []rune{0x1E5D},                 // Case map
	0x1E5E:  []rune{0x1E5F},                 // Case map
	0x1E60:  []rune{0x1E61},                 // Case map
	0x1E62:  []rune{0x1E63},                 // Case map
	0x1E64:  []rune{0x1E65},                 // Case map
	0x1E66:  []rune{0x1E67},                 // Case map
	0x1E68:  []rune{0x1E69},                 // Case map
	0x1E6A:  []rune{0x1E6B},                 // Case map
	0x1E6C:  []rune{0x1E6D},                 // Case map
	0x1E6E:  []rune{0x1E6F},                 // Case map
	0x1E70:  []rune{0x1E71},                 // Case map
	0x1E72:  []rune{0x1E73},                 // Case map
	0x1E74:  []rune{0x1E75},                 // Case map
	0x1E76:  []rune{0x1E77},                 // Case map
	0x1E78:  []rune{0x1E79},                 // Case map
	0x1E7A:  []rune{0x1E7B},                 // Case map
	0x1E7C:  []rune{0x1E7D},                 // Case map
	0x1E7E:  []rune{0x1E7F},                 // Case map
	0x1E80:  []rune{0x1E81},                 // Case map
	0x1E82:  []rune{0x1E83},                 // Case map
	0x1E84:  []rune{0x1E85},                 // Case map
	0x1E86:  []rune{0x1E87},                 // Case map
	0x1E88:  []rune{0x1E89},                 // Case map
	0x1E8A:  []rune{0x1E8B},                 // Case map
	0x1E8C:  []rune{0x1E8D},                 // Case map
	0x1E8E:  []rune{0x1E8F},                 // Case map
	0x1E90:  []rune{0x1E91},                 // Case map
	0x1E92:  []rune{0x1E93},                 // Case map
	0x1E94:  []rune{0x1E95},                 // Case map
	0x1E96:  []rune{0x0068, 0x0331},         // Case map
	0x1E97:  []rune{0x0074, 0x0308},         // Case map
	0x1E98:  []rune{0x0077, 0x030A},         // Case map
	0x1E99:  []rune{0x0079, 0x030A},         // Case map
	0x1E9A:  []rune{0x0061, 0x02BE},         // Case map
	0x1E9B:  []rune{0x1E61},                 // Case map
	0x1EA0:  []rune{0x1EA1},                 // Case map
	0x1EA2:  []rune{0x1EA3},                 // Case map
	0x1EA4:  []rune{0x1EA5},                 // Case map
	0x1EA6:  []rune{0x1EA7},                 // Case map
	0x1EA8:  []rune{0x1EA9},                 // Case map
	0x1EAA:  []rune{0x1EAB},                 // Case map
	0x1EAC:  []rune{0x1EAD},                 // Case map
	0x1EAE:  []rune{0x1EAF},                 // Case map
	0x1EB0:  []rune{0x1EB1},                 // Case map
	0x1EB2:  []rune{0x1EB3},                 // Case map
	0x1EB4:  []rune{0x1EB5},                 // Case map
	0x1EB6:  []rune{0x1EB7},                 // Case map
	0x1EB8:  []rune{0x1EB9},                 // Case map
	0x1EBA:  []rune{0x1EBB},                 // Case map
	0x1EBC:  []rune{0x1EBD},                 // Case map
	0x1EBE:  []rune{0x1EBF},                 // Case map
	0x1EC0:  []rune{0x1EC1},                 // Case map
	0x1EC2:  []rune{0x1EC3},                 // Case map
	0x1EC4:  []rune{0x1EC5},                 // Case map
	0x1EC6:  []rune{0x1EC7},                 // Case map
	0x1EC8:  []rune{0x1EC9},                 // Case map
	0x1ECA:  []rune{0x1ECB},                 // Case map
	0x1ECC:  []rune{0x1ECD},                 // Case map
	0x1ECE:  []rune{0x1ECF},                 // Case map
	0x1ED0:  []rune{0x1ED1},                 // Case map
	0x1ED2:  []rune{0x1ED3},                 // Case map
	0x1ED4:  []rune{0x1ED5},                 // Case map
	0x1ED6:  []rune{0x1ED7},                 // Case map
	0x1ED8:  []rune{0x1ED9},                 // Case map
	0x1EDA:  []rune{0x1EDB},                 // Case map
	0x1EDC:  []rune{0x1EDD},                 // Case map
	0x1EDE:  []rune{0x1EDF},                 // Case map
	0x1EE0:  []rune{0x1EE1},                 // Case map
	0x1EE2:  []rune{0x1EE3},                 // Case map
	0x1EE4:  []rune{0x1EE5},                 // Case map
	0x1EE6:  []rune{0x1EE7},                 // Case map
	0x1EE8:  []rune{0x1EE9},                 // Case map
	0x1EEA:  []rune{0x1EEB},                 // Case map
	0x1EEC:  []rune{0x1EED},                 // Case map
	0x1EEE:  []rune{0x1EEF},                 // Case map
	0x1EF0:  []rune{0x1EF1},                 // Case map
	0x1EF2:  []rune{0x1EF3},                 // Case map
	0x1EF4:  []rune{0x1EF5},                 // Case map
	0x1EF6:  []rune{0x1EF7},                 // Case map
	0x1EF8:  []rune{0x1EF9},                 // Case map
	0x1F08:  []rune{0x1F00},                 // Case map
	0x1F09:  []rune{0x1F01},                 // Case map
	0x1F0A:  []rune{0x1F02},                 // Case map
	0x1F0B:  []rune{0x1F03},                 // Case map
	0x1F0C:  []rune{0x1F04},                 // Case map
	0x1F0D:  []rune{0x1F05},                 // Case map
	0x1F0E:  []rune{0x1F06},                 // Case map
	0x1F0F:  []rune{0x1F07},                 // Case map
	0x1F18:  []rune{0x1F10},                 // Case map
	0x1F19:  []rune{0x1F11},                 // Case map
	0x1F1A:  []rune{0x1F12},                 // Case map
	0x1F1B:  []rune{0x1F13},                 // Case map
	0x1F1C:  []rune{0x1F14},                 // Case map
	0x1F1D:  []rune{0x1F15},                 // Case map
	0x1F28:  []rune{0x1F20},                 // Case map
	0x1F29:  []rune{0x1F21},                 // Case map
	0x1F2A:  []rune{0x1F22},                 // Case map
	0x1F2B:  []rune{0x1F23},                 // Case map
	0x1F2C:  []rune{0x1F24},                 // Case map
	0x1F2D:  []rune{0x1F25},                 // Case map
	0x1F2E:  []rune{0x1F26},                 // Case map
	0x1F2F:  []rune{0x1F27},                 // Case map
	0x1F38:  []rune{0x1F30},                 // Case map
	0x1F39:  []rune{0x1F31},                 // Case map
	0x1F3A:  []rune{0x1F32},                 // Case map
	0x1F3B:  []rune{0x1F33},                 // Case map
	0x1F3C:  []rune{0x1F34},                 // Case map
	0x1F3D:  []rune{0x1F35},                 // Case map
	0x1F3E:  []rune{0x1F36},                 // Case map
	0x1F3F:  []rune{0x1F37},                 // Case map
	0x1F48:  []rune{0x1F40},                 // Case map
	0x1F49:  []rune{0x1F41},                 // Case map
	0x1F4A:  []rune{0x1F42},                 // Case map
	0x1F4B:  []rune{0x1F43},                 // Case map
	0x1F4C:  []rune{0x1F44},                 // Case map
	0x1F4D:  []rune{0x1F45},                 // Case map
	0x1F50:  []rune{0x03C5, 0x0313},         // Case map
	0x1F52:  []rune{0x03C5, 0x0313, 0x0300}, // Case map
	0x1F54:  []rune{0x03C5, 0x0313, 0x0301}, // Case map
	0x1F56:  []rune{0x03C5, 0x0313, 0x0342}, // Case map
	0x1F59:  []rune{0x1F51},                 // Case map
	0x1F5B:  []rune{0x1F53},                 // Case map
	0x1F5D:  []rune{0x1F55},                 // Case map
	0x1F5F:  []rune{0x1F57},                 // Case map
	0x1F68:  []rune{0x1F60},                 // Case map
	0x1F69:  []rune{0x1F61},                 // Case map
	0x1F6A:  []rune{0x1F62},                 // Case map
	0x1F6B:  []rune{0x1F63},                 // Case map
	0x1F6C:  []rune{0x1F64},                 // Case map
	0x1F6D:  []rune{0x1F65},                 // Case map
	0x1F6E:  []rune{0x1F66},                 // Case map
	0x1F6F:  []rune{0x1F67},                 // Case map
	0x1F80:  []rune{0x1F00, 0x03B9},         // Case map
	0x1F81:  []rune{0x1F01, 0x03B9},         // Case map
	0x1F82:  []rune{0x1F02, 0x03B9},         // Case map
	0x1F83:  []rune{0x1F03, 0x03B9},         // Case map
	0x1F84:  []rune{0x1F04, 0x03B9},         // Case map
	0x1F85:  []rune{0x1F05, 0x03B9},         // Case map
	0x1F86:  []rune{0x1F06, 0x03B9},         // Case map
	0x1F87:  []rune{0x1F07, 0x03B9},         // Case map
	0x1F88:  []rune{0x1F00, 0x03B9},         // Case map
	0x1F89:  []rune{0x1F01, 0x03B9},         // Case map
	0x1F8A:  []rune{0x1F02, 0x03B9},         // Case map
	0x1F8B:  []rune{0x1F03, 0x03B9},         // Case map
	0x1F8C:  []rune{0x1F04, 0x03B9},         // Case map
	0x1F8D:  []rune{0x1F05, 0x03B9},         // Case map
	0x1F8E:  []rune{0x1F06, 0x03B9},         // Case map
	0x1F8F:  []rune{0x1F07, 0x03B9},         // Case map
	0x1F90:  []rune{0x1F20, 0x03B9},         // Case map
	0x1F91:  []rune{0x1F21, 0x03B9},         // Case map
	0x1F92:  []rune{0x1F22, 0x03B9},         // Case map
	0x1F93:  []rune{0x1F23, 0x03B9},         // Case map
	0x1F94:  []rune{0x1F24, 0x03B9},         // Case map
	0x1F95:  []rune{0x1F25, 0x03B9},         // Case map
	0x1F96:  []rune{0x1F26, 0x03B9},         // Case map
	0x1F97:  []rune{0x1F27, 0x03B9},         // Case map
	0x1F98:  []rune{0x1F20, 0x03B9},         // Case map
	0x1F99:  []rune{0x1F21, 0x03B9},         // Case map
	0x1F9A:  []rune{0x1F22, 0x03B9},         // Case map
	0x1F9B:  []rune{0x1F23, 0x03B9},         // Case map
	0x1F9C:  []rune{0x1F24, 0x03B9},         // Case map
	0x1F9D:  []rune{0x1F25, 0x03B9},         // Case map
	0x1F9E:  []rune{0x1F26, 0x03B9},         // Case map
	0x1F9F:  []rune{0x1F27, 0x03B9},         // Case map
	0x1FA0:  []rune{0x1F60, 0x03B9},         // Case map
	0x1FA1:  []rune{0x1F61, 0x03B9},         // Case map
	0x1FA2:  []rune{0x1F62, 0x03B9},         // Case map
	0x1FA3:  []rune{0x1F63, 0x03B9},         // Case map
	0x1FA4:  []rune{0x1F64, 0x03B9},         // Case map
	0x1FA5:  []rune{0x1F65, 0x03B9},         // Case map
	0x1FA6:  []rune{0x1F66, 0x03B9},         // Case map
	0x1FA7:  []rune{0x1F67, 0x03B9},         // Case map
	0x1FA8:  []rune{0x1F60, 0x03B9},         // Case map
	0x1FA9:  []rune{0x1F61, 0x03B9},         // Case map
	0x1FAA:  []rune{0x1F62, 0x03B9},         // Case map
	0x1FAB:  []rune{0x1F63, 0x03B9},         // Case map
	0x1FAC:  []rune{0x1F64, 0x03B9},         // Case map
	0x1FAD:  []rune{0x1F65, 0x03B9},         // Case map
	0x1FAE:  []rune{0x1F66, 0x03B9},         // Case map
	0x1FAF:  []rune{0x1F67, 0x03B9},         // Case map
	0x1FB2:  []rune{0x1F70, 0x03B9},         // Case map
	0x1FB3:  []rune{0x03B1, 0x03B9},         // Case map
	0x1FB4:  []rune{0x03AC, 0x03B9},         // Case map
	0x1FB6:  []rune{0x03B1, 0x0342},         // Case map
	0x1FB7:  []rune{0x03B1, 0x0342, 0x03B9}, // Case map
	0x1FB8:  []rune{0x1FB0},                 // Case map
	0x1FB9:  []rune{0x1FB1},                 // Case map
	0x1FBA:  []rune{0x1F70},                 // Case map
	0x1FBB:  []rune{0x1F71},                 // Case map
	0x1FBC:  []rune{0x03B1, 0x03B9},         // Case map
	0x1FBE:  []rune{0x03B9},                 // Case map
	0x1FC2:  []rune{0x1F74, 0x03B9},         // Case map
	0x1FC3:  []rune{0x03B7, 0x03B9},         // Case map
	0x1FC4:  []rune{0x03AE, 0x03B9},         // Case map
	0x1FC6:  []rune{0x03B7, 0x0342},         // Case map
	0x1FC7:  []rune{0x03B7, 0x0342, 0x03B9}, // Case map
	0x1FC8:  []rune{0x1F72},                 // Case map
	0x1FC9:  []rune{0x1F73},                 // Case map
	0x1FCA:  []rune{0x1F74},                 // Case map
	0x1FCB:  []rune{0x1F75},                 // Case map
	0x1FCC:  []rune{0x03B7, 0x03B9},         // Case map
	0x1FD2:  []rune{0x03B9, 0x0308, 0x0300}, // Case map
	0x1FD3:  []rune{0x03B9, 0x0308, 0x0301}, // Case map
	0x1FD6:  []rune{0x03B9, 0x0342},         // Case map
	0x1FD7:  []rune{0x03B9, 0x0308, 0x0342}, // Case map
	0x1FD8:  []rune{0x1FD0},                 // Case map
	0x1FD9:  []rune{0x1FD1},                 // Case map
	0x1FDA:  []rune{0x1F76},                 // Case map
	0x1FDB:  []rune{0x1F77},                 // Case map
	0x1FE2:  []rune{0x03C5, 0x0308, 0x0300}, // Case map
	0x1FE3:  []rune{0x03C5, 0x0308, 0x0301}, // Case map
	0x1FE4:  []rune{0x03C1, 0x0313},         // Case map
	0x1FE6:  []rune{0x03C5, 0x0342},         // Case map
	0x1FE7:  []rune{0x03C5, 0x0308, 0x0342}, // Case map
	0x1FE8:  []rune{0x1FE0},                 // Case map
	0x1FE9:  []rune{0x1FE1},                 // Case map
	0x1FEA:  []rune{0x1F7A},                 // Case map
	0x1FEB:  []rune{0x1F7B},                 // Case map
	0x1FEC:  []rune{0x1FE5},                 // Case map
	0x1FF2:  []rune{0x1F7C, 0x03B9},         // Case map
	0x1FF3:  []rune{0x03C9, 0x03B9},         // Case map
	0x1FF4:  []rune{0x03CE, 0x03B9},         // Case map
	0x1FF6:  []rune{0x03C9, 0x0342},         // Case map
	0x1FF7:  []rune{0x03C9, 0x0342, 0x03B9}, // Case map
	0x1FF8:  []rune{0x1F78},                 // Case map
	0x1FF9:  []rune{0x1F79},                 // Case map
	0x1FFA:  []rune{0x1F7C},                 // Case map
	0x1FFB:  []rune{0x1F7D},                 // Case map
	0x1FFC:  []rune{0x03C9, 0x03B9},         // Case map
	0x2126:  []rune{0x03C9},                 // Case map
	0x212A:  []rune{0x006B},                 // Case map
	0x212B:  []rune{0x00E5},                 // Case map
	0x2160:  []rune{0x2170},                 // Case map
	0x2161:  []rune{0x2171},                 // Case map
	0x2162:  []rune{0x2172},                 // Case map
	0x2163:  []rune{0x2173},                 // Case map
	0x2164:  []rune{0x2174},                 // Case map
	0x2165:  []rune{0x2175},                 // Case map
	0x2166:  []rune{0x2176},                 // Case map
	0x2167:  []rune{0x2177},                 // Case map
	0x2168:  []rune{0x2178},                 // Case map
	0x2169:  []rune{0x2179},                 // Case map
	0x216A:  []rune{0x217A},                 // Case map
	0x216B:  []rune{0x217B},                 // Case map
	0x216C:  []rune{0x217C},                 // Case map
	0x216D:  []rune{0x217D},                 // Case map
	0x216E:  []rune{0x217E},                 // Case map
	0x216F:  []rune{0x217F},                 // Case map
	0x24B6:  []rune{0x24D0},                 // Case map
	0x24B7:  []rune{0x24D1},                 // Case map
	0x24B8:  []rune{0x24D2},                 // Case map
	0x24B9:  []rune{0x24D3},                 // Case map
	0x24BA:  []rune{0x24D4},                 // Case map
	0x24BB:  []rune{0x24D5},                 // Case map
	0x24BC:  []rune{0x24D6},                 // Case map
	0x24BD:  []rune{0x24D7},                 // Case map
	0x24BE:  []rune{0x24D8},                 // Case map
	0x24BF:  []rune{0x24D9},                 // Case map
	0x24C0:  []rune{0x24DA},                 // Case map
	0x24C1:  []rune{0x24DB},                 // Case map
	0x24C2:  []rune{0x24DC},                 // Case map
	0x24C3:  []rune{0x24DD},                 // Case map
	0x24C4:  []rune{0x24DE},                 // Case map
	0x24C5:  []rune{0x24DF},                 // Case map
	0x24C6:  []rune{0x24E0},                 // Case map
	0x24C7:  []rune{0x24E1},                 // Case map
	0x24C8:  []rune{0x24E2},                 // Case map
	0x24C9:  []rune{0x24E3},                 // Case map
	0x24CA:  []rune{0x24E4},                 // Case map
	0x24CB:  []rune{0x24E5},                 // Case map
	0x24CC:  []rune{0x24E6},                 // Case map
	0x24CD:  []rune{0x24E7},                 // Case map
	0x24CE:  []rune{0x24E8},                 // Case map
	0x24CF:  []rune{0x24E9},                 // Case map
	0xFB00:  []rune{0x0066, 0x0066},         // Case map
	0xFB01:  []rune{0x0066, 0x0069},         // Case map
	0xFB02:  []rune{0x0066, 0x006C},         // Case map
	0xFB03:  []rune{0x0066, 0x0066, 0x0069}, // Case map
	0xFB04:  []rune{0x0066, 0x0066, 0x006C}, // Case map
	0xFB05:  []rune{0x0073, 0x0074},         // Case map
	0xFB06:  []rune{0x0073, 0x0074},         // Case map
	0xFB13:  []rune{0x0574, 0x0576},         // Case map
	0xFB14:  []rune{0x0574, 0x0565},         // Case map
	0xFB15:  []rune{0x0574, 0x056B},         // Case map
	0xFB16:  []rune{0x057E, 0x0576},         // Case map
	0xFB17:  []rune{0x0574, 0x056D},         // Case map
	0xFF21:  []rune{0xFF41},                 // Case map
	0xFF22:  []rune{0xFF42},                 // Case map
	0xFF23:  []rune{0xFF43},                 // Case map
	0xFF24:  []rune{0xFF44},                 // Case map
	0xFF25:  []rune{0xFF45},                 // Case map
	0xFF26:  []rune{0xFF46},                 // Case map
	0xFF27:  []rune{0xFF47},                 // Case map
	0xFF28:  []rune{0xFF48},                 // Case map
	0xFF29:  []rune{0xFF49},                 // Case map
	0xFF2A:  []rune{0xFF4A},                 // Case map
	0xFF2B:  []rune{0xFF4B},                 // Case map
	0xFF2C:  []rune{0xFF4C},                 // Case map
	0xFF2D:  []rune{0xFF4D},                 // Case map
	0xFF2E:  []rune{0xFF4E},                 // Case map
	0xFF2F:  []rune{0xFF4F},                 // Case map
	0xFF30:  []rune{0xFF50},                 // Case map
	0xFF31:  []rune{0xFF51},                 // Case map
	0xFF32:  []rune{0xFF52},                 // Case map
	0xFF33:  []rune{0xFF53},                 // Case map
	0xFF34:  []rune{0xFF54},                 // Case map
	0xFF35:  []rune{0xFF55},                 // Case map
	0xFF36:  []rune{0xFF56},                 // Case map
	0xFF37:  []rune{0xFF57},                 // Case map
	0xFF38:  []rune{0xFF58},                 // Case map
	0xFF39:  []rune{0xFF59},                 // Case map
	0xFF3A:  []rune{0xFF5A},                 // Case map
	0x10400: []rune{0x10428},                // Case map
	0x10401: []rune{0x10429},                // Case map
	0x10402: []rune{0x1042A},                // Case map
	0x10403: []rune{0x1042B},                // Case map
	0x10404: []rune{0x1042C},                // Case map
	0x10405: []rune{0x1042D},                // Case map
	0x10406: []rune{0x1042E},                // Case map
	0x10407: []rune{0x1042F},                // Case map
	0x10408: []rune{0x10430},                // Case map
	0x10409: []rune{0x10431},                // Case map
	0x1040A: []rune{0x10432},                // Case map
	0x1040B: []rune{0x10433},                // Case map
	0x1040C: []rune{0x10434},                // Case map
	0x1040D: []rune{0x10435},                // Case map
	0x1040E: []rune{0x10436},                // Case map
	0x1040F: []rune{0x10437},                // Case map
	0x10410: []rune{0x10438},                // Case map
	0x10411: []rune{0x10439},                // Case map
	0x10412: []rune{0x1043A},                // Case map
	0x10413: []rune{0x1043B},                // Case map
	0x10414: []rune{0x1043C},                // Case map
	0x10415: []rune{0x1043D},                // Case map
	0x10416: []rune{0x1043E},                // Case map
	0x10417: []rune{0x1043F},                // Case map
	0x10418: []rune{0x10440},                // Case map
	0x10419: []rune{0x10441},                // Case map
	0x1041A: []rune{0x10442},                // Case map
	0x1041B: []rune{0x10443},                // Case map
	0x1041C: []rune{0x10444},                // Case map
	0x1041D: []rune{0x10445},                // Case map
	0x1041E: []rune{0x10446},                // Case map
	0x1041F: []rune{0x10447},                // Case map
	0x10420: []rune{0x10448},                // Case map
	0x10421: []rune{0x10449},                // Case map
	0x10422: []rune{0x1044A},                // Case map
	0x10423: []rune{0x1044B},                // Case map
	0x10424: []rune{0x1044C},                // Case map
	0x10425: []rune{0x1044D},                // Case map
}

// TableB3 represents RFC-3454 Table B.3.
var TableB3 Mapping = tableB3

var tableC1_1 = Set{
	RuneRange{0x0020, 0x0020}, // SPACE
}

// TableC1_1 represents RFC-3454 Table C.1.1.
var TableC1_1 Set = tableC1_1

var tableC1_2 = Set{
	RuneRange{0x00A0, 0x00A0}, // NO-BREAK SPACE
	RuneRange{0x1680, 0x1680}, // OGHAM SPACE MARK
	RuneRange{0x2000, 0x2000}, // EN QUAD
	RuneRange{0x2001, 0x2001}, // EM QUAD
	RuneRange{0x2002, 0x2002}, // EN SPACE
	RuneRange{0x2003, 0x2003}, // EM SPACE
	RuneRange{0x2004, 0x2004}, // THREE-PER-EM SPACE
	RuneRange{0x2005, 0x2005}, // FOUR-PER-EM SPACE
	RuneRange{0x2006, 0x2006}, // SIX-PER-EM SPACE
	RuneRange{0x2007, 0x2007}, // FIGURE SPACE
	RuneRange{0x2008, 0x2008}, // PUNCTUATION SPACE
	RuneRange{0x2009, 0x2009}, // THIN SPACE
	RuneRange{0x200A, 0x200A}, // HAIR SPACE
	RuneRange{0x200B, 0x200B}, // ZERO WIDTH SPACE
	RuneRange{0x202F, 0x202F}, // NARROW NO-BREAK SPACE
	RuneRange{0x205F, 0x205F}, // MEDIUM MATHEMATICAL SPACE
	RuneRange{0x3000, 0x3000}, // IDEOGRAPHIC SPACE
}

// TableC1_2 represents RFC-3454 Table C.1.2.
var TableC1_2 Set = tableC1_2

var tableC2_1 = Set{
	RuneRange{0x0000, 0x001F}, // [CONTROL CHARACTERS]
	RuneRange{0x007F, 0x007F}, // DELETE
}

// TableC2_1 represents RFC-3454 Table C.2.1.
var TableC2_1 Set = tableC2_1

var tableC2_2 = Set{
	RuneRange{0x0080, 0x009F},   // [CONTROL CHARACTERS]
	RuneRange{0x06DD, 0x06DD},   // ARABIC END OF AYAH
	RuneRange{0x070F, 0x070F},   // SYRIAC ABBREVIATION MARK
	RuneRange{0x180E, 0x180E},   // MONGOLIAN VOWEL SEPARATOR
	RuneRange{0x200C, 0x200C},   // ZERO WIDTH NON-JOINER
	RuneRange{0x200D, 0x200D},   // ZERO WIDTH JOINER
	RuneRange{0x2028, 0x2028},   // LINE SEPARATOR
	RuneRange{0x2029, 0x2029},   // PARAGRAPH SEPARATOR
	RuneRange{0x2060, 0x2060},   // WORD JOINER
	RuneRange{0x2061, 0x2061},   // FUNCTION APPLICATION
	RuneRange{0x2062, 0x2062},   // INVISIBLE TIMES
	RuneRange{0x2063, 0x2063},   // INVISIBLE SEPARATOR
	RuneRange{0x206A, 0x206F},   // [CONTROL CHARACTERS]
	RuneRange{0xFEFF, 0xFEFF},   // ZERO WIDTH NO-BREAK SPACE
	RuneRange{0xFFF9, 0xFFFC},   // [CONTROL CHARACTERS]
	RuneRange{0x1D173, 0x1D17A}, // [MUSICAL CONTROL CHARACTERS]
}

// TableC2_2 represents RFC-3454 Table C.2.2.
var TableC2_2 Set = tableC2_2

var tableC3 = Set{
	RuneRange{0xE000, 0xF8FF},     // [PRIVATE USE, PLANE 0]
	RuneRange{0xF0000, 0xFFFFD},   // [PRIVATE USE, PLANE 15]
	RuneRange{0x100000, 0x10FFFD}, // [PRIVATE USE, PLANE 16]
}

// TableC3 represents RFC-3454 Table C.3.
var TableC3 Set = tableC3

var tableC4 = Set{
	RuneRange{0xFDD0, 0xFDEF},     // [NONCHARACTER CODE POINTS]
	RuneRange{0xFFFE, 0xFFFF},     // [NONCHARACTER CODE POINTS]
	RuneRange{0x1FFFE, 0x1FFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0x2FFFE, 0x2FFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0x3FFFE, 0x3FFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0x4FFFE, 0x4FFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0x5FFFE, 0x5FFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0x6FFFE, 0x6FFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0x7FFFE, 0x7FFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0x8FFFE, 0x8FFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0x9FFFE, 0x9FFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0xAFFFE, 0xAFFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0xBFFFE, 0xBFFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0xCFFFE, 0xCFFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0xDFFFE, 0xDFFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0xEFFFE, 0xEFFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0xFFFFE, 0xFFFFF},   // [NONCHARACTER CODE POINTS]
	RuneRange{0x10FFFE, 0x10FFFF}, // [NONCHARACTER CODE POINTS]
}

// TableC4 represents RFC-3454 Table C.4.
var TableC4 Set = tableC4

var tableC5 = Set{
	RuneRange{0xD800, 0xDFFF}, // [SURROGATE CODES]
}

// TableC5 represents RFC-3454 Table C.5.
var TableC5 Set = tableC5

var tableC6 = Set{
	RuneRange{0xFFF9, 0xFFF9}, // INTERLINEAR ANNOTATION ANCHOR
	RuneRange{0xFFFA, 0xFFFA}, // INTERLINEAR ANNOTATION SEPARATOR
	RuneRange{0xFFFB, 0xFFFB}, // INTERLINEAR ANNOTATION TERMINATOR
	RuneRange{0xFFFC, 0xFFFC}, // OBJECT REPLACEMENT CHARACTER
	RuneRange{0xFFFD, 0xFFFD}, // REPLACEMENT CHARACTER
}

// TableC6 represents RFC-3454 Table C.6.
var TableC6 Set = tableC6

var tableC7 = Set{
	RuneRange{0x2FF0, 0x2FFB}, // [IDEOGRAPHIC DESCRIPTION CHARACTERS]
}

// TableC7 represents RFC-3454 Table C.7.
var TableC7 Set = tableC7

var tableC8 = Set{
	RuneRange{0x0340, 0x0340}, // COMBINING GRAVE TONE MARK
	RuneRange{0x0341, 0x0341}, // COMBINING ACUTE TONE MARK
	RuneRange{0x200E, 0x200E}, // LEFT-TO-RIGHT MARK
	RuneRange{0x200F, 0x200F}, // RIGHT-TO-LEFT MARK
	RuneRange{0x202A, 0x202A}, // LEFT-TO-RIGHT EMBEDDING
	RuneRange{0x202B, 0x202B}, // RIGHT-TO-LEFT EMBEDDING
	RuneRange{0x202C, 0x202C}, // POP DIRECTIONAL FORMATTING
	RuneRange{0x202D, 0x202D}, // LEFT-TO-RIGHT OVERRIDE
	RuneRange{0x202E, 0x202E}, // RIGHT-TO-LEFT OVERRIDE
	RuneRange{0x206A, 0x206A}, // INHIBIT SYMMETRIC SWAPPING
	RuneRange{0x206B, 0x206B}, // ACTIVATE SYMMETRIC SWAPPING
	RuneRange{0x206C, 0x206C}, // INHIBIT ARABIC FORM SHAPING
	RuneRange{0x206D, 0x206D}, // ACTIVATE ARABIC FORM SHAPING
	RuneRange{0x206E, 0x206E}, // NATIONAL DIGIT SHAPES
	RuneRange{0x206F, 0x206F}, // NOMINAL DIGIT SHAPES
}

// TableC8 represents RFC-3454 Table C.8.
var TableC8 Set = tableC8

var tableC9 = Set{
	RuneRange{0xE0001, 0xE0001}, // LANGUAGE TAG
	RuneRange{0xE0020, 0xE007F}, // [TAGGING CHARACTERS]
}

// TableC9 represents RFC-3454 Table C.9.
var TableC9 Set = tableC9

var tableD1 = Set{
	RuneRange{0x05BE, 0x05BE},
	RuneRange{0x05C0, 0x05C0},
	RuneRange{0x05C3, 0x05C3},
	RuneRange{0x05D0, 0x05EA},
	RuneRange{0x05F0, 0x05F4},
	RuneRange{0x061B, 0x061B},
	RuneRange{0x061F, 0x061F},
	RuneRange{0x0621, 0x063A},
	RuneRange{0x0640, 0x064A},
	RuneRange{0x066D, 0x066F},
	RuneRange{0x0671, 0x06D5},
	RuneRange{0x06DD, 0x06DD},
	RuneRange{0x06E5, 0x06E6},
	RuneRange{0x06FA, 0x06FE},
	RuneRange{0x0700, 0x070D},
	RuneRange{0x0710, 0x0710},
	RuneRange{0x0712, 0x072C},
	RuneRange{0x0780, 0x07A5},
	RuneRange{0x07B1, 0x07B1},
	RuneRange{0x200F, 0x200F},
	RuneRange{0xFB1D, 0xFB1D},
	RuneRange{0xFB1F, 0xFB28},
	RuneRange{0xFB2A, 0xFB36},
	RuneRange{0xFB38, 0xFB3C},
	RuneRange{0xFB3E, 0xFB3E},
	RuneRange{0xFB40, 0xFB41},
	RuneRange{0xFB43, 0xFB44},
	RuneRange{0xFB46, 0xFBB1},
	RuneRange{0xFBD3, 0xFD3D},
	RuneRange{0xFD50, 0xFD8F},
	RuneRange{0xFD92, 0xFDC7},
	RuneRange{0xFDF0, 0xFDFC},
	RuneRange{0xFE70, 0xFE74},
	RuneRange{0xFE76, 0xFEFC},
}

// TableD1 represents RFC-3454 Table D.1.
var TableD1 Set = tableD1

var tableD2 = Set{
	RuneRange{0x0041, 0x005A},
	RuneRange{0x0061, 0x007A},
	RuneRange{0x00AA, 0x00AA},
	RuneRange{0x00B5, 0x00B5},
	RuneRange{0x00BA, 0x00BA},
	RuneRange{0x00C0, 0x00D6},
	RuneRange{0x00D8, 0x00F6},
	RuneRange{0x00F8, 0x0220},
	RuneRange{0x0222, 0x0233},
	RuneRange{0x0250, 0x02AD},
	RuneRange{0x02B0, 0x02B8},
	RuneRange{0x02BB, 0x02C1},
	RuneRange{0x02D0, 0x02D1},
	RuneRange{0x02E0, 0x02E4},
	RuneRange{0x02EE, 0x02EE},
	RuneRange{0x037A, 0x037A},
	RuneRange{0x0386, 0x0386},
	RuneRange{0x0388, 0x038A},
	RuneRange{0x038C, 0x038C},
	RuneRange{0x038E, 0x03A1},
	RuneRange{0x03A3, 0x03CE},
	RuneRange{0x03D0, 0x03F5},
	RuneRange{0x0400, 0x0482},
	RuneRange{0x048A, 0x04CE},
	RuneRange{0x04D0, 0x04F5},
	RuneRange{0x04F8, 0x04F9},
	RuneRange{0x0500, 0x050F},
	RuneRange{0x0531, 0x0556},
	RuneRange{0x0559, 0x055F},
	RuneRange{0x0561, 0x0587},
	RuneRange{0x0589, 0x0589},
	RuneRange{0x0903, 0x0903},
	RuneRange{0x0905, 0x0939},
	RuneRange{0x093D, 0x0940},
	RuneRange{0x0949, 0x094C},
	RuneRange{0x0950, 0x0950},
	RuneRange{0x0958, 0x0961},
	RuneRange{0x0964, 0x0970},
	RuneRange{0x0982, 0x0983},
	RuneRange{0x0985, 0x098C},
	RuneRange{0x098F, 0x0990},
	RuneRange{0x0993, 0x09A8},
	RuneRange{0x09AA, 0x09B0},
	RuneRange{0x09B2, 0x09B2},
	RuneRange{0x09B6, 0x09B9},
	RuneRange{0x09BE, 0x09C0},
	RuneRange{0x09C7, 0x09C8},
	RuneRange{0x09CB, 0x09CC},
	RuneRange{0x09D7, 0x09D7},
	RuneRange{0x09DC, 0x09DD},
	RuneRange{0x09DF, 0x09E1},
	RuneRange{0x09E6, 0x09F1},
	RuneRange{0x09F4, 0x09FA},
	RuneRange{0x0A05, 0x0A0A},
	RuneRange{0x0A0F, 0x0A10},
	RuneRange{0x0A13, 0x0A28},
	RuneRange{0x0A2A, 0x0A30},
	RuneRange{0x0A32, 0x0A33},
	RuneRange{0x0A35, 0x0A36},
	RuneRange{0x0A38, 0x0A39},
	RuneRange{0x0A3E, 0x0A40},
	RuneRange{0x0A59, 0x0A5C},
	RuneRange{0x0A5E, 0x0A5E},
	RuneRange{0x0A66, 0x0A6F},
	RuneRange{0x0A72, 0x0A74},
	RuneRange{0x0A83, 0x0A83},
	RuneRange{0x0A85, 0x0A8B},
	RuneRange{0x0A8D, 0x0A8D},
	RuneRange{0x0A8F, 0x0A91},
	RuneRange{0x0A93, 0x0AA8},
	RuneRange{0x0AAA, 0x0AB0},
	RuneRange{0x0AB2, 0x0AB3},
	RuneRange{0x0AB5, 0x0AB9},
	RuneRange{0x0ABD, 0x0AC0},
	RuneRange{0x0AC9, 0x0AC9},
	RuneRange{0x0ACB, 0x0ACC},
	RuneRange{0x0AD0, 0x0AD0},
	RuneRange{0x0AE0, 0x0AE0},
	RuneRange{0x0AE6, 0x0AEF},
	RuneRange{0x0B02, 0x0B03},
	RuneRange{0x0B05, 0x0B0C},
	RuneRange{0x0B0F, 0x0B10},
	RuneRange{0x0B13, 0x0B28},
	RuneRange{0x0B2A, 0x0B30},
	RuneRange{0x0B32, 0x0B33},
	RuneRange{0x0B36, 0x0B39},
	RuneRange{0x0B3D, 0x0B3E},
	RuneRange{0x0B40, 0x0B40},
	RuneRange{0x0B47, 0x0B48},
	RuneRange{0x0B4B, 0x0B4C},
	RuneRange{0x0B57, 0x0B57},
	RuneRange{0x0B5C, 0x0B5D},
	RuneRange{0x0B5F, 0x0B61},
	RuneRange{0x0B66, 0x0B70},
	RuneRange{0x0B83, 0x0B83},
	RuneRange{0x0B85, 0x0B8A},
	RuneRange{0x0B8E, 0x0B90},
	RuneRange{0x0B92, 0x0B95},
	RuneRange{0x0B99, 0x0B9A},
	RuneRange{0x0B9C, 0x0B9C},
	RuneRange{0x0B9E, 0x0B9F},
	RuneRange{0x0BA3, 0x0BA4},
	RuneRange{0x0BA8, 0x0BAA},
	RuneRange{0x0BAE, 0x0BB5},
	RuneRange{0x0BB7, 0x0BB9},
	RuneRange{0x0BBE, 0x0BBF},
	RuneRange{0x0BC1, 0x0BC2},
	RuneRange{0x0BC6, 0x0BC8},
	RuneRange{0x0BCA, 0x0BCC},
	RuneRange{0x0BD7, 0x0BD7},
	RuneRange{0x0BE7, 0x0BF2},
	RuneRange{0x0C01, 0x0C03},
	RuneRange{0x0C05, 0x0C0C},
	RuneRange{0x0C0E, 0x0C10},
	RuneRange{0x0C12, 0x0C28},
	RuneRange{0x0C2A, 0x0C33},
	RuneRange{0x0C35, 0x0C39},
	RuneRange{0x0C41, 0x0C44},
	RuneRange{0x0C60, 0x0C61},
	RuneRange{0x0C66, 0x0C6F},
	RuneRange{0x0C82, 0x0C83},
	RuneRange{0x0C85, 0x0C8C},
	RuneRange{0x0C8E, 0x0C90},
	RuneRange{0x0C92, 0x0CA8},
	RuneRange{0x0CAA, 0x0CB3},
	RuneRange{0x0CB5, 0x0CB9},
	RuneRange{0x0CBE, 0x0CBE},
	RuneRange{0x0CC0, 0x0CC4},
	RuneRange{0x0CC7, 0x0CC8},
	RuneRange{0x0CCA, 0x0CCB},
	RuneRange{0x0CD5, 0x0CD6},
	RuneRange{0x0CDE, 0x0CDE},
	RuneRange{0x0CE0, 0x0CE1},
	RuneRange{0x0CE6, 0x0CEF},
	RuneRange{0x0D02, 0x0D03},
	RuneRange{0x0D05, 0x0D0C},
	RuneRange{0x0D0E, 0x0D10},
	RuneRange{0x0D12, 0x0D28},
	RuneRange{0x0D2A, 0x0D39},
	RuneRange{0x0D3E, 0x0D40},
	RuneRange{0x0D46, 0x0D48},
	RuneRange{0x0D4A, 0x0D4C},
	RuneRange{0x0D57, 0x0D57},
	RuneRange{0x0D60, 0x0D61},
	RuneRange{0x0D66, 0x0D6F},
	RuneRange{0x0D82, 0x0D83},
	RuneRange{0x0D85, 0x0D96},
	RuneRange{0x0D9A, 0x0DB1},
	RuneRange{0x0DB3, 0x0DBB},
	RuneRange{0x0DBD, 0x0DBD},
	RuneRange{0x0DC0, 0x0DC6},
	RuneRange{0x0DCF, 0x0DD1},
	RuneRange{0x0DD8, 0x0DDF},
	RuneRange{0x0DF2, 0x0DF4},
	RuneRange{0x0E01, 0x0E30},
	RuneRange{0x0E32, 0x0E33},
	RuneRange{0x0E40, 0x0E46},
	RuneRange{0x0E4F, 0x0E5B},
	RuneRange{0x0E81, 0x0E82},
	RuneRange{0x0E84, 0x0E84},
	RuneRange{0x0E87, 0x0E88},
	RuneRange{0x0E8A, 0x0E8A},
	RuneRange{0x0E8D, 0x0E8D},
	RuneRange{0x0E94, 0x0E97},
	RuneRange{0x0E99, 0x0E9F},
	RuneRange{0x0EA1, 0x0EA3},
	RuneRange{0x0EA5, 0x0EA5},
	RuneRange{0x0EA7, 0x0EA7},
	RuneRange{0x0EAA, 0x0EAB},
	RuneRange{0x0EAD, 0x0EB0},
	RuneRange{0x0EB2, 0x0EB3},
	RuneRange{0x0EBD, 0x0EBD},
	RuneRange{0x0EC0, 0x0EC4},
	RuneRange{0x0EC6, 0x0EC6},
	RuneRange{0x0ED0, 0x0ED9},
	RuneRange{0x0EDC, 0x0EDD},
	RuneRange{0x0F00, 0x0F17},
	RuneRange{0x0F1A, 0x0F34},
	RuneRange{0x0F36, 0x0F36},
	RuneRange{0x0F38, 0x0F38},
	RuneRange{0x0F3E, 0x0F47},
	RuneRange{0x0F49, 0x0F6A},
	RuneRange{0x0F7F, 0x0F7F},
	RuneRange{0x0F85, 0x0F85},
	RuneRange{0x0F88, 0x0F8B},
	RuneRange{0x0FBE, 0x0FC5},
	RuneRange{0x0FC7, 0x0FCC},
	RuneRange{0x0FCF, 0x0FCF},
	RuneRange{0x1000, 0x1021},
	RuneRange{0x1023, 0x1027},
	RuneRange{0x1029, 0x102A},
	RuneRange{0x102C, 0x102C},
	RuneRange{0x1031, 0x1031},
	RuneRange{0x1038, 0x1038},
	RuneRange{0x1040, 0x1057},
	RuneRange{0x10A0, 0x10C5},
	RuneRange{0x10D0, 0x10F8},
	RuneRange{0x10FB, 0x10FB},
	RuneRange{0x1100, 0x1159},
	RuneRange{0x115F, 0x11A2},
	RuneRange{0x11A8, 0x11F9},
	RuneRange{0x1200, 0x1206},
	RuneRange{0x1208, 0x1246},
	RuneRange{0x1248, 0x1248},
	RuneRange{0x124A, 0x124D},
	RuneRange{0x1250, 0x1256},
	RuneRange{0x1258, 0x1258},
	RuneRange{0x125A, 0x125D},
	RuneRange{0x1260, 0x1286},
	RuneRange{0x1288, 0x1288},
	RuneRange{0x128A, 0x128D},
	RuneRange{0x1290, 0x12AE},
	RuneRange{0x12B0, 0x12B0},
	RuneRange{0x12B2, 0x12B5},
	RuneRange{0x12B8, 0x12BE},
	RuneRange{0x12C0, 0x12C0},
	RuneRange{0x12C2, 0x12C5},
	RuneRange{0x12C8, 0x12CE},
	RuneRange{0x12D0, 0x12D6},
	RuneRange{0x12D8, 0x12EE},
	RuneRange{0x12F0, 0x130E},
	RuneRange{0x1310, 0x1310},
	RuneRange{0x1312, 0x1315},
	RuneRange{0x1318, 0x131E},
	RuneRange{0x1320, 0x1346},
	RuneRange{0x1348, 0x135A},
	RuneRange{0x1361, 0x137C},
	RuneRange{0x13A0, 0x13F4},
	RuneRange{0x1401, 0x1676},
	RuneRange{0x1681, 0x169A},
	RuneRange{0x16A0, 0x16F0},
	RuneRange{0x1700, 0x170C},
	RuneRange{0x170E, 0x1711},
	RuneRange{0x1720, 0x1731},
	RuneRange{0x1735, 0x1736},
	RuneRange{0x1740, 0x1751},
	RuneRange{0x1760, 0x176C},
	RuneRange{0x176E, 0x1770},
	RuneRange{0x1780, 0x17B6},
	RuneRange{0x17BE, 0x17C5},
	RuneRange{0x17C7, 0x17C8},
	RuneRange{0x17D4, 0x17DA},
	RuneRange{0x17DC, 0x17DC},
	RuneRange{0x17E0, 0x17E9},
	RuneRange{0x1810, 0x1819},
	RuneRange{0x1820, 0x1877},
	RuneRange{0x1880, 0x18A8},
	RuneRange{0x1E00, 0x1E9B},
	RuneRange{0x1EA0, 0x1EF9},
	RuneRange{0x1F00, 0x1F15},
	RuneRange{0x1F18, 0x1F1D},
	RuneRange{0x1F20, 0x1F45},
	RuneRange{0x1F48, 0x1F4D},
	RuneRange{0x1F50, 0x1F57},
	RuneRange{0x1F59, 0x1F59},
	RuneRange{0x1F5B, 0x1F5B},
	RuneRange{0x1F5D, 0x1F5D},
	RuneRange{0x1F5F, 0x1F7D},
	RuneRange{0x1F80, 0x1FB4},
	RuneRange{0x1FB6, 0x1FBC},
	RuneRange{0x1FBE, 0x1FBE},
	RuneRange{0x1FC2, 0x1FC4},
	RuneRange{0x1FC6, 0x1FCC},
	RuneRange{0x1FD0, 0x1FD3},
	RuneRange{0x1FD6, 0x1FDB},
	RuneRange{0x1FE0, 0x1FEC},
	RuneRange{0x1FF2, 0x1FF4},
	RuneRange{0x1FF6, 0x1FFC},
	RuneRange{0x200E, 0x200E},
	RuneRange{0x2071, 0x2071},
	RuneRange{0x207F, 0x207F},
	RuneRange{0x2102, 0x2102},
	RuneRange{0x2107, 0x2107},
	RuneRange{0x210A, 0x2113},
	RuneRange{0x2115, 0x2115},
	RuneRange{0x2119, 0x211D},
	RuneRange{0x2124, 0x2124},
	RuneRange{0x2126, 0x2126},
	RuneRange{0x2128, 0x2128},
	RuneRange{0x212A, 0x212D},
	RuneRange{0x212F, 0x2131},
	RuneRange{0x2133, 0x2139},
	RuneRange{0x213D, 0x213F},
	RuneRange{0x2145, 0x2149},
	RuneRange{0x2160, 0x2183},
	RuneRange{0x2336, 0x237A},
	RuneRange{0x2395, 0x2395},
	RuneRange{0x249C, 0x24E9},
	RuneRange{0x3005, 0x3007},
	RuneRange{0x3021, 0x3029},
	RuneRange{0x3031, 0x3035},
	RuneRange{0x3038, 0x303C},
	RuneRange{0x3041, 0x3096},
	RuneRange{0x309D, 0x309F},
	RuneRange{0x30A1, 0x30FA},
	RuneRange{0x30FC, 0x30FF},
	RuneRange{0x3105, 0x312C},
	RuneRange{0x3131, 0x318E},
	RuneRange{0x3190, 0x31B7},
	RuneRange{0x31F0, 0x321C},
	RuneRange{0x3220, 0x3243},
	RuneRange{0x3260, 0x327B},
	RuneRange{0x327F, 0x32B0},
	RuneRange{0x32C0, 0x32CB},
	RuneRange{0x32D0, 0x32FE},
	RuneRange{0x3300, 0x3376},
	RuneRange{0x337B, 0x33DD},
	RuneRange{0x33E0, 0x33FE},
	RuneRange{0x3400, 0x4DB5},
	RuneRange{0x4E00, 0x9FA5},
	RuneRange{0xA000, 0xA48C},
	RuneRange{0xAC00, 0xD7A3},
	RuneRange{0xD800, 0xFA2D},
	RuneRange{0xFA30, 0xFA6A},
	RuneRange{0xFB00, 0xFB06},
	RuneRange{0xFB13, 0xFB17},
	RuneRange{0xFF21, 0xFF3A},
	RuneRange{0xFF41, 0xFF5A},
	RuneRange{0xFF66, 0xFFBE},
	RuneRange{0xFFC2, 0xFFC7},
	RuneRange{0xFFCA, 0xFFCF},
	RuneRange{0xFFD2, 0xFFD7},
	RuneRange{0xFFDA, 0xFFDC},
	RuneRange{0x10300, 0x1031E},
	RuneRange{0x10320, 0x10323},
	RuneRange{0x10330, 0x1034A},
	RuneRange{0x10400, 0x10425},
	RuneRange{0x10428, 0x1044D},
	RuneRange{0x1D000, 0x1D0F5},
	RuneRange{0x1D100, 0x1D126},
	RuneRange{0x1D12A, 0x1D166},
	RuneRange{0x1D16A, 0x1D172},
	RuneRange{0x1D183, 0x1D184},
	RuneRange{0x1D18C, 0x1D1A9},
	RuneRange{0x1D1AE, 0x1D1DD},
	RuneRange{0x1D400, 0x1D454},
	RuneRange{0x1D456, 0x1D49C},
	RuneRange{0x1D49E, 0x1D49F},
	RuneRange{0x1D4A2, 0x1D4A2},
	RuneRange{0x1D4A5, 0x1D4A6},
	RuneRange{0x1D4A9, 0x1D4AC},
	RuneRange{0x1D4AE, 0x1D4B9},
	RuneRange{0x1D4BB, 0x1D4BB},
	RuneRange{0x1D4BD, 0x1D4C0},
	RuneRange{0x1D4C2, 0x1D4C3},
	RuneRange{0x1D4C5, 0x1D505},
	RuneRange{0x1D507, 0x1D50A},
	RuneRange{0x1D50D, 0x1D514},
	RuneRange{0x1D516, 0x1D51C},
	RuneRange{0x1D51E, 0x1D539},
	RuneRange{0x1D53B, 0x1D53E},
	RuneRange{0x1D540, 0x1D544},
	RuneRange{0x1D546, 0x1D546},
	RuneRange{0x1D54A, 0x1D550},
	RuneRange{0x1D552, 0x1D6A3},
	RuneRange{0x1D6A8, 0x1D7C9},
	RuneRange{0x20000, 0x2A6D6},
	RuneRange{0x2F800, 0x2FA1D},
	RuneRange{0xF0000, 0xFFFFD},
	RuneRange{0x100000, 0x10FFFD},
}

// TableD2 represents RFC-3454 Table D.2.
var TableD2 Set = tableD2
