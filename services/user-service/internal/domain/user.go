package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID        string
	Name      string
	Title     string
	Company   string
	Level     int32
	LevelText string
	JoinedAt  string
	Location  string
	Team      string
}

func (u *User) GetID() string        { return u.ID }
func (u *User) GetName() string      { return u.Name }
func (u *User) GetTitle() string     { return u.Title }
func (u *User) GetCompany() string   { return u.Company }
func (u *User) GetLevel() int32      { return u.Level }
func (u *User) GetLevelText() string { return u.LevelText }
func (u *User) GetJoinedAt() string  { return u.JoinedAt }
func (u *User) GetLocation() string  { return u.Location }
func (u *User) GetTeam() string      { return u.Team }

type Repository interface {
	GetCurrentUser(userID string) (*User, error)
	BatchGetUsers(userIDs []string) ([]*User, error)
}
