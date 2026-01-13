package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func sqlString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func sqlInt32(i int32) sql.NullInt32 {
	return sql.NullInt32{Int32: i, Valid: true}
}

func sqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func sqlNullInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func sqlNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u == nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: *u, Valid: true}
}

func sqlNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func nullInt32ToPtr(ni sql.NullInt32) *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}

func nullUUIDToPtr(nu uuid.NullUUID) *uuid.UUID {
	if !nu.Valid {
		return nil
	}
	return &nu.UUID
}

func nullTimeToTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func sqlBoolPtr(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

func sqlNullStringPtr(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func sqlNullInt32Ptr(i int32) sql.NullInt32 {
	return sql.NullInt32{Int32: i, Valid: true}
}

func sqlNullUUIDPtr(u uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: u, Valid: true}
}

func sqlNullBoolPtr(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

func sqlNullTimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func int32PtrToInt(i int32) *int {
	val := int(i)
	return &val
}

func sqlNullBoolFromPtr(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}
