# redis-lua-scripts

![Go](https://github.com/jimako1989/redis-lua-scripts/workflows/Go/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A curated library of Redis Lua scripts

### Namespace
 - The first number of the lua script filename indicates the number of KEYS.
 Example: ```2_HSETXP.lua``` takes 2 keys.

### Examples

#### HASHES_XP
Group ```HASHES_XP``` is Hashes with EXPIREAT by each field.
Although you can set ```EXPIRE``` to ```HASHES```, hash set will expire by key so that you’ll discard all fields associated with the key as well.
To avoid that, ```HASHES_XP``` creates two hash sets, one of those is normal hash set containing the key and fields, the other one hash set has time of the field expire at, whose key of the hash set is “<key>.EXPIREAT”.
```go
// To load lua script, execute GetScript
script, err := GetScript("HASHES_XP/2_HSETXP")

// HSETXP key expire(secs) field value
_, err := redis.Bool(script.Do(conn, "key", "10", "field", "value"))

// check it works! b is true
b, err := redis.Bool(conn.Do("HEXISTS", "key", "field"))
```

#### SORTED_SETS_XP
Group ```SORTED_SETS_XP``` is SortedSets with EXPIREAT by each member.
```SORTED_SETS_XP``` creates two sorted sets as well, one of those is normal sorted set containing the key and fields, the other one sorted set has time of the field expire at, whose key of the hash set is “<key>.EXPIREAT”.
```go
// To load lua script, execute GetScript
script, err := GetScript("SORTED_SETS_XP/2_ZADDXP")

// ZADDXP key expire(secs) field value
_, err := redis.Bool(script.Do(conn, "key", "10", "1.3", "member"))

// check it works! f is 1.3
f, err := redis.Float64(conn.Do("ZSCORE", "key", "member"))
```

#### GetAllScripts
func ```GetAllScripts``` enables loading all scripts of the group specified.
```go
scripts, err := GetAllScripts("SORTED_SETS_XP")
```
