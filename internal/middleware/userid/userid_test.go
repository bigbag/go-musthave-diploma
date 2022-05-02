package userid

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

var (
	secretKey  = "secret"
	cookieName = "auth"
	contextKey = "userid"
)

var authMiddleware = New(Config{
	CookieName: cookieName,
	ContextKey: contextKey,
	Secret:     secretKey,
})

func makeToken(userID int) string {
	expiresTime := time.Now().Add(time.Hour * 24)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(userID),
		ExpiresAt: expiresTime.Unix(),
	})

	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return ""
	}
	return token
}

func Test_User_Cookie(t *testing.T) {
	const userID = 12
	app := fiber.New()

	app.Use(New())

	app.Get("/", func(c *fiber.Ctx) error {
		userID := c.Locals(contextKey).(string)
		return c.SendString(userID)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: makeToken(userID),
	})

	resp, err := app.Test(req)

	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")

	body, err := ioutil.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)

	utils.AssertEqual(t, strconv.Itoa(userID), string(body))

	defer resp.Body.Close()
}

func Test_User_Cookie_Next_Error(t *testing.T) {
	app := fiber.New()

	app.Use(New())

	app.Get("/", func(c *fiber.Ctx) error {
		return errors.New("next error")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: makeToken(12),
	})

	resp, err := app.Test(req)

	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 500, resp.StatusCode, "Status code")
	utils.AssertEqual(t, "", resp.Header.Get(fiber.HeaderContentEncoding))

	body, err := ioutil.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, "next error", string(body))

	defer resp.Body.Close()
}

func Test_User_Cookie_Next(t *testing.T) {
	app := fiber.New()
	app.Use(New(Config{
		Next: func(_ *fiber.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusNotFound, resp.StatusCode)

	defer resp.Body.Close()
}
