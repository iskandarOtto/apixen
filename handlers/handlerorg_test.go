package handlerorg

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	modelorg "github.com/maszuari/apixen/models"
	"github.com/stretchr/testify/assert"
)

type (
	OrgModelStub struct{}
)

func (om *OrgModelStub) SaveComment(orgname string, com string) error {
	list := []*modelorg.CommentOrg{}
	if orgname == "" {
		return errors.New("Oganization name is required")
	}

	if com == "" {
		return errors.New("Comment is required")
	}
	list = append(list, &modelorg.CommentOrg{100, com, orgname, "2019-06-20"})
	return nil
}

func (om *OrgModelStub) IsOrgNameAvailable(orgname string) (bool, error) {
	return true, nil
}

func (om *OrgModelStub) GetCommentsByOrgName(orgname string) ([]*modelorg.CommentOrg, error) {
	list := []*modelorg.CommentOrg{}
	list = append(list, &modelorg.CommentOrg{90, "test 111", "acme", "2019-06-20"})
	list = append(list, &modelorg.CommentOrg{91, "test abc", "acme", "2019-06-21"})
	return list, nil
}

func (om *OrgModelStub) GetMembersByOrgName(orgname string) ([]*modelorg.Member, error) {
	members := []*modelorg.Member{
		{Email: "abby@acme.com", Username: "abby", Avatarurl: "/avatar/abby.png", Followerno: 4, Followingno: 5},
		{Email: "ronan@acme.com", Username: "ronan", Avatarurl: "/avatar/ronan.png", Followerno: 6, Followingno: 7},
	}
	return members, nil
}

func (om *OrgModelStub) DeleteCommentsByOrgName(orgname string) error {

	if orgname != "acme" {
		return errors.New("Invalid organization name")
	}
	return nil
}
func (om *OrgModelStub) FindCommentByID(id int) (modelorg.Comment, error) {
	c := modelorg.Comment{}
	return c, nil
}

func TestGetMembersByOrgName(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/orgs/acme/members/", nil)
	if err != nil {
		t.Fatal(err)
	}

	o := &OrgModelStub{}
	h := NewHandler(o)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/orgs/{orgname}/members/", h.GetMembersByOrgName)
	router.ServeHTTP(rr, req)

	expect, err := createJSONMemberResponse("n", "Success")
	assert.Equal(t, http.StatusOK, rr.Code)
	actual := strings.TrimSuffix(rr.Body.String(), "\n")
	assert.Equal(t, expect, actual)
}

func TestDeleteCommentsByOrgName(t *testing.T) {

	req, err := http.NewRequest("DELETE", "/orgs/acme/comments", nil)
	if err != nil {
		t.Fatal(err)
	}

	o := &OrgModelStub{}
	h := NewHandler(o)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/orgs/{orgname}/comments", h.DeleteCommentsByOrgName)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expect, err := createJSONGenericResponse("n", "Success")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, rr.Code)
	actual := strings.TrimSuffix(rr.Body.String(), "\n")
	assert.Equal(t, expect, actual)
}

func TestSaveComment(t *testing.T) {

	jsonBody := map[string]interface{}{
		"comment": "This is a comment",
	}
	body, _ := json.Marshal(jsonBody)

	req, err := http.NewRequest("POST", "/orgs/acme/comments/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	o := &OrgModelStub{}
	h := NewHandler(o)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/orgs/{orgname}/comments/", h.SaveComment)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expect, err := createJSONGenericResponse("n", "Success")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, rr.Code)
	actual := strings.TrimSuffix(rr.Body.String(), "\n")
	assert.Equal(t, expect, actual)
}

func TestGetCommentsByOrgName(t *testing.T) {

	req, err := http.NewRequest("GET", "/orgs/acme/comments/", nil)
	if err != nil {
		t.Fatal(err)
	}

	o := &OrgModelStub{}
	h := NewHandler(o)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/orgs/{orgname}/comments/", h.GetCommentsByOrgName)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expect := `{"error":"n","message":"Success","comments":[{"id":90,"comment":"test 111","orgname":"acme","created":"2019-06-20"},{"id":91,"comment":"test abc","orgname":"acme","created":"2019-06-21"}]}`
	assert.Equal(t, http.StatusOK, rr.Code)
	actual := strings.TrimSuffix(rr.Body.String(), "\n")
	assert.Equal(t, expect, actual)
}

func TestHello(t *testing.T) {

	req, err := http.NewRequest("GET", "/hello/abc", nil)
	if err != nil {
		t.Fatal(err)
	}

	o := &OrgModelStub{}
	h := NewHandler(o)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/hello/{name}", h.Hello)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expect := "Hello abc"

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expect, rr.Body.String())

}

func createJSONGenericResponse(e string, m string) (string, error) {
	gr := GenericRes{Error: e, Message: m}
	out, err := json.Marshal(gr)
	if err != nil {
		return "none", err
	}
	return string(out), nil
}

func createJSONMemberResponse(e string, m string) (string, error) {
	mr := MembersRes{Error: e, Message: m}
	var list []*modelorg.Member
	if e == "n" {
		list = []*modelorg.Member{
			{Email: "abby@acme.com", Username: "abby", Avatarurl: "/avatar/abby.png", Followerno: 4, Followingno: 5},
			{Email: "ronan@acme.com", Username: "ronan", Avatarurl: "/avatar/ronan.png", Followerno: 6, Followingno: 7},
		}
	} else {
		list = []*modelorg.Member{}
	}
	mr.Members = list
	out, err := json.Marshal(mr)
	if err != nil {
		return "none", err
	}
	return string(out), nil
}
