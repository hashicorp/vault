package vault

//go:generate go run github.com/hashicorp/vault/tools/stubmaker "myfunc() {}"

func testmyfunc() {
	myfunc()
}
