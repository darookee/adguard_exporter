// Copyright (C) 2016 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adguard

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/common/log"
)

type adguardserver struct {
	*httptest.Server
}

func handler(server *adguardserver, uri string, filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(b)
	}
}

func newadguardServer(uri, filename string) *adguardserver {
	h := &adguardserver{}
	h.Server = httptest.NewServer(handler(h, uri, filename))
	return h
}

func getClientAndServer(t *testing.T, uri, filename string) (*adguardserver, *Client) {
	h := newadguardServer(uri, filename)
	client, err := NewClient(h.URL)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return h, client
}

func TestadguardGetMetrics(t *testing.T) {
	server, client := getClientAndServer(t, "", "testdata/stats.json")
	defer server.Close()
	metrics, err := client.GetMetrics()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Infof("Metrics response: %s", metrics)
	if metrics.dnsQueries != 100000 {
		log.Fatalf("Invalid metrics response: %s", metrics)
	}
}
