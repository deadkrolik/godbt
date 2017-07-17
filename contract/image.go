package contract

//Image - base type for set of rows
type Image []Row

//Row - one row in table
type Row struct {
	Table string
	Data  map[string]string
}
