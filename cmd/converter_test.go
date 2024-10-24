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

func Test_strToParams(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name       string
		args       args
		wantParams []string
		wantErr    bool
	}{
		{
			name:       "Default usage",
			args:       args{str: "name=value"},
			wantParams: []string{"name", "value"},
			wantErr:    false,
		},
		{
			name:       "Multy params usage",
			args:       args{str: "name1=value1 name2=value2"},
			wantParams: []string{"name1", "value1", "name2", "value2"},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParams, err := strToParams(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("strToParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotParams, tt.wantParams) {
				t.Errorf("strToParams() = %v, want %v", gotParams, tt.wantParams)
			}
		})
	}
}
