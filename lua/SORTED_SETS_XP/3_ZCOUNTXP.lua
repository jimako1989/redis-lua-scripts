-- SORTED SETS with EXPIRE by member
-- ZCOUNTXP key min max
---- This command returns the number of keys still alive between min and max

local ZSET_EXPIREAT_KEY = KEYS[1]..".EXPIREAT"
local now = tonumber(redis.call('TIME')[1])

local c = 0
for v in redis.call('ZRANGE', KEYS[1], KEYS[2], KEYS[3]) do
    if tonumber(redis.call('ZSCORE', ZSET_EXPIREAT_KEY, v)) < now then
        redis.call('ZREM', KEYS[1], v)
        redis.call('ZREM', ZSET_EXPIREAT_KEY, v)
    else
        c = c + 1
    end
end

return c
