package docs

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Chat Server API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    body { margin: 0; background: #f7f7f7; }
    #auth-bar {
      max-width: 1200px;
      margin: 0 auto;
      padding: 12px 16px;
      background: #fff;
      border-bottom: 1px solid #eee;
      font-family: Arial, sans-serif;
      font-size: 14px;
    }
    #auth-token {
      width: 60%;
      max-width: 720px;
      padding: 6px 8px;
      border: 1px solid #ddd;
      border-radius: 4px;
    }
    #swagger-ui { max-width: 1200px; margin: 0 auto; }
  </style>
</head>
<body>
  <div id="auth-bar">
    Authorization: <input id="auth-token" placeholder="Bearer <token> (leave empty for public endpoints)" />
  </div>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    function getAuthToken() {
      var input = document.getElementById("auth-token");
      if (!input) return "";
      return input.value.trim();
    }

    window.onload = function () {
      SwaggerUIBundle({
        url: "/docs/swagger.json",
        dom_id: "#swagger-ui",
        deepLinking: true,
        persistAuthorization: true,
        requestInterceptor: function (req) {
          req.headers = req.headers || {};
          req.headers["Connect-Protocol-Version"] = "1";
          var token = getAuthToken();
          if (token) {
            if (token.indexOf("Bearer ") !== 0) {
              token = "Bearer " + token;
            }
            req.headers["Authorization"] = token;
          }
          return req;
        }
      });
    };
  </script>
</body>
</html>`

func Register(router *gin.Engine) {
	router.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerHTML))
	})
	router.GET("/docs/swagger.json", func(c *gin.Context) {
		if swaggerPath, ok := findSwaggerPath(); ok {
			c.File(swaggerPath)
			return
		}
		c.String(http.StatusNotFound, "swagger json not found, run: buf generate")
	})
}

func findSwaggerPath() (string, bool) {
	candidates := []string{
		filepath.Join("docs", "swagger", "chat_server.swagger.json"),
		filepath.Join("..", "docs", "swagger", "chat_server.swagger.json"),
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}
