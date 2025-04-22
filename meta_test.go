package logr

import (
	"reflect"
	"testing"
)

func TestM(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		args args
		want Meta
	}{
		{
			name: "New Meta",
			args: args{
				key:   "some-key",
				value: "some-value",
			},
			want: Meta{"some-key": "some-value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := M(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("M() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeta_Copy(t *testing.T) {
	tests := []struct {
		name string
		m    Meta
		want Meta
	}{
		{
			name: "New Meta",
			m:    Meta{"some-key": "some-value"},
			want: Meta{"some-key": "some-value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Copy()
			tt.m.With("another-key", "another-value")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeta_With(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		m    Meta
		args args
		want Meta
	}{
		{
			name: "New Meta",
			m:    Meta{"some-key": "some-value"},
			want: Meta{"some-key": "some-value", "another-key": "another-value"},
			args: args{
				key:   "another-key",
				value: "another-value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.With(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("With() = %v, want %v", got, tt.want)
			}
		})
	}
}
