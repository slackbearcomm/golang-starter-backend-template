package dbhelpers

import (
	"context"
	"fmt"
	"gogql/app/models"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ResolveFilterSort generates a query string for db query with filters
func ResolveFilterSort(tableName dbTable, filter models.SearchFilter) string {
	// Resolve filter.SortBy
	// - Date Created = "DateCreated"
	// - Date Updated = "DateUpdated"
	// - Alphabetical

	var orderBy string
	switch filter.SortBy {
	case "Alphabetical":
		orderBy = fmt.Sprintf("%s.id", tableName)
	case "DateCreated":
		orderBy = fmt.Sprintf("%s.created_at", tableName)
	case "DateUpdated":
		fallthrough
	default:
		orderBy = fmt.Sprintf("%s.updated_at", tableName)
	}

	// Resolve filter.SortDir
	// - Descending
	// - Ascending

	var orderDir string
	switch filter.SortDir {
	case "Ascending":
		orderDir = "ASC"
	case "Descending":
		orderDir = "DESC"
	default:
		orderDir = "ASC"
	}

	output := fmt.Sprintf("ORDER BY %s %s OFFSET %d LIMIT %d", orderBy, orderDir, filter.Offset, filter.Limit)

	return output
}

// GetGlobalTotal returns the count
func GetGlobalTotal(
	ctx context.Context,
	conn *pgxpool.Pool,
	tableName dbTable,
	conditionsQuery string,
	queryArgs []interface{},
) (int, error) {
	var total int

	selectQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	queryStmt := fmt.Sprintf("%s %s", selectQuery, conditionsQuery)

	row := conn.QueryRow(ctx, queryStmt, queryArgs...)
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func UniqueKeys(list []string) []string {
	keys := []string{}

	for _, e := range list {
		if !containsKey(keys[:], e) {
			keys = append(keys, e)
		}
	}
	return keys
}

func containsKey(list []string, key string) bool {
	for _, e := range list {
		if e == key {
			return true
		}
	}
	return false
}

func UniqueUIDs(list []*uuid.UUID) []*uuid.UUID {
	uids := []*uuid.UUID{}

	for _, e := range list {
		if !containsUID(uids[:], e) {
			uids = append(uids, e)
		}
	}
	return uids
}

func containsUID(list []*uuid.UUID, uid *uuid.UUID) bool {
	for _, e := range list {
		if e == uid {
			return true
		}
	}
	return false
}
