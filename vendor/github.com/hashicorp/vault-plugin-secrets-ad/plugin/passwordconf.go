package plugin

type passwordConf struct {
	TTL       int    `json:"ttl"`
	MaxTTL    int    `json:"max_ttl"`
	Length    int    `json:"length"`
	Formatter string `json:"formatter"`
}

func (c *passwordConf) Map() map[string]interface{} {
	return map[string]interface{}{
		"ttl":       c.TTL,
		"max_ttl":   c.MaxTTL,
		"length":    c.Length,
		"formatter": c.Formatter,
	}
}
