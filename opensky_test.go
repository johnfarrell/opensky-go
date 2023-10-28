package opensky

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	type args struct {
		opts []ClientOption
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "Default Client",
			args: args{
				opts: []ClientOption{},
			},
			want: &Client{
				httpClient: http.Client{
					Timeout: 1 * time.Second,
				},
				authHeader: "",
				baseUrl:    "https://opensky-network.org/api",
			},
			wantErr: false,
		},
		{
			name: "Custom timeout",
			args: args{
				opts: []ClientOption{
					WithTimeout(5 * time.Second),
				},
			},
			want: &Client{
				httpClient: http.Client{
					Timeout: 5 * time.Second,
				},
				authHeader: "",
				baseUrl:    "https://opensky-network.org/api",
			},
			wantErr: false,
		},
		{
			name: "Authentication",
			args: args{
				opts: []ClientOption{
					WithAuthentication("username", "password"),
				},
			},
			want: &Client{
				httpClient: http.Client{
					Timeout: 1 * time.Second,
				},
				authHeader: "username:password",
				baseUrl:    "https://opensky-network.org/api",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetStates(t *testing.T) {

	cli, _ := NewClient(WithTimeout(120 * time.Second))
	tests := []struct {
		name    string
		want    StateResponse
		wantErr bool
	}{
		{
			name:    "Basic response",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cli.GetStates(&BoundingBox{
				LongitudeMin: -122.4845,
				LongitudeMax: -122.1693,
				LatitudeMin:  47.5014,
				LatitudeMax:  47.6954,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
