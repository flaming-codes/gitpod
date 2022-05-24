// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package testdb

import (
	"github.com/gitpod-io/gitpod/common-go/db"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
)

func Connect(t *testing.T) *gorm.DB {
	t.Helper()

	conn, err := db.Connect(db.ConnectionParams{
		User:     "gitpod",
		Password: "test",
		Host:     "localhost:3306",
		Database: "gitpod",
	})
	require.NoError(t, err, "must establish a db connection")

	t.Cleanup(func() {
		rawConn, err := conn.DB()
		require.NoError(t, err)

		require.NoError(t, rawConn.Close(), "must close database connection")
	})

	return conn
}
