-- Set HSET with EXPIRE
redis.call('MULTI')
redis.call('HSET', KEYS[1], ARGV[1], ARGV[2])
redis.call('EXPIRE', KEYS[1], ARGV[3])
redis.call('EXEC')
