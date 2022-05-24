// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package proxy

import (
	"net/url"
	"testing"
)

func TestRewriteURL(t *testing.T) {
	type input struct {
		proxy    *Proxy
		u        url.URL
		fromRepo string
		toRepo   string
		host     string
		tag      string
	}
	tests := []struct {
		Name string
		in   input
		u    url.URL
	}{
		{
			Name: "updates alias reference when it is cross blob mount",
			in: input{
				proxy: &Proxy{
					Host: url.URL{Host: "localhost:8080", Scheme: "http"},
					Aliases: map[string]Repo{
						"base": {
							Repo: "/gitpod/base-images",
						},
						"target": {
							Repo: "/gitpod/workspace-images",
						},
					},
					proxies: nil,
				},
				fromRepo: "base",
				toRepo:   "/gitpod/base-images",
				host:     "registry.gitlab.com",
				tag:      "tag12345",
				u: url.URL{
					Host:     "localhost.com",
					RawQuery: "/mounts/uploads/?mount=sha:12345?from=base",
				},
			},
			u: url.URL{
				Host:     "registry.gitlab.com",
				RawQuery: "/mounts/uploads/?mount=sha:12345?from=/gitpod/base-images",
			},
		},
		{
			Name: "does not update alias reference when it is not cross blob mount",
			in: input{
				proxy: &Proxy{
					Host: url.URL{Host: "localhost:8080", Scheme: "http"},
					Aliases: map[string]Repo{
						"base": {
							Repo: "/gitpod/base-images",
						},
						"target": {
							Repo: "/gitpod/workspace-images",
						},
					},
				},
				fromRepo: "base",
				toRepo:   "/gitpod/base-images",
				host:     "registry.gitlab.com",
				tag:      "tag12345",
				u: url.URL{
					Host:     "localhost.com",
					RawQuery: "/blobs/sha:12345base",
				},
			},
			u: url.URL{
				Host:     "registry.gitlab.com",
				RawQuery: "/blobs/sha:12345base",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.in.proxy.rewriteURL(&test.in.u, test.in.fromRepo, test.in.toRepo, test.in.host, test.in.tag)
			if test.in.u.RawQuery != test.u.RawQuery {
				t.Errorf("expected raw query: %s but got %s", test.u.RawQuery, test.in.u.RawQuery)
			}
		})
	}

}
