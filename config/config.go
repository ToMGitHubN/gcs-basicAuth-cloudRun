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
package config

import (
	"context"
	"net/http"

	"github.com/ToMGitHubN/gcs-basicAuth-cloudRun/backends/gcs"
	"github.com/ToMGitHubN/gcs-basicAuth-cloudRun/backends/proxy"
)

// Setup will be called once at the start of the program.
func Setup() error {
	return gcs.Setup()
}

// GET will be called in main.go for GET requests
func GET(ctx context.Context, output http.ResponseWriter, input *http.Request) {
	gcs.Read(ctx, output, input, LoggingOnly)
}

// HEAD will be called in main.go for HEAD requests
func HEAD(ctx context.Context, output http.ResponseWriter, input *http.Request) {
	gcs.ReadMetadata(ctx, output, input, LoggingOnly)
}

// func POST

// func DELETE

// OPTIONS will be called in main.go for OPTIONS requests
func OPTIONS(ctx context.Context, output http.ResponseWriter, input *http.Request) {
	proxy.SendOptions(ctx, output, input, LoggingOnly)
}
