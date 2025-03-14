package Giphy

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGiphy_PullFromGiphy(t *testing.T) {
	type fields struct {
		apiKey string
	}
	type args struct {
		searchTerm string
	}
	tests := []struct {
		name    string
		server  *httptest.Server
		fields  fields
		args    args
		wantErr error
		want    *Response
	}{
		{
			name: "empty test input",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("{\"data\":[]}"))
			})),
			fields: fields{
				apiKey: "test",
			},
			args: args{
				searchTerm: "test",
			},
			want: &Response{
				Data: []struct {
					Id     string `json:"id"`
					Images struct {
						Preview struct {
							URL string `json:"url"`
						} `json:"downsized_small"`
						Original struct {
							URL string `json:"url"`
						} `json:"original"`
					} `json:"images"`
				}{},
			},
		},
		{
			name: "test against giphy with a search term",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("{\"data\":[{\"images\":{\"downsized_small\":{\"url\":\"http://test.com\"}, \"original\":{\"url\":\"http://test.com\"}}}]}"))
			})),
			fields: fields{
				apiKey: "test",
			},
			args: args{
				searchTerm: "test",
			},
			want: &Response{
				Data: []struct {
					Id     string `json:"id"`
					Images struct {
						Preview struct {
							URL string `json:"url"`
						} `json:"downsized_small"`
						Original struct {
							URL string `json:"url"`
						} `json:"original"`
					} `json:"images"`
				}{
					{
						Images: struct {
							Preview struct {
								URL string `json:"url"`
							} `json:"downsized_small"`
							Original struct {
								URL string `json:"url"`
							} `json:"original"`
						}{
							Preview: struct {
								URL string `json:"url"`
							}{
								URL: "http://test.com",
							},
							Original: struct {
								URL string `json:"url"`
							}{
								URL: "http://test.com",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Giphy{
				apiKey: tt.fields.apiKey,
				host:   tt.server.URL,
				limit:  5,
			}
			gotErr, got := g.PullFromGiphy(tt.args.searchTerm)
			if tt.wantErr != nil && gotErr != tt.wantErr {
				t.Errorf("PullFromGiphy() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
			if tt.wantErr == nil && gotErr != nil {
				t.Errorf("PullFromGiphy() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullFromGiphy() got = %v, want %v", got, tt.want)
			}
		})
	}
}
