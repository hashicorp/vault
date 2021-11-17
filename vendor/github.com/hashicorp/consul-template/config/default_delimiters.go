package config

// DefaultDelims is used to configure the default delimiters used for all templates
type DefaultDelims struct {
	// Left is the left delimiter for templating
	Left *string `mapstructure:"left"`

	// Right is the right delimiter for templating
	Right *string `mapstructure:"right"`
}

// DefaultDefaultDelims returns the default DefaultDelims
func DefaultDefaultDelims() *DefaultDelims {
	return &DefaultDelims{}
}

// Copy returns a copy of the DefaultDelims
func (c *DefaultDelims) Copy() *DefaultDelims {
	if c == nil {
		return nil
	}

	return &DefaultDelims{
		Left:  c.Left,
		Right: c.Right,
	}
}

// Merge merges the DefaultDelims
func (c *DefaultDelims) Merge(o *DefaultDelims) *DefaultDelims {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	if o.Left != nil {
		r.Left = o.Left
	}

	if o.Right != nil {
		r.Right = o.Right
	}

	return r
}
