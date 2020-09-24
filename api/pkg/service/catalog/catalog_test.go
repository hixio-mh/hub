// Copyright © 2020 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package catalog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tektoncd/hub/api/gen/catalog"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

// NewServiceTest returns the catalog service implementation for test.
func NewServiceTest(api app.Config) catalog.Service {
	svc := auth.NewService(api, "catalog")
	wq := newCatalogSyncer(api)

	s := &service{
		svc,
		wq,
	}
	return s
}

func TestRefresh(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	catalogSvc := NewServiceTest(tc)
	ctx := auth.WithUserID(context.Background(), 11)

	payload := &catalog.RefreshPayload{Name: "catalog-official", Org: "tektoncd"}
	res, err := catalogSvc.Refresh(ctx, payload)
	assert.NoError(t, err)
	assert.Equal(t, 10001, int(res.ID))
	assert.Equal(t, "queued", res.Status)
}

func TestRefreshAgain(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	catalogSvc := NewServiceTest(tc)
	ctx := auth.WithUserID(context.Background(), 11)

	payload := &catalog.RefreshPayload{Name: "catalog-official", Org: "tektoncd"}
	res, err := catalogSvc.Refresh(ctx, payload)
	assert.NoError(t, err)
	assert.Equal(t, 10001, int(res.ID))
	assert.Equal(t, "queued", res.Status)

	res, err = catalogSvc.Refresh(ctx, payload)
	assert.NoError(t, err)
	assert.Equal(t, 10001, int(res.ID))
	assert.Equal(t, "queued", res.Status)
}
