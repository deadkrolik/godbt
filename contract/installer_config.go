package contract

//InstallerConfig - config fo DB
type InstallerConfig struct {
	//Engine for installing dumps, located in /installers/
	Type string
	//DSN string for database
	ConnString string
	//How to clear tables before installing dump
	ClearMethod int
	//Check connection when creating installer instance
	DisableConnCheck bool
}

const (
	//ClearMethodTruncate - clears table using: TRUNCATE table
	//Do not use with transaction because of https://dev.mysql.com/doc/refman/5.5/en/implicit-commit.html
	ClearMethodTruncate = 0
	//ClearMethodNoClear - do not clear table
	ClearMethodNoClear = 1
	//ClearMethodDeleteAll - clears table using: DELETE FROM table
	ClearMethodDeleteAll = 2
)
