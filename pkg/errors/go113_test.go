// Copyright 2021 dairongpeng <dairongpeng@foxmail.com>. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build go1.13

package errors

import (
	stderrors "errors"
	"fmt"
	"reflect"
	"testing"
)

func TestErrorChainCompat(t *testing.T) {
	err := stderrors.New("error that gets wrapped")
	wrapped := Wrap(err, "wrapped up")
	if !stderrors.Is(wrapped, err) {
		t.Errorf("Wrap does not support Go 1.13 error chains")
	}
}

func TestIs(t *testing.T) {
	err := New("test")

	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with stack",
			args: args{
				err:    WithStack(err),
				target: err,
			},
			want: true,
		},
		{
			name: "with message",
			args: args{
				err:    WithMessage(err, "test"),
				target: err,
			},
			want: true,
		},
		{
			name: "with message format",
			args: args{
				err:    WithMessagef(err, "%s", "test"),
				target: err,
			},
			want: true,
		},
		{
			name: "std errors compatibility",
			args: args{
				err:    fmt.Errorf("wrap it: %w", err),
				target: err,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

type customErr struct {
	msg string
}

func (c customErr) Error() string { return c.msg }

func TestAs(t *testing.T) {
	var err = customErr{msg: "test message"}

	type args struct {
		err    error
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with stack",
			args: args{
				err:    WithStack(err),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "with message",
			args: args{
				err:    WithMessage(err, "test"),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "with message format",
			args: args{
				err:    WithMessagef(err, "%s", "test"),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "std errors compatibility",
			args: args{
				err:    fmt.Errorf("wrap it: %w", err),
				target: new(customErr),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := As(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("As() = %v, want %v", got, tt.want)
			}

			ce := tt.args.target.(*customErr)
			if !reflect.DeepEqual(err, *ce) {
				t.Errorf("set target error failed, target error is %v", *ce)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	err := New("test")

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "with stack",
			args: args{err: WithStack(err)},
			want: err,
		},
		{
			name: "with message",
			args: args{err: WithMessage(err, "test")},
			want: err,
		},
		{
			name: "with message format",
			args: args{err: WithMessagef(err, "%s", "test")},
			want: err,
		},
		{
			name: "std errors compatibility",
			args: args{err: fmt.Errorf("wrap: %w", err)},
			want: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unwrap(tt.args.err); !reflect.DeepEqual(err, tt.want) {
				t.Errorf("Unwrap() error = %v, want %v", err, tt.want)
			}
		})
	}
}
