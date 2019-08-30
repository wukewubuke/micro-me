package model

import (
	"github.com/go-xorm/xorm"
	"time"
)

type (
	Members struct {
		Id int64 `json:"id"`
		Token string `json:"token" xorm:"varchar(100) notnull 'token'"`
		Username string `json:"username" xorm:"varchar(30) notnull 'username'"`
		Password string `json:"password" xorm:"varchar(40) notnull 'password'"`
		CreateTime time.Time `json:"create_time" xorm:"DateTime 'create_time'"`
		UpdateTime time.Time `json:"update_time" xorm:"DateTime 'update_time'"`
	}
	MembersModel struct {
		mysql *xorm.Engine
	}
)


func NewMembersModel(mysql *xorm.Engine) *MembersModel{
	return &MembersModel{
		mysql: mysql,
	}
}


func (m *MembersModel)FindByToken(token string) (*Members, error){
	member := new(Members)
	if _, err := m.mysql.Where("token = ?", token).Get(member); err != nil {
		return nil,err
	}
	return member,nil
}

func (m *MembersModel)FindById(id int64) (*Members, error){
	member := new(Members)
	if _, err := m.mysql.Where("id = ?", id).Get(member); err != nil {
		return nil,err
	}
	return member,nil
}

func (m *MembersModel)FindByUsername(username string) (*Members, error){
	member := new(Members)
	if _, err := m.mysql.Where("username = ?", username).Get(member); err != nil {
		return nil,err
	}
	return member,nil
}


func (m *MembersModel)InsertMember(member *Members) error {
	if _, err := m.mysql.Insert(member); err != nil {
		return err
	}
	return nil
}