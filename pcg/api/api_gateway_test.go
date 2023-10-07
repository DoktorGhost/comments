package api_gateway_test

import (
	"APIGateway/db"
	api_gateway "APIGateway/pcg/api"
	"APIGateway/pcg/comments"
	"APIGateway/pcg/news"
	"APIGateway/pcg/types"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAggregateRequests(t *testing.T) {

	db.InitDB()
	// Инициализация тестовой базы данных и добавление данных
	newsID, _ := news.AddNews("Test News", "This is a test news.")
	commentId1, _ := comments.AddComment(newsID, "Comment 1", 0)
	commentId2, _ := comments.AddComment(newsID, "Comment 2", 0)

	// Проведем тест функции AggregateRequests
	newsInfo, commentsList, err := api_gateway.AggregateRequests(newsID)

	// Проверим, что ошибка равна nil
	if err != nil {
		t.Fatalf("Ошибка в AggregateRequests: %v", err)
	}

	// Проверим, что newsInfo имеет ожидаемый заголовок
	expectedTitle := "Test News"
	if newsInfo.Title != expectedTitle {
		t.Fatalf("Заголовок новости не соответствует ожидаемому: ожидаемый=%s, полученный=%s", expectedTitle, newsInfo.Title)
	} else {
		log.Printf("Заголовок новости соответствует ожидаемому: ожидаемый=%s, полученный=%s", expectedTitle, newsInfo.Title)
	}

	// Проверим, что commentsList содержит два комментария
	if len(commentsList) != 2 {
		t.Fatalf("Количество комментариев не соответствует ожидаемому: ожидаемый=2, полученный=%d", len(commentsList))
	} else {
		log.Printf("Количество комментариев соответствует ожидаемому: ожидаемый=2, полученный=%d", len(commentsList))
	}

	//теперь удалим новость и комменатрии
	err = news.DeleteNews(newsID)
	if err != nil {
		t.Fatalf("Ошибка при удалении новости: %v", err)
	}

	err = comments.DeleteComment(commentId1)

	if err != nil {
		t.Fatalf("Ошибка при удалении комментария: %v", err)
	} else {
		log.Printf("Комментарий с ID=%v удален", commentId1)
	}

	err = comments.DeleteComment(commentId2)

	if err != nil {
		t.Fatalf("Ошибка при удалении комментария: %v", err)
	} else {
		log.Printf("Комментарий с ID=%v удален", commentId2)
	}

}

func TestCreateCommentHandler(t *testing.T) {
	// Создаем тестовый HTTP сервер
	server := httptest.NewServer(http.HandlerFunc(api_gateway.CreateCommentHandler))
	defer server.Close()

	db.InitDB()
	// Инициализация тестовой базы данных и добавление данных
	newsID, _ := news.AddNews("Test News", "This is a test news.")

	// Подготавливаем JSON-запрос
	commentData := struct {
		NewsID   int    `json:"news_id"`
		Text     string `json:"text"`
		ParentID int    `json:"parent_id"`
	}{
		NewsID:   newsID,
		Text:     "Test comment",
		ParentID: 0,
	}

	commentJSON, err := json.Marshal(commentData)
	if err != nil {
		t.Fatalf("Ошибка при кодировании JSON-запроса: %v", err)
	}

	// Отправляем POST-запрос на тестовый сервер
	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(commentJSON))
	if err != nil {
		t.Fatalf("Ошибка при отправке POST-запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем код состояния HTTP-ответа
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался код состояния 200, получен: %v", resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка при чтении тела ответа: %v", err)
	}

	//fmt.Println("Response Body:", string(bodyBytes))

	// Декодируем JSON-ответ
	var createdComment types.Comment

	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
	err = decoder.Decode(&createdComment)
	//fmt.Println("JSON Response:", createdComment)
	if err != nil {
		t.Fatalf("Ошибка при декодировании JSON-ответа: %v", err)
	}

	// Проверяем созданный комментарий
	if createdComment.ID == 0 {
		t.Fatalf("Ожидался ненулевой ID для созданного комментария")
	}

	//теперь удалим новость и комменатрии
	err = news.DeleteNews(newsID)
	if err != nil {
		t.Fatalf("Ошибка при удалении новости: %v", err)
	}

	err = comments.DeleteComment(createdComment.ID)

	if err != nil {
		t.Fatalf("Ошибка при удалении комментария: %v", err)
	}

}

func TestSearchNewsHandler(t *testing.T) {
	// Создаем тестовый HTTP сервер с обработчиком SearchNewsByTitleHandler
	server := httptest.NewServer(http.HandlerFunc(api_gateway.SearchNewsHandler))
	defer server.Close()

	db.InitDB()

	// Добавим несколько новостей
	newsID_1, _ := news.AddNews("Golang один из лучших языков programming", "This is a test news.")
	newsID_2, _ := news.AddNews("pRograMMing without brain", "This is a test news.")
	newsID_3, _ := news.AddNews("а в этой новости нет поискового запроса!", "This is a test news.")

	defer func() {
		_ = news.DeleteNews(newsID_1)
		_ = news.DeleteNews(newsID_2)
		_ = news.DeleteNews(newsID_3)
	}()

	// Подготавливаем тестовые запросы
	// Тест 1: Поиск с валидными ключевыми словами
	query1 := "golang"
	url1 := server.URL + "/news/search?page=1&pageSize=8&query=" + url.QueryEscape(query1)

	// Тест 2: Поиск с ключевыми словами, которых в новостях нет
	query2 := "apple"
	url2 := server.URL + "/news/search?page=1&pageSize=8&query=" + url.QueryEscape(query2)

	// Тест 3: Пустой запрос
	query3 := ""
	url3 := server.URL + "/news/search?page=1&pageSize=8&query=" + url.QueryEscape(query3)

	// Отправляем GET-запросы на тестовый сервер и проверяем результаты

	// Тест 1: Ожидаем результаты поиска
	resp1, err1 := http.Get(url1)
	if err1 != nil {
		t.Fatalf("Ошибка при выполнении GET-запроса: %v", err1)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался код состояния 200, получен: %v", resp1.Status)
	}

	var response1 []types.News
	decoder1 := json.NewDecoder(resp1.Body)
	if err := decoder1.Decode(&response1); err != nil {
		t.Fatalf("Ошибка при декодировании JSON-ответа: %v", err)
	}

	if len(response1) == 0 {
		t.Fatal("Ожидалось получить результаты поиска, но получен пустой ответ")
	}

	// Проверяем, что результаты соответствуют ожиданиям (поиск по ключевым словам)
	for _, news := range response1 {
		containsKeywords := false
		if strings.Contains(strings.ToLower(news.Title), query1) {
			containsKeywords = true
			break
		}

		if !containsKeywords {
			t.Fatalf("Полученная новость не содержит одно из ключевых слов: %v", news.Title)
		}
	}

	// Тест 2: Ожидаем пустой результат
	resp2, err2 := http.Get(url2)
	if err2 != nil {
		t.Fatalf("Ошибка при выполнении GET-запроса: %v", err2)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался код состояния 200, получен: %v", resp2.Status)
	}

	var response2 []types.News
	decoder2 := json.NewDecoder(resp2.Body)
	if err := decoder2.Decode(&response2); err != nil {
		t.Fatalf("Ошибка при декодировании JSON-ответа: %v", err)
	}

	if len(response2) > 0 {
		t.Fatal("Ожидалось получить пустой результат, но получены результаты поиска")
	}

	// Тест 3: Пустой запрос - ожидаем все новости
	resp3, err3 := http.Get(url3)
	if err3 != nil {
		t.Fatalf("Ошибка при выполнении GET-запроса: %v", err3)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался код состояния 200, получен: %v", resp3.Status)
	}

	var response3 []types.News
	decoder3 := json.NewDecoder(resp3.Body)
	if err := decoder3.Decode(&response3); err != nil {
		t.Fatalf("Ошибка при декодировании JSON-ответа: %v", err)
	}
	if len(response3) != 3 {
		t.Fatal("Ожидалось получить 3 новостей, а получено ", len(response3))
	}
}
