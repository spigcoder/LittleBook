local key = KEYS[1]
local cntKey = key..":cnt"
local expectCode = ARGV[1]
local code = redis.call("get", key)
local cntValue = redis.call("get", cntKey)
local cnt = tonumber(cntValue) or -1 

--为什么不用验证验证码以过期，因为过期了，那么就返回验证码错误就好了
if cnt <= 0 then
    -- 验证码已经发送超过5次了, 用户一直输入错误，有人搞你
    return -1
end
if code == expectCode then
    -- 验证码正确, 之后就不可以再使用了
    redis.call("set", cntKey, -1)
    return 0
else 
    -- 验证码错误
    redis.call("decr", cntKey)
    return -2
end
