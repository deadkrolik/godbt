package contract

//Modifier - changing attrs values
type Modifier func(table string, key string, value string) string

//ModifiersList - list of functions
type ModifiersList map[string]Modifier
