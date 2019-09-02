package main

import (
	"reflect"
	"testing"
)

func Test_doCheckService(t *testing.T) {
	type args struct {
		service *Service
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "GET OK",
			args: args{
				service: &Service{
					Name:     "google.com.br",
					Protocol: "http",
					URI:      "www.google.com.br",
					Port:     80,
					Method:   "GET",
				},
			},
			want: 200,
		},
		{
			name: "GET NOT OK",
			args: args{
				service: &Service{
					Name:     "DUMMY",
					Protocol: "http",
					URI:      "192.192.192.192",
					Port:     801,
					Method:   "GET",
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doCheckService(tt.args.service)
			if !reflect.DeepEqual(tt.args.service.StatusCode, tt.want) {
				t.Errorf("=> Got %v wanted %v", tt.args.service.StatusCode, tt.want)
			}
		})
	}
}
