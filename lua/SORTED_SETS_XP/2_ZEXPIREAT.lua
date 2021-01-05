-- SORTED SETS with EXPIRE by member
-- ZEXPIREAT key member
---- This command returns the expire time of the member, or -1 if the member is already expired.

local ZSET_EXPIREAT_KEY = KEYS[1]..".EXPIREAT"
local now = tonumber(redis.call('TIME')[1])

local expireAt = redis.call('ZSCORE', ZSET_EXPIREAT_KEY, KEYS[2])

if tonumber(expireAt) < now then
    redis.call('ZREM', KEYS[1], KEYS[2])
    redis.call('ZREM', ZSET_EXPIREAT_KEY, KEYS[2])
    return -1
else
    return expireAt
end
