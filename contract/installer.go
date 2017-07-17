package contract

//Installer - executes requests based in Image structure in real DB
type Installer interface {
	InstallImage(Image) error
	GetTableRowsCount(string) (int64, error)
	GetTableImage(string, ...interface{}) (Image, error)
	WithTransaction() error
	Rollback() error
	SetClearMethod(int) Installer
}

const (
	//SortAsc - for method GetTableImage
	SortAsc = 0
	//SortDesc - for method GetTableImage
	SortDesc = 1
)
