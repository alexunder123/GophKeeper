package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	// user test
	os.Remove("user_config.json")
	tests := []struct {
		name string
		want *UserConfig
	}{
		{
			name: "userTest1 create file",
			want: &UserConfig{
				RunAddress: "127.0.0.1:3200",
			},
		},
		{
			name: "userTest2 read file",
			want: &UserConfig{
				RunAddress: "127.0.0.1:3200",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserConfig()
			require.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
	
}
