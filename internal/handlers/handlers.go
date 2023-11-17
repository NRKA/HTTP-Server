package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NRKA/HTTP-Server/internal/kafka"
	"github.com/NRKA/HTTP-Server/internal/repository"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	queryParamKey = "key"
	topic         = "TOPIC"
)

type articleInterface interface {
	Create(ctx context.Context, article repository.Article) (int64, error)
	GetByID(ctx context.Context, id int64) (repository.Article, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, article repository.Article) error
}
type articleHandler struct {
	repo        articleInterface
	producer    kafka.KafkaInterface
	currentTime func() time.Time
}

func NewArticleHandler(repo articleInterface, producer kafka.KafkaInterface) *articleHandler {
	return &articleHandler{
		repo:        repo,
		producer:    producer,
		currentTime: time.Now,
	}
}
func (handler *articleHandler) SetCustomTimeFunc(timeFunc func() time.Time) {
	handler.currentTime = timeFunc
}

func CreateRouter(handler *articleHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/article", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			err, status := handler.Create(w, req)
			if err != nil {
				handler.handleError(w, status, err.Error())
			}
		case http.MethodPut:
			err, status := handler.Update(w, req)
			if err == nil {
				handler.handleError(w, status, err.Error())
			}
		}
	}).Methods(http.MethodPost, http.MethodPut)

	router.HandleFunc(fmt.Sprintf("/article/{%s:[0-9]+}", queryParamKey),
		func(w http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodGet:
				err, status := handler.GetByID(w, req)
				if err != nil {
					handler.handleError(w, status, err.Error())
				}

			case http.MethodDelete:
				err, status := handler.Delete(w, req)
				if err != nil {
					handler.handleError(w, status, err.Error())
				}
			}
		}).Methods(http.MethodGet, http.MethodDelete)

	return router
}
func (handler *articleHandler) handleError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(errorMessage))
}

func (handler *articleHandler) Create(w http.ResponseWriter, req *http.Request) (error, int) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf(errReadReqBody + err.Error()), http.StatusInternalServerError
	}
	var articleData repository.Article
	if err = json.Unmarshal(body, &articleData); err != nil {
		return fmt.Errorf(errParseJson + err.Error()), http.StatusInternalServerError
	}
	if articleData.Name == "" || articleData.Rating < 1 {
		return fmt.Errorf(errInvalidData), http.StatusBadRequest
	}
	id, err := handler.repo.Create(req.Context(), articleData)
	if err != nil {
		return fmt.Errorf(err.Error()), http.StatusInternalServerError
	}
	articleData.ID = id
	articleJson, err := json.Marshal(articleData)
	if err != nil {
		return fmt.Errorf(errCreateJson + err.Error()), http.StatusInternalServerError
	}

	err = handler.producer.SendEvent(os.Getenv(topic), kafka.Event{
		TimeStamp:   handler.currentTime(),
		Type:        req.Method,
		RequestBody: string(body),
	})
	if err != nil {
		log.Println(errSendEvent)
	}
	handler.handleError(w, http.StatusOK, string(articleJson))
	return nil, http.StatusOK
}

func (handler *articleHandler) GetByID(w http.ResponseWriter, req *http.Request) (error, int) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		return fmt.Errorf(errQueryParamKey), http.StatusBadRequest
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return fmt.Errorf(errParseInt + err.Error()), http.StatusBadRequest
	}
	article, err := handler.repo.GetByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrArticalNotFound) {
			return fmt.Errorf(err.Error()), http.StatusNotFound
		}
		return fmt.Errorf(errArticleGetById + err.Error()), http.StatusInternalServerError
	}

	articleJson, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf(errCreateJson + err.Error()), http.StatusInternalServerError
	}
	err = handler.producer.SendEvent(os.Getenv(topic), kafka.Event{
		TimeStamp:   handler.currentTime(),
		Type:        req.Method,
		RequestBody: "",
	})
	if err != nil {
		log.Println(errSendEvent)
	}
	handler.handleError(w, http.StatusOK, string(articleJson))
	return nil, http.StatusOK
}

func (handler *articleHandler) Delete(w http.ResponseWriter, req *http.Request) (error, int) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		return fmt.Errorf(errQueryParamKey), http.StatusBadRequest
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return fmt.Errorf(errParseInt + err.Error()), http.StatusBadRequest
	}

	err = handler.repo.Delete(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrArticalNotFound) {
			return fmt.Errorf(errArticleNotFound + err.Error()), http.StatusNotFound
		}
		return fmt.Errorf(errArticleDelete + err.Error()), http.StatusInternalServerError
	}

	err = handler.producer.SendEvent(os.Getenv(topic), kafka.Event{
		TimeStamp:   handler.currentTime(),
		Type:        req.Method,
		RequestBody: "",
	})
	if err != nil {
		log.Println(errSendEvent)
	}
	handler.handleError(w, http.StatusOK, "")
	return nil, http.StatusOK
}

func (handler *articleHandler) Update(w http.ResponseWriter, req *http.Request) (error, int) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf(errReadReqBody + err.Error()), http.StatusInternalServerError
	}
	var articleData repository.Article
	if err = json.Unmarshal(body, &articleData); err != nil {
		return fmt.Errorf(errParseJson + err.Error()), http.StatusInternalServerError
	}
	if articleData.Name == "" || articleData.Rating < 1 {
		return fmt.Errorf(errInvalidData), http.StatusBadRequest
	}
	err = handler.repo.Update(req.Context(), articleData)
	if err != nil {
		if errors.Is(err, repository.ErrArticalNotFound) {
			return fmt.Errorf(errArticleNotFound + err.Error()), http.StatusNotFound
		}
		return fmt.Errorf(errArticleUpdate + err.Error()), http.StatusInternalServerError
	}

	err = handler.producer.SendEvent(os.Getenv(topic), kafka.Event{
		TimeStamp:   handler.currentTime(),
		Type:        req.Method,
		RequestBody: string(body),
	})
	if err != nil {
		log.Println(errSendEvent)
	}

	handler.handleError(w, http.StatusOK, "")
	return nil, http.StatusOK
}
