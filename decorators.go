package pouch

// RedisTableable entities are able to return a key structure
// that can be used in conjunction with identifiable fields
// to deterministically generate the identifying key in redis
// for the given entity.
type RedisDecorated interface {
	KeyFormula() string
	SetFieldFromString(string, string) error
}
