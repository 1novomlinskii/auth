# JWT-авторизация: схема работы

Два сервиса: **Auth Service** (выдаёт токены) и **URL Shortener** (принимает запросы).

## 1. Логин

```
Клиент → POST /v1/auth/login { email, password }
                              ↓
                         Auth Service проверяет пароль
                              ↓
                   Создаёт access token (JWT, 15мин)
                   + refresh token (7 дней, хеш в БД)
                              ↓
Клиент ← { access_token, refresh_token }
```

Access token (JWT, RS256):
```json
{
  "sub": "user-uuid",
  "tenant_id": "...",
  "role": "user",
  "iat": 1712345678,
  "exp": 1712346578
}
```

## 2. Запрос к другому сервису

```
Клиент → POST /api/shorten { url: "..." }
          Header: Authorization: Bearer <access_token>
                              ↓
                ┌── URL Shortener ──────────┐
                │  1. Проверяет JWT         │
                │     (публичный ключ RSA)  │
                │  2. Извлекает sub (user)  │
                │  3. Создаёт ссылку        │
                └──────────────────────────┘
```

URL Shortener **не ходит в Auth Service** — проверяет JWT локально, stateless.

## 3. Если access token протух

```
Клиент → POST /api/shorten { ... }
          Header: Bearer <expired_token>
                              ↓
          URL Shortener: "token expired" (401)
                              ↓
Клиент → POST /v1/auth/refresh { refresh_token }
                              ↓
          Auth Service проверяет хеш в БД (user_sessions)
                              ↓
Клиент ← { new_access_token, new_refresh_token }
                              ↓
Клиент → POST /api/shorten { ... }  ← с новым токеном
          Header: Bearer <new_token>
```

## 4. Проверка JWT в сервисе-потребителе (пример)

```go
func AuthMiddleware(publicKey *rsa.PublicKey) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tokenStr := extractBearer(r.Header.Get("Authorization"))

            claims := &Claims{}
            token, err := jwt.ParseWithClaims(tokenStr, claims,
                func(t *jwt.Token) (any, error) { return publicKey, nil })

            if err != nil || !token.Valid {
                http.Error(w, "unauthorized", 401)
                return
            }

            ctx := context.WithValue(r.Context(), "user_id", claims.Sub)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

## 5. Схема целиком

```
┌──────────┐  JWT in Header   ┌────────────────┐
│  Client  │ ───────────────→  │ URL Shortener  │
│          │                  │                │
│          │                  │ Проверяет JWT  │
│          │                  │ локально (RSA) │
└──────────┘                  └────────────────┘
     │                               │
     │ POST /v1/auth/login           │
     ▼                               │
┌──────────┐                         │
│  Auth    │  Выдаёт JWT             │
│  Service │  (RS256, приватный ключ)│
│          │                         │
│  Хранит  │  Публичный ключ RSA    │
│  юзеров  │  → в конфиг URL        │
│  пароли  │    Shortener'а         │
│  сессии  │                         │
└──────────┘                         │
                                     │
┌──────────┐                         │
│  Redis   │  Blacklist отозванных   │
│          │  access токенов         │
└──────────┘                         │
```

## Ключевые моменты

| Токен | Время жизни | Где проверяется | Где хранится |
|-------|------------|-----------------|-------------|
| Access token (JWT) | 15 мин | Локально, в сервисе-потребителе | Нет (stateless) |
| Refresh token | 7 дней | Только в Auth Service | user_sessions (SHA-256 хеш) |

- Access token — **без сети**, быстрая проверка
- Refresh token — **через Auth Service**, можно отозвать
- Если access украли — живёт 15 минут
- Если refresh украли — Auth Service видит второй запрос и отзывает все сессии пользователя