-- HASH SETS with EXPIRE by fields
-- HSETXPAT key expire_at field value [field value...]

local HSET_EXPIREAT_KEY = KEYS[1]..".EXPIREAT"
local expireAt = tonumber(KEYS[2])

-- get the list of fields in arguments, also values
local k = ""
for i, v in ipairs(ARGV) do
    if i % 2 == 0 then
        k = v
        -- set / update expireAt whether the field exists or not
        redis.call('HSET', HSET_EXPIREAT_KEY, k, expireAt)
    else
        redis.call('HSET', KEYS[1], k, v)
    end
end
