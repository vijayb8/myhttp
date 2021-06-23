package main

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func Test_getMD5(t *testing.T) {
	ch := make(chan string, 1)
	t.Run("not valid url", func(t *testing.T) {
		getMD5("google.com", ch)
		resp := <- ch
		assert.Equal(t, resp, "")
	})
	t.Run("valid url", func(t *testing.T) {
		getMD5("http://google.com", ch)
		resp := <- ch
		//check if the size is 32
		assert.Equal(t, len(resp), 32)
	})
}

func Test_getHash(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"get valid hash",
			args{
				content: []byte("test"),
			},
			"098f6bcd4621d373cade4e832627b4f6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHash(tt.args.content); got != tt.want {
				t.Errorf("getHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getValidUrl(t *testing.T) {
	type args struct {
		rawUrl string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"get valid url with scheme",
			args{
				rawUrl: "google.com",
			},
			"http://google.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValidUrl(tt.args.rawUrl); got != tt.want {
				t.Errorf("getValidUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}