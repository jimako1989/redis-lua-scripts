-- SORTED LISTS with EXPIRE by member
-- ZLLENXP key
return redis.call('LLEN', KEYS[1])
