-- SORTED SETS with EXPIRE by member
-- ZCARDXP key
---- This command returns the number of keys still alive

local ZSET_EXPIREAT_KEY = KEYS[1]..".EXPIREAT"
local now = tonumber(redis.call('TIME')[1])

local c = 0
local m = ""
local skip = false
for i, v in ipairs(redis.call('ZRANGE', ZSET_EXPIREAT_KEY, 0, -1, "WITHSCORES")) do
    if i % 2 == 1 then
        m = v
        if tonumber(v) < now then
            redis.call('ZREM', KEYS[1], m)
            redis.call('ZREM', ZSET_EXPIREAT_KEY, m)
            skip = true
        end
    else
        if not skip then
            c = c + 1
        end
        skip = false
    end
end

return c
