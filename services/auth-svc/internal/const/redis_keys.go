package _const

import (
	"fmt"
	"strings"
)

var (
	RedisKeyRefreshToken    = RedisKey{PrefixKey: "rf_token"}
	RedisKeyRolePermission  = RedisKey{PrefixKey: "role_perm_1"}
	RedisKeyUserCompanyRole = RedisKey{PrefixKey: "usr_comp_role_1"}
	RedisKeyLoginRateLimit  = RedisKey{PrefixKey: "rate_limit_login"}
	RedisKeyBlockLogin      = RedisKey{PrefixKey: "block_login"}

	RedisKeyWhitelistIP = RedisKey{PrefixKey: "authsvc-v1:whitelist_ip"}

	RedisKeyBackendDomainsTTL = 24 * 3600
)

var (
	RedisKeyUserDBCacheById      = RedisKey{PrefixKey: "v1_user_db_cache_by_id"}
	RedisKeyUserDBCacheByAccount = RedisKey{PrefixKey: "v1_user_db_cache_by_account"}
)

type RedisKey struct {
	PrefixKey string
}

// Key will gen
// Eg: <Rediskey>.Key(abc-bcs)
func (k *RedisKey) Key(key ...string) string {

	keyString := strings.Join(key, "-")
	return fmt.Sprintf("%v:%v", k.PrefixKey, keyString)
}
