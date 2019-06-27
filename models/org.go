package modelorg

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type CommentOrg struct {
	ID      uint64 `db:"id" json:"id"`
	Comment string `db:"comment" json:"comment"`
	Orgname string `db:"orgname" json:"orgname"`
	Created string `db:"created" json:"created"`
}

type Comment struct {
	ID      uint64 `db:"id"`
	Comment string `db:"comment"`
	OrgId   uint64 `db:"orgid"`
	Created string `db:"created"`
}

type Member struct {
	Email       string
	Username    string
	Avatarurl   string
	Followerno  uint64
	Followingno uint64
}

type OrgModelImp interface {
	SaveComment(orgname string, com string) error
	IsOrgNameAvailable(orgname string) (bool, error)
	GetCommentsByOrgName(orgname string) ([]*CommentOrg, error)
	GetMembersByOrgName(orgname string) ([]*Member, error)
	DeleteCommentsByOrgName(orgname string) error
	FindCommentByID(id int) (Comment, error)
}

type OrgModel struct {
	db *sqlx.DB
}

func NewOrgModel(db *sqlx.DB) *OrgModel {
	return &OrgModel{db: db}
}

func (om *OrgModel) SaveComment(orgname string, com string) error {
	var orgid int
	err := om.db.Get(&orgid, "SELECT id FROM organizations WHERE shortname = $1", orgname)
	if err != nil {
		return err
	}
	//i := strconv.Itoa(orgid)
	tx := om.db.MustBegin()
	tx.MustExec("INSERT INTO comments (comment, orgid) VALUES ($1, $2)", com, orgid)
	tx.Commit()
	return nil
}

func (om *OrgModel) IsOrgNameAvailable(orgname string) (bool, error) {

	var count = 0
	err := om.db.Get(&count, "SELECT COUNT(id) FROM organizations WHERE shortname = $1", orgname)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (om *OrgModel) GetCommentsByOrgName(orgname string) ([]*CommentOrg, error) {
	list := []*CommentOrg{}
	err := om.db.Select(&list, "SELECT a.id AS id, a.comment AS comment, b.shortname AS orgname, a.created AS created FROM comments a INNER JOIN organizations b ON a.orgid = b.id WHERE b.shortname = $1 ORDER BY a.created DESC", orgname)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (om *OrgModel) GetMembersByOrgName(orgname string) ([]*Member, error) {
	list := []*Member{}
	err := om.db.Select(&list, "SELECT email, username, avatarurl, followerno, followingno FROM members a INNER JOIN organizations b ON a.orgid=b.id WHERE b.shortname=$1 ORDER BY followerno DESC", orgname)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (om *OrgModel) DeleteCommentsByOrgName(orgname string) error {

	list := []*Comment{}
	err := om.db.Select(&list, "SELECT a.id, a.comment, a.orgid , a.created FROM comments a INNER JOIN organizations b ON a.orgid = b.id WHERE b.shortname = $1", orgname)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return nil
	}

	orgId := list[0].OrgId
	tx := om.db.MustBegin()
	for _, item := range list {
		tx.MustExec("INSERT INTO deletedcomments (comment, orgid, created ) VALUES ($1, $2, $3)", item.Comment, item.OrgId, item.Created)
	}

	tx.MustExec("DELETE FROM comments WHERE orgid = $1", orgId)
	tx.Commit()

	return nil
}

func (om *OrgModel) FindCommentByID(id int) (Comment, error) {

	c := Comment{}
	err := om.db.Get(&c, "SELECT * FROM comments WHERE id = $1", id)
	if err != nil {
		return c, err
	}
	return c, nil
}
