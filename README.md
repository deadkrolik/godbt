# GODBT
Golang package for database testing (like PHP DBunit)

## Start using it

1. Download and install it:

```sh
$ go get -u github.com/deadkrolik/godbt
```

2. Import it in your code:

```go
package test

import (
    "github.com/deadkrolik/godbt"
    "github.com/deadkrolik/godbt/contract"
)
```
3. Create test object, shared for all test functions and configure mysql conection:

```go
var (
    tester *godbt.Tester
    err    error
)

func init() {
    tester, err = godbt.GetTester(contract.InstallerConfig{
        Type:        "mysql",
        ConnString:  "user:password@/database?charset=utf8&parseTime=True&loc=Local",
        ClearMethod: contract.ClearMethodTruncate,
    })
    if err != nil {
        panic(err.Error())
    }
}
```
## Loading images

Image - set of rows you can load from different sources.

1. XML-file as a source:

```go
image, err := tester.GetImageManager().LoadImage("path/to/file.xml")
if err != nil {
    //...
}
```

2. XML-formatted string as a source:

```go
image, err := tester.GetImageManager().LoadImage(`
<?xml version="1.0" ?>
<dataset>
    <data id="1" k1="1" k2="2"/>
    <data id="2" k1="2" k2="3"/>
</dataset>`)
if err != nil {
    //...
}
```

Root tag should be "dataset", "data" in example is a table name, "k1", "k2" - columns names.

3. Second optional parameter for LoadImage is a list of replacement functions:

```go
contract.ModifiersList{
    "value1": func(table string, key string, value string) string {
        return "modified_" + value
    },
})
```

If we have XML-file like that:

```xml
<?xml version="1.0" ?>
<dataset>
    <data k1="value1" k2="value2" k3="---value1---"/>
</dataset>
```

The result Image will have k1="modified_value1", k2="value2", k3="---modified_value1---". So replacement
function will be called if attribute value contains your string. It can be used to deal with templates (replace
"{CURRENT_DATE}" to real date, for example).

## Installing images

1. Here is a code, that inserts dataset to real database: 

```go
err = tester.GetInstaller().InstallImage(image)
if err != nil {
    //...
}
```

2. It can be done in transaction (it's faster, but you can't use ClearMethodTruncate):

```go
installer := tester.GetInstaller()
err = installer.WithTransaction()
if err != nil {
    //...
}

err = installer.SetClearMethod(contract.ClearMethodDeleteAll).InstallImage(yourImage)
if err != nil {
    //...
}

//here is you app code

err = installer.Rollback()
if err != nil {
    //...
}
```

## Checking test results

1. After loading and installing image you should run your real app code, that works with database.
 
2. Then you can check if you database state correct and the same as you predefined state, that can be described as another Image.

3. Checking count of rows in real table:

```go
count, err := tester.GetInstaller().GetTableRowsCount("data")
if err != nil {
    //...
}
//count assert ...
```

4. Checking state:

```go
realImage, err := tester.GetInstaller().GetTableImage(
    "data",//table name
    []string{"column1", "column2", "column3"},//only columns
    map[string]int{//rows order is important
        "column1": contract.SortAsc,
        "column2": contract.SortDesc,
    },
)
if err != nil {
    //...
}

testImage, err := tester.GetImageManager().LoadImage(`
     <?xml version="1.0" ?>
     <dataset>
         <data column1="1" column2="2" column3="1"/>
         <data column1="2" column2="2" column3="2"/>
     </dataset>`
 )
 if err != nil {
     //...
 }

diffs := tester.GetImageManager().GetImagesDiff(testImage, realImage)
t.Log(diffs)
```

If images a same *diffs* will have zero length. If not - will contain error messages.
