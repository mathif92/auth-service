package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func Update(ctx context.Context, tx *sqlx.Tx, table string, attributes map[string]interface{}, filters map[string]interface{}) error {
	// Create SET values UPDATE statement section
	var attributesSections []string
	var attributeValues []interface{}
	for column, value := range attributes {
		attributesSections = append(attributesSections, fmt.Sprintf("%s = %v", column, value))
		attributeValues = append(attributeValues, value)
	}
	attributesStr := strings.Join(attributesSections, ", ")

	// Create WHERE values UPDATE statement section
	var filtersSections []string
	var filterValues []interface{}
	for column, value := range filters {
		filtersSections = append(filtersSections, fmt.Sprintf("%s = %v", column, value))
		filterValues = append(filterValues, value)
	}
	filtersStr := strings.Join(filtersSections, "AND ")

	// Put it all together
	updateQuery := fmt.Sprintf("UPDATE %s SET %s", table, attributesStr)
	if len(filters) > 0 {
		updateQuery = fmt.Sprintf("%s WHERE %s", updateQuery, filtersStr)
	}

	queryArgs := append(attributeValues, filterValues...)

	_, err := tx.ExecContext(ctx, updateQuery, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}
