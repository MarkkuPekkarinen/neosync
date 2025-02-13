package dbconnectconfig

import (
	"fmt"
	"testing"

	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	"github.com/nucleuscloud/neosync/internal/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	pgconnectionFixture = &mgmtv1alpha1.PostgresConnection{
		Host:    "localhost",
		Port:    5432,
		Name:    "postgres",
		User:    "test-user",
		Pass:    "test-pass",
		SslMode: ptr("verify"),
	}
)

func Test_NewFromPostgresConnection(t *testing.T) {
	t.Run("Connection", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			actual, err := NewFromPostgresConnection(
				&mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_Connection{
							Connection: pgconnectionFixture,
						},
					},
				},
				&testConnectionTimeout,
				testutil.GetTestLogger(t),
			)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(
				t,
				"postgres://test-user:test-pass@localhost:5432/postgres?connect_timeout=5&sslmode=verify",
				actual.String(),
			)
			assert.Equal(t, "test-user", actual.GetUser())
		})
		t.Run("ok_no_timeout", func(t *testing.T) {
			actual, err := NewFromPostgresConnection(
				&mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_Connection{
							Connection: pgconnectionFixture,
						},
					},
				},
				nil,
				testutil.GetTestLogger(t),
			)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(
				t,
				"postgres://test-user:test-pass@localhost:5432/postgres?sslmode=verify",
				actual.String(),
			)
			assert.Equal(t, "test-user", actual.GetUser())
		})
		t.Run("ok_no_port", func(t *testing.T) {
			actual, err := NewFromPostgresConnection(
				&mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_Connection{
							Connection: &mgmtv1alpha1.PostgresConnection{
								Host: "localhost",
								// Port:    5432,
								Name:    "postgres",
								User:    "test-user",
								Pass:    "test-pass",
								SslMode: ptr("verify"),
							},
						},
					},
				},
				&testConnectionTimeout,
				testutil.GetTestLogger(t),
			)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(
				t,
				"postgres://test-user:test-pass@localhost/postgres?connect_timeout=5&sslmode=verify",
				actual.String(),
			)
			assert.Equal(t, "test-user", actual.GetUser())
		})
		t.Run("ok_no_pass", func(t *testing.T) {
			actual, err := NewFromPostgresConnection(
				&mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_Connection{
							Connection: &mgmtv1alpha1.PostgresConnection{
								Host: "localhost",
								Port: 5432,
								Name: "postgres",
								User: "test-user",
								// Pass:    "test-pass",
								SslMode: ptr("verify"),
							},
						},
					},
				},
				&testConnectionTimeout,
				testutil.GetTestLogger(t),
			)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(
				t,
				"postgres://test-user@localhost:5432/postgres?connect_timeout=5&sslmode=verify",
				actual.String(),
			)
			assert.Equal(t, "test-user", actual.GetUser())
		})
		t.Run("ok_no_creds", func(t *testing.T) {
			actual, err := NewFromPostgresConnection(
				&mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_Connection{
							Connection: &mgmtv1alpha1.PostgresConnection{
								Host: "localhost",
								Port: 5432,
								Name: "postgres",
								// User:    "test-user",
								// Pass:    "test-pass",
								SslMode: ptr("verify"),
							},
						},
					},
				},
				&testConnectionTimeout,
				testutil.GetTestLogger(t),
			)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(
				t,
				"postgres://localhost:5432/postgres?connect_timeout=5&sslmode=verify",
				actual.String(),
			)
			assert.Equal(t, "", actual.GetUser())
		})
	})

	t.Run("URL", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			actual, err := NewFromPostgresConnection(
				&mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_Url{
							Url: "postgres://test-user:test-pass@localhost:5432/postgres?sslmode=disable",
						},
					},
				},
				&testConnectionTimeout,
				testutil.GetTestLogger(t),
			)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(
				t,
				"postgres://test-user:test-pass@localhost:5432/postgres?connect_timeout=5&sslmode=disable",
				actual.String(),
			)
			assert.Equal(t, "test-user", actual.GetUser())
		})
		t.Run("ok_no_timeout", func(t *testing.T) {
			actual, err := NewFromPostgresConnection(
				&mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_Url{
							Url: "postgres://test-user:test-pass@localhost:5432/postgres",
						},
					},
				},
				nil,
				testutil.GetTestLogger(t),
			)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(
				t,
				"postgres://test-user:test-pass@localhost:5432/postgres",
				actual.String(),
			)
			assert.Equal(t, "test-user", actual.GetUser())
		})
		t.Run("ok_user_provided_timeout", func(t *testing.T) {
			actual, err := NewFromPostgresConnection(
				&mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_Url{
							Url: "postgres://test-user:test-pass@localhost:5432/postgres?connect_timeout=10",
						},
					},
				},
				&testConnectionTimeout,
				testutil.GetTestLogger(t),
			)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(
				t,
				"postgres://test-user:test-pass@localhost:5432/postgres?connect_timeout=10",
				actual.String(),
			)
			assert.Equal(t, "test-user", actual.GetUser())
		})
	})

	t.Run("URL from Env", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			viper.Set(fmt.Sprintf("%s%s", userDefinedEnvPrefix, "PG_URL"), "postgres://test-user:testpass@localhost:3309/mydb")
			actual, err := NewFromPostgresConnection(
				&mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_UrlFromEnv{
							UrlFromEnv: "USER_DEFINED_PG_URL",
						},
					},
				},
				&testConnectionTimeout,
				testutil.GetTestLogger(t),
			)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(
				t,
				"postgres://test-user:testpass@localhost:3309/mydb?connect_timeout=5",
				actual.String(),
			)
			assert.Equal(t, "test-user", actual.GetUser())
		})
		t.Run("error_no_prefix", func(t *testing.T) {
			viper.Set("PG_URL_NO_PREFIX", "postgres://test-user:testpass@localhost:3309/mydb")

			_, err := NewFromPostgresConnection(
				&mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_UrlFromEnv{
							UrlFromEnv: "PG_URL_NO_PREFIX",
						},
					},
				},
				&testConnectionTimeout,
				testutil.GetTestLogger(t),
			)
			require.Error(t, err)
		})
	})
}
