// Package radix implements all functionality needed to work with redis and all
// things related to it, including redis cluster, pubsub, sentinel, scanning,
// lua scripting, and more.
//
// This package has extensive examples documenting advanced behavior not covered
// here.
//
// Creating a Client
//
// For a single redis instance use PoolConfig to create a connection pool. The
// connection pool implements the Client interface. It is thread-safe and will
// automatically create, reuse, and recreate connections as needed:
//
//	client, err := (radix.PoolConfig{}).New("tcp", "127.0.0.1:6379")
//	if err != nil {
//		// handle error
//	}
//
// If you're using sentinel or cluster you should use SentinelConfig or
// ClusterConfig (respectively) to create your Client instead.
//
// Commands
//
// Any redis command can be performed by passing a Cmd into a Client's Do
// method. Each Cmd instance should only be used once. The return from the Cmd
// can be captured into any appopriate go primitive type, or a slice, map, or
// struct, if the command returns an array.
//
//	// discard the result
//	err := client.Do(ctx, radix.Cmd(nil, "SET", "foo", "someval"))
//
//	var fooVal string
//	err := client.Do(ctx, radix.Cmd(&fooVal, "GET", "foo"))
//
//	var fooValB []byte
//	err := client.Do(ctx, radix.Cmd(&fooValB, "GET", "foo"))
//
//	var barI int
//	err := client.Do(ctx, radix.Cmd(&barI, "INCR", "bar"))
//
//	var bazEls []string
//	err := client.Do(ctx, radix.Cmd(&bazEls, "LRANGE", "baz", "0", "-1"))
//
//	var buzMap map[string]string
//	err := client.Do(ctx, radix.Cmd(&buzMap, "HGETALL", "buz"))
//
// FlatCmd can also be used if you wish to use non-string arguments like
// integers, slices, maps, or structs, and have them automatically be flattened
// into a single string slice.
//
// Other Actions
//
// Cmd and FlatCmd both implement the Action interface. Other Actions include
// Pipeline, WithConn, and EvalScript.Cmd. Any of these may be passed into any
// Client's Do method.
//
//	var fooVal string
//	p := radix.NewPipeline()
//	p.Append(radix.FlatCmd(nil, "SET", "foo", 1))
//	p.Append(radix.Cmd(&fooVal, "GET", "foo"))
//
//	if err := client.Do(p); err != nil {
//		panic(err)
//	}
//	fmt.Printf("fooVal: %q\n", fooVal)
//
// Transactions
//
// There are two ways to perform transactions in redis. The first is with the
// MULTI/EXEC commands, which can be done using the WithConn Action (see its
// example). The second is using EVAL with lua scripting, which can be done
// using the EvalScript Action (again, see its example).
//
// EVAL with lua scripting is recommended in almost all cases. It only requires
// a single round-trip, it's infinitely more flexible than MULTI/EXEC, it's
// simpler to code, and for complex transactions, which would otherwise need a
// WATCH statement with MULTI/EXEC, it's significantly faster.
//
// AUTH and other settings via Dialer
//
// Dialer has fields like AuthPass and SelectDB which can be used to configure
// Conns at creation.
//
// PoolConfig takes a Dialer as one of its fields, so that all Conns the Pool
// creates will be created with those settings.
//
// Other Clients which create their own Pools, like Cluster and Sentinel, will
// take in a PoolConfig which can be used to configure the Pools they create.
//
// For example, to create a Cluster instance which uses a particular AUTH
// password on all Conns:
//
//	cfg := radix.ClusterConfig{
//		PoolConfig: radix.PoolConfig{
//			Dialer: radix.Dialer{
//				AuthPass: "mySuperSecretPassword",
//			},
//		},
//	}
//
//	client, err := cfg.New(ctx, []string{redisAddr1, redisAddr2, redisAddr3})
//
// Custom implementations
//
// All interfaces in this package were designed such that they could have custom
// implementations. There is no dependency within radix that demands any
// interface be implemented by a particular underlying type, so feel free to
// create your own Pools or Conns or Actions or whatever makes your life easier.
//
// Errors
//
// Errors returned from redis can be explicitly checked for using the the
// resp3.SimpleError type. Note that the errors.As or errors.Is functions,
// introduced in go 1.13, should be used.
//
//	var redisErr resp3.SimpleError
//	err := client.Do(ctx, radix.Cmd(nil, "AUTH", "wrong password"))
//	if errors.As(err, &redisErr) {
//		log.Printf("redis error returned: %s", redisErr.S)
//	}
//
package radix
