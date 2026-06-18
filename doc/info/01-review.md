план действий

1 // healthHandler responds with the database connectivity status.
func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) { 
перенести в internal/delivery/http

2 добавь  .golangci.yml в корень

3. docker-compose.yml — обновить версию:
postgres:
  image: postgres:18-alpine
4. Все 7 миграций — заменить:
DEFAULT gen_random_uuid()  →  DEFAULT uuidv7()
Больше ничего не нужно. Расширение pg_uuidv7 не требуется — uuidv7() встроена в PG18.
Вниму внимание: down-миграции не трогаем — DROP TABLE не зависит от того, какой функцией генерировался UUID.

5 global_handle —переименовать  login

6 переименовать user_local_passwords просто user_passwords

7 добавить поля для сброса пароля
 reset_token_hash TEXT;
reset_token_expires_at TIMESTAMPTZ;

8 сделать карту БД какой нибудь BD-INFO.MD
где расписать что значит каждое поле БД и пример что там будет лежать 






