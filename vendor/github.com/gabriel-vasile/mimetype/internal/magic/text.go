package magic

import (
	"bytes"

	"github.com/gabriel-vasile/mimetype/internal/charset"
	"github.com/gabriel-vasile/mimetype/internal/json"
)

var (
	// Html matches a Hypertext Markup Language file.
	Html = markup(
		[]byte("<!DOCTYPE HTML"),
		[]byte("<HTML"),
		[]byte("<HEAD"),
		[]byte("<SCRIPT"),
		[]byte("<IFRAME"),
		[]byte("<H1"),
		[]byte("<DIV"),
		[]byte("<FONT"),
		[]byte("<TABLE"),
		[]byte("<A"),
		[]byte("<STYLE"),
		[]byte("<TITLE"),
		[]byte("<B"),
		[]byte("<BODY"),
		[]byte("<BR"),
		[]byte("<P"),
		[]byte("<!--"),
	)
	// Xml matches an Extensible Markup Language file.
	Xml = markup([]byte("<?XML"))
	// Owl2 matches an Owl ontology file.
	Owl2 = xml(newXmlSig("Ontology", `xmlns="http://www.w3.org/2002/07/owl#"`))
	// Rss matches a Rich Site Summary file.
	Rss = xml(newXmlSig("rss", ""))
	// Atom matches an Atom Syndication Format file.
	Atom = xml(newXmlSig("feed", `xmlns="http://www.w3.org/2005/Atom"`))
	// Kml matches a Keyhole Markup Language file.
	Kml = xml(
		newXmlSig("kml", `xmlns="http://www.opengis.net/kml/2.2"`),
		newXmlSig("kml", `xmlns="http://earth.google.com/kml/2.0"`),
		newXmlSig("kml", `xmlns="http://earth.google.com/kml/2.1"`),
		newXmlSig("kml", `xmlns="http://earth.google.com/kml/2.2"`),
	)
	// Xliff matches a XML Localization Interchange File Format file.
	Xliff = xml(newXmlSig("xliff", `xmlns="urn:oasis:names:tc:xliff:document:1.2"`))
	// Collada matches a COLLAborative Design Activity file.
	Collada = xml(newXmlSig("COLLADA", `xmlns="http://www.collada.org/2005/11/COLLADASchema"`))
	// Gml matches a Geography Markup Language file.
	Gml = xml(
		newXmlSig("", `xmlns:gml="http://www.opengis.net/gml"`),
		newXmlSig("", `xmlns:gml="http://www.opengis.net/gml/3.2"`),
		newXmlSig("", `xmlns:gml="http://www.opengis.net/gml/3.3/exr"`),
	)
	// Gpx matches a GPS Exchange Format file.
	Gpx = xml(newXmlSig("gpx", `xmlns="http://www.topografix.com/GPX/1/1"`))
	// Tcx matches a Training Center XML file.
	Tcx = xml(newXmlSig("TrainingCenterDatabase", `xmlns="http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2"`))
	// X3d matches an Extensible 3D Graphics file.
	X3d = xml(newXmlSig("X3D", `xmlns:xsd="http://www.w3.org/2001/XMLSchema-instance"`))
	// Amf matches an Additive Manufacturing XML file.
	Amf = xml(newXmlSig("amf", ""))
	// Threemf matches a 3D Manufacturing Format file.
	Threemf = xml(newXmlSig("model", `xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02"`))
	// Xfdf matches a XML Forms Data Format file.
	Xfdf = xml(newXmlSig("xfdf", `xmlns="http://ns.adobe.com/xfdf/"`))
	// VCard matches a Virtual Contact File.
	VCard = ciPrefix([]byte("BEGIN:VCARD\n"), []byte("BEGIN:VCARD\r\n"))
	// ICalendar matches a iCalendar file.
	ICalendar = ciPrefix([]byte("BEGIN:VCALENDAR\n"), []byte("BEGIN:VCALENDAR\r\n"))
	phpPageF  = ciPrefix(
		[]byte("<?PHP"),
		[]byte("<?\n"),
		[]byte("<?\r"),
		[]byte("<? "),
	)
	phpScriptF = shebang(
		[]byte("/usr/local/bin/php"),
		[]byte("/usr/bin/php"),
		[]byte("/usr/bin/env php"),
	)
	// Js matches a Javascript file.
	Js = shebang(
		[]byte("/bin/node"),
		[]byte("/usr/bin/node"),
		[]byte("/bin/nodejs"),
		[]byte("/usr/bin/nodejs"),
		[]byte("/usr/bin/env node"),
		[]byte("/usr/bin/env nodejs"),
	)
	// Lua matches a Lua programming language file.
	Lua = shebang(
		[]byte("/usr/bin/lua"),
		[]byte("/usr/local/bin/lua"),
		[]byte("/usr/bin/env lua"),
	)
	// Perl matches a Perl programming language file.
	Perl = shebang(
		[]byte("/usr/bin/perl"),
		[]byte("/usr/bin/env perl"),
	)
	// Python matches a Python programming language file.
	Python = shebang(
		[]byte("/usr/bin/python"),
		[]byte("/usr/local/bin/python"),
		[]byte("/usr/bin/env python"),
	)
	// Tcl matches a Tcl programming language file.
	Tcl = shebang(
		[]byte("/usr/bin/tcl"),
		[]byte("/usr/local/bin/tcl"),
		[]byte("/usr/bin/env tcl"),
		[]byte("/usr/bin/tclsh"),
		[]byte("/usr/local/bin/tclsh"),
		[]byte("/usr/bin/env tclsh"),
		[]byte("/usr/bin/wish"),
		[]byte("/usr/local/bin/wish"),
		[]byte("/usr/bin/env wish"),
	)
	// Rtf matches a Rich Text Format file.
	Rtf = prefix([]byte("{\\rtf1"))
)

// Text matches a plain text file.
//
// TODO: This function does not parse BOM-less UTF16 and UTF32 files. Not really
// sure it should. Linux file utility also requires a BOM for UTF16 and UTF32.
func Text(raw []byte, limit uint32) bool {
	// First look for BOM.
	if cset := charset.FromBOM(raw); cset != "" {
		return true
	}
	return isText(raw)
}

// Php matches a PHP: Hypertext Preprocessor file.
func Php(raw []byte, limit uint32) bool {
	if res := phpPageF(raw, limit); res {
		return res
	}
	return phpScriptF(raw, limit)
}

// Json matches a JavaScript Object Notation file.
func Json(raw []byte, limit uint32) bool {
	parsed, err := json.Scan(raw)
	if len(raw) < int(limit) {
		return err == nil
	}

	return parsed == len(raw)
}

// GeoJson matches a RFC 7946 GeoJSON file.
//
// GeoJson detection implies searching for key:value pairs like: `"type": "Feature"`
// in the input.
// BUG(gabriel-vasile): The "type" key should be searched for in the root object.
func GeoJson(raw []byte, limit uint32) bool {
	raw = trimLWS(raw)
	if len(raw) == 0 {
		return false
	}
	// GeoJSON is always a JSON object, not a JSON array.
	if raw[0] != '{' {
		return false
	}

	s := []byte(`"type"`)
	si, sl := bytes.Index(raw, s), len(s)

	if si == -1 {
		return false
	}

	// If the "type" string is the suffix of the input,
	// there is no need to search for the value of the key.
	if si+sl == len(raw) {
		return false
	}
	// Skip the "type" part.
	raw = raw[si+sl:]
	// Skip any whitespace before the colon.
	raw = trimLWS(raw)
	// Check for colon.
	if len(raw) == 0 || raw[0] != ':' {
		return false
	}
	// Skip any whitespace after the colon.
	raw = trimLWS(raw[1:])

	geoJsonTypes := [][]byte{
		[]byte(`"Feature"`),
		[]byte(`"FeatureCollection"`),
		[]byte(`"Point"`),
		[]byte(`"LineString"`),
		[]byte(`"Polygon"`),
		[]byte(`"MultiPoint"`),
		[]byte(`"MultiLineString"`),
		[]byte(`"MultiPolygon"`),
		[]byte(`"GeometryCollection"`),
	}
	for _, t := range geoJsonTypes {
		if bytes.HasPrefix(raw, t) {
			return true
		}
	}

	return false
}

// NdJson matches a Newline delimited JSON file.
func NdJson(raw []byte, limit uint32) bool {
	// Separator with carriage return and new line `\r\n`.
	srn := []byte{0x0D, 0x0A}

	// Separator with only new line `\n`.
	sn := []byte{0x0A}

	// Total bytes scanned.
	parsed := 0

	// Split by `srn`.
	for rni, insrn := range bytes.Split(raw, srn) {
		// Separator byte count should be added only after the first split.
		if rni != 0 {
			// Add two as `\r\n` is used for split.
			parsed += 2
		}
		// Split again by `sn`.
		for ni, insn := range bytes.Split(insrn, sn) {
			// Separator byte count should be added only after the first split.
			if ni != 0 {
				// Add one as `\n` is used for split.
				parsed++
			}
			// Empty line is valid.
			if len(insn) == 0 {
				continue
			}
			p, err := json.Scan(insn)
			parsed += p
			if parsed < int(limit) && err != nil {
				return false
			}
		}
	}

	// Empty inputs should not pass as valid NDJSON with 0 lines.
	return parsed > 2 && parsed == len(raw)
}

// Svg matches a SVG file.
func Svg(raw []byte, limit uint32) bool {
	return bytes.Contains(raw, []byte("<svg"))
}

// isText considers any file containing null bytes as a binary file.
// There is plenty room for disagreement regarding what should be considered a
// text file. This approach is used by diff, cat, and other linux utilities.
func isText(raw []byte) bool {
	l := 8000
	if len(raw) > l {
		raw = raw[:l]
	}
	return bytes.IndexByte(raw, 0) == -1
}
