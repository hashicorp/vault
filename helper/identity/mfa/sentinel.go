package mfa

func (c *Config) SentinelGet(key string) (interface{}, error) {
	if c == nil {
		return nil, nil
	}
	switch key {
	case "type":
		return c.Type, nil
	case "name":
		return c.Name, nil
	case "mount_accessor":
		return c.MountAccessor, nil
	}

	return nil, nil
}

func (c *Config) SentinelKeys() []string {
	return []string{
		"type",
		"name",
		"mount_accessor",
	}
}
