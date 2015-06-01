package gore

// Cluster consists of fix number of shards, with each shard holds a portion
// of the keyset. Cluster can be created by adding shards, or using sentinel.
type Cluster struct {
	addresses     []*addressWithPassword
	shards        []*Pool
	sentinel      bool
	ShardStrategy func(string, int) int
}

type addressWithPassword struct {
	address  string
	password string
}

// NewCluster creates new cluster. You must add shards to this cluster manually
func NewCluster() *Cluster {
	return &Cluster{
		sentinel:      false,
		ShardStrategy: DefaultShardStrategy,
	}
}

// AddShard add a list of shards to the cluster.
func (c *Cluster) AddShard(addresses ...string) {
	if !c.sentinel {
		for _, address := range addresses {
			c.addresses = append(c.addresses, &addressWithPassword{address, ""})
		}
	}
}

// AddShardWithPassword add a password-protected shard
func (c *Cluster) AddShardWithPassword(address, password string) {
	if !c.sentinel {
		c.addresses = append(c.addresses, &addressWithPassword{address, password})
	}
}

// Dial connects the cluster to all shards. If one shard cannot be connected, the whole
// operation will fail.
func (c *Cluster) Dial() (err error) {
	if c.sentinel {
		return nil
	}
	if len(c.addresses) == 0 {
		return ErrNoShard
	}
	defer func() {
		if err != nil {
			for _, pool := range c.shards {
				pool.Close()
			}
		}
	}()
	for _, address := range c.addresses {
		pool := &Pool{Password: address.password}
		err = pool.Dial(address.address)
		if err != nil {
			return err
		}
		c.shards = append(c.shards, pool)
	}
	return nil
}

// Execute runs a command on the cluster. The command will be send to appropriate shard
// based on its key. If the command has no key (PING, INFO), this function returns
// ErrNoKey
func (c *Cluster) Execute(cmd *Command) (*Reply, error) {
	if len(cmd.args) < 1 {
		return nil, ErrNoKey
	}
	key := string(convertString(cmd.args[0]))
	pool := c.shards[DefaultShardStrategy(key, len(c.shards))]
	conn, err := pool.Acquire()
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, ErrNotConnected
	}
	defer pool.Release(conn)
	return cmd.Run(conn)
}

// DefaultShardStrategy converts a string key into number and takes modulo
// with the size of cluster
func DefaultShardStrategy(key string, size int) int {
	r := 0
	for _, c := range key {
		r = r*256 + int(c)
		r = r % size
	}
	return r
}
