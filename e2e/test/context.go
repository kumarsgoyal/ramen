// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

// nolint: thelper // for using dt *testing.T and keeping test code idiomatic.
package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ramendr/ramen/e2e/dractions"
	"github.com/ramendr/ramen/e2e/types"
	"go.uber.org/zap"
)

type Context struct {
	workload types.Workload
	deployer types.Deployer
	name     string
	logger   *zap.SugaredLogger
}

func NewContext(w types.Workload, d types.Deployer, log *zap.SugaredLogger) Context {
	name := strings.ToLower(d.GetName() + "-" + w.GetName() + "-" + w.GetAppName())

	return Context{
		workload: w,
		deployer: d,
		name:     name,
		logger:   log.Named(name),
	}
}

func (c *Context) Deployer() types.Deployer {
	return c.deployer
}

func (c *Context) Workload() types.Workload {
	return c.workload
}

func (c *Context) Name() string {
	return c.name
}

func (c *Context) Namespace() string {
	if ns := c.deployer.GetNamespace(); ns != "" {
		return ns
	}

	return c.name
}

func (c *Context) Logger() *zap.SugaredLogger {
	return c.logger
}

// Validated return an error if the combination of deployer and workload is not supported.
// TODO: validate that the workload is compatible with the clusters.
func (c *Context) Validate() error {
	if !c.deployer.IsWorkloadSupported(c.workload) {
		return fmt.Errorf("workload %q not supported by deployer %q", c.workload.GetName(), c.deployer.GetName())
	}

	return nil
}

func (c *Context) Deploy(dt *testing.T) {
	t := WithLog(dt, c.logger)
	t.Helper()

	if err := c.deployer.Deploy(c); err != nil {
		t.Fatalf("Failed to deploy workload: %s", err)
	}
}

func (c *Context) Undeploy(dt *testing.T) {
	t := WithLog(dt, c.logger)
	t.Helper()

	if err := c.deployer.Undeploy(c); err != nil {
		t.Fatalf("Failed to undeploy workload: %s", err)
	}
}

func (c *Context) Enable(dt *testing.T) {
	t := WithLog(dt, c.logger)
	t.Helper()

	if err := dractions.EnableProtection(c); err != nil {
		t.Fatalf("Failed to enable protection for workload: %s", err)
	}
}

func (c *Context) Disable(dt *testing.T) {
	t := WithLog(dt, c.logger)
	t.Helper()

	if err := dractions.DisableProtection(c); err != nil {
		t.Fatalf("Failed to disable protection for workload: %s", err)
	}
}

func (c *Context) Failover(dt *testing.T) {
	t := WithLog(dt, c.logger)
	t.Helper()

	if err := dractions.Failover(c); err != nil {
		t.Fatalf("Failed to failover workload: %s", err)
	}
}

func (c *Context) Relocate(dt *testing.T) {
	t := WithLog(dt, c.logger)
	t.Helper()

	if err := dractions.Relocate(c); err != nil {
		t.Fatalf("Failed to relocate workload: %s", err)
	}
}
