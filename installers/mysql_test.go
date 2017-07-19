package installers

import (
	"strings"
	"testing"

	"github.com/deadkrolik/godbt/contract"
)

func TestGetImageQuery(t *testing.T) {
	var (
		err   error
		query string
	)
	i, _ := GetInstallerMysql(contract.InstallerConfig{
		DisableConnCheck: true,
	})

	query, err = i.getImageQuery("test")
	if !strings.Contains(query, "SELECT * FROM `test`") {
		t.Fatalf("query should has only select pattern, not `%s`", query)
	}

	query, err = i.getImageQuery("test", 123)
	if err == nil {
		t.Fatal("getImageQuery second param should not be int")
	}

	query, err = i.getImageQuery("test", []string{"a1", "a2"})
	if !strings.Contains(query, "SELECT `a1`, `a2` FROM `test`") {
		t.Fatalf("query should contain fields to select, not `%s`", query)
	}

	query, err = i.getImageQuery("test", []string{})
	if !strings.Contains(query, "SELECT * FROM `test`") {
		t.Fatalf("query should contain * if fields are not set, not `%s`", query)
	}

	query, err = i.getImageQuery("test", []string{}, 123)
	if err == nil {
		t.Fatal("getImageQuery third param should not be int")
	}

	query, err = i.getImageQuery("test", []string{}, map[string]int{
		"a": contract.SortAsc,
		"b": contract.SortDesc,
	})
	if !strings.Contains(query, " ORDER BY") || !strings.Contains(query, "`a` ASC") || !strings.Contains(query, "`b` DESC") {
		t.Fatalf("query should contain ORDER BY statement, not `%s`", query)
	}
}
