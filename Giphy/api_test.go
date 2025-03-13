package Giphy

import (
	"encoding/json"
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
		name     string
		fields   fields
		args     args
		want     error
		wantJson string
	}{
		{
			name: "empty test input",
			fields: fields{
				apiKey: "test",
			},
			args: args{
				searchTerm: "test",
			},
			wantJson: "{\"data\":[]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Giphy{
				apiKey: tt.fields.apiKey,
			}
			got, got1 := g.PullFromGiphy(tt.args.searchTerm)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullFromGiphy() got = %v, want %v", got, tt.want)
			}
			gotJson, err := json.Marshal(got1)
			if err != nil {
				t.Errorf("Error marshalling JSON: %v", err)
			}
			if string(gotJson) != tt.wantJson {
				t.Errorf("PullFromGiphy() got1 = %v, want %v", string(gotJson), tt.wantJson)
			}
		})
	}
}
