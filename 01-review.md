перенести миграции в апп
можно ли сделать версию миграций

что за настройка cfg.Log.Pretty 

обясни каждую строчку

	// Wrap with request logging.
	var handler http.Handler = mux
	handler = hlog.RequestIDHandler("req_id", "X-Request-ID")(handler)
	handler = hlog.NewHandler(a.logger)(handler)
	handler = hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("request")
	})(handler)

		addr := ":" + a.cfg.HTTP.Portтут же только порт? а хост?

переделать  make-local-up  b=1  по аналогии с lint:

// healthHandler responds with the database connectivity status.
func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) { перенести в internal/delivery/http


----------------------------------
миграции

gen_random_uuid(), v7

users global_handle что за поле что там лежать будет
может добавить отдельно ФИО
индекс на user_id email_primary для поиска по email или  user_id, например для востановления пароль надо натйи пользователя по email

users tenant_id убрать дефолт


user_local_passwords просто user_passwords
поле для временного хеша при сборсе и востановлении пароля

CREATE INDEX idx_sessions_user ON user_sessions (user_id) WHERE revoked_at IS NULL;
то есть если я добавляю revoked_at то индекс пропаает?



------------------
 общие вопросы

как в общес будет работать, например призодит запрос на прилоеднеи создать короткую ссылку, я иду в авторизацию, я хочу использовать JWT что приложение отправит в Auth