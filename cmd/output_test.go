package cmd

import (
	"os"
	"testing"
)

func TestFileOutput_Write(t *testing.T) {
	type fields struct {
		filePath string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Create file with text",
			fields:  fields{"./"},
			args:    args{[]byte("Hello RIDOS")},
			wantErr: false,
		},
		{
			name:    "Don`t find path of peremision deniet",
			fields:  fields{"//"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileOutput := &FileOutput{
				filePath: tt.fields.filePath,
			}
			if err := fileOutput.Write(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("FileOutput.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			os.Remove(tt.fields.filePath)
		})
	}
}
