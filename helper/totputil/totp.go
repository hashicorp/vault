package totputil

import (
	"github.com/hashicorp/vault/logical/framework"
	cache "github.com/patrickmn/go-cache"
)

type Backend struct {
	*framework.Backend

	UsedCodes *cache.Cache
}
