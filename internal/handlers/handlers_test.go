package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NRKA/HTTP-Server/internal/kafka"
	mock_kafka_interface "github.com/NRKA/HTTP-Server/internal/kafka/mocks"
	"github.com/NRKA/HTTP-Server/internal/repository"
	mock_repository "github.com/NRKA/HTTP-Server/internal/repository/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestArticleHandler_Create(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name             string
		mockArguments    repository.Article
		mockReturnValue  int64
		mockError        error
		expectedCode     int
		expectedResponse repository.Article
		mockKafka        func(*gomock.Controller, kafka.Event) kafka.KafkaInterface
	}{{
		name:             "success",
		mockArguments:    repository.Article{ID: 1, Name: "name", Rating: 10},
		mockReturnValue:  int64(1),
		mockError:        nil,
		expectedCode:     http.StatusOK,
		expectedResponse: repository.Article{ID: 1, Name: "name", Rating: 10},
		mockKafka: func(controller *gomock.Controller, event kafka.Event) kafka.KafkaInterface {
			mockProducer := mock_kafka_interface.NewMockKafkaInterface(controller)
			mockProducer.EXPECT().SendEvent("crud", event).Return(nil)
			return mockProducer
		},
	},
		{
			name:             "internal server error",
			mockArguments:    repository.Article{ID: 1, Name: "name", Rating: 10},
			mockReturnValue:  int64(1),
			mockError:        fmt.Errorf("failed to create article: internal server error"),
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: repository.Article{},
			mockKafka: func(controller *gomock.Controller, event kafka.Event) kafka.KafkaInterface {
				return mock_kafka_interface.NewMockKafkaInterface(controller)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			data, err := json.Marshal(tc.mockArguments)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/article", bytes.NewReader(data))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockArticleInterface(ctrl)
			mockKafka := tc.mockKafka(ctrl, kafka.Event{
				TimeStamp:   time.Date(2023, 10, 22, 22, 22, 22, 22, time.Local),
				Type:        req.Method,
				RequestBody: string(data),
			})
			handler := NewArticleHandler(mockRepo, mockKafka)
			mockRepo.EXPECT().Create(gomock.Any(), repository.Article{
				ID:        1,
				Name:      "name",
				Rating:    10,
				CreatedAt: time.Time{},
			}).Return(tc.mockReturnValue, tc.mockError)
			defer ctrl.Finish()

			handler.SetCustomTimeFunc(func() time.Time {
				return time.Date(2023, 10, 22, 22, 22, 22, 22, time.Local)
			})
			err, status := handler.Create(rr, req)
			require.Equal(t, tc.mockError, err)

			assert.Equal(t, tc.expectedCode, status)
			if status != http.StatusOK {
				assert.Equal(t, tc.mockError, err)
				return
			}
			var actual repository.Article
			err = json.Unmarshal(rr.Body.Bytes(), &actual)
			require.NoError(t, err)
			assert.Equal(t, actual, tc.expectedResponse)
		})
	}
}

func TestArticleHandler_GetByID(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name              string
		mockArgument      int64
		mockReturnArticle repository.Article
		mockError         error
		expectedCode      int
		expectedResponse  repository.Article
		mockKafka         func(*gomock.Controller, kafka.Event) kafka.KafkaInterface
	}{{
		name:              "success",
		mockArgument:      1,
		mockReturnArticle: repository.Article{ID: 1, Name: "name", Rating: 10},
		mockError:         nil,
		expectedCode:      http.StatusOK,
		expectedResponse:  repository.Article{ID: 1, Name: "name", Rating: 10},
		mockKafka: func(controller *gomock.Controller, event kafka.Event) kafka.KafkaInterface {
			mockProducer := mock_kafka_interface.NewMockKafkaInterface(controller)
			mockProducer.EXPECT().SendEvent("crud", event).Return(nil)
			return mockProducer
		}}, {
		name:              "article not found",
		mockArgument:      99,
		mockReturnArticle: repository.Article{},
		mockError:         repository.ErrArticalNotFound,
		expectedCode:      http.StatusNotFound,
		expectedResponse:  repository.Article{},
		mockKafka: func(controller *gomock.Controller, event kafka.Event) kafka.KafkaInterface {
			return mock_kafka_interface.NewMockKafkaInterface(controller)
		}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockArticleInterface(ctrl)
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/article/%d", tc.mockArgument), bytes.NewReader([]byte{}))
			require.NoError(t, err)

			mockKafka := tc.mockKafka(ctrl, kafka.Event{
				TimeStamp:   time.Date(2023, 10, 22, 22, 22, 22, 22, time.Local),
				Type:        req.Method,
				RequestBody: "",
			})
			handler := NewArticleHandler(mockRepo, mockKafka)
			mockRepo.EXPECT().GetByID(gomock.Any(), tc.mockArgument).Return(tc.mockReturnArticle, tc.mockError)
			defer ctrl.Finish()

			// act
			req = mux.SetURLVars(req, map[string]string{queryParamKey: strconv.Itoa(int(tc.mockArgument))})
			rr := httptest.NewRecorder()

			handler.SetCustomTimeFunc(func() time.Time {
				return time.Date(2023, 10, 22, 22, 22, 22, 22, time.Local)
			})
			err, status := handler.GetByID(rr, req)
			assert.Equal(t, tc.expectedCode, status)

			if status != http.StatusOK {
				assert.Equal(t, tc.mockError, err)
				return
			}
			var actual repository.Article
			err = json.Unmarshal(rr.Body.Bytes(), &actual)
			require.NoError(t, err)
			assert.Equal(t, actual, tc.expectedResponse)
		})
	}
}

func TestArticleHandler_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		mockArgument int64
		mockError    error
		expectedCode int
		mockKafka    func(*gomock.Controller, kafka.Event) kafka.KafkaInterface
	}{{
		name:         "success",
		mockArgument: 1,
		mockError:    nil,
		expectedCode: http.StatusOK,
		mockKafka: func(controller *gomock.Controller, event kafka.Event) kafka.KafkaInterface {
			mockProducer := mock_kafka_interface.NewMockKafkaInterface(controller)
			mockProducer.EXPECT().SendEvent("crud", event).Return(nil)
			return mockProducer
		},
	}, {
		name:         "article not found",
		mockArgument: 9999,
		mockError:    repository.ErrArticalNotFound,
		expectedCode: http.StatusNotFound,
		mockKafka: func(controller *gomock.Controller, event kafka.Event) kafka.KafkaInterface {
			return mock_kafka_interface.NewMockKafkaInterface(controller)
		},
	}, {
		name:         "internal server error",
		mockArgument: 9999,
		mockError:    fmt.Errorf("failed to delete article"),
		expectedCode: http.StatusInternalServerError,
		mockKafka: func(controller *gomock.Controller, event kafka.Event) kafka.KafkaInterface {
			return mock_kafka_interface.NewMockKafkaInterface(controller)
		},
	},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockArticleInterface(ctrl)
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/article/%d", tc.mockArgument), bytes.NewReader([]byte{}))
			require.NoError(t, err)
			mockKafka := tc.mockKafka(ctrl, kafka.Event{
				TimeStamp:   time.Date(2023, 10, 22, 22, 22, 22, 22, time.Local),
				Type:        req.Method,
				RequestBody: "",
			})
			handler := NewArticleHandler(mockRepo, mockKafka)
			mockRepo.EXPECT().Delete(gomock.Any(), tc.mockArgument).Return(tc.mockError)
			defer ctrl.Finish()

			req = mux.SetURLVars(req, map[string]string{queryParamKey: strconv.Itoa(int(tc.mockArgument))})
			rr := httptest.NewRecorder()

			handler.SetCustomTimeFunc(func() time.Time {
				return time.Date(2023, 10, 22, 22, 22, 22, 22, time.Local)
			})
			err, status := handler.Delete(rr, req)
			assert.Equal(t, tc.expectedCode, status)
		})
	}
}

func TestArticleHandler_Update(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		mockArguments repository.Article
		mockError     error
		expectedCode  int
		mockKafka     func(*gomock.Controller, kafka.Event) kafka.KafkaInterface
	}{{
		name:          "success",
		mockArguments: repository.Article{ID: 123, Name: "name", Rating: 10},
		mockError:     nil,
		expectedCode:  http.StatusOK,
		mockKafka: func(controller *gomock.Controller, event kafka.Event) kafka.KafkaInterface {
			mockProducer := mock_kafka_interface.NewMockKafkaInterface(controller)
			mockProducer.EXPECT().SendEvent("crud", event).Return(nil)
			return mockProducer
		},
	}, {
		name:          "article not found",
		mockArguments: repository.Article{ID: 123123, Name: "name", Rating: 10},
		mockError:     repository.ErrArticalNotFound,
		expectedCode:  http.StatusNotFound,
		mockKafka: func(controller *gomock.Controller, event kafka.Event) kafka.KafkaInterface {
			return mock_kafka_interface.NewMockKafkaInterface(controller)
		},
	},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockArticleInterface(ctrl)
			data, err := json.Marshal(tc.mockArguments)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPut, "/article", bytes.NewReader(data))
			require.NoError(t, err)

			mockKafka := tc.mockKafka(ctrl, kafka.Event{
				TimeStamp:   time.Date(2023, 10, 22, 22, 22, 22, 22, time.Local),
				Type:        req.Method,
				RequestBody: string(data),
			})
			handler := NewArticleHandler(mockRepo, mockKafka)
			mockRepo.EXPECT().Update(gomock.Any(), tc.mockArguments).Return(tc.mockError)
			defer ctrl.Finish()

			rr := httptest.NewRecorder()

			handler.SetCustomTimeFunc(func() time.Time {
				return time.Date(2023, 10, 22, 22, 22, 22, 22, time.Local)
			})
			//act
			err, status := handler.Update(rr, req)

			//assert
			assert.Equal(t, tc.expectedCode, status)
		})
	}
}
