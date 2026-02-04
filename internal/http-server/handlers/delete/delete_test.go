package delete_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/napryag/url-shortener/internal/http-server/handlers/delete"
	"github.com/napryag/url-shortener/internal/http-server/handlers/delete/mocks"
	"github.com/napryag/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		url       string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			url:   "https://www.google.com/?hl=RU",
			alias: "test_alias",
		},
		//TODO: more test cases
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			input := fmt.Sprintf(`{"alias":"%s"}`, tc.alias)

			req, err := http.NewRequest(http.MethodDelete, "/url", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp delete.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
