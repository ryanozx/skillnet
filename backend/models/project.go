package models

import (
	"fmt"
)

type ProjectsArray struct {
	Projects    []ProjectMinimal `json:"projects"`
	NextPageURL string
}

func (pArray *ProjectsArray) TestFormat() *ProjectsArray {
	if len(pArray.Projects) == 0 {
		return &ProjectsArray{
			NextPageURL: pArray.NextPageURL,
		}
	}
	output := ProjectsArray{
		Projects:    []ProjectMinimal{},
		NextPageURL: pArray.NextPageURL,
	}
	for _, project := range pArray.Projects {
		output.Projects = append(output.Projects, *project.TestFormat())
	}
	return &output
}

type ProjectMinimal struct {
	ID            uint
	Name          string
	URL           string    `gorm:"-:all"`
	CommunityID   uint      `gorm:"<-:create; not null"`
	Community     Community `json:"-"`
	CommunityName string    `gorm:"-:all"`
	ProjectImgURI string
	ProjectInfo   string
}

func (pm *ProjectMinimal) TestFormat() *ProjectMinimal {
	return pm
}

type Project struct {
	ProjectMinimal `gorm:"embedded"`
	OwnerID        string              `json:"-" gorm:"<-:create; not null"`
	User           User                `json:"-" gorm:"foreignKey:OwnerID"`
	Members        []ProjectMembership `json:"-"`
	PublicCanPost  bool
	Posts          []Post `json:"-" gorm:"constraint:OnDelete:CASCADE"`
}

func (p *Project) TestFormat() *Project {
	output := Project{
		ProjectMinimal: *p.ProjectMinimal.TestFormat(),
		User:           *p.User.TestFormat(),
		PublicCanPost:  p.PublicCanPost,
	}
	return &output
}

type ProjectView struct {
	ProjectMinimal
	Owner         UserMinimal
	PublicCanPost bool
	IsOwner       bool
}

func (p *Project) ProjectView(userID string) *ProjectView {
	output := ProjectView{
		ProjectMinimal: *p.GetProjectMinimal(),
		Owner:          *p.User.GetUserMinimal(),
		PublicCanPost:  p.PublicCanPost,
		IsOwner:        userID == p.OwnerID,
	}
	return &output
}

func (p *Project) GetProjectMinimal() *ProjectMinimal {
	p.URL = GenerateProjectURL(p)
	p.CommunityName = p.Community.Name
	return &p.ProjectMinimal
}

func (p *Project) GetUserID() string {
	return p.OwnerID
}

type ProjectMembership struct {
	UserID    string
	User      User
	ProjectID uint
	Project   Project
}

func GenerateProjectURL(project *Project) string {
	url := fmt.Sprintf("%s/projects/%d", ClientAddress, project.ID)
	return url
}
