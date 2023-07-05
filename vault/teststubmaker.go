package vault

import "github.com/hashicorp/vault/vault/cluster"

func testmyfunc() {
	myfunc(cluster.Listener{})
}
