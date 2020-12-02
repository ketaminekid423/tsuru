// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package router

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrDynamicRouterNotFound = errors.New("dynamic router not found")
)

type DynamicRouter struct {
	Name   string
	Type   string
	Config map[string]interface{}
}

type DynamicRouterService interface {
	Get(ctx context.Context, name string) (*DynamicRouter, error)
	List(context.Context) ([]DynamicRouter, error)
	Remove(ctx context.Context, name string) error
	Create(context.Context, DynamicRouter) error
	Update(context.Context, DynamicRouter) error
}

type DynamicRouterStorage interface {
	Save(context.Context, DynamicRouter) error
	Get(ctx context.Context, name string) (*DynamicRouter, error)
	List(context.Context) ([]DynamicRouter, error)
	Remove(ctx context.Context, name string) error
}

type HealthcheckData struct {
	Path    string
	Status  int
	Body    string
	TCPOnly bool
}

func (hc *HealthcheckData) String() string {
	if hc.TCPOnly {
		return "tcp only"
	}
	status := ""
	if hc.Status != 0 {
		status = fmt.Sprintf(", status: %d", hc.Status)
	}
	path := hc.Path
	if path == "" {
		path = "/"
	}
	body := hc.Body
	if body != "" {
		body = fmt.Sprintf(", body: %q", body)
	}
	return fmt.Sprintf("path: %q%s%s", path, status, body)
}
