package localware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/labstack/echo/v4"
)

type ClerkDbContext struct {
	DbContext
	Claims     *clerk.SessionClaims
	Clerk_User *clerk.User
}

type AuthorizationParams struct {
	jwt.VerifyParams
	// AuthorizationFailureHandler gets executed when request authorization
	// fails. Pass a custom http.Handler to control the http.Response for
	// invalid authorization. The default is a Response with an empty body
	// and 401 Unauthorized status.
	AuthorizationFailureHandler http.Handler
	// JWKSClient is the jwks.Client that will be used to fetch the
	// JSON Web Key Set. A default client will be used if none is
	// provided.
	JWKSClient *jwks.Client
}

// Rewritten from: https://github.com/clerk/clerk-sdk-go/blob/v2.0.4/http/middleware.go#L38
// The syntax for this middleware does not seem to
// be supported by echo. To solve this problem, the syntax
// on the website and api want there to be an
// Authorization: Bearer <token>... this does not
// exist from what I can tell. What does exist is
// the __session token (on the cookie). This token
// needs to be verified. If it is, then we have
// authentication.
//
// WithHeaderAuthorization checks the Authorization request header
// for a valid Clerk authorization JWT. The token is parsed and verified
// and the active session claims are written to the http.Request context.
// The middleware uses Bearer authentication, so the Authorization header
// is expected to have the following format:
// Authorization: Bearer <token>
//

const session_key = "__session"

func WithHeaderAuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Assert that the context is of type *CustomContext
		cc, ok := c.(*DbContext)

		if !ok {
			err := errors.New("could not validate DbContext")
			c.Logger().Error(err)
			return c.HTML(http.StatusBadRequest, Make404Html())
		}

		params := &AuthorizationParams{}

		if params.Clock == nil {
			params.Clock = clerk.NewClock()
		}

		cookies := c.Request().Cookies()

		authorization := ""

		for i := len(cookies) - 1; i >= 0; i-- {
			if cookies[i].Name == session_key {
				authorization = cookies[i].Value
			}
		}

		if authorization == "" {
			err := errors.New("authorization is nil")
			cc.Logger().Error(err)
			return c.HTML(http.StatusBadRequest, Make404Html())
		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		decoded, err := jwt.Decode(cc.Request().Context(), &jwt.DecodeParams{Token: token})
		if err != nil {
			cc.Logger().Error(err)
			return cc.HTML(http.StatusBadRequest, Make404Html())
		}

		params.Token = token
		if params.JWK == nil {
			params.JWK, err = getJWK(cc.Request().Context(), params.JWKSClient, decoded.KeyID, params.Clock)
			if err != nil {
				cc.Logger().Error(err)
				return cc.HTML(http.StatusBadRequest, Make404Html())
			}
		}

		claims, err := jwt.Verify(cc.Request().Context(), &params.VerifyParams)
		if err != nil {
			cc.Logger().Error(err)
			return cc.HTML(http.StatusBadRequest, Make404Html())
		}

		usr, err := user.Get(cc.Request().Context(), claims.Subject)
		if err != nil {
			cc.Logger().Error(err)
			return cc.HTML(http.StatusBadRequest, Make404Html())
		}

		// Token was verified. Add the session claims to the request context.
		newCtx := &ClerkDbContext{*cc, claims, usr}

		return next(newCtx)
	}
}

// Retrieve the JSON web key for the provided token from the JWKS set.

// Tries a cached value first, but if there's no value or the entry
// has expired, it will fetch the JWK set from the API and cache the
// value.
func getJWK(ctx context.Context, jwksClient *jwks.Client, kid string, clock clerk.Clock) (*clerk.JSONWebKey, error) {
	if kid == "" {
		return nil, fmt.Errorf("missing jwt kid header claim")
	}

	jwk := getCache().get(kid)
	if jwk == nil || !getCache().isValid(kid, clock.Now().UTC()) {
		var err error
		jwk, err = jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
			KeyID:      kid,
			JWKSClient: jwksClient,
		})
		if err != nil {
			return nil, err
		}
	}
	getCache().set(kid, jwk, clock.Now().UTC().Add(time.Hour))
	return jwk, nil
}

// A cache to store JSON Web Keys.
type jwkCache struct {
	mu      sync.RWMutex
	entries map[string]*cacheEntry
}

// Each entry in the JWK cache has a value and an expiration date.
type cacheEntry struct {
	value     *clerk.JSONWebKey
	expiresAt time.Time
}

// IsValid returns true if a non-expired entry exists in the cache
// for the provided key, false otherwise.
func (c *jwkCache) isValid(key string, t time.Time) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	return ok && entry != nil && entry.expiresAt.After(t)
}

// Get fetches the JSON Web Key for the provided key, unless the
// entry has expired.
func (c *jwkCache) get(key string) *clerk.JSONWebKey {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok || entry == nil {
		return nil
	}
	return entry.value
}

// Set stores the JSON Web Key in the provided key and sets the
// expiration date.
func (c *jwkCache) set(key string, value *clerk.JSONWebKey, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = &cacheEntry{
		value:     value,
		expiresAt: expiresAt,
	}
}

var cacheInit sync.Once

// A "singleton" JWK cache for the package.
var cache *jwkCache

// getCache returns the library's default cache singleton.
// Please note that the returned Cache is a package-level variable.
// Using the package with more than one Clerk API secret keys might
// require to use different Clients with their own Cache.
func getCache() *jwkCache {
	cacheInit.Do(func() {
		cache = &jwkCache{
			entries: map[string]*cacheEntry{},
		}
	})
	return cache
}

// Private method to return 404 page html
func Make404Html() string {
	return `
          <!DOCTYPE html>
          <html>

          <head>
            <title>
             404 
            </title>
            <meta charset="UTF-8" />

			      <script async crossorigin="anonymous" data-clerk-publishable-key="pk_test_Zmx5aW5nLWNvcmFsLTgyLmNsZXJrLmFjY291bnRzLmRldiQ" src="https://flying-coral-82.clerk.accounts.dev/npm/@clerk/clerk-js@latest/dist/clerk.browser.js" type="text/javascript"></script>
            <link href="https://cdn.jsdelivr.net/npm/flowbite@3.1.2/dist/flowbite.min.css" rel="stylesheet"/>
			      <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/toastify-js/src/toastify.min.css"/>

            <meta name="viewport" content="width=device-width, initial-scale=1" />
            <!-- Clerk Publishable key and Frontend API URL -->
            <link href="/assets/css/output.css" rel="stylesheet" />
            <!-- This is the flowbite darkmode script -->
          </head>

          <body>

						<button type="button" id="user-menu-button" aria-expanded="false"></button>
            <section class="bg-white dark:bg-gray-900">
              <div class="py-8 px-4 mx-auto max-w-screen-xl text-center lg:py-16">
                <h1
                  class="mb-4 text-4xl font-extrabold tracking-tight leading-none text-gray-900 md:text-5xl lg:text-6xl dark:text-white">
                  Keep Searchin Matey!
                </h1>
                <p class="mb-8 text-lg font-normal text-gray-500 lg:text-xl sm:px-16 lg:px-48 dark:text-gray-400">
                  There be trecherous waters here, route yourself back to safety below!
                </p>
                <div class="flex flex-col space-y-4 sm:flex-row sm:justify-center sm:space-y-0">
                  <a href="https://prate.pro"
                    class="inline-flex justify-center items-center py-3 px-5 text-base font-medium text-center text-white rounded-lg bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 dark:focus:ring-blue-900">
                     Safety 
                    <svg class="w-3.5 h-3.5 ms-2 rtl:rotate-180" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none"
                      viewBox="0 0 14 10">
                      <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M1 5h12m0 0L9 1m4 4L9 9"></path>
                    </svg>
                  </a>
                </div>
              </div>
            </section>
          </body>
          </html>
          <script src="/assets/js/index.js"></script>
			    <script src="/assets/js/dashboard.js"></script>

          `
}
