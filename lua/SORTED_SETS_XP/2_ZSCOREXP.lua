-- SORTED SETS with EXPIRE by member
-- ZSCOREXP key member
---- This command returns the score of the member

local ZSET_EXPIREAT_KEY = KEYS[1]..".EXPIREAT"
local now = tonumber(redis.call('TIME')[1])

local score = redis.call('ZSCORE', KEYS[1], KEYS[2])

if tonumber(redis.call('ZSCORE', ZSET_EXPIREAT_KEY, KEYS[2])) < now then
    redis.call('ZREM', KEYS[1], KEYS[2])
    redis.call('ZREM', ZSET_EXPIREAT_KEY, KEYS[2])
    return ""
else
    return score
end
