package helpers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/wildanfaz/store-management/internal/models"
)

func MySQLUpdateQueryValues(payload any, table, tag string) (string, []any, error) {
	vo := reflect.ValueOf(payload)
	to := reflect.TypeOf(payload)

	if vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
		to = to.Elem()
	}

	if vo.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("Expected %v, got %v", reflect.Struct, vo.Kind())
	}

	query := fmt.Sprintf("UPDATE %s \nSET ", table)

	values := []any{}
	whereKeys := []string{}
	whereValues := []any{}

	for i := 0; i < vo.NumField(); i++ {
		v := vo.Field(i)

		if v.IsZero() {
			continue
		}

		column := to.Field(i).Tag.Get("column")

		if strings.Contains(column, "where") {
			whereValues = append(whereValues, v.Interface())
			whereKeys = append(whereKeys, strings.Split(column, ":")[1])
			continue
		}
		query = fmt.Sprintf(`%s%s = ?, `, query, column)
		values = append(values, v.Interface())
	}
	query = strings.TrimSuffix(query, ", ")

	for i := 0; i < len(whereKeys); i++ {
		if i == 0 {
			query += "\nWHERE "
		}
		query = fmt.Sprintf("%s%s = ? AND ", query, whereKeys[i])
		values = append(values, whereValues[i])
	}
	query = strings.TrimSuffix(query, "AND ")

	return query, values, nil
}

func PostgreSQLQueryListProducts(payload models.Product) (string, []any) {
	var (
		counter = 1
		values  []any
	)

	var query = `
	SELECT id, name, description, price, quantity, created_at, updated_at
	FROM products
	`

	if payload.Name != "" {
		payload.Name = "%" + payload.Name + "%"
		values = append(values, payload.Name)
		query = fmt.Sprintf("%s \nWHERE name ILIKE $%d", query, counter)
		counter++
	}

	query += "\nORDER BY created_at DESC"

	pagination := reflect.ValueOf(payload.Pagination)

	if !pagination.IsZero() {
		if pagination.FieldByName("PerPage").Int() < 1 {
			payload.Pagination.PerPage = 10
		}

		values = append(values, payload.PerPage)
		query = fmt.Sprintf("%s \nLIMIT $%d", query, counter)
		counter++

		if pagination.FieldByName("Page").Int() > 0 {
			values = append(values, (payload.Page-1)*payload.PerPage)
			query = fmt.Sprintf("%s \nOFFSET $%d", query, counter)
		}
	}

	return query, values
}

func PostgreSQLQueryUpdateProduct(id int, payload models.Product) (string, []any) {
	var (
		counter = 1
		values  []any
	)

	var query = `
	UPDATE products
	SET
	`

	v := reflect.ValueOf(payload)
	t := reflect.TypeOf(payload)

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			query = fmt.Sprintf("%s %s = $%d, ", query, t.Field(i).Tag.Get("column"), counter)
			values = append(values, v.Field(i).Interface())
			counter++
		}
	}

	query = strings.TrimSuffix(query, ", ")
	query = fmt.Sprintf("%s \nWHERE id = $%d", query, counter)
	values = append(values, id)

	return query, values
}

func PostgreSQLQueryListOrders(payload models.Order) (string, []any) {
	var (
		counter = 1
		values  []any
	)

	var query = `
	SELECT 
	o.id, o.status, o.quantity, o.created_at, o.updated_at,
	p.id, p.name, p.description, p.price, p.quantity, p.created_at, p.updated_at
	FROM orders o
	JOIN products p ON p.id = o.product_id
	JOIN users u ON u.id = o.user_id
	WHERE u.id = $1
	`
	values = append(values, payload.UserId)
	counter++

	if payload.Status != "" {
		values = append(values, payload.Status)
		query = fmt.Sprintf("%s AND o.status = $%d", query, counter)
		counter++
	}

	query += "\nORDER BY o.created_at DESC"

	pagination := reflect.ValueOf(payload.Pagination)

	if !pagination.IsZero() {
		if pagination.FieldByName("PerPage").Int() < 1 {
			payload.Pagination.PerPage = 10
		}

		values = append(values, payload.PerPage)
		query = fmt.Sprintf("%s \nLIMIT $%d", query, counter)
		counter++

		if pagination.FieldByName("Page").Int() > 0 {
			values = append(values, (payload.Page-1)*payload.PerPage)
			query = fmt.Sprintf("%s \nOFFSET $%d", query, counter)
		}
	}

	return query, values
}
