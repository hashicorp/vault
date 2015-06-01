// Copyright 2014 keimoon. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package gore is a full feature Redis client for Go:
  - Convenient command building and reply parsing
  - Pipeline, multi-exec, LUA scripting
  - Pubsub
  - Connection pool
  - Redis sentinel
  - Client implementation of sharding

Connections

Gore only supports TCP connection for Redis. The connection is thread-safe and can be auto-repaired
with or without sentinel.

  conn, err := gore.Dial("localhost:6379") //Connect to redis server at localhost:6379
  if err != nil {
    return
  }
  defer conn.Close()

Command

Redis command is built with NewCommand

  gore.NewCommand("SET", "kirisame", "marisa") // SET kirisame marisa
  gore.NewCommand("ZADD", "magician", 1337, "alice") // ZADD magician 1337 alice
  gore.NewCommand("HSET", "sdm", "sakuya", 99) // HSET smd sakuya 99

In the last command, the value stored by redis will be the string "99", not the integer 99.
  Integer and float values are converted to string using strconv
  Boolean values are convert to "1" and "0"
  Nil values are stored as zero length string
  Other types are converted to string using standard fmt.Sprint

To efficiently store integer, you can use gore.FixInt or gore.VarInt

Compact integer

Gore supports compacting integer to reduce memory used by redis. There are 2 ways of compacting integer:
  gore.FixInt stores an integer as a fixed 8 bytes []byte.
  gore.VarInt encodes an integer with variable length []byte.

  gore.NewCommand("SET", "fixint", gore.FixInt(1337)) // Set fixint as an 8 bytes []byte
  gore.NewCommand("SET", "varint", gore.VarInt(1337)) // varint only takes 3 bytes

Reply

A redis reply is return when the command is run on a connection

  rep, err := gore.NewCommand("GET", "kirisame").Run(conn)

Parsing the reply is straightforward:

  s, _ := rep.String()  // Return string value if reply is simple string (status) or bulk string
  b, _ := rep.Bytes()   // Return a byte array
  x, _ := rep.Integer() // Return integer value if reply type is integer (INCR, DEL)
  e, _ := rep.Error()   // Return error message if reply type is error
  a, _ := rep.Array()   // Return reply list if reply type is array (MGET, ZRANGE)

Reply converting

Reply support convenient methods to convert to other types

  x, _ := rep.Int()    // Convert string value to int64. This method is different from rep.Integer()
  f, _ := rep.Float()  // Convert string value to float64
  t, _ := rep.Bool()   // Convert string value to boolean, where "1" is true and "0" is false
  x, _ := rep.FixInt() // Convert string value to FixInt
  x, _ := rep.VarInt() // Convert string value to VarInt

To convert an array reply to a slice, you can use Slice method:

  s := []int
  err := rep.Slice(&s) // Convert an array reply to a slice of integer

The following slice element types are supported:
  - integer (int, int64)
  - float (float64)
  - string and []byte
  - FixInt and VarInt
  - *gore.Pair for converting map data from HGETALL or ZRANGE WITHSCORES

Reply returns from HGETALL or SENTINEL master can be converted into a map
using Map:

  m, err:= rep.Map()

Pipeline

Gore supports pipelining using gore.Pipeline:

  p := gore.NewPipeline()
  p.Add(gore.NewCommand("SET", "kirisame", "marisa"))
  p.Add(gore.NewCommand("SET", "alice", "margatroid"))
  replies, _ := p.Run(conn)
  for _, r := range replies {
      // Deal with individual reply here
  }

Script

Script can be set from a string or read from a file, and can be executed over
a connection. Gore makes sure to use EVALSHA before using EVAL to save bandwidth.

  s := gore.NewScript()
  s.SetBody("return redis.call('SET', KEYS[1], ARGV[1])")
  rep, err := s.Execute(conn, 1, "kirisame", "marisa")

Script can be loaded from a file:

  s := gore.NewScript()
  s.ReadFromFile("scripts/set.lua")
  rep, err := s.Execute(conn, 1, "kirisame", "marisa")

Script map

If your application use a lot of script files, you can manage them through ScriptMap

  gore.LoadScripts("scripts", ".*\\.lua") // Load all .lua file from scripts folder
  s := gore.GetScripts("set.lua") // Get script from set.lua file
  rep, err := s.Execute(conn, 1, "kirisame", "marisa") // And execute

Pubsub

Publish message to a channel is easy, you can use gore.Command to issue a PUBLISH
over a connection, or use gore.Publish method:

  gore.Publish(conn, "touhou", "Hello!")

To handle subscriptions, you should allocate a dedicated connection and assign it
to gore.Subscriptions:

  subs := gore.NewSubscriptions(conn)
  subs.Subscribe("test")
  subs.PSubscribe("tou*")

To receive messages, the subcriber should spawn a new goroutine and use
Subscriptions Message channel:

  go func() {
      for message := range subs.Message() {
          if message == nil {
               break
          }
          fmt.Println("Got message from %s, originate from %s: %s", message.Channel, message.OriginalChannel, message.Message)
      }
  }()

Connection pool

To use connection pool, a Pool should be created when application startup. The Dial() method
of the pool should be called to make initial connection to the redis server. If Dial() fail,
it is up to the application to decide to fail fast, or wait and connect again later.

  pool := &gore.Pool{
      InitialConn: 5,  // Initial number of connections to open
      MaximumConn: 10, // Maximum number of connections to open
  }
  err := pool.Dial("localhost:6379")
  if err != nil {
      log.Error(err)
      return
  }
  ...

In each goroutine, a connection from the pool can be get by Acquire() method. Release() method
should always be called later to return the connection to the pool, even in error situation.

  // Inside a goroutine
  conn, err := pool.Acquire()
  if err != nil {
      // Error can happens when goroutine try to acquire a conn
      // from the pool. Application should fail fast here.
      return
  }
  defer pool.Release(conn)
  if conn == nil {
      // This happens when the pool was closed. Application should
      // fail here.
      return
  }
  // Do every thing with the conn, exclusively.
  ...

To gracefully close the pool, call Close() method anywhere in your program.

Transaction

Transaction is implemented using MULTI, EXEC and WATCH. Using transaction
directly with a Conn is not goroutine-safe, so transaction should be used
with connection pool only.

  tr := gore.NewTransaction(conn)
  tr.Watch("a key") // Watch a key
  tr.Watch("another key")
  rep, _ := NewCommand("GET", "a key").Run(conn)
  value, _ := rep.Int()
  tr.Add(NewCommand("SET", "a key", value + 1)) // Add a command to the transaction
  _, err := tr.Commit() // Commit the transaction
  if err == nil {
       // Transaction OK!!!
  } else if err == gore.ErrKeyChanged {
       // Watched key has been changed, transaction should be started over.
  } else {
       // Other errors, transaction should be aborted
  }

Authentication

Gore supports Redis authentication in single connection, pool, sentinel.

Authentication with a single connection can be done by sending "AUTH" command to redis
server as normal, but this can be trouble some when the server is down, and is reconnected
after that. To deal with this problem, gore provides conn.Auth() method:

  conn.Auth("secret password")

This method should be called when the connection is initialized. By calling Auth(), when 
gore tries to reconnect, is will also attempt to send AUTH command to redis server right
after the connection is made.

To configure Auth password with gore.Pool, you can set pool.Password before calling pool.Dial().
Like gore.Conn, gore.Pool also automatically send AUTH command when reconnected.

If you are using sentinel to retrieve Pool or Cluster, instead of using GetPool or GetCluster method,
you can use GetPoolWithPassword or GetClusterWithPassword to connect with a password-protected
pool/cluster

Sentinel

Redis Sentinel is a system that monitors other Redis instance, notify application
when something is wrong with monitored Redis instance, and do automatic failover.
Please note that Redis Sentinel is still in beta stage, and only supported fully
in Redis version 2.8 and above. For more information about setting up Redis Sentinel,
please refer to the official document at http://redis.io/topics/sentinel

Using Redis Sentinel with gore is simple:

  // First, you need to create a Sentinel object:
  s := gore.NewSentinel()
  // Add some Sentinel servers to this object.
  // In production environment, you should have at least 3 Sentinel Servers
  s.AddServer("127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381")
  // Initialize the Sentinel
  err := s.Dial()
  if err != nil {
      return
  }
  // Now, the Sentinel is ready, you can get a monitored pool of connection from the
  // sentinel by using one function:
  pool, err := s.GetPool("mymaster")

The name of the pool ("mymaster") must be an already monitored instance name, otherwise,
the function will return ErrNil. The application also should not call GetPool function
repeatedly because internal locking may cause dropping in performance. It should assign
and reuse the pool variable instead. Because the GetPool function is normally used when
the application starts up, it will fail immediately if the redis instance is still down.
Application can use a for loop and sleep to retry to connect.

Sharding

Gore supports simple sharding strategy: a fixed number of Redis instances are grouped into
"cluster", each instance holds a portion of the cluster keyset.

When a single command needed to be execute on the cluster, gore will redirect the command
to approriate instance based on the key. Gore makes sure that each key will be redirected
to only one instance consistently. Because of the nature of the fixed-sharding, the number
of Redis instances in the cluster should never change, and pipeline or transaction is not
supported.

Gore provides two ways to connect to a cluster.

The first way is using Sentinel. All Redis instances in the same cluster should have the
same prefix, and the suffix should be a number. For example: "mycluster1", "mycluster2",
..., "mycluster20". Using Sentinel, you can get a cluster relatively easy:

  s := NewSentinel()
  s.AddServer("127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381")
  err := s.Dial()
  if err != nil {
      return
  }
  c, err := s.GetCluster("mycluster")

The second way is to add shard to a cluster manually:

  c := NewCluster()
  c.AddShard("127.0.0.1:6379", "127.0.0.1:6380")
  err := c.Dial()
  if err != nil {
      return
  }

Using cluster

A single command can be ran on the cluster with Execute:

  rep, err := c.Execute(NewCommand("SET", "kirisame", "marisa"))
  if err != nil || !rep.IsOk() {
      return
  }
  rep, err := c.Execute(NewCommand("GET", "kirisame"))
  if err != nil {
      return
  }
  value, _ := rep.String() // value should be "marisa"

*/
package gore
