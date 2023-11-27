package pages

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

type Page struct {
    ServerWebSocketUrl string
}

func Index(c echo.Context) error {
    fmt.Println(os.Getenv("SERVER_URL")) 

    page := Page{
        ServerWebSocketUrl: "ws://localhost:127.0.0.1:5173/ws",
    }

	return c.Render(200, "index.html", page)
}
