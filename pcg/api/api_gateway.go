/*
package api_gateway


import (
	"CommentsService/pcg/comments"
	"CommentsService/pcg/logs"
	"CommentsService/pcg/news"
	"CommentsService/pcg/types"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

func AggregateRequests(newsID int) (types.News, []types.Comment, error) {
	// Создайте канал и WaitGroup
	resultCh := make(chan types.Result, 2)
	var wg sync.WaitGroup

	// Горутина для запроса информации о новости
	wg.Add(1)
	go func() {
		defer wg.Done()

		newsInfo, err := news.GetNews(newsID)
		resultCh <- types.Result{NewsInfo: newsInfo, Err: err}
	}()

	// Горутина для запроса списка комментариев
	wg.Add(1)
	go func() {
		defer wg.Done()

		comments, err := comments.GetCommentsByNewsID(newsID)
		resultCh <- types.Result{Comments: comments, Err: err}
	}()

	// Дождитесь завершения всех горутин
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Прочитайте результаты и ошибки
	var newsInfo types.News
	var comments []types.Comment
	var firstErr error

	for result := range resultCh {
		if result.Err != nil && firstErr == nil {
			firstErr = result.Err
		}

		if result.NewsInfo.ID != 0 {
			newsInfo = result.NewsInfo
		}

		comments = append(comments, result.Comments...)
	}

	if firstErr != nil {
		return types.News{}, nil, firstErr
	}

	return newsInfo, comments, nil
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса POST
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
	requestID := uuid.New().String()

	// Устанавливаем Request ID в заголовке запроса
	r.Header.Set("X-Request-ID", requestID)

	logs.LogRequest(r.Header.Get("X-Request-ID"), r.RemoteAddr, http.StatusOK)

	// Декодируем JSON-запрос от клиента в структуру Comment
	var comment types.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Ошибка при чтении JSON-запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Выполняем запрос к сервису комментариев, используя структуру comment
	createdCommentID, err := comments.AddComment(comment.NewsID, comment.Text, comment.ParentCommentID)
	if err != nil {
		http.Error(w, "Ошибка при создании комментария: "+err.Error(), http.StatusInternalServerError)
		return
	}
	createdComment, err := comments.GetComment(createdCommentID)
	if err != nil {
		http.Error(w, "Ошибка при извлечении комментария: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Кодируем созданный комментарий в формат JSON
	responseJSON, err := json.Marshal(createdComment)
	if err != nil {
		http.Error(w, "Ошибка при кодировании ответа: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type и записываем ответ в тело HTTP-ответа
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func SearchNewsHandler(w http.ResponseWriter, r *http.Request) {
	page, pageSize, query := extractParametersFromRequest(r)
	requestID := uuid.New().String()

	// Устанавливаем Request ID в заголовке запроса
	r.Header.Set("X-Request-ID", requestID)

	logs.LogRequest(r.Header.Get("X-Request-ID"), r.RemoteAddr, http.StatusOK)

	matchingNews, err := news.SearchNews(page, pageSize, query)
	if err != nil {
		http.Error(w, "Ошибка при поиске новостей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(matchingNews)
	if err != nil {
		http.Error(w, "Ошибка при кодировании ответа: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func extractParametersFromRequest(r *http.Request) (int, int, string) {
	// Получаем параметры запроса из URL
	values := r.URL.Query()

	// Извлекаем параметр "page"
	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		// Обработка ошибки или установка значения по умолчанию
		page = 1
	}

	// Извлекаем параметр "pageSize"
	pageSize, err := strconv.Atoi(values.Get("pageSize"))
	if err != nil {
		// Обработка ошибки или установка значения по умолчанию
		pageSize = 10
	}

	// Извлекаем параметр "query"
	query := values.Get("query")

	return page, pageSize, query
}
*/
