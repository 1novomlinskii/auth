# Database Schema — Auth Service

Все таблицы используют `uuidv7()` для генерации UUID.  
`uuidv7()` встроена в PostgreSQL 18 — расширение не требуется.

---

## `users` — Пользователи (ядро)

| Поле | Тип | Описание | Пример |
|------|-----|----------|--------|
| `user_id` | UUID PK | Идентификатор пользователя | `0194c5e0-7b3a-7f00-8000-000000000001` |
| `tenant_id` | UUID | Идентификатор тенанта (multi-tenant) | `00000000-0000-0000-0000-000000000001` |
| `login` | VARCHAR(50) UNIQUE | Уникальный логин пользователя | `john_doe` |
| `display_name` | VARCHAR(255) | Отображаемое имя | `John Doe` |
| `email_primary` | VARCHAR(255) UNIQUE | Основной email | `john@example.com` |
| `locale` | VARCHAR(10) | Локаль пользователя | `ru_RU` |
| `timezone` | VARCHAR(50) | Часовой пояс | `Europe/Moscow` |
| `created_at` | TIMESTAMPTZ | Дата регистрации | `2026-06-18T10:00:00Z` |
| `updated_at` | TIMESTAMPTZ | Дата последнего обновления профиля | `2026-06-18T12:30:00Z` |
| `is_active` | BOOLEAN | Активен ли пользователь | `true` |

---

## `user_identities` — OAuth2/SSO провайдеры (Google, GitHub и т.д.)

| Поле | Тип | Описание | Пример |
|------|-----|----------|--------|
| `identity_id` | UUID PK | Идентификатор связки | `0194c5e0-8b4a-7f00-8000-000000000002` |
| `user_id` | UUID FK → users | Ссылка на пользователя | `0194c5e0-7b3a-7f00-8000-000000000001` |
| `provider` | VARCHAR(20) | Провайдер | `google` |
| `provider_user_id` | VARCHAR(255) | ID пользователя у провайдера | `1234567890` (sub от Google) |
| `provider_email` | VARCHAR(255) | Email от провайдера | `john@gmail.com` |
| `provider_name` | VARCHAR(255) | Имя от провайдера | `John Doe` |
| `avatar_url` | TEXT | Ссылка на аватар | `https://avatars.example.com/123` |
| `last_login_at` | TIMESTAMPTZ | Последний вход через провайдера | `2026-06-18T10:00:00Z` |
| `created_at` | TIMESTAMPTZ | Дата привязки | `2026-06-18T10:00:00Z` |

Уникальность: `(provider, provider_user_id)`

---

## `user_password_hash` — Хэши паролей (только для local-регистрации)

| Поле | Тип | Описание | Пример |
|------|-----|----------|--------|
| `user_id` | UUID PK FK → users | Ссылка на пользователя | `0194c5e0-7b3a-7f00-8000-000000000001` |
| `password_hash` | TEXT | bcrypt-хэш пароля (cost=12) | `$2a$12$LJ3m4ys3Lk...` |
| `reset_token_hash` | TEXT | SHA-256 хэш токена сброса пароля | `a1b2c3d4e5f6...` (64 hex символа) |
| `reset_token_expires_at` | TIMESTAMPTZ | Срок действия токена сброса | `2026-06-18T11:00:00Z` |
| `created_at` | TIMESTAMPTZ | Дата установки пароля | `2026-06-18T10:00:00Z` |
| `updated_at` | TIMESTAMPTZ | Дата последней смены пароля | `2026-06-18T12:00:00Z` |

Только для пользователей, зарегистрированных через email+password.  
OAuth-пользователи не имеют записи в этой таблице.

---

## `user_sessions` — Сессии (refresh токены)

| Поле | Тип | Описание | Пример |
|------|-----|----------|--------|
| `session_id` | UUID PK | Идентификатор сессии | `0194c5e0-9c5a-7f00-8000-000000000003` |
| `user_id` | UUID FK → users | Владелец сессии | `0194c5e0-7b3a-7f00-8000-000000000001` |
| `identity_id` | UUID FK → user_identities | Через какой OAuth-провайдер вошёл (NULL для local) | `0194c5e0-8b4a-7f00-8000-000000000002` |
| `refresh_token` | TEXT UNIQUE | SHA-256 хэш refresh токена | `a1b2c3d4e5f6...` |
| `device_info` | TEXT | Информация об устройстве | `Mozilla/5.0 ...` |
| `ip_address` | INET | IP адрес клиента | `192.168.1.1` |
| `expires_at` | TIMESTAMPTZ | Дата истечения refresh токена | `2026-07-18T10:00:00Z` |
| `created_at` | TIMESTAMPTZ | Начало сессии | `2026-06-18T10:00:00Z` |
| `revoked_at` | TIMESTAMPTZ | Дата отзыва сессии (NULL = активна) | `2026-06-18T12:00:00Z` |

Индекс: `idx_sessions_user` — partial index `(user_id) WHERE revoked_at IS NULL`

---

## `api_keys` — Ключи для service-to-service авторизации

| Поле | Тип | Описание | Пример |
|------|-----|----------|--------|
| `api_key_id` | UUID PK | Идентификатор ключа | `0194c5e0-ad6a-7f00-8000-000000000004` |
| `tenant_id` | UUID | Тенант, которому принадлежит ключ | `00000000-0000-0000-0000-000000000001` |
| `service_name` | TEXT | Название сервиса | `notification-service` |
| `key_prefix` | VARCHAR(8) | Префикс ключа (для идентификации) | `sk_auth_` |
| `key_hash` | TEXT UNIQUE | SHA-256 хэш полного ключа | `a1b2c3d4e5f6...` |
| `permissions` | TEXT[] | Массив разрешений | `{users:read,users:write}` |
| `created_by` | UUID FK → users | Кто создал ключ | `0194c5e0-7b3a-7f00-8000-000000000001` |
| `expires_at` | TIMESTAMPTZ | Дата истечения | `2027-06-18T10:00:00Z` |
| `is_active` | BOOLEAN | Активен ли ключ | `true` |
| `created_at` | TIMESTAMPTZ | Дата создания | `2026-06-18T10:00:00Z` |

---

## `roles` — Роли (RBAC)

| Поле | Тип | Описание | Пример |
|------|-----|----------|--------|
| `role_id` | UUID PK | Идентификатор роли | `0194c5e0-be7a-7f00-8000-000000000005` |
| `tenant_id` | UUID | Тенант, которому принадлежит роль | `00000000-0000-0000-0000-000000000001` |
| `name` | TEXT | Название роли | `admin` |
| `permissions` | TEXT[] | Массив разрешений | `{users:*}` |
| `is_system` | BOOLEAN | Системная роль (нельзя удалить) | `true` |
| `description` | TEXT | Описание роли | `Администратор системы` |
| `created_at` | TIMESTAMPTZ | Дата создания | `2026-06-18T10:00:00Z` |
| `updated_at` | TIMESTAMPTZ | Дата последнего изменения | `2026-06-18T12:00:00Z` |

Уникальность: `(tenant_id, name)`

---

## `user_roles` — Связь пользователей и ролей (M2M)

| Поле | Тип | Описание | Пример |
|------|-----|----------|--------|
| `user_id` | UUID FK → users | Пользователь | `0194c5e0-7b3a-7f00-8000-000000000001` |
| `role_id` | UUID FK → roles | Роль | `0194c5e0-be7a-7f00-8000-000000000005` |
| `assigned_at` | TIMESTAMPTZ | Дата назначения | `2026-06-18T10:00:00Z` |

PK: `(user_id, role_id)`