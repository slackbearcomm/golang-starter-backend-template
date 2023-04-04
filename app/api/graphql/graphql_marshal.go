package graphql

import (
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
)

// MarshalNullString to handle custom type
func MarshalNullString(ns null.String) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalString(ns.String)
}

// UnmarshalNullString to handle custom type
func UnmarshalNullString(v interface{}) (null.String, error) {
	if v == nil {
		return null.String{Valid: false}, nil
	}
	s, err := graphql.UnmarshalString(v)

	if err != nil {
		return null.StringFrom(s), fmt.Errorf("graphql marshal error")
	}

	return null.StringFrom(s), nil
}

// MarshalNullInt to handle custom type
func MarshalNullInt(ns null.Int) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalInt(ns.Int)
}

// UnmarshalNullInt to handle custom type
func UnmarshalNullInt(v interface{}) (null.Int, error) {
	if v == nil {
		return null.Int{Valid: false}, nil
	}
	s, err := graphql.UnmarshalInt(v)
	if err != nil {
		return null.IntFrom(s), fmt.Errorf("graphql marshal error")
	}

	return null.IntFrom(s), nil
}

// MarshalNullInt64 to handle custom type
func MarshalNullInt64(ns null.Int64) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalInt64(ns.Int64)
}

// UnmarshalNullInt64 to handle custom type
func UnmarshalNullInt64(v interface{}) (null.Int64, error) {
	if v == nil {
		return null.Int64{Valid: false}, nil
	}
	s, err := graphql.UnmarshalInt64(v)
	if err != nil {
		return null.Int64From(s), fmt.Errorf("graphql marshal error")
	}

	return null.Int64From(s), nil
}

// MarshalNullFloat to handle custom type
func MarshalNullFloat(ns null.Float64) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalFloat(ns.Float64)
}

// UnmarshalNullFloat to handle custom type
func UnmarshalNullFloat(v interface{}) (null.Float64, error) {
	if v == nil {
		return null.Float64{Valid: false}, nil
	}
	s, err := graphql.UnmarshalFloat(v)
	if err != nil {
		return null.Float64From(s), fmt.Errorf("graphql marshal error")
	}

	return null.Float64From(s), nil
}

// MarshalNullTime to handle custom type
func MarshalNullTime(ns null.Time) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalTime(ns.Time)
}

// UnmarshalNullTime to handle custom type
func UnmarshalNullTime(v interface{}) (null.Time, error) {
	if v == nil {
		return null.Time{Valid: false}, nil
	}
	s, err := graphql.UnmarshalTime(v)
	if err != nil {
		return null.TimeFrom(s), fmt.Errorf("graphql marshal error")
	}
	return null.TimeFrom(s), nil
}

// MarshalNullBool to handle custom type
func MarshalNullBool(ns null.Bool) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalBoolean(ns.Bool)
}

// UnmarshalNullBool to handle custom type
func UnmarshalNullBool(v interface{}) (null.Bool, error) {
	if v == nil {
		return null.Bool{Valid: false}, nil
	}
	s, err := graphql.UnmarshalBoolean(v)
	if err != nil {
		return null.BoolFrom(s), fmt.Errorf("graphql marshal error")
	}
	return null.BoolFrom(s), nil
}

// MarshalUUID allows uuid to be marshalled by graphql
func MarshalUUID(id uuid.UUID) graphql.Marshaler {
	return graphql.MarshalString(id.String())
}

// UnmarshalUUID allows uuid to be unmarshalled by graphql
func UnmarshalUUID(v interface{}) (uuid.UUID, error) {
	idAsString, ok := v.(string)
	if !ok {
		return uuid.Nil, errors.New("id should be a valid UUID")
	}
	return uuid.FromString(idAsString)
}

// MarshalNullUUID allows uuid to be marshalled by graphql
func MarshalNullUUID(id uuid.NullUUID) graphql.Marshaler {
	if !id.Valid {
		return graphql.Null
	}
	return graphql.MarshalString(id.UUID.String())
}

// UnmarshalNullUUID allows uuid to be unmarshalled by graphql
func UnmarshalNullUUID(v interface{}) (uuid.NullUUID, error) {
	if v == nil {
		return uuid.NullUUID{Valid: false}, nil
	}
	s, err := graphql.UnmarshalString(v)
	if err != nil {
		return uuid.NullUUID{Valid: false, UUID: uuid.FromStringOrNil(s)}, fmt.Errorf("graphql marshal error")
	}

	uid, err := uuid.FromString(s)
	if err != nil {
		return uuid.NullUUID{Valid: false, UUID: uuid.FromStringOrNil(s)}, fmt.Errorf("graphql marshal error")
	}
	return uuid.NullUUID{Valid: true, UUID: uid}, nil
}
