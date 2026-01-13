package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSqlNullString(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected sql.NullString
	}{
		{
			name:     "non-nil string",
			input:    stringPtr("test"),
			expected: sql.NullString{String: "test", Valid: true},
		},
		{
			name:     "nil string",
			input:    nil,
			expected: sql.NullString{Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sqlNullString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNullStringToPtr(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected *string
	}{
		{
			name:     "valid string",
			input:    sql.NullString{String: "test", Valid: true},
			expected: stringPtr("test"),
		},
		{
			name:     "invalid string",
			input:    sql.NullString{Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nullStringToPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSqlNullInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    *int32
		expected sql.NullInt32
	}{
		{
			name:     "non-nil int32",
			input:    int32Ptr(42),
			expected: sql.NullInt32{Int32: 42, Valid: true},
		},
		{
			name:     "nil int32",
			input:    nil,
			expected: sql.NullInt32{Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sqlNullInt32(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNullInt32ToPtr(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt32
		expected *int32
	}{
		{
			name:     "valid int32",
			input:    sql.NullInt32{Int32: 42, Valid: true},
			expected: int32Ptr(42),
		},
		{
			name:     "invalid int32",
			input:    sql.NullInt32{Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nullInt32ToPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSqlNullUUID(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name     string
		input    *uuid.UUID
		expected uuid.NullUUID
	}{
		{
			name:     "non-nil uuid",
			input:    &testUUID,
			expected: uuid.NullUUID{UUID: testUUID, Valid: true},
		},
		{
			name:     "nil uuid",
			input:    nil,
			expected: uuid.NullUUID{Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sqlNullUUID(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNullUUIDToPtr(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name     string
		input    uuid.NullUUID
		expected *uuid.UUID
	}{
		{
			name:     "valid uuid",
			input:    uuid.NullUUID{UUID: testUUID, Valid: true},
			expected: &testUUID,
		},
		{
			name:     "invalid uuid",
			input:    uuid.NullUUID{Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nullUUIDToPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNullTimeToTimePtr(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    sql.NullTime
		expected *time.Time
	}{
		{
			name:     "valid time",
			input:    sql.NullTime{Time: now, Valid: true},
			expected: &now,
		},
		{
			name:     "invalid time",
			input:    sql.NullTime{Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nullTimeToTimePtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions for testing
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}
