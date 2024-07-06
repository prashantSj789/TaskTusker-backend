package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)


type CreateuserRequest struct{
  Name string `json:"name"`
  UserName string `json:"userName"`
  Password string `json:"password"`
  Email string `json:"email"`
}
type createuserMessage struct{
	Message string `json:"message"`
}
type User struct {
  ID int `json:"id"`
  Name string `json:"name"`
  UserName string `json:"userName"`
  Password string `json:"password"`
  Email string `json:"email"`
  Is_Premium bool `json:"is_premium"`
  Created_at time.Time `json:"created_at"`
}
type UserLoginRequest struct {
  UserName string `json:"userName"`
  Password string `json:"password"`
}
type LoginMessage struct {
  Message string `json:"message"`
  Token string `json:"token"`
}
type Issue struct {
  ID int `json:"id"`
  Problem string `json:"problem"`
  Created_by string `json:"created_by"`
  Organization string `json:"organization"`
  Created_at time.Time `json:"created_ar"`
}

type Organization struct {
  ID int `json:"id"`
  Name string `json:"name"`
  Creator string `json:"creator"`
  Team []string `json:"team"`
  CreateadAt time.Time `json:"time"`
}
type CreateOrganizationRequest struct {
  Name string `json:"name"`
}
type InviteRequest struct {
  Organization string `json:"organization"`
  Team []string `json:"team"`
}
type InviteRequestMessage struct{
  Message string `json:"message"`
}
type DeleteMemberRequest struct{
  Organization string `json:"organization"`
  Member string `json:"member"`
}
type GetAllOrg struct{
  Organizations []*Organization `json:"organizations"`
}
type  CreateIssueRequest struct {
  Problem string `json:"problem"`
  Organization string `json:"organization"`
}
type getallissue struct{
  Issues []*Issue `json:"issues"`
}
type Comment struct{
  ID int `json:"id"`
  Content string `json:"content"`
  PostId int `json:"post_id"`
  CreatedBy string `json:"created_by"`
  Created_at time.Time `json:"created_at"`
}
type DispComment struct{
   ID int `json:"id"`
   Content string `json:"content"`
   CreatedBy string `json:"created_by"`
   Created_at time.Time `json:"created_at"`
}
type getissue struct{
  Issue *Issue `json:"issue"`
  Comments []*DispComment `json:"comments"`
}
type CommentRequest struct{
  Content string `json:"content"`
}
func NewUser(name, username, password, email string)(*User,error){
  loc, _ := time.LoadLocation("Asia/Kolkata")
  pin,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
  if err!=nil{
		return nil,err
	}
  return &User{
    Name: name,
    UserName: username,
    Email: email,
    Password: string(pin),
    Is_Premium: false,
    Created_at: time.Now().In(loc),
  },nil
}
func NewOrganization(name,creater string )(*Organization,error){
  loc, _ := time.LoadLocation("Asia/Kolkata")
  return &Organization{
    Name: name,
    Creator: creater,
    CreateadAt: time.Now().In(loc),  
    },nil
}
func NewIssue(username,Problem,Organization string) (*Issue,error){
  loc, _ := time.LoadLocation("Asia/Kolkata")
  return &Issue{
   Problem: Problem,
   Created_by: username,
   Organization: Organization,
   Created_at: time.Now().In(loc),
  },nil
}
func NewComment(user,content string ,post int ) (*Comment,error){ 
  loc, _ := time.LoadLocation("Asia/Kolkata")
  return &Comment{
    Content: content,
    CreatedBy: user,
    PostId: post,
    Created_at: time.Now().In(loc),
  },nil
}
func (u *User) ValidatePassword(password string) bool{
  return bcrypt.CompareHashAndPassword([]byte(u.Password),[]byte(password))==nil
}