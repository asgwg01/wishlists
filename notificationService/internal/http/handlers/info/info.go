package info

import (
	"context"
	"log/slog"
	"net/http"
	dowork "notificationService/internal/services/doWork"
)

func NewHandler(log *slog.Logger, srv dowork.IServiceWorkSome) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const logPrefix = "handlers.info"
		log := log.With(
			slog.String("where", logPrefix),
		)
		log.Debug("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

		ok, _ := srv.WorkSome(context.Background(), "some work")
		if ok {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, err := w.Write([]byte(infoPage))
		if err != nil {
			log.Error("Error write request", slog.String("err", err.Error()))
		}

		log.Info("End of work")
	}
}

const infoPage = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Сервер</title>
  <style>
    body {margin:0;height:100vh;display:flex;align-items:center;justify-content:center;background:#121212;color:#e0e0e0;font:14px Arial}
    .c {width:360px;padding:20px;background:#1e1e1e;border-radius:8px;text-align:center}
    .i {text-align:left;margin:12px 0}
    .l {color:#8a8a8a;font-size:13px}
    .v {font-weight:500}
  </style>
</head>
<body>
  <div class="c">
    <h1 style="color:#bb86fc;margin:0 0 16px">Информация о приложении</h1>
    <div class="i"><div class="l">Название</div><div class="v">Приложение</div></div>
    <div class="i"><div class="l">Версия</div><div class="v">1.0.0</div></div>
    <div class="i"><div class="l">IP</div><div class="v" id="ip">...</div></div>
    <div class="i"><div class="l">Порт</div><div class="v" id="port">...</div></div>
  </div>
  <script>
    const url = new URL(location.href);
    document.getElementById('ip').textContent = url.hostname;
    document.getElementById('port').textContent = url.port || (url.protocol === 'https:' ? '443' : '80');
  </script>
</body>
</html>

`
