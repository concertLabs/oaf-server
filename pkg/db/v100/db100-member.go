package db100

import (
	"errors"
)

//Member models the m:n relation between Users and Sections
type Member struct {
	SectionID int `json:"sectionid" db:"SectionID"`
	UserID    int `json:"userid" db:"UserID"`
	Rights    int `json:"rights" db:"Rights"`
}

func (m *Member) getIDs() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, m.SectionID)
	interfaceSlice = append(interfaceSlice, m.UserID)
	return interfaceSlice
}

func (m *Member) getTablename() string {
	return "Members"
}

func (m *Member) getIDColumns() []string {
	return []string{"SectionID", "UserID"}
}

func (m *Member) getInsertColumns() []string {
	result := m.getUpdateColumns()
	result = append(result, "SectionID")
	result = append(result, "UserID")
	return result
}

func (m *Member) getInsertFields() []interface{} {
	var interfaceSlice = m.getUpdateFields()
	interfaceSlice = append(interfaceSlice, m.SectionID)
	interfaceSlice = append(interfaceSlice, m.UserID)
	return interfaceSlice
}

func (m *Member) getUpdateColumns() []string {
	return []string{"Rights"}
}

func (m *Member) getUpdateFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, m.Rights)
	return interfaceSlice
}

//Insert inserts a new Member into the database
func (m *Member) Insert() error {
	var err error
	_, err = insertDBO(m)
	if err != nil {
		return errors.New("Error inserting Member:" + err.Error())
	}
	return nil
}

//GetMembers gives back all Members in the Database
func GetMembers(orgid int) ([]Member, error) {
	var m []Member
	var err error
	if orgid < 1 {
		err = db.Select(&m, `SELECT * FROM "Members"`)
	} else {
		query := `SELECT "Members"."SectionID" AS "SectionID", "Members"."UserID" as "UserID", "Members"."Rights" as "Rights"
		FROM "Members", "Sections" WHERE "Members"."SectionID" = "Sections"."SectionID" and "Sections"."OrganizationID" = ?`
		query = db.Rebind(query)
		err = db.Select(&m, query, orgid)
	}
	if err != nil {
		return m, errors.New("Error getting Member:" + err.Error())
	}
	return m, nil
}

//GetDetails takes a Member struct with only the UserID and SectionID and tries to fetch the remaining infos
func (m *Member) GetDetails() error {
	err := getDetailsDBO(m)
	if err != nil {
		return errors.New("Error getting Member details:" + err.Error())
	}
	return nil
}

//Patch patches a Member with new Info from a second struct
func (m *Member) Patch(mm Member) error {
	m.Rights = mm.Rights
	return nil
}

//Update updates the Right Field of a Member in the Database
func (m *Member) Update() error {
	err := updateDBO(m)
	if err != nil {
		return errors.New("Error updating Members: " + err.Error())
	}
	return nil
}

//DeleteMember deletes a Member with the given UserID and SectionID
func DeleteMember(UserID int, SectionID int) error {
	query := db.Rebind(`DELETE FROM "Members" WHERE "UserID" = ? and "SectionID" = ?`)
	_, err := db.Exec(query, UserID, SectionID)
	if err != nil {
		return errors.New("Error deleting Member: " + err.Error())
	}
	return nil
}
