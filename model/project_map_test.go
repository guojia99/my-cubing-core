package model

import (
	"fmt"
	"testing"

	"github.com/guojia99/go-tables/table"
)

func TestTableProjectsItems(t *testing.T) {
	tb := table.DefaultSimpleTable(projectsItems)
	fmt.Println(tb)
}
