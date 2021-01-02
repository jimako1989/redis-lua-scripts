-- Set HINCRBY with PEXPIRE
redis.call('HINCRBY', KEYS[1], ARGV[1], ARGV[2])
redis.call('PEXPIRE', KEYS[1], ARGV[3])
