package db

import "strings"

// RawOrderBySQL generates an SQL ORDER BY clause based on the provided list of orders.
//
// Parameters:
// - o: a list of orders of type ListOrders.
// - separator: a string used to split the orders into sort fields and sort directions.
//
// Return:
// - a string representing the generated SQL ORDER BY clause.
func RawOrderBySQL[ListOrders ~[]string](o ListOrders, separator string) string {
	if len(o) == 0 {
		return ""
	}
	orders := []string{}

	for _, v := range o {
		sort := strings.Split(v, separator)
		if len(sort) != 2 {
			orders = append(orders, v)
			continue
		}
		switch strings.ToLower(sort[1]) {
		case "desc":
			orders = append(orders, sort[0]+" DESC")
		default:
			orders = append(orders, sort[0]+" ASC")
		}
	}

	return "ORDER BY " + strings.Join(orders, ", ")
}
