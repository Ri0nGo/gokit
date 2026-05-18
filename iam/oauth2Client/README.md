# oauth2Client

`oauth2Client` 是 IAM 服务的 OAuth2 客户端包，用于内部系统接入 IAM：生成授权地址、使用 code 换取 `access_token`/`refresh_token`/`openid`、校验 token、刷新 token、获取用户信息。

## 安装

```bash
go get github.com/Ri0nGo/gokit
```

## 初始化客户端

```go
client, err := oauth2Client.New(oauth2Client.Config{
	BaseURL:      "https://iam.example.com",
	ClientID:     "your-client-id",
	ClientSecret: "your-secret",
	RedirectURI:  "https://your-system.example.com/oauth/callback",
	Scopes:       []string{"snsapi_userinfo"},
	Timeout:      3 * time.Second,
})
if err != nil {
	return err
}
```

`ClientID` 对应 IAM 认证应用的 `client_id`。`ClientSecret` 对应 IAM 认证应用的 `secret_key`，请求 token 时会作为 `secret` 发送。

客户端不支持从外部传入 `http.Client`。如未配置 `Timeout`，默认超时时间为 3 秒。

IAM OAuth2 接口返回统一响应结构：`code/msg/data`。客户端方法会直接返回 `data` 中的业务结构。

## 1. 获取授权码 URL

```go
authURL, err := client.AuthCodeURL("state-value")
if err != nil {
	return err
}
```

用户访问 `authURL` 并完成登录后，IAM 会跳转到 `RedirectURI`，并在 query 中带上 `code` 和 `state`。

## 2. code 换取 token

```go
token, err := client.GetUserToken(ctx, code)
if err != nil {
	return err
}

fmt.Println(token.AccessToken)
fmt.Println(token.RefreshToken)
fmt.Println(token.OpenID)
fmt.Println(token.ExpiresIn)
fmt.Println(token.Scope)
```

内部请求：

```text
GET /api/v1/oauth/token?client_id=CLIENT_ID&secret=SECRET&code=CODE&grant_type=authorization_code
```

## 3. 校验 access_token

```go
err := client.CheckUserToken(ctx, token.AccessToken, token.OpenID)
if err != nil {
	return err
}
```

内部请求：

```text
GET /api/v1/oauth/auth?access_token=ACCESS_TOKEN&openid=OPENID
```

## 4. 刷新 access_token

```go
refreshed, err := client.RefreshUserToken(ctx, token.RefreshToken)
if err != nil {
	return err
}

fmt.Println(refreshed.AccessToken)
fmt.Println(refreshed.RefreshToken)
fmt.Println(refreshed.OpenID)
```

内部请求：

```text
GET /api/v1/oauth/refresh_token?client_id=CLIENT_ID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN
```

## 5. 获取用户信息

```go
user, err := client.GetUserInfo(ctx, token.AccessToken, token.OpenID)
if err != nil {
	return err
}

fmt.Println(user.OpenID)
fmt.Println(user.Username)
fmt.Println(user.DisplayName)
fmt.Println(user.Roles)
```

内部请求：

```text
GET /api/v1/oauth/userinfo?access_token=ACCESS_TOKEN&openid=OPENID
```

## 注意事项

不使用 `Authorization` Header。`access_token` 和 `openid` 都通过 query 参数传递。

`openid` 在 code 换 token 阶段返回，不在授权回调阶段返回。后续校验 token 和获取用户信息都需要携带 `openid`。

OAuth2 用户信息不暴露 IAM 内部用户 ID，只返回 `openid`、`username`、`display_name`、`status`、`roles`。
