---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: JWT 
---
JWT token 是一个字符串，有三部分组成，头部、载荷与签名，中间用 `.` 分隔

##### 头部（header）

头部通常由两部分组成：令牌的类型和正在使用的签名算法

```json
{
    "alg":"HS256",
    "typ:"JWT"
}
```

然后使用 `base64url` 编码得到头部

##### 载荷（payload）

载荷中放置 `token` 的一些基本信息。

载荷属性分为三类：

- 预定义（Registered）
- 公有（public）
- 私有（private）

```json
{
    "sub":"1",	// subject 主题
    "iss":"https://localhost:8080/auth/login",	// issuer 签发人
    "iat":145158154,	// issued at 签发时间
    "exp":145464552,	// expiration time 过期时间
    "nbf":1451568945,	// not before 生效时间，此时间之前是无效的
    "jti":"3451545fsdfsdfas25245153",	// 编号
    "aud":"dev"	// audience 受众
}
```

> Tips：上 7 个字段都是由官方所定义的，即预定义

##### 签名（signature）

签名时需要用到前面编码过的两个字符串，如果以 `HMACSHA256` 加密，

```
HMACSHA256(
base64UrlEncode(headher)+'.'+base64UrlEncode(payload),secret
)
```

加密后再进行 `base64Url` 编码最后得到字符串就是 `token`。

签名的作用是保证 jwt 没有被篡改过

> HMAC 算法是不可逆的算法，类似 MD5 和 hash ，但多了个密钥，密钥由服务端持有，客户端把 token 发给服务端，服务端可以把其中的头部和载荷再加上事先共享的 secret 再进行一次 HMAC 加密，得到的结构和 token 的第三段进行对比，如果一样则表明数据没有被篡改

JWT 使用两种方式：

- 加到 `URL`：`?token=jwt_token`
- 加到 `header`，`Authorization：Bearer JWT_token`

JWT 在客户端的存储方式

- LocalStorage
- SessionStorage
- Cookie

推荐第三种，第一二中存在跨域读取限制

##### Cookie 的跨域策略

子可以读父，父不可以读子，兄弟之间不可以互相访问

##### JWT 无状态

因为 JWT 有效期完全和载荷中的编码过期时间一致，服务端不维护任何状态，无状态的优势：

- 节省服务器资源，，因为服务器无需唯一个状态
- 适合分布式，因为服务端无需维护状态，因此如果服务端是多台服务器组成的分布式集群，那么无需像『有状态』一样互相同步各自的状态
- 时间换空间，因为 token 校验时通过签名校验，签名校验消耗的是 CPU 时间，而『有状态』是通过客户端提供的凭据对服务端现有的状态进行一次查询，消耗的是 I/O 和内存、磁盘空间。对于一个 web 服务器来说，其属于 I/O 密集型，因此通过时间换空间，可以提高整体的硬件使用率。

