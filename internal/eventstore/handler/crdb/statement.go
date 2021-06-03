package crdb

import (
	"strconv"
	"strings"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

func NewCreateStatement(aggregateType eventstore.AggregateType, sequence, previousSequence uint64, values []handler.Column) handler.Statement {
	cols, params, args := columnsToQuery(values)
	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")

	return handler.Statement{
		AggregateType:    aggregateType,
		Sequence:         sequence,
		PreviousSequence: previousSequence,
		Execute: func(ex handler.Executer, projectionName string) error {
			if aggregateType == "" {
				return handler.ErrNoAggregateType
			}
			if projectionName == "" {
				return handler.ErrNoTable
			}
			if previousSequence >= sequence {
				return handler.ErrPrevSeqGtSeq
			}
			if len(values) == 0 {
				return handler.ErrNoValues
			}
			query := "INSERT INTO " + projectionName + " (" + columnNames + ") VALUES (" + valuesPlaceholder + ")"
			_, err := ex.Exec(query, args...)
			return err
		},
	}
}

func NewUpsertStatement(aggregateType eventstore.AggregateType, sequence, previousSequence uint64, values []handler.Column) handler.Statement {
	cols, params, args := columnsToQuery(values)
	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")

	return handler.Statement{
		AggregateType:    aggregateType,
		Sequence:         sequence,
		PreviousSequence: previousSequence,
		Execute: func(ex handler.Executer, projectionName string) error {
			if aggregateType == "" {
				return handler.ErrNoAggregateType
			}
			if projectionName == "" {
				return handler.ErrNoTable
			}
			if previousSequence >= sequence {
				return handler.ErrPrevSeqGtSeq
			}
			if len(values) == 0 {
				return handler.ErrNoValues
			}
			query := "UPSERT INTO " + projectionName + " (" + columnNames + ") VALUES (" + valuesPlaceholder + ")"
			_, err := ex.Exec(query, args...)
			return err
		},
	}
}

func NewUpdateStatement(aggregateType eventstore.AggregateType, sequence, previousSequence uint64, values, conditions []handler.Column) handler.Statement {
	cols, params, args := columnsToQuery(values)
	wheres, whereArgs := columnsToWhere(conditions, len(params))
	args = append(args, whereArgs...)

	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")
	wheresPlaceholders := strings.Join(wheres, " AND ")

	return handler.Statement{
		Sequence:         sequence,
		PreviousSequence: previousSequence,
		Execute: func(ex handler.Executer, projectionName string) error {
			if aggregateType == "" {
				return handler.ErrNoAggregateType
			}
			if projectionName == "" {
				return handler.ErrNoTable
			}
			if previousSequence >= sequence {
				return handler.ErrPrevSeqGtSeq
			}
			if len(values) == 0 {
				return handler.ErrNoValues
			}
			if len(conditions) == 0 {
				return handler.ErrNoCondition
			}
			query := "UPDATE " + projectionName + " SET (" + columnNames + ") = (" + valuesPlaceholder + ") WHERE " + wheresPlaceholders
			_, err := ex.Exec(query, args...)
			return err
		},
	}
}

func NewDeleteStatement(aggregateType eventstore.AggregateType, sequence, previousSequence uint64, conditions []handler.Column) handler.Statement {
	wheres, args := columnsToWhere(conditions, 0)

	wheresPlaceholders := strings.Join(wheres, " AND ")

	return handler.Statement{
		AggregateType:    aggregateType,
		Sequence:         sequence,
		PreviousSequence: previousSequence,
		Execute: func(ex handler.Executer, projectionName string) error {
			if aggregateType == "" {
				return handler.ErrNoAggregateType
			}
			if projectionName == "" {
				return handler.ErrNoTable
			}
			if previousSequence >= sequence {
				return handler.ErrPrevSeqGtSeq
			}
			if len(conditions) == 0 {
				return handler.ErrNoCondition
			}
			query := "DELETE FROM " + projectionName + " WHERE " + wheresPlaceholders
			_, err := ex.Exec(query, args...)
			return err
		},
	}
}

func NewNoOpStatement(aggregateType eventstore.AggregateType, sequence, previousSequence uint64) handler.Statement {
	return handler.Statement{
		AggregateType:    aggregateType,
		Sequence:         sequence,
		PreviousSequence: previousSequence,
	}
}

func columnsToQuery(cols []handler.Column) (names []string, parameters []string, values []interface{}) {
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

func columnsToWhere(cols []handler.Column, paramOffset int) (wheres []string, values []interface{}) {
	wheres = make([]string, len(cols))
	values = make([]interface{}, len(cols))

	for i, col := range cols {
		wheres[i] = "(" + col.Name + " = $" + strconv.Itoa(i+1+paramOffset) + ")"
		values[i] = col.Value
	}

	return wheres, values
}
