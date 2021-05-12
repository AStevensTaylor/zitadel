package handler

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Statement struct {
	Sequence         uint64
	PreviousSequence uint64
	TableName        string

	execute func(*sql.Tx) error
}

func NewCreateStatement(table string, values []Column, sequence, previousSequence uint64) Statement {
	cols, params, args := columnsToQuery(values)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ", "), strings.Join(params, ", "))

	return Statement{
		TableName:        table,
		Sequence:         sequence,
		PreviousSequence: previousSequence,
		execute: func(tx *sql.Tx) error {
			_, err := tx.Exec(query, args...)
			return err
		},
	}
}

func NewUpdateStatement(table string, pk []Column, values []Column, sequence, previousSequence uint64) Statement {
	cols, params, args := columnsToQuery(values)
	wheres, whereArgs := columnsToWhere(pk, len(params))
	args = append(args, whereArgs...)
	query := fmt.Sprintf("UPDATE %s SET (%s) = (%s) WHERE %s", table, strings.Join(cols, ", "), strings.Join(params, ", "), strings.Join(wheres, " AND "))

	return Statement{
		TableName:        table,
		Sequence:         sequence,
		PreviousSequence: previousSequence,
		execute: func(tx *sql.Tx) error {
			_, err := tx.Exec(query, args...)
			return err
		},
	}
}

func NewDeleteStatement(table string, conditions []Column, sequence, previousSequence uint64) Statement {
	wheres, args := columnsToWhere(conditions, 0)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, strings.Join(wheres, " AND "))

	return Statement{
		TableName:        table,
		Sequence:         sequence,
		PreviousSequence: previousSequence,
		execute: func(tx *sql.Tx) error {
			_, err := tx.Exec(query, args...)
			return err
		},
	}
}

func NewNoOpStatement(table string, sequence, previousSequence uint64) Statement {
	return Statement{
		TableName:        table,
		Sequence:         sequence,
		PreviousSequence: previousSequence,
	}
}

func (stmt *Statement) Execute(tx *sql.Tx) error {
	if stmt.execute == nil {
		return nil
	}
	return stmt.execute(tx)
}

type Column struct {
	Name  string
	Value interface{}
}

func columnsToQuery(cols []Column) (names []string, parameters []string, values []interface{}) {
	names = make([]string, len(cols))
	values = make([]interface{}, len(cols))
	parameters = make([]string, len(cols))
	for i, col := range cols {
		names[i] = col.Name
		values[i] = col.Value
		parameters[i] = "$" + strconv.Itoa(i+1)

	}
	return names, parameters, values
}

func columnsToWhere(cols []Column, paramOffset int) (wheres []string, values []interface{}) {
	wheres = make([]string, len(cols))
	values = make([]interface{}, len(cols))

	for i, col := range cols {
		wheres[i] = "(" + col.Name + " = $" + strconv.Itoa(i+1+paramOffset) + ")"
		values[i] = col.Value
	}

	return wheres, values
}
