package db100

import (
	"errors"
)

type databaseObject interface {
	getIDs() []interface{}
	getTablename() string
	getIDColumns() []string
	getInsertColumns() []string
	getInsertFields() []interface{}
	getUpdateColumns() []string
	getUpdateFields() []interface{}
}

func queryLetter(i int) string {
	var result string
	if i == 0 {
		result = "( "
	} else {
		result = ", "
	}
	return result
}

func addColumnsToQuery(query string, columns []string, where bool) string {
	for i, c := range columns {
		query = query + ` "` + c + `" = ?`
		if i != (len(columns) - 1) {
			if where {
				query = query + " AND"
			} else {
				query = query + ","
			}
		}
	}
	return query
}

func insertDBO(dbo databaseObject) (int, error) {
	var err error
	var nid int
	query := `INSERT INTO "` + dbo.getTablename() + `" `
	columns := dbo.getInsertColumns()
	for i, c := range columns {
		query = query + queryLetter(i) + `"` + c + `"`
	}
	query = query + ") VALUES "
	for i := 0; i < len(columns); i++ {
		query = query + queryLetter(i) + "?"
	}
	query = query + ")"
	query = db.Rebind(query)
	if db.DriverName() == pgDriverName {
		idcols := dbo.getIDColumns()
		if len(idcols) > 1 {
			nid, err = insertDBOPG(query, "", dbo.getInsertFields())
		} else {
			nid, err = insertDBOPG(query, dbo.getIDColumns()[0], dbo.getInsertFields())
		}
	} else {
		nid, err = insertDBOOther(query, dbo.getInsertFields())
	}
	return nid, err
}

func insertDBOPG(query string, returning string, a []interface{}) (int, error) {
	var newid int
	if returning != "" {
		query = query + ` RETURNING "` + returning + `"`
	}
	tx := db.MustBegin()
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return -1, errors.New("Error preparing Statement:" + err.Error())
	}
	if returning != "" {
		stmt.QueryRow(a...).Scan(newid)
		if err != nil {
			tx.Rollback()
			return -1, errors.New("Error executing Statement:" + err.Error())
		}
	} else {
		stmt.QueryRow(a...)
	}
	err = tx.Commit()
	if err != nil {
		return -1, errors.New("Error executing Commit:" + err.Error())
	}
	return newid, nil
}

func insertDBOOther(query string, a []interface{}) (int, error) {
	res, err := db.Exec(query, a...)
	if err != nil {
		return -1, errors.New("Error inserting: " + err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, errors.New("Error fetching new ID: " + err.Error())
	}
	return int(id), nil
}

func updateDBO(dbo databaseObject) error {
	query := `UPDATE "` + dbo.getTablename() + `" SET`
	columns := dbo.getUpdateColumns()
	query = addColumnsToQuery(query, columns, false)
	columns = dbo.getIDColumns()
	query = query + ` WHERE`
	query = addColumnsToQuery(query, columns, true)
	query = db.Rebind(query)
	fields := dbo.getUpdateFields()
	fields = append(fields, dbo.getIDs()...)
	_, err := db.Exec(query, fields...)
	if err != nil {
		return errors.New("Error updating: " + err.Error())
	}
	return nil
}

func buildSelectQuery(dbo databaseObject, whereColumn string, whereValue int) (string, []interface{}) {
	var is []interface{}
	query := `SELECT * FROM "` + dbo.getTablename() + `"`
	useWhere := (whereValue > 0)
	if useWhere {
		query = query + ` WHERE "` + whereColumn + `" = ?`
		is = append(is, whereValue)
	}
	query = db.Rebind(query)
	return query, is
}

func getDetailsDBO(dest interface{}) error {
	dbo := dest.(databaseObject)
	query := `SELECT * FROM "` + dbo.getTablename() + `" WHERE`
	query = addColumnsToQuery(query, dbo.getIDColumns(), true)
	query = query + ` LIMIT 1`
	query = db.Rebind(query)
	err := db.Get(dest, query, dbo.getIDs()...)
	if err != nil {
		return errors.New("Error Selecting:" + err.Error())
	}
	return nil
}
