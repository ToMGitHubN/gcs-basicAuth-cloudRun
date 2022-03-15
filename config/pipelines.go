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
	"time"

	"github.com/ToMGitHubN/gcs-basicAuth-cloudRun/filter"
	gocache "github.com/patrickmn/go-cache"
)

// DEFAULT: A proxy that simply logs requests.
var LoggingOnly = filter.Pipeline{
	filter.LogRequest,
}

// EXAMPLE: No funny stuff.
var NoFilters = filter.Pipeline{}

// EXAMPLE: Send everything compressed.
var ZippingProxy = filter.Pipeline{
	filter.GZip,
	filter.LogRequest,
}

// EXAMPLE: Cache media in the proxy's memory.
var CacheMedia = filter.Pipeline{
	cacheMedia,
	filter.LogRequest,
}

// mediaCache is a cache for media.
// TODO(domz): Gonna need some memory bounds here.
var mediaCache = gocache.New(5*time.Minute, 10*time.Minute)

// cacheSetter matches the filter.CacheSet type.
func cacheSetter(k string, b []byte, d time.Duration) {
	mediaCache.Set(k, b, d)
}

// cacheSetter matches the gcs.CacheGet type.
// Basically, we have to deal with the conversion from ifc/nil to []byte here.
func cacheGetter(k string) ([]byte, bool) {
	ifc, hit := mediaCache.Get(k)
	if hit {
		return ifc.([]byte), true
	}
	return []byte{}, false
}

// cacheMedia applies mediaCache to the FillCache filter.
func cacheMedia(c context.Context, mfh filter.MediaFilterHandle) error {
	return filter.FillCache(c, mfh, cacheSetter)
}
