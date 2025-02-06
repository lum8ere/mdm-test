package run_processor

import (
	"encoding/json"
	"io"
	"net/http"

	"mdm/libs/4_common/smart_context"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// AppHandler определяет сигнатуру обработчика.
type AppHandler func(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error)

// JSONResponseMiddleware оборачивает вызов AppHandler в http.HandlerFunc.
func JSONResponseMiddleware(sctx smart_context.ISmartContext, handler AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Декодируем JSON-тело запроса и объединяем его с URL-параметрами.
		data, err := parseJSONBody(r)
		if err != nil {
			handleError(w, err, sctx.GetLogger())
			return
		}

		// Пример: извлекаем параметр "id" из URL и добавляем его в data
		if id := chi.URLParam(r, "id"); id != "" {
			data["id"] = id
		}

		// Вызываем обработчик с распарсенными данными.
		response, err := handler(sctx, data)
		if err != nil {
			handleError(w, err, sctx.GetLogger())
			return
		}

		// Отправляем JSON-ответ.
		w.Header().Set("Content-Type", "application/json")
		if response != nil {
			if err := json.NewEncoder(w).Encode(response); err != nil {
				sctx.Error("Failed to encode JSON response", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

// parseJSONBody пытается декодировать тело запроса как JSON в map[string]interface{}.
func parseJSONBody(r *http.Request) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if r.Body == nil {
		return data, nil
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return data, nil
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// handleError отправляет JSON-ответ с сообщением об ошибке.
func handleError(w http.ResponseWriter, err error, log *zap.Logger) {
	log.Error("Handler error", zap.Error(err))
	w.Header().Set("Content-Type", "application/json")
	statusCode := http.StatusInternalServerError
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
