package copyutil

import (
	"reflect"
	"testing"

	"github.com/samber/lo"
)

func TestCopy(t *testing.T) {
	type from struct {
		Name string
		Age  *int
	}
	type to struct {
		Name string
		Age  *int
	}
	type args struct {
		from *from
	}
	tests := []struct {
		name    string
		args    args
		want    *to
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				from: &from{
					Name: "foo",
					Age:  lo.ToPtr(0),
				},
			},
			want: &to{
				Name: "foo",
				Age:  lo.ToPtr(0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Copy[from, to](tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}
