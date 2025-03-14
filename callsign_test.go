package main

import (
	"github.com/pyrousnet/slash_commands/Callook"
	"github.com/pyrousnet/slash_commands/HamDb"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_getCallSignInfo(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		r      *http.Request
		server *httptest.Server
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test getCallSignInfo pulling from HamDb",
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{
					// make url.Url that matches http://localhost:8080/?text=K1LNX
					URL: &url.URL{Scheme: "http", Host: "localhost:8080", RawQuery: "text=K1LNX"},
				},
				server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"hamdb":{"version":"1.0","callsign":{"call":"K1LNX","class":"E","expires":"2022-06-06","status":"Active","grid":"FN42","lat":"42.1234","lon":"-71.1234","fname":"John","mi":"Q","name":"Doe","suffix":"","addr1":"123 Main St","addr2":"","state":"MA","zip":"12345","country":"USA"},"messages":{"status":"OK"}}}`))
				})),
			},
			want: `{"response_type":"in_channel","text":"| Data | Value |\n| :------ | :-------|\n| Callsign | K1LNX |\n| Name | John Q Doe |\n| City | , MA 12345 |\n| Last3 | LNX |\n| Class | Extra |\n| Status | Unknown |"}`,
		},
		{
			name: "Test getCallSignInfo pulling from Callook",
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{
					// make url.Url that matches http://localhost:8080/?text=K1LNX
					URL: &url.URL{Scheme: "http", Host: "localhost:8080", RawQuery: "text=K1LNX"},
				},
				server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"status":"INVALID","type":"callsign","current":{"callsign":"K1LNX","operClass":"E"},"previous":{"callsign":"K1LNX","operClass":"E"},"trustee":{"callsign":"K1LNX","name":"John Doe"},"name":"John Doe","address":{"line1":"123 Main St","line2":"","attn":""},"location":{"latitude":"42.1234","longitude":"-71.1234","gridsquare":"FN42"},"otherInfo":{"grantDate":"2022-06-06","expiryDate":"2022-06-06","lastActionDate":"2022-06-06","frn":"1234567890","ulsUrl":"http://callook.info"}}`))
				})),
			},
			want: `{"response_type":"in_channel","text":"| Data | Value |\n| :------ | :-------|\n| Callsign | K1LNX |\n| Name | John Doe |\n| City |  |\n| Last3 | LNX |\n| Class | E |\n| Status | INVALID |"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HamDb.HamDbUrl = tt.args.server.URL + "/"
			Callook.CallookUrl = tt.args.server.URL + "/"
			getCallSignInfo(tt.args.w, tt.args.r)

			// get the response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			if resp.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.Code)
			}

			// check data against expected
			got := resp.Body.String()
			if got != tt.want {
				t.Errorf("getCallSignInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
