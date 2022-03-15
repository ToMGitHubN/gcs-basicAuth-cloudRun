//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/ToMGitHubN/gcs-basicAuth-cloudRun/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	// initialize
	log.Print("starting server...")
	handler := http.HandlerFunc(ProxyHTTPGCS)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Warn().Msgf("defaulting to port %s", port)
	}

	// Initialize
	if err := config.Setup(); err != nil {
		log.Fatal().Msgf("main setup: %v", err)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	http2server := &http2.Server{}
	if err := http.ListenAndServe(":"+port, h2c.NewHandler(handler, http2server)); err != nil {
		log.Fatal().Msgf("main: %v", err)
	}
}

// Basic認証
func CheckAuth(input *http.Request) bool {
	// 認証機能を有効にしているか
	if os.Getenv("BASIC_ATUH_ENABLE") != "true" {
		return true
	}

	// 認証情報取得
	clientID, clientSecret, ok := input.BasicAuth()
	if ok == false {
		return false
	}
	return clientID == os.Getenv("BASIC_ATUH_ID") && clientSecret == os.Getenv("BASIC_ATUH_PASSWORD")
}

// ProxyHTTPGCS is the entry point for the cloud function, providing a proxy that
// permits HTTP protocol usage of a GCS bucket's contents.
func ProxyHTTPGCS(output http.ResponseWriter, input *http.Request) {

	// Basic認証チェック
	if CheckAuth(input) == false {
		// 認証失敗時
		output.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		output.WriteHeader(http.StatusUnauthorized) // 401
		http.Error(output, "Unauthorized", 401)
		return
	}

	ctx := context.Background()

	// route HTTP methods to appropriate handlers.
	switch input.Method {
	case http.MethodGet:
		config.GET(ctx, output, input)
	case http.MethodHead:
		config.HEAD(ctx, output, input)
	case http.MethodOptions:
		config.OPTIONS(ctx, output, input)
	default:
		http.Error(output, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
