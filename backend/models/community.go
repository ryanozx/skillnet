package models

import (
	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name     string `gorm:"<-:create; not null"`
	OwnerID  string `json:"-" gorm:"<-:create; not null"`
	User     User   `json:"-" gorm:"foreignKey:OwnerID"`
	About    null.String
	Projects []Project `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Posts    []Post    `json:"-" gorm:"constraint:OnDelete:CASCADE"`
}

func (c *Community) TestFormat() *Community {
	output := Community{
		Model: c.Model,
		Name:  c.Name,
		About: c.About,
	}
	return &output
}

func (c *Community) GetUserID() string {
	return c.OwnerID
}

type CommunityView struct {
	Community Community
	IsOwner   bool
}

func (cv *CommunityView) TestFormat() *CommunityView {
	output := CommunityView{
		Community: *cv.Community.TestFormat(),
		IsOwner:   cv.IsOwner,
	}
	return &output
}

func (c *Community) CommunityView(userID string) *CommunityView {
	output := &CommunityView{
		Community: *c,
		IsOwner:   c.OwnerID == userID,
	}
	return output
}

type CommunityArray struct {
	Communities []CommunityView
	NextPageURL string
}

func (ca *CommunityArray) TestFormat() *CommunityArray {
	if len(ca.Communities) == 0 {
		return &CommunityArray{
			NextPageURL: ca.NextPageURL,
		}
	}
	output := &CommunityArray{
		Communities: []CommunityView{},
		NextPageURL: ca.NextPageURL,
	}
	for _, cv := range ca.Communities {
		output.Communities = append(output.Communities, *cv.TestFormat())
	}
	return output
}
