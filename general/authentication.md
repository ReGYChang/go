- [JWT](#jwt)
  - [What is JWT](#what-is-jwt)
  - [JWT Principal](#jwt-principal)
    - [Header](#header)
    - [Payload](#payload)
    - [Signature](#signature)
  - [How does JWT Work](#how-does-jwt-work)

# JWT

## What is JWT

`JWT(JSON Web Token)` 是一個非常輕巧的規範, 其允許使用者使用 `JWT` 在 client 及 server 之間傳遞安全可靠的資訊, 利用 `HTTP` 通訊過程中進行身份認證

在以下場景中非常適合使用 `JWT`:
- Authorization: 其為使用 `JWT` 最頻繁的場景, 一旦使用者登錄, 後續每個請求都將包含 `JWT`, 允許使用者訪問該 token 允許的 routes, services 及 resources, `SSO` 是現在廣泛使用 `JWT` 的一個應用, 因其輕量且可輕鬆跨域使用
- Information Exchange: 對於點對點之間安全的傳輸訊息而言, `JWT` 是一種很好的方式, 因其可以被簽名, 可以正確驗證訊息發送者資訊及訊息內容是否被篡改

## JWT Principal

`JWT` 由三個部分組成, 其之間使用 `.` 連接, 分別為:

- Header
- Payload
- Signature

一個典型的 `JWT` 如下:

```
xxxxx.yyyyy.zzzzz
```

### Header

`Header` 一般由兩部分組成: token 類型及加密演算法名稱(HMACSHA256 或 RSA)

如:

```json
{
    'alg': "HS256",
    'typ': "JWT"
}
```

再使用 `BASE64` 對此 JSON 進行編碼後就得到 `JWT` 的第一個部分

### Payload

`Payload` 中包含聲明(要求), 關於實體(通常為使用者)和其他資料的聲明

聲明有三種類型:

- Registered Claims: 有一組預先定義的聲明, 其不為強制, 但推薦, 如: iss(issuer), exp(expiration time), sub(subject), aud(audience) 等
- Public: 可以隨意定義
- Private: 用於許可使用點對點之間共享資訊, 且不是註冊或公開的聲明

下面為一個例子:

```json
{
    "sub": '1234567890',
    "name": 'regy',
    "admin":true
}
```

對 `Payload` 進行 `BASE64` 編碼即為 `JWT` 第二部分

>❗️ 需注意不要在 JWT Payload 或 Header 中存放敏感資訊, 除非其透過加密

### Signature

> 為了計算出 signature, 必須有編碼後的 Header 和 Payload, 及一個 private key, 並使用 Header 中指定的加密演算法對其簽名即可

如:

```java
HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload), secret)
```

`Signature` 主要用於校驗訊息在傳遞過程中是否遭到竄改, 且對於使用 private key 簽名的 token, 還可以驗證 JWT 的發送方是否偽造身份

## How does JWT Work

- Client 通過 username/password 登陸 Server
- Server 對 Client 身份進行驗證
- Server 對該 user 生成 token 並返回 Client
- Client 發起 request 攜帶此 token
- Server 收到 request 先驗證 token 正確性再返回資料
- Client 將 token 保存至本地瀏覽器, 一般保存在 `cookie`

> Server 不需保存 token, 只需對 token 中攜帶的資訊進行驗證即可; 無論 Client 訪問 Server 中的哪台機器, 只需通過使用者資訊驗證即可


