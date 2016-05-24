package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"strconv"

	"database/sql"
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
)

// GetObservedProperty returns an ObservedProperty by id
func (gdb *GostDatabase) GetObservedProperty(id string) (*entities.ObservedProperty, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select id, name, definition, description FROM %s.observedproperty where id = $1", gdb.Schema)
	observedProperty, err := processObservedProperty(gdb.Db, sql, intID)
	if err != nil {
		return nil, err
	}

	return observedProperty, nil
}

// GetObservedPropertyByDatastream returns an ObservedProperty by id
func (gdb *GostDatabase) GetObservedPropertyByDatastream(id string) (*entities.ObservedProperty, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select observedproperty.id, observedproperty.name, observedproperty.definition, observedproperty.description FROM %s.observedproperty inner join %s.datastream on datastream.observedproperty_id = observedproperty.id where datastream.id = $1", gdb.Schema, gdb.Schema)
	observedProperty, err := processObservedProperty(gdb.Db, sql, intID)
	if err != nil {
		return nil, err
	}

	return observedProperty, nil
}

// GetObservedProperties returns all observed properties
func (gdb *GostDatabase) GetObservedProperties() ([]*entities.ObservedProperty, error) {
	sql := fmt.Sprintf("select id, name, definition, description FROM %s.observedproperty", gdb.Schema)
	return processObservedProperties(gdb.Db, sql)
}

func processObservedProperty(db *sql.DB, sql string, args ...interface{}) (*entities.ObservedProperty, error) {
	observedProperties, err := processObservedProperties(db, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(observedProperties) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("ObservedProperty not found"))
	}

	return observedProperties[0], nil
}

func processObservedProperties(db *sql.DB, sql string, args ...interface{}) ([]*entities.ObservedProperty, error) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var observedProperties = []*entities.ObservedProperty{}

	for rows.Next() {
		var opID int
		var name string
		var definition string
		var description string
		err2 := rows.Scan(&opID, &name, &definition, &description)
		if err2 != nil {
			return nil, err2
		}

		op := entities.ObservedProperty{
			ID:          strconv.Itoa(opID),
			Name:        name,
			Definition:  definition,
			Description: description,
		}

		observedProperties = append(observedProperties, &op)
	}

	return observedProperties, nil
}

// PostObservedProperty adds an ObservedProperty to the database
func (gdb *GostDatabase) PostObservedProperty(op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	var opID int
	sql := fmt.Sprintf("INSERT INTO %s.observedproperty (name, definition, description) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, op.Name, op.Definition, op.Description).Scan(&opID)
	if err != nil {
		return nil, err
	}

	op.ID = strconv.Itoa(opID)
	return op, nil
}

// ObservedPropertyExists checks if a ObservedProperty is present in the database based on a given id.
func (gdb *GostDatabase) ObservedPropertyExists(thingID int) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.observedproperty WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}