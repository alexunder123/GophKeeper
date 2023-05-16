package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	os.Remove("config.json")
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "test1 create file",
			want: &Config{
				RunAddress:        "127.0.0.1:3200",
				DatabaseDirectory: "/users/",
				SQLDatabase:       "postgres://postgres:1@localhost:5432/postgres?sslmode=disable",
				Expires:           2,
				LenghtSesionID:    16,
				LenghtUserID:      12,
				LockingTime:       15,
			},
		},
		{
			name: "test2 read file",
			want: &Config{
				RunAddress:        "127.0.0.1:3200",
				DatabaseDirectory: "/users/",
				SQLDatabase:       "postgres://postgres:1@localhost:5432/postgres?sslmode=disable",
				Expires:           2,
				LenghtSesionID:    16,
				LenghtUserID:      12,
				LockingTime:       15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig()
			require.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}

}
