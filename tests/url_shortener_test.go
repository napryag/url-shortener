package tests

import (
	"net/http"
	"net/url"
	"path"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gavv/httpexpect/v2"
	"github.com/napryag/url-shortener/internal/http-server/handlers/url/save"
	"github.com/napryag/url-shortener/internal/lib/api"
	"github.com/napryag/url-shortener/internal/lib/random"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8080"
)

func TestURLShortener_HappyPAth(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/url").WithJSON(save.Request{
		URL:   gofakeit.URL(),
		Alias: random.NewRandomString(10),
	}).
		WithBasicAuth("myuser", "mypassword").
		Expect().Status(200).
		JSON().Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirectDelete(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).WithBasicAuth("myuser", "mypassword").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			// Redirect

			testRedirect(t, alias, tc.url)

			// Delete

			req := e.DELETE("/"+path.Join("url", alias)).
				WithBasicAuth("myuser", "mypassword").
				Expect().Status(http.StatusOK).
				JSON().Object()

			req.Value("status").String().IsEqual("OK")

			// Redirect

			testRedirectNotFound(t, alias)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}

func testRedirectNotFound(t *testing.T, alias string) {

	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	_, err := api.GetRedirect(u.String())
	require.ErrorIs(t, err, api.ErrInvalidStatusCode)
}
