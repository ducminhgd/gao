package number

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxNumberInt32(t *testing.T) {
	tests := []struct {
		name    string
		args    []int32
		want    int32
		wantErr error
	}{
		{
			name:    "empty arguments",
			args:    []int32{},
			want:    0,
			wantErr: errors.New("arguments are required"),
		},
		{
			name:    "only one element",
			args:    []int32{1},
			want:    1,
			wantErr: nil,
		},
		{
			name:    "positive and incremental",
			args:    []int32{1, 2, 3, 4, 100},
			want:    100,
			wantErr: nil,
		},
		{
			name:    "positive and decremental",
			args:    []int32{100, 99, 98, 5, 0},
			want:    100,
			wantErr: nil,
		},
		{
			name:    "negative and random",
			args:    []int32{-100, 5, 99, 98, 0},
			want:    99,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MaxNumber(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMaxNumberInt64(t *testing.T) {
	tests := []struct {
		name    string
		args    []int64
		want    int64
		wantErr error
	}{
		{
			name:    "empty arguments",
			args:    []int64{},
			want:    0,
			wantErr: errors.New("arguments are required"),
		},
		{
			name:    "only one element",
			args:    []int64{1},
			want:    1,
			wantErr: nil,
		},
		{
			name:    "positive and incremental",
			args:    []int64{1, 2, 3, 4, 100},
			want:    100,
			wantErr: nil,
		},
		{
			name:    "positive and decremental",
			args:    []int64{100, 99, 98, 5, 0},
			want:    100,
			wantErr: nil,
		},
		{
			name:    "negative and random",
			args:    []int64{-100, 5, 99, 98, 0},
			want:    99,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MaxNumber(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMaxNumberFloat32(t *testing.T) {
	tests := []struct {
		name    string
		args    []float32
		want    float32
		wantErr error
	}{
		{
			name:    "empty arguments",
			args:    []float32{},
			want:    0,
			wantErr: errors.New("arguments are required"),
		},
		{
			name:    "only one element",
			args:    []float32{1.1},
			want:    1.1,
			wantErr: nil,
		},
		{
			name:    "positive and incremental",
			args:    []float32{1.1, 2.1, 3.3, 4.4, 100.1},
			want:    100.1,
			wantErr: nil,
		},
		{
			name:    "positive and decremental",
			args:    []float32{100.1, 99.99, 98.9, 5.5, 0.1},
			want:    100.1,
			wantErr: nil,
		},
		{
			name:    "negative and random",
			args:    []float32{-100.1, 5.5, 99, 98.2, 0.4},
			want:    99,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MaxNumber(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMaxNumberFloat64(t *testing.T) {
	tests := []struct {
		name    string
		args    []float64
		want    float64
		wantErr error
	}{
		{
			name:    "empty arguments",
			args:    []float64{},
			want:    0,
			wantErr: errors.New("arguments are required"),
		},
		{
			name:    "only one element",
			args:    []float64{1.1},
			want:    1.1,
			wantErr: nil,
		},
		{
			name:    "positive and incremental",
			args:    []float64{1.1, 2.1, 3.3, 4.4, 100.1},
			want:    100.1,
			wantErr: nil,
		},
		{
			name:    "positive and decremental",
			args:    []float64{100.1, 99.99, 98.9, 5.5, 0.1},
			want:    100.1,
			wantErr: nil,
		},
		{
			name:    "negative and random",
			args:    []float64{-100.1, 5.5, 99, 98.2, 0.4},
			want:    99,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MaxNumber(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMinNumberInt32(t *testing.T) {
	tests := []struct {
		name    string
		args    []int32
		want    int32
		wantErr error
	}{
		{
			name:    "empty arguments",
			args:    []int32{},
			want:    0,
			wantErr: errors.New("arguments are required"),
		},
		{
			name:    "only one element",
			args:    []int32{1},
			want:    1,
			wantErr: nil,
		},
		{
			name:    "positive and incremental",
			args:    []int32{1, 2, 3, 4, 100},
			want:    1,
			wantErr: nil,
		},
		{
			name:    "positive and decremental",
			args:    []int32{100, 99, 98, 5, 0},
			want:    0,
			wantErr: nil,
		},
		{
			name:    "negative and random",
			args:    []int32{-100, 5, 99, 98, 0},
			want:    -100,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MinNumber(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMinNumberInt64(t *testing.T) {
	tests := []struct {
		name    string
		args    []int64
		want    int64
		wantErr error
	}{
		{
			name:    "empty arguments",
			args:    []int64{},
			want:    0,
			wantErr: errors.New("arguments are required"),
		},
		{
			name:    "only one element",
			args:    []int64{1},
			want:    1,
			wantErr: nil,
		},
		{
			name:    "positive and incremental",
			args:    []int64{1, 2, 3, 4, 100},
			want:    1,
			wantErr: nil,
		},
		{
			name:    "positive and decremental",
			args:    []int64{100, 99, 98, 5, 0},
			want:    0,
			wantErr: nil,
		},
		{
			name:    "negative and random",
			args:    []int64{-100, 5, 99, 98, 0},
			want:    -100,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MinNumber(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMinNumberFloat32(t *testing.T) {
	tests := []struct {
		name    string
		args    []float32
		want    float32
		wantErr error
	}{
		{
			name:    "empty arguments",
			args:    []float32{},
			want:    0,
			wantErr: errors.New("arguments are required"),
		},
		{
			name:    "only one element",
			args:    []float32{1.1},
			want:    1.1,
			wantErr: nil,
		},
		{
			name:    "positive and incremental",
			args:    []float32{1.1, 2.2, 3.3, 4.4, 100.1},
			want:    1.1,
			wantErr: nil,
		},
		{
			name:    "positive and decremental",
			args:    []float32{100.1, 99.99, 98.1, 5.45, 0},
			want:    0,
			wantErr: nil,
		},
		{
			name:    "negative and random",
			args:    []float32{-100, 5.5, 99.1, 98.2, 0},
			want:    -100,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MinNumber(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMinNumberFloat64(t *testing.T) {
	tests := []struct {
		name    string
		args    []float64
		want    float64
		wantErr error
	}{
		{
			name:    "empty arguments",
			args:    []float64{},
			want:    0,
			wantErr: errors.New("arguments are required"),
		},
		{
			name:    "only one element",
			args:    []float64{1.1},
			want:    1.1,
			wantErr: nil,
		},
		{
			name:    "positive and incremental",
			args:    []float64{1.1, 2.2, 3.3, 4.4, 100.1},
			want:    1.1,
			wantErr: nil,
		},
		{
			name:    "positive and decremental",
			args:    []float64{100.1, 99.99, 98.1, 5.45, 0},
			want:    0,
			wantErr: nil,
		},
		{
			name:    "negative and random",
			args:    []float64{-100, 5.5, 99.1, 98.2, 0},
			want:    -100,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MinNumber(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
