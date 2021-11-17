package magic

// Torrent has bencoded text in the beginning.
var Torrent = prefix([]byte("d8:announce"))
