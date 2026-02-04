package redirect_test

import (
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/napryag/url-shortener/internal/http-server/handlers/redirect"
	"github.com/napryag/url-shortener/internal/http-server/handlers/redirect/mocks"
	"github.com/napryag/url-shortener/internal/lib/api"
	"github.com/napryag/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
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
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			require.Equal(t, tc.url, redirectedToURL)
		})
	}
}
