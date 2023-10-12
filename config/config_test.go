package config

import (
	"reflect"
	"testing"
)

func Test_readConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{
			name: "success",
			want: &Config{
				PostgresDriver:   "postgres",
				PostgresUser:     "posts",
				PostgresPassword: "p0stgr3s",
				PostgresHost:     "localhost",
				PostgresPort:     "5432",
				PostgresPath:     "posts",
				PostgresSslmode:  "disable",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig error = %v, want %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig got %v, want %v", got, tt.want)
			}
		})
	}
}
