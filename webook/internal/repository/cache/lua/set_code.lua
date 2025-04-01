-- 发送到的key，也就是code:义务:手机号
local key = KEYS[1]
-- 发送次数
local cntKey = key..":cnt"
local val = ARGV[1]
-- 验证码的有效时间是10分钟
local ttl = tonumber(redis.call("ttl", key))
-- 没有过期时间，可能是其他人误操作
if ttl == -1 then
    return -2
elseif ttl == -2 or ttl < 540 then
    -- 没有发送过验证码，直接设置
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 5)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 这里是发送的次数太多了，一分钟内要求发送超过两次
    return -1
end