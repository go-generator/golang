package main

import (
	"database/sql"
	"fmt"
)

func PrepareStatementPicker(dbName string, count int) string {
	switch dbName {
	case "mysql":
		return "?"
	case "postgres":
		return fmt.Sprintf("$%v", count)
	case "sqlserver":
		return fmt.Sprintf("@p%v", count)
	}
	return ""
}
func InsertIntoSyntaxPicker(dbName string, k string) string {
	switch dbName {
	case "mysql":
		return "`" + k + "`"
	case "postgres", "sqlserver":
		return k
	}
	return ""
}
func (t SqlType) GetById() error {
	db2, err := sql.Open(t.Info.DriverName, t.Info.DataSourceName)
	if err != nil {
		return err
	}
	defer db2.Close()
	query := "SELECT * FROM " + t.Info.Source + " WHERE _id =" + PrepareStatementPicker(t.Info.DriverName, 1)
	rows, err := db2.Query(query, t.IdParam)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	tmp2 := ""
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}
		var value string
		tmp := "{"
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			tmp += "\"" + columns[i] + "\": \"" + value + "\","
		}
		tmp2 += tmp[:len(tmp)-1] + "}"
	}
	*t.OutputString = tmp2
	return nil
}
func (t SqlType) GetAll() error {
	db2, err := sql.Open(t.Info.DriverName, t.Info.DataSourceName)
	if err != nil {
		return err
	}
	defer db2.Close()
	query := "SELECT * FROM " + t.Info.Source
	rows, err := db2.Query(query)
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	tmp2 := ""
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}
		var value string
		tmp := "{"
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			tmp += "\"" + columns[i] + "\": \"" + value + "\","
		}
		tmp2 += tmp[:len(tmp)-1] + "}"
	}
	*t.OutputString = tmp2
	return nil
}
func (t SqlType) Create() error {
	db2, err := sql.Open(t.Info.DriverName, t.Info.DataSourceName)
	if err != nil {
		return err
	}
	defer db2.Close()
	query := "INSERT INTO " + t.Info.Source + "("
	query2 := ") VALUES ("
	s1 := []interface{}{}
	count := 1
	for k, v := range t.InputMap {
		query += InsertIntoSyntaxPicker(t.Info.DriverName, k) + ","
		s1 = append(s1, v)
		query2 += PrepareStatementPicker(t.Info.DriverName, count) + ","
		count++
	}
	query2 = query2[:len(query2)-1] + ")"
	query = query[:len(query)-1] + query2
	pre, err := db2.Prepare(query)
	if err != nil {
		return err
	}
	_, err = pre.Exec(s1...)
	if err != nil {
		return err
	}

	*t.OutputString = "Created Successfully"
	return nil
}
func (t SqlType) Update() error {
	db2, err := sql.Open(t.Info.DriverName, t.Info.DataSourceName)
	if err != nil {
		return err
	}
	defer db2.Close()
	query := "UPDATE " + t.Info.Source + " SET "
	count := 1
	s1 := []interface{}{}
	for k, v := range t.InputMap {
		query += InsertIntoSyntaxPicker(t.Info.DriverName, k) + "= " + PrepareStatementPicker(t.Info.DriverName, count) + ","
		count++
		s1 = append(s1, v)
	}
	query = query[:len(query)-1] + " WHERE (_id = " + PrepareStatementPicker(t.Info.DriverName, count) + ")"
	s1 = append(s1, t.IdParam)
	_, err = db2.Query(query, s1...)
	if err != nil {
		return err
	}
	*t.OutputString = "Updated Successfully"
	return nil
}
func (t SqlType) Delete() error {
	db2, err := sql.Open(t.Info.DriverName, t.Info.DataSourceName)
	if err != nil {
		return err
	}
	defer db2.Close()
	query := "DELETE FROM " + t.Info.Source + " WHERE _id=" + PrepareStatementPicker(t.Info.DriverName, 1)
	_, err = db2.Query(query, t.IdParam)
	if err != nil {
		return err
	}
	*t.OutputString = "Delete Successfully"
	return nil
}
