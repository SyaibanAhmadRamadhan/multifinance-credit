package pagination_test

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetOffsetValue(t *testing.T) {
	type args struct {
		page     int64
		pageSize int64
	}
	tests := []struct {
		name     string
		args     args
		expected int64
	}{
		{
			name: "default zero value",
			args: args{
				page:     0,
				pageSize: 5,
			},
			expected: 0,
		},
		{
			name: "basic test",
			args: args{
				page:     3,
				pageSize: 5,
			},
			expected: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			// WHEN
			actual := pagination.GetOffsetValue(tt.args.page, tt.args.pageSize)
			// THEN
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestGetPageCount(t *testing.T) {
	type args struct {
		pageSize  int64
		totalData int64
	}
	tests := []struct {
		name     string
		args     args
		expected int64
	}{
		{
			name: "default zero value",
			args: args{
				pageSize:  0,
				totalData: 10,
			},
			expected: 1,
		},
		{
			name: "basic test",
			args: args{
				pageSize:  5,
				totalData: 10,
			},
			expected: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			// WHEN
			actual := pagination.GetPageCount(tt.args.pageSize, tt.args.totalData)
			// THEN
			require.Equal(t, tt.expected, actual)
		})
	}
}
