package cmd

import (
	"reflect"
	"testing"
)

func TestParseAndReformatJson(t *testing.T) {
	type args struct {
		jsonData string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Good string",
			args: args{
				jsonData: `{"name": "richard", "age": 25}`,
			},
			want:    []byte(`{"age":25,"name":"richard"}`),
			wantErr: false,
		},
		{
			name: "Bad walue",
			args: args{
				jsonData: `name: "richard" "age": 25`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAndReformatJson(tt.args.jsonData)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAndReformatJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAndReformatJson() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
