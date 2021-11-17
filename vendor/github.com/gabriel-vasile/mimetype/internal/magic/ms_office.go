package magic

import (
	"bytes"
	"encoding/binary"
	"strings"
)

var (
	xlsxSigFiles = []string{
		"xl/worksheets/",
		"xl/drawings/",
		"xl/theme/",
		"xl/_rels/",
		"xl/styles.xml",
		"xl/workbook.xml",
		"xl/sharedStrings.xml",
	}
	docxSigFiles = []string{
		"word/media/",
		"word/_rels/document.xml.rels",
		"word/document.xml",
		"word/styles.xml",
		"word/fontTable.xml",
		"word/settings.xml",
		"word/numbering.xml",
		"word/header",
		"word/footer",
	}
	pptxSigFiles = []string{
		"ppt/slides/",
		"ppt/media/",
		"ppt/slideLayouts/",
		"ppt/theme/",
		"ppt/slideMasters/",
		"ppt/tags/",
		"ppt/notesMasters/",
		"ppt/_rels/",
		"ppt/handoutMasters/",
		"ppt/notesSlides/",
		"ppt/presentation.xml",
		"ppt/tableStyles.xml",
		"ppt/presProps.xml",
		"ppt/viewProps.xml",
	}
)

// zipTokenizer holds the source zip file and scanned index.
type zipTokenizer struct {
	in []byte
	i  int // current index
}

// next returns the next file name from the zip headers.
// https://web.archive.org/web/20191129114319/https://users.cs.jmu.edu/buchhofp/forensics/formats/pkzip.html
func (t *zipTokenizer) next() (fileName string) {
	if t.i > len(t.in) {
		return
	}
	in := t.in[t.i:]
	// pkSig is the signature of the zip local file header.
	pkSig := []byte("PK\003\004")
	pkIndex := bytes.Index(in, pkSig)
	// 30 is the offset of the file name in the header.
	fNameOffset := pkIndex + 30
	// end if signature not found or file name offset outside of file.
	if pkIndex == -1 || fNameOffset > len(in) {
		return
	}

	fNameLen := int(binary.LittleEndian.Uint16(in[pkIndex+26 : pkIndex+28]))
	if fNameLen <= 0 || fNameOffset+fNameLen > len(in) {
		return
	}
	t.i += fNameOffset + fNameLen
	return string(in[fNameOffset : fNameOffset+fNameLen])
}

// msoXML reads at most first 10 local headers and returns whether the input
// looks like a Microsoft Office file.
func msoXML(in []byte, prefixes ...string) bool {
	t := zipTokenizer{in: in}
	for i, tok := 0, t.next(); i < 10 && tok != ""; i, tok = i+1, t.next() {
		for p := range prefixes {
			if strings.HasPrefix(tok, prefixes[p]) {
				return true
			}
		}
	}

	return false
}

// Xlsx matches a Microsoft Excel 2007 file.
func Xlsx(raw []byte, limit uint32) bool {
	return msoXML(raw, xlsxSigFiles...)
}

// Docx matches a Microsoft Word 2007 file.
func Docx(raw []byte, limit uint32) bool {
	return msoXML(raw, docxSigFiles...)
}

// Pptx matches a Microsoft PowerPoint 2007 file.
func Pptx(raw []byte, limit uint32) bool {
	return msoXML(raw, pptxSigFiles...)
}

// Ole matches an Open Linking and Embedding file.
//
// https://en.wikipedia.org/wiki/Object_Linking_and_Embedding
func Ole(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1})
}

// Aaf matches an Advanced Authoring Format file.
// See: https://pyaaf.readthedocs.io/en/latest/about.html
// See: https://en.wikipedia.org/wiki/Advanced_Authoring_Format
func Aaf(raw []byte, limit uint32) bool {
	if len(raw) < 31 {
		return false
	}
	return bytes.HasPrefix(raw[8:], []byte{0x41, 0x41, 0x46, 0x42, 0x0D, 0x00, 0x4F, 0x4D}) &&
		(raw[30] == 0x09 || raw[30] == 0x0C)
}

// Doc matches a Microsoft Word 97-2003 file.
//
// BUG(gabriel-vasile): Doc should look for subheaders like Ppt and Xls does.
//
// Ole is a container for Doc, Ppt, Pub and Xls.
// Right now, when an Ole file is detected, it is considered to be a Doc file
// if the checks for Ppt, Pub and Xls failed.
func Doc(raw []byte, limit uint32) bool {
	return true
}

// Ppt matches a Microsoft PowerPoint 97-2003 file or a PowerPoint 95 presentation.
func Ppt(raw []byte, limit uint32) bool {
	// Root CLSID test is the safest way to detect identify OLE, however, the format
	// often places the root CLSID at the end of the file.
	if matchOleClsid(raw, []byte{
		0x10, 0x8d, 0x81, 0x64, 0x9b, 0x4f, 0xcf, 0x11,
		0x86, 0xea, 0x00, 0xaa, 0x00, 0xb9, 0x29, 0xe8,
	}) || matchOleClsid(raw, []byte{
		0x70, 0xae, 0x7b, 0xea, 0x3b, 0xfb, 0xcd, 0x11,
		0xa9, 0x03, 0x00, 0xaa, 0x00, 0x51, 0x0e, 0xa3,
	}) {
		return true
	}

	lin := len(raw)
	if lin < 520 {
		return false
	}
	pptSubHeaders := [][]byte{
		{0xA0, 0x46, 0x1D, 0xF0},
		{0x00, 0x6E, 0x1E, 0xF0},
		{0x0F, 0x00, 0xE8, 0x03},
	}
	for _, h := range pptSubHeaders {
		if bytes.HasPrefix(raw[512:], h) {
			return true
		}
	}

	if bytes.HasPrefix(raw[512:], []byte{0xFD, 0xFF, 0xFF, 0xFF}) &&
		raw[518] == 0x00 && raw[519] == 0x00 {
		return true
	}

	return lin > 1152 && bytes.Contains(raw[1152:min(4096, lin)],
		[]byte("P\x00o\x00w\x00e\x00r\x00P\x00o\x00i\x00n\x00t\x00 D\x00o\x00c\x00u\x00m\x00e\x00n\x00t"))
}

// Xls matches a Microsoft Excel 97-2003 file.
func Xls(raw []byte, limit uint32) bool {
	// Root CLSID test is the safest way to detect identify OLE, however, the format
	// often places the root CLSID at the end of the file.
	if matchOleClsid(raw, []byte{
		0x10, 0x08, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
	}) || matchOleClsid(raw, []byte{
		0x20, 0x08, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
	}) {
		return true
	}

	lin := len(raw)
	if lin < 520 {
		return false
	}
	xlsSubHeaders := [][]byte{
		{0x09, 0x08, 0x10, 0x00, 0x00, 0x06, 0x05, 0x00},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x10},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x1F},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x22},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x23},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x28},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x29},
	}
	for _, h := range xlsSubHeaders {
		if bytes.HasPrefix(raw[512:], h) {
			return true
		}
	}

	return lin > 1152 && bytes.Contains(raw[1152:min(4096, lin)],
		[]byte("W\x00k\x00s\x00S\x00S\x00W\x00o\x00r\x00k\x00B\x00o\x00o\x00k"))
}

// Pub matches a Microsoft Publisher file.
func Pub(raw []byte, limit uint32) bool {
	return matchOleClsid(raw, []byte{
		0x01, 0x12, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46,
	})
}

// Msg matches a Microsoft Outlook email file.
func Msg(raw []byte, limit uint32) bool {
	return matchOleClsid(raw, []byte{
		0x0B, 0x0D, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46,
	})
}

// Helper to match by a specific CLSID of a compound file.
//
// http://fileformats.archiveteam.org/wiki/Microsoft_Compound_File
func matchOleClsid(in []byte, clsid []byte) bool {
	if len(in) <= 512 {
		return false
	}

	// SecID of first sector of the directory stream
	firstSecID := int(binary.LittleEndian.Uint32(in[48:52]))

	// Expected offset of CLSID for root storage object
	clsidOffset := 512*(1+firstSecID) + 80

	if len(in) <= clsidOffset+16 {
		return false
	}

	return bytes.HasPrefix(in[clsidOffset:], clsid)
}
