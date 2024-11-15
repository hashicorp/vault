package serialization

// ParsableAction Encapsulates a method with a single Parsable parameter
type ParsableAction func(Parsable) error

// ParsableWriter  Encapsulates a method that receives a Parsable and SerializationWriter as parameters
type ParsableWriter func(Parsable, SerializationWriter) error

// SerializationWriterProxyFactory factory that allows the composition of before and after callbacks on existing factories.
type SerializationWriterProxyFactory struct {
	factory              SerializationWriterFactory
	onBeforeAction       ParsableAction
	onAfterAction        ParsableAction
	onSerializationStart ParsableWriter
}

// NewSerializationWriterProxyFactory constructs a new instance of SerializationWriterProxyFactory
func NewSerializationWriterProxyFactory(
	factory SerializationWriterFactory,
	onBeforeAction ParsableAction,
	onAfterAction ParsableAction,
	onSerializationStart ParsableWriter,
) *SerializationWriterProxyFactory {
	return &SerializationWriterProxyFactory{
		factory:              factory,
		onBeforeAction:       onBeforeAction,
		onAfterAction:        onAfterAction,
		onSerializationStart: onSerializationStart,
	}
}

func (s *SerializationWriterProxyFactory) GetValidContentType() (string, error) {
	return s.factory.GetValidContentType()
}

func (s *SerializationWriterProxyFactory) GetSerializationWriter(contentType string) (SerializationWriter, error) {
	writer, err := s.factory.GetSerializationWriter(contentType)
	if err != nil {
		return nil, err
	}

	originalBefore := writer.GetOnBeforeSerialization()
	err = writer.SetOnBeforeSerialization(func(parsable Parsable) error {
		if s != nil {
			err := s.onBeforeAction(parsable)
			if err != nil {
				return err
			}
		}
		if originalBefore != nil {
			err := originalBefore(parsable)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	originalAfter := writer.GetOnAfterObjectSerialization()
	err = writer.SetOnAfterObjectSerialization(func(parsable Parsable) error {
		if s != nil {
			err := s.onAfterAction(parsable)
			if err != nil {
				return err
			}
		}
		if originalAfter != nil {
			err := originalAfter(parsable)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	originalStart := writer.GetOnStartObjectSerialization()
	err = writer.SetOnStartObjectSerialization(func(parsable Parsable, writer SerializationWriter) error {
		if s != nil {
			err := s.onSerializationStart(parsable, writer)
			if err != nil {
				return err
			}
		}
		if originalBefore != nil {
			err := originalStart(parsable, writer)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return writer, nil
}
