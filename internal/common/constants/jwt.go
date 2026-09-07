package constants

// JWTContextKey 是 echo-jwt 把解析后的 *jwt.Token 写入 Echo context 的键。
// 中间件写入、controller 读取必须使用同一个值。
const JWTContextKey = "user"
