// Package mgo offers a rich MongoDB driver for Go.
//
// #########################################################
//
// THIS DRIVER IS UNMAINTAINED! See here for details:
//
// https://github.com/go-mgo/mgo/blob/v2-unstable/README.md
//
// #########################################################
//
// Usage of the driver revolves around the concept of sessions.  To
// get started, obtain a session using the Dial function:
//
//     session, err := mgo.Dial(url)
//
// This will establish one or more connections with the cluster of
// servers defined by the url parameter.  From then on, the cluster
// may be queried with multiple consistency rules (see SetMode) and
// documents retrieved with statements such as:
//
//     c := session.DB(database).C(collection)
//     err := c.Find(query).One(&result)
//
// New sessions are typically created by calling session.Copy on the
// initial session obtained at dial time. These new sessions will share
// the same cluster information and connection pool, and may be easily
// handed into other methods and functions for organizing logic.
// Every session created must have its Close method called at the end
// of its life time, so its resources may be put back in the pool or
// collected, depending on the case.
//
// For more details, see the documentation for the types and methods.
//
package mgo
