package abstractions

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

// MultipartBody represents a multipart body for a request or a response.
type MultipartBody interface {
	serialization.Parsable
	// AddOrReplacePart adds or replaces a part in the multipart body.
	AddOrReplacePart(name string, contentType string, content any) error
	// GetPartValue gets the value of a part in the multipart body.
	GetPartValue(name string) (any, error)
	// RemovePart removes a part from the multipart body.
	RemovePart(name string) error
	// SetRequestAdapter sets the request adapter to use for serialization.
	SetRequestAdapter(requestAdapter RequestAdapter)
	// GetRequestAdapter gets the request adapter to use for serialization.
	GetRequestAdapter() RequestAdapter
	// GetBoundary returns the boundary used in the multipart body.
	GetBoundary() string
}
type multipartBody struct {
	parts            map[string]multipartEntry
	originalNamesMap map[string]string
	boundary         string
	requestAdapter   RequestAdapter
}

func NewMultipartBody() MultipartBody {
	return &multipartBody{
		parts:            make(map[string]multipartEntry),
		originalNamesMap: make(map[string]string),
		boundary:         strings.ReplaceAll(uuid.New().String(), "-", ""),
	}
}
func normalizePartName(original string) string {
	return strings.ToLower(original)
}
func stringReference(original string) *string {
	return &original
}

// AddOrReplacePart adds or replaces a part in the multipart body.
func (m *multipartBody) AddOrReplacePart(name string, contentType string, content any) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if contentType == "" {
		return errors.New("contentType cannot be empty")
	}
	if content == nil {
		return errors.New("content cannot be nil")
	}
	normalizedName := normalizePartName(name)
	m.parts[normalizedName] = multipartEntry{
		ContentType: contentType,
		Content:     content,
	}
	m.originalNamesMap[normalizedName] = name

	return nil
}

// GetPartValue gets the value of a part in the multipart body.
func (m *multipartBody) GetPartValue(name string) (any, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	normalizedName := normalizePartName(name)
	if part, ok := m.parts[normalizedName]; ok {
		return part.Content, nil
	}
	return nil, nil
}

// RemovePart removes a part from the multipart body.
func (m *multipartBody) RemovePart(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	normalizedName := normalizePartName(name)
	delete(m.parts, normalizedName)
	delete(m.originalNamesMap, normalizedName)
	return nil
}

// Serialize writes the objects properties to the current writer.
func (m *multipartBody) Serialize(writer serialization.SerializationWriter) error {
	if writer == nil {
		return errors.New("writer cannot be nil")
	}
	if m.requestAdapter == nil {
		return errors.New("requestAdapter cannot be nil")
	}
	serializationWriterFactory := m.requestAdapter.GetSerializationWriterFactory()
	if serializationWriterFactory == nil {
		return errors.New("serializationWriterFactory cannot be nil")
	}
	if len(m.parts) == 0 {
		return errors.New("no parts to serialize")
	}

	first := true
	for partName, part := range m.parts {
		if first {
			first = false
		} else {
			if err := writer.WriteStringValue("", stringReference("")); err != nil {
				return err
			}
		}
		if err := writer.WriteStringValue("", stringReference("--"+m.boundary)); err != nil {
			return err
		}
		if err := writer.WriteStringValue("Content-Type", stringReference(part.ContentType)); err != nil {
			return err
		}
		partOriginalName := m.originalNamesMap[partName]
		if err := writer.WriteStringValue("Content-Disposition", stringReference("form-data; name=\""+partOriginalName+"\"")); err != nil {
			return err
		}
		if err := writer.WriteStringValue("", stringReference("")); err != nil {
			return err
		}
		if parsable, ok := part.Content.(serialization.Parsable); ok {
			partWriter, error := serializationWriterFactory.GetSerializationWriter(part.ContentType)
			defer partWriter.Close()
			if error != nil {
				return error
			}
			if error = partWriter.WriteObjectValue("", parsable); error != nil {
				return error
			}
			partContent, error := partWriter.GetSerializedContent()
			if error != nil {
				return error
			}
			if error = writer.WriteByteArrayValue("", partContent); error != nil {
				return error
			}
		} else if str, ok := part.Content.(string); ok {
			if error := writer.WriteStringValue("", stringReference(str)); error != nil {
				return error
			}
		} else if byteArray, ok := part.Content.([]byte); ok {
			if error := writer.WriteByteArrayValue("", byteArray); error != nil {
				return error
			}
		} else {
			return errors.New("unsupported part type")
		}
	}
	if err := writer.WriteStringValue("", stringReference("")); err != nil {
		return err
	}
	if err := writer.WriteStringValue("", stringReference("--"+m.boundary+"--")); err != nil {
		return err
	}

	return nil
}

// GetFieldDeserializers returns the deserialization information for this object.
func (m *multipartBody) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	panic("not implemented")
}

// GetRequestAdapter gets the request adapter to use for serialization.
func (m *multipartBody) GetRequestAdapter() RequestAdapter {
	return m.requestAdapter
}

// SetRequestAdapter sets the request adapter to use for serialization.
func (m *multipartBody) SetRequestAdapter(requestAdapter RequestAdapter) {
	m.requestAdapter = requestAdapter
}

// GetBoundary returns the boundary used in the multipart body.
func (m *multipartBody) GetBoundary() string {
	return m.boundary
}

type multipartEntry struct {
	ContentType string
	Content     any
}
