package test

import (
	"github.com/stretchr/testify/assert"
    "gopkg.in/DATA-DOG/go-sqlmock.v2"
	"testing"
	"github.com/maszuari/apixen/models"
)

func TestFindDatabaseRecord(t *testing.T) {
	
	mockDB, mock, sqlxDB := MockDB(t)
	defer mockDB.Close()
	
	var cols []string = []string{"id", "comment", "orgid", "created"}
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(1, "test 123", 1, "2019-06-19 09:30:12 am"))

	om := modelorg.NewOrgModel(sqlxDB)
	rs, err := om.FindCommentByID(1)
	if err!=nil{
		t.Fatalf("An error '%s' was not expecting", err)
	}
	expect := modelorg.Comment{ID:1, Comment:"test 123", OrgId: 1, Created: "2019-06-19 09:30:12 am"}
	assert.Equal(t, expect, rs)

}