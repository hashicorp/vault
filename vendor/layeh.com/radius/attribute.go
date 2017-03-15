package radius // import "layeh.com/radius"

// Attribute is a RADIUS attribute, which is part of a RADIUS packet.
type Attribute struct {
	Type  byte
	Value interface{}
}

// AttributeCodec defines how an Attribute is encoded and decoded to and from
// wire data.
type AttributeCodec interface {
	// Note: do not store wire; make a copy of it.
	Decode(packet *Packet, wire []byte) (interface{}, error)
	Encode(packet *Packet, value interface{}) ([]byte, error)
}
