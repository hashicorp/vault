package magic

var (
	// Mp4 matches an MP4 file.
	Mp4 = ftyp(
		[]byte("avc1"), []byte("dash"), []byte("iso2"), []byte("iso3"),
		[]byte("iso4"), []byte("iso5"), []byte("iso6"), []byte("isom"),
		[]byte("mmp4"), []byte("mp41"), []byte("mp42"), []byte("mp4v"),
		[]byte("mp71"), []byte("MSNV"), []byte("NDAS"), []byte("NDSC"),
		[]byte("NSDC"), []byte("NSDH"), []byte("NDSM"), []byte("NDSP"),
		[]byte("NDSS"), []byte("NDXC"), []byte("NDXH"), []byte("NDXM"),
		[]byte("NDXP"), []byte("NDXS"), []byte("F4V "), []byte("F4P "),
	)
	// ThreeGP matches a 3GPP file.
	ThreeGP = ftyp(
		[]byte("3gp1"), []byte("3gp2"), []byte("3gp3"), []byte("3gp4"),
		[]byte("3gp5"), []byte("3gp6"), []byte("3gp7"), []byte("3gs7"),
		[]byte("3ge6"), []byte("3ge7"), []byte("3gg6"),
	)
	// ThreeG2 matches a 3GPP2 file.
	ThreeG2 = ftyp(
		[]byte("3g24"), []byte("3g25"), []byte("3g26"), []byte("3g2a"),
		[]byte("3g2b"), []byte("3g2c"), []byte("KDDI"),
	)
	// AMp4 matches an audio MP4 file.
	AMp4 = ftyp(
		// audio for Adobe Flash Player 9+
		[]byte("F4A "), []byte("F4B "),
		// Apple iTunes AAC-LC (.M4A) Audio
		[]byte("M4B "), []byte("M4P "),
		// MPEG-4 (.MP4) for SonyPSP
		[]byte("MSNV"),
		// Nero Digital AAC Audio
		[]byte("NDAS"),
	)
	// QuickTime matches a QuickTime File Format file.
	QuickTime = ftyp([]byte("qt  "), []byte("moov"))
	// Mqv matches a Sony / Mobile QuickTime  file.
	Mqv = ftyp([]byte("mqt "))
	// M4a matches an audio M4A file.
	M4a = ftyp([]byte("M4A "))
	// M4v matches an Appl4 M4V video file.
	M4v = ftyp([]byte("M4V "), []byte("M4VH"), []byte("M4VP"))
	// Heic matches a High Efficiency Image Coding (HEIC) file.
	Heic = ftyp([]byte("heic"), []byte("heix"))
	// HeicSequence matches a High Efficiency Image Coding (HEIC) file sequence.
	HeicSequence = ftyp([]byte("hevc"), []byte("hevx"))
	// Heif matches a High Efficiency Image File Format (HEIF) file.
	Heif = ftyp([]byte("mif1"), []byte("heim"), []byte("heis"), []byte("avic"))
	// HeifSequence matches a High Efficiency Image File Format (HEIF) file sequence.
	HeifSequence = ftyp([]byte("msf1"), []byte("hevm"), []byte("hevs"), []byte("avcs"))
	// TODO: add support for remaining video formats at ftyps.com.
)
