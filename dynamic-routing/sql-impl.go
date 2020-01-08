package main

import (
	"database/sql"
)

func (t SqlType) GetById() error {
	db2, err := sql.Open(t.Info.DriverName, t.Info.DataSourceName)
	if err != nil {
		return err
	}
	defer db2.Close()
	query := "SELECT * FROM " + t.Info.Source + " WHERE _id ='" + t.IdParam + "'"
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
	for k, v := range t.InputMap {
		query += k + ","
		tmp, ok := v.(string)
		if ok {
			query2 += "'" + tmp + "',"
		}
	}
	query2 = query2[:len(query2)-1] + ")"
	query = query[:len(query)-1] + query2
	_, err = db2.Query(query)
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
	for k, v := range t.InputMap {
		tmp, ok := v.(string)
		if ok {
			query += k + "= '" + tmp + "',"
		}

	}
	query = query[:len(query)-1] + " WHERE (_id = '" + t.IdParam + "')"
	_, err = db2.Query(query)
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
	query := "DELETE FROM " + t.Info.Source + " WHERE _id= '" + t.IdParam + "'"
	_, err = db2.Query(query)
	if err != nil {
		return err
	}
	*t.OutputString = "Delete Successfully"
	return nil
}
