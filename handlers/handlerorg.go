package handlerorg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	modelorg "github.com/maszuari/apixen/models"
)

type CommentOrgsRes struct {
	Error       string                 `json:"error"`
	Message     string                 `json:"message"`
	CommentOrgs []*modelorg.CommentOrg `json:"comments"`
}

type MembersRes struct {
	Error   string             `json:"error"`
	Message string             `json:"message"`
	Members []*modelorg.Member `json:"members"`
}

type GenericRes struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type UserComment struct {
	Comment string `json:"comment"`
}

type handler struct {
	OrgModel modelorg.OrgModelImp
}

func NewHandler(om modelorg.OrgModelImp) *handler {
	return &handler{om}
}

func (h *handler) GetMembersByOrgName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgname := vars["orgname"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	mr := &MembersRes{}
	list, err := h.OrgModel.GetMembersByOrgName(orgname)
	if err != nil {
		log.Println(err)
		mr.Error = "y"
		mr.Message = "Database error"
	} else {
		mr.Error = "n"
		mr.Message = "Success"
		mr.Members = list
	}
	json.NewEncoder(w).Encode(mr)
}

func (h *handler) DeleteCommentsByOrgName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgname := vars["orgname"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	gr := &GenericRes{}
	err := h.OrgModel.DeleteCommentsByOrgName(orgname)
	if err != nil {
		log.Println(err)
		gr.Error = "y"
		gr.Message = "Database error"
	} else {
		gr.Error = "n"
		gr.Message = "Success"
	}

	json.NewEncoder(w).Encode(gr)
}

func (h *handler) SaveComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgname := vars["orgname"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dec := json.NewDecoder(r.Body)
	var uc UserComment
	err := dec.Decode(&uc)
	r.Body.Close()
	gr := &GenericRes{}
	if err != nil {
		log.Println(err)
		gr.Error = "y"
		gr.Message = "Invalid JSON format"
	} else {
		err = h.OrgModel.SaveComment(orgname, uc.Comment)
		if err != nil {
			log.Println(err)
			gr.Error = "y"
			gr.Message = "Database error"
		} else {
			gr.Error = "n"
			gr.Message = "Success"
		}
	}

	json.NewEncoder(w).Encode(gr)

}

func (h *handler) GetCommentsByOrgName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgname := vars["orgname"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	c, err := h.OrgModel.GetCommentsByOrgName(orgname)
	res := &CommentOrgsRes{}
	if err != nil {
		log.Println(err)
		res.Error = "y"
		res.Message = "Database error"

	} else {
		res.Error = "n"
		res.Message = "Success"
		res.CommentOrgs = c
	}
	json.NewEncoder(w).Encode(res)
}

func (h *handler) Hello(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello "+name)
}
