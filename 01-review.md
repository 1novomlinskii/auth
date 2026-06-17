план действий

// healthHandler responds with the database connectivity status.
func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) { перенести в internal/delivery/http


добавь  .golangci.yml в корень

План
1. docker-compose.yml — обновить версию:
postgres:
  image: postgres:18-alpine
2. Все 7 миграций — заменить:
DEFAULT gen_random_uuid()  →  DEFAULT uuidv7()
Больше ничего не нужно. Расширение pg_uuidv7 не требуется — uuidv7() встроена в PG18.
Вниму внимание: down-миграции не трогаем — DROP TABLE не зависит от того, какой функцией генерировался UUID.

global_handle — уникальный публичный username пользователя. lo
display_name = то, что показывается другим пользователям: "John Doe", "Анна Иванова", "Cat Lover 42".

global_handle - login
display_name -  ИФ

first_name VARCHAR(100),
last_name  VARCHAR(100),

user_local_passwords просто user_passwords

ALTER TABLE user_local_passwords ADD COLUMN reset_token_hash TEXT;
ALTER TABLE user_local_passwords ADD COLUMN reset_token_expires_at TIMESTAMPTZ;


Есть несколько вариантов.
Вариант A: колонки в user_local_passwords
ALTER TABLE user_local_passwords ADD COLUMN reset_token_hash TEXT;
ALTER TABLE user_local_passwords ADD COLUMN reset_token_expires_at TIMESTAMPTZ;
Плюс: просто. Минус: можно иметь только один активный reset-токен.
Вариант B: отдельная таблица (рекомендую)
CREATE TABLE password_reset_tokens (
    token_id    UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id     UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL,                  -- SHA-256 хеш токена
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,                    -- NULL = ещё не использован
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_reset_tokens_user ON password_reset_tokens (user_id);
Плюсы:
- Можно запросить несколько сбросов (старые автоматически протухают)
- Аудит — видно кто и когда запрашивал сброс
- Легко чистить просроченные через фоновую задачу
Вариант C: Redis
Хранить user_id → token_hash + TTL 15min в Redis. Просто и быстро, но если Redis лёг — токены пропадут.
Для pet project я бы выбрал вариант B — он чище по архитектуре, не привязывает к Redis, и ложится в стиль остальных таблиц. Добавить его как 008_password_reset_tokens.sql.


-------------------------------
Ты прав, тут нюанс. Давай разложу.
Когда используется этот partial index
Запрос: "Показать все активные сессии пользователя"
SELECT * FROM user_sessions WHERE user_id = 'u1' AND revoked_at IS NULL;
Индекс содержит только revoked_at IS NULL → находит за O(log N) среди 10 записей, а не 10 миллионов.
Когда НЕ используется
- Проверка refresh token при входе — ищем по refresh_token (у него UNIQUE — свой индекс)
- Поиск протухшей сессии — ищем по expires_at < NOW(), для этого нужен отдельный индекс на expires_at, а не на user_id
Твой сценарий: "сессия протухла, надо проверить"
Логика валидации токена:
-- 1. Ищем по хешу токена (индекс UNIQUE на refresh_token)
SELECT * FROM user_sessions WHERE refresh_token = 'hash...';

-- 2. В коде проверяем:
if session.revoked_at != nil || session.expires_at.Before(time.Now()) {
    return "token invalid"
}
Никакого seq scan по всей таблице — первый шаг уже нашёл строку по индексу refresh_token. Partial index на user_id тут вообще не участвует.

--------------------


Hardcoded tenant UUID в миграции 001_users — '00000000-0000-0000-0000-000000000001'. Если это single-tenant — норм, если multi-tenant — будет больно. Стоит хотя бы замитировать TODO в миграции.


------------------
 общие вопросы


audit
