# 1.0.1 (April 12th, 2019)

IMPROVEMENTS

* InMemTransport: Add timeout for sending a message [[GH-313](https://github.com/hashicorp/raft/pull/313)]
* ensure 'make deps' downloads test dependencies like testify [[GH-310](https://github.com/hashicorp/raft/pull/310)]
* Clarifies function of CommitTimeout [[GH-309](https://github.com/hashicorp/raft/pull/309)]
* Add additional metrics regarding log dispatching and committal [[GH-316](https://github.com/hashicorp/raft/pull/316)]

# 1.0.0 (October 3rd, 2017)

v1.0.0 takes the changes that were staged in the library-v2-stage-one branch. This version manages server identities using a UUID, so introduces some breaking API changes. It also versions the Raft protocol, and requires some special steps when interoperating with Raft servers running older versions of the library (see the detailed comment in config.go about version compatibility). You can reference https://github.com/hashicorp/consul/pull/2222 for an idea of what was required to port Consul to these new interfaces.

# 0.1.0 (September 29th, 2017)

v0.1.0 is the original stable version of the library that was in master and has been maintained with no breaking API changes. This was in use by Consul prior to version 0.7.0.
