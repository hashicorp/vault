/*
Package megos provides a client library for accessing information of a Apache Mesos cluster.

Construct a new megos client, then use the various functions on the client to
access different information of Mesos HTTP endpoints.
For example to identify the leader node:

	node1, _ := url.Parse("http://192.168.1.120:5050/")
	node2, _ := url.Parse("http://192.168.1.122:5050/")

	mesos := megos.NewClient([]*url.URL{node1, node2})
	leader, err := mesos.DetermineLeader()
	if err != nil {
	panic(err)
	}

	fmt.Println(leader)
	// Output:
	// master@192.168.1.122:5050

More examples are available in the README.md on github: https://github.com/andygrunwald/megos
*/
package megos
