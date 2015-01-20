package piface

// pouch interface to storage services (for mocking/testing/etc)

type RedisFace interface {
	// Connection functions
	Auth(string) (bool, error)
	Echo(string) (string, error)
	Ping() (string, error)
	Quit() (bool, error)
	Select(int) (bool, error)

	// Hash functions
	Hdel(string, string, ...string) (int, error)
	Hexists(string, string) (bool, error)
	Hget(string, string) (*string, error)
	Hgetall(string) ([]string, error)
	Hincrby(string, string, int) (int64, error)
	Hincrbyfloat(string, string, float64) (float64, error)
	Hkeys(string) ([]string, error)
	Hlen(string) (int64, error)
	Hmget(string, string, ...string) ([]string, error)
	Hmset(string, map[string]string) error
	Hset(string, string, interface{}) (int64, error)
	Hsetnx(string, string, interface{}) (int64, error)
	Hvals(string) ([]string, error)
	Hscan(string, int64) (int64, []string, error)
}
