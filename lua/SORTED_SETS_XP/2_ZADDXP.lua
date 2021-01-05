-- SORTED SETS with EXPIRE by member
-- ZADDXP key expire(sec) score member [score member...]
---- This command creates a another SORTED SETS whose key is KEYS[1]+'.EXPIREAT' contains EXPIREAT TIME by member.

local ZSET_EXPIREAT_KEY = KEYS[1]..".EXPIREAT"
local expireAt = tonumber(redis.call('TIME')[1]) + tonumber(KEYS[2])

local k = ""
for i, v in ipairs(ARGV) do
    if i % 2 == 1 then
        k = v
        -- set / update expireAt whether the field exists or not
        redis.call('ZADD', ZSET_EXPIREAT_KEY, k, expireAt)
    else
        redis.call('ZADD', KEYS[1], k, v)
    end
end

return true
