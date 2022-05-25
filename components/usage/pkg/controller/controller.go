// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package controller

import (
	"fmt"
	"github.com/gitpod-io/gitpod/common-go/log"
	"github.com/robfig/cron"
	"sync"
	"time"
)

func New(schedule string, reconciler Reconciler) (*Controller, error) {
	sched, err := cron.Parse(schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cron schedule: %w", err)
	}

	return &Controller{
		schedule:   sched,
		reconciler: reconciler,
		scheduler:  cron.NewWithLocation(time.UTC),
	}, nil
}

type Controller struct {
	schedule   cron.Schedule
	reconciler Reconciler

	scheduler *cron.Cron

	runningJobs sync.WaitGroup
}

func (c *Controller) Start() {
	log.
		c.scheduler.Schedule(c.schedule, cron.FuncJob(func() {
		c.runningJobs.Add(1)
		defer c.runningJobs.Done()

		err := c.reconciler.Reconcile()
		if err != nil {
			log.WithError(err).Errorf("Controller run failed.")
		}
	}))

	c.scheduler.Start()
}

// Stop terminates the Controller and awaits for all running jobs to complete.
func (c *Controller) Stop() {
	// Stop any new jobs from running
	c.scheduler.Stop()

	// Wait for existing jobs to finish
	c.runningJobs.Wait()
}
