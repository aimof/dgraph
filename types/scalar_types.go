/*
 * Copyright 2016-2018 Dgraph Labs, Inc. and Contributors
 *
 * This file is available under the Apache License, Version 2.0,
 * with the Commons Clause restriction.
 */

package types

import (
	"time"

	"github.com/dgraph-io/dgraph/protos/intern"
	geom "github.com/twpayne/go-geom"
)

const (
	nanoSecondsInSec = 1000000000
)

// Note: These ids are stored in the posting lists to indicate the type
// of the data. The order *cannot* be changed without breaking existing
// data. When adding a new type *always* add to the end of this list.
// Never delete anything from this list even if it becomes unused.
const (
	BinaryID   = TypeID(intern.Posting_BINARY)
	IntID      = TypeID(intern.Posting_INT)
	FloatID    = TypeID(intern.Posting_FLOAT)
	BoolID     = TypeID(intern.Posting_BOOL)
	DateTimeID = TypeID(intern.Posting_DATETIME)
	StringID   = TypeID(intern.Posting_STRING)
	GeoID      = TypeID(intern.Posting_GEO)
	UidID      = TypeID(intern.Posting_UID)
	PasswordID = TypeID(intern.Posting_PASSWORD)
	DefaultID  = TypeID(intern.Posting_DEFAULT)
)

var typeNameMap = map[string]TypeID{
	"int":      IntID,
	"float":    FloatID,
	"string":   StringID,
	"bool":     BoolID,
	"datetime": DateTimeID,
	"geo":      GeoID,
	"uid":      UidID,
	"password": PasswordID,
	"default":  DefaultID,
}

type TypeID intern.Posting_ValType

func (t TypeID) Enum() intern.Posting_ValType {
	return intern.Posting_ValType(t)
}

func (t TypeID) Name() string {
	switch t {
	case IntID:
		return "int"
	case FloatID:
		return "float"
	case BoolID:
		return "bool"
	case StringID:
		return "string"
	case DateTimeID:
		return "datetime"
	case GeoID:
		return "geo"
	case UidID:
		return "uid"
	case PasswordID:
		return "password"
	case DefaultID:
		return "default"
	case BinaryID:
		return "binary"
	}
	return ""
}

// Val is a value with type information.
type Val struct {
	Tid   TypeID
	Value interface{}
}

// TypeForName returns the type corresponding to the given name.
// If name is not recognized, it returns nil.
func TypeForName(name string) (TypeID, bool) {
	t, ok := typeNameMap[name]
	return t, ok
}

func (t TypeID) IsScalar() bool {
	return t != UidID
}

// ValueForType returns the zero value for a type id
func ValueForType(id TypeID) Val {
	switch id {
	case BinaryID:
		var b []byte
		return Val{BinaryID, &b}

	case IntID:
		var i int64
		return Val{IntID, &i}

	case FloatID:
		var f float64
		return Val{FloatID, &f}

	case BoolID:
		var b bool
		return Val{BoolID, &b}

	case DateTimeID:
		var t time.Time
		return Val{DateTimeID, &t}

	case StringID:
		var s string
		return Val{StringID, s}

	case DefaultID:
		var s string
		return Val{DefaultID, s}

	case GeoID:
		var g geom.T
		return Val{GeoID, &g}

	case UidID:
		var i uint64
		return Val{UidID, &i}

	case PasswordID:
		var p string
		return Val{PasswordID, p}

	default:
		return Val{}
	}
}

func createDate(y int, m time.Month, d int) time.Time {
	var dt time.Time
	dt = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return dt
}

func ParseTime(val string) (time.Time, error) {
	var t time.Time
	if err := t.UnmarshalText([]byte(val)); err == nil {
		return t, err
	}
	// try without timezone
	if t, err := time.Parse(dateTimeFormat, val); err == nil {
		return t, err
	}
	if t, err := time.Parse(dateFormatYMD, val); err == nil {
		return t, err
	}
	if t, err := time.Parse(dateFormatYM, val); err == nil {
		return t, err
	}
	return time.Parse(dateFormatY, val)
}

const dateFormatYMD = "2006-01-02"
const dateFormatYM = "2006-01"
const dateFormatY = "2006"
const dateTimeFormat = "2006-01-02T15:04:05"
