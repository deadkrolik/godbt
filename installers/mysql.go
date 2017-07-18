package installers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/deadkrolik/godbt/contract"
	"strings"
	//need that
	_ "github.com/go-sql-driver/mysql"
)

//Mysql - install image to MySQL
type Mysql struct {
	connection  *sql.DB
	handler     queryAble
	clearMethod int
}

//https://github.com/golang/go/issues/14468
type queryAble interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

//GetInstallerMysql - installer instance
func GetInstallerMysql(config contract.InstallerConfig) (*Mysql, error) {
	db, err := sql.Open("mysql", config.ConnString)
	if err != nil {
		return nil, err
	}

	if !config.DisableConnCheck {
		err = db.Ping()
		if err != nil {
			return nil, err
		}
	}

	return &Mysql{
		connection:  db,
		clearMethod: config.ClearMethod,
		handler:     db,
	}, nil
}

//InstallImage - install Image
func (m *Mysql) InstallImage(image contract.Image) error {
	err := m.clearTables(image)
	if err != nil {
		return err
	}

	for _, row := range image {
		if len(row.Data) == 0 {
			continue
		}

		var values []interface{}
		var keys []string
		for k, v := range row.Data {
			values = append(values, v)
			keys = append(keys, "`"+k+"` = ?")
		}

		query := fmt.Sprintf("INSERT INTO `%s` SET %s", row.Table, strings.Join(keys, ", "))
		stmt, err := m.db().Prepare(query)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	return nil
}

//GetTableRowsCount - table rows count
func (m *Mysql) GetTableRowsCount(table string) (int64, error) {
	var count int64

	row := m.db().QueryRow("SELECT COUNT(*) FROM `" + table + "`")
	err := row.Scan(&count)

	return count, err
}

//GetTableImage - table to Image
func (m *Mysql) GetTableImage(table string, args ...interface{}) (contract.Image, error) {
	query, err := m.getImageQuery(table, args...)
	if err != nil {
		return contract.Image{}, err
	}

	rows, err := m.db().Query(query)
	if err != nil {
		return contract.Image{}, err
	}
	defer func() {
		_ = rows.Close()
	}()

	columns, err := rows.Columns()
	if err != nil {
		return contract.Image{}, err
	}
	columnCount := len(columns)

	var image contract.Image

	for rows.Next() {
		values, err := m.getColumnsValues(rows, columnCount)
		if err != nil {
			return contract.Image{}, err
		}

		stmt := contract.Row{
			Table: table,
		}
		stmt.Data = make(map[string]string, len(columns))
		for i := range columns {
			stmt.Data[columns[i]] = string(values[i])
		}
		image = append(image, stmt)
	}

	err = rows.Err()
	if err != nil {
		return contract.Image{}, err
	}

	return image, nil
}

//WithTransaction - start transaction
func (m *Mysql) WithTransaction() error {
	tx, err := m.connection.Begin()
	if err != nil {
		return err
	}

	m.handler = tx
	return nil
}

//Rollback - rollback previous transaction
func (m *Mysql) Rollback() error {
	tx, ok := m.handler.(*sql.Tx)
	if !ok {
		return errors.New("Can't cast handler to sql.Tx")
	}

	return tx.Rollback()
}

//db - real DB handler
func (m *Mysql) db() queryAble {
	return m.handler
}

//SetClearMethod - rewrite clear method
func (m *Mysql) SetClearMethod(method int) contract.Installer {
	m.clearMethod = method
	return m
}

//getImageQuery - parsing
//second param - columns list, third - sort for columns
func (m *Mysql) getImageQuery(table string, args ...interface{}) (string, error) {
	columnsString := "*"
	if len(args) > 0 {
		columnsList, ok := args[0].([]string)
		if !ok {
			return "", errors.New("Second param for GetTableImage should be []string")
		}

		for i := range columnsList {
			columnsList[i] = "`" + columnsList[i] + "`"
		}
		if len(columnsList) > 0 {
			columnsString = strings.Join(columnsList, ", ")
		}
	}

	orderString := ""
	if len(args) > 1 {
		orders, ok := args[1].(map[string]int)
		if !ok {
			return "", errors.New("Third param for GetTableImage should be map[string]int")
		}

		var ordersList []string
		for key, order := range orders {
			sort := "ASC"
			if order == contract.SortDesc {
				sort = "DESC"
			}
			ordersList = append(ordersList, fmt.Sprintf("`%s` %s", key, sort))
		}
		if len(ordersList) > 0 {
			orderString = "ORDER BY " + strings.Join(ordersList, ", ")
		}
	}

	return fmt.Sprintf("SELECT %s FROM `%s` %s", columnsString, table, orderString), nil
}

//getColumnsValues - DB-row to bytes (golang-nuts/-9h9UwrsX7Q)
func (m *Mysql) getColumnsValues(rows *sql.Rows, count int) ([][]byte, error) {
	var (
		values   [][]byte
		pointers []interface{}
	)

	pointers = make([]interface{}, count)
	values = make([][]byte, count)
	for i := range pointers {
		pointers[i] = &values[i]
	}

	err := rows.Scan(pointers...)
	if err != nil {
		return [][]byte{}, err
	}

	return values, nil
}

//clearTables - clear tables before insert
func (m *Mysql) clearTables(image contract.Image) error {
	if m.clearMethod == contract.ClearMethodNoClear {
		return nil
	}

	tables := make(map[string]bool)
	for _, row := range image {
		tables[row.Table] = true
	}

	for table := range tables {
		query := "TRUNCATE TABLE `%s`"
		if m.clearMethod == contract.ClearMethodDeleteAll {
			query = "DELETE FROM `%s`"
		}

		_, err := m.db().Exec(fmt.Sprintf(query, table))
		if err != nil {
			return err
		}
	}

	return nil
}
