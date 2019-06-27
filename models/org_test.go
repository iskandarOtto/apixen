package modelorg

import (
	"testing"

	test "github.com/maszuari/apixen/tests"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
)

func TestDeleteCommentsByOrgName(t *testing.T) {
	mockDB, mock, sqlxDB := test.MockDB(t)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "comment", "orgid", "created"}).
		AddRow(1, "Comment 1", 1, "2019-06-18").
		AddRow(2, "Comment 2", 1, "2019-06-19")

	mock.ExpectQuery("SELECT (.+) FROM comments a INNER JOIN organizations b (.+)").
		WithArgs("acme").
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT (.+)").
		WithArgs("Comment 1", 1, "2019-06-18").   // comment, orgid, created
		WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row
	mock.ExpectExec("INSERT (.+)").
		WithArgs("Comment 2", 1, "2019-06-19").   // comment, orgid, created
		WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row
	mock.ExpectExec("DELETE (.+)").
		WithArgs(1).                              //orgid
		WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row
	mock.ExpectCommit()

	om := NewOrgModel(sqlxDB)
	err := om.DeleteCommentsByOrgName("acme")

	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetMembersByOrgName(t *testing.T) {
	mockDB, mock, sqlxDB := test.MockDB(t)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"email", "username", "avatarurl", "followerno", "followingno"}).
		AddRow("abe@acme.com", "abbie", "/avatar/abbie.png", 3, 1).
		AddRow("bob@acme.com", "bobby", "/avatar/abbie.png", 8, 9)

	mock.ExpectQuery("SELECT (.+) FROM members a INNER JOIN organizations b (.+)").
		WithArgs("acme").
		WillReturnRows(rows)

	om := NewOrgModel(sqlxDB)
	rs, err := om.GetMembersByOrgName("acme")

	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}

	data := []*Member{}
	m1 := &Member{"abe@acme.com", "abbie", "/avatar/abbie.png", 3, 1}
	m2 := &Member{"bob@acme.com", "bobby", "/avatar/abbie.png", 8, 9}
	data = append(data, m1)
	data = append(data, m2)
	expect := data
	assert.Equal(t, expect, rs)
}

func TestGetCommentsByOrgName(t *testing.T) {
	mockDB, mock, sqlxDB := test.MockDB(t)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "comment", "orgname", "created"}).
		AddRow(1, "Comment 1", "acme", "2019-06-18").
		AddRow(2, "Comment 2", "acme", "2019-06-19")

	mock.ExpectQuery("SELECT (.+) FROM comments a INNER JOIN organizations b (.+)").
		WithArgs("acme").
		WillReturnRows(rows)

	om := NewOrgModel(sqlxDB)
	rs, err := om.GetCommentsByOrgName("acme")

	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}

	frows := []*CommentOrg{}
	c1 := &CommentOrg{1, "Comment 1", "acme", "2019-06-18"}
	c2 := &CommentOrg{2, "Comment 2", "acme", "2019-06-19"}
	frows = append(frows, c1)
	frows = append(frows, c2)
	expect := frows
	assert.Equal(t, expect, rs)
}

func TestSaveComment(t *testing.T) {
	mockDB, mock, sqlxDB := test.MockDB(t)
	defer mockDB.Close()

	var cols []string = []string{"id"}
	orgid := 1
	mock.ExpectQuery("SELECT * ").
		WithArgs("acme").
		WillReturnRows(sqlmock.NewRows(cols).
			AddRow(orgid))

	mock.ExpectBegin()
	mock.ExpectExec("INSERT (.+)").
		WithArgs("test 123", orgid).              // comment, orgid
		WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row
	mock.ExpectCommit()

	om := NewOrgModel(sqlxDB)
	err := om.SaveComment("acme", "test 123")
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestIsOrgNameAvailable(t *testing.T) {

	mockDB, mock, sqlxDB := test.MockDB(t)
	defer mockDB.Close()

	var cols []string = []string{"count"}
	mock.ExpectQuery("SELECT * ").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(1))

	om := NewOrgModel(sqlxDB)
	rs, err := om.IsOrgNameAvailable("acme")
	if err != nil {
		t.Fatalf("An error '%s' was not expecting", err)
	}

	expect := true
	assert.Equal(t, expect, rs)

}
