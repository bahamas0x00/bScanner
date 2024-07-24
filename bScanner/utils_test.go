package bScanner

import "testing"

func Test_urlCheck(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"t1", args{"baidu.com"}, "http://baidu.com/"},
		{"t2", args{"baidu.com/"}, "http://baidu.com/"},
		{"t3", args{"http://baidu.com"}, "http://baidu.com/"},
		{"t4", args{"https://baidu.com"}, "https://baidu.com/"},
		{"t5", args{"https://baidu.com/"}, "https://baidu.com/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := urlCheck(tt.args.url); got != tt.want {
				t.Errorf(" urlCheck() = %v, want %v", got, tt.want)

			}
		})
	}
}
