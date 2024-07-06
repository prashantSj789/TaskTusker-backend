package main

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	_"github.com/pelletier/go-toml/query"
)

type PostgressStore struct{
	db *sql.DB
}
type storage interface{
  CreateUser(*User)error
  CreateOrganization(*Organization)error
  GetUserbyUserName(string)(*User,error)
  GetOrganizationByName(string)(*Organization,error)
  AddMenber(*Organization,string) error
  RemoveMember(*Organization,string)error
  GetOrganizations(string)([]*Organization,error)
  CreateCard(*Card)error
  GetCards(string)([]*Card,error)
  GetCard(int)(*Card,error)
  UpdateCard(*Card) error
  DeleteCard(int) error
  CreateIssue(*Issue) error
  GetIssuses(string) ([]*Issue,error)
  GetIssuebyId(id int)(*Issue,[]*DispComment,error)
  CreateComment(*Comment) error
}
func NewPostgressStore() (*PostgressStore,error){
	conStr:= "user=postgres dbname=postgres host=database-3.c5i8mwka8jer.us-east-1.rds.amazonaws.com port=5432  password=my-go-jira "
	db, err:= sql.Open("postgres",conStr)
	if err!=nil{
		panic(err)
	}
	if err:=db.Ping();err!=nil{
		return nil,err
	}
	return &PostgressStore{
		db: db,
	},nil
}
func (s *PostgressStore) init() (error,error,error,error,error){
	return s.CreateUserTable(),s.CreateOrganizationTable(),s.CreateCardTable(),s.CreateCommentsTable(),s.CreateIssueTable()

}

func (s *PostgressStore) CreateUserTable() error{
	query:= `create table if not exists Users ( 
		id serial primary key, 
		name varchar(50),
		user_name varchar(100),  
		password varchar(200), 
		email varchar(50),
		is_premium boolean,
		createad_at timestamp 
		)`
		_,err := s.db.Exec(query)
		return err	
}
func (s *PostgressStore) CreateOrganizationTable() error{
	Query:= `create table if not exists Organizations (
	    id serial primary key,
	    name varchar(50),
		Created_at timestamp,
		Creater varchar(50),
		team text[]
		)`
		_,err:= s.db.Exec(Query)
		return err
}
func (s *PostgressStore) CreateCardTable() error {
    Query:=`create table if not exists Cards (
	id serial primary key,
	name varchar(50),
	task text,
	organization varchar(50),
	issued_to varchar(50),
	issues_at timestamp,
	state varchar(50)
	)`
	_,err:=s.db.Exec(Query)
	return err
}
func (s *PostgressStore) CreateIssueTable() error{
	Query:=`create table if not exists Issues(
	id serial primary key,
	content text,
	created_by varchar(50),
	organization varchar(50),
	created_at timestamp
	)`
	_,err:=s.db.Exec(Query)
	return err
}
func (s *PostgressStore) CreateCommentsTable() error{
	Query:=`create table if not exists Comments(
	id serial primary key,
	content text,
	created_by varchar(50),
	post_id integer,
	created_at timestamp
	)`
	_,err:=s.db.Exec(Query)
	return err
}
func (s *PostgressStore) CreateUser(u *User) error{
	query:=`insert into Users
	(name,user_name,password,email,is_premium,createad_at)
	values($1,$2,$3,$4,$5,$6)`
	_,err:=s.db.Query(
		query,
		u.Name,
		u.UserName,
		u.Password,
		u.Email,
		u.Is_Premium,
		u.Created_at,
	)
	if  err!=nil {
		return err;
	}
	return nil
}
func (s *PostgressStore) CreateOrganization(o *Organization) error{
	queryy:=`insert into Organizations
	(name,Created_at,Creater,team)
	values($1,$2,$3,$4)`

	_,err:=s.db.Query(
		queryy,
		o.Name,
		o.CreateadAt,
		o.Creator,
		pq.Array(o.Team),
	)
	if  err!=nil {
		return err;
	}
	return nil
}
func (s *PostgressStore) CreateCard(c *Card) error {
	query:=`insert into Cards
	(name,task,organization,issued_to,issues_at,state)
	values($1,$2,$3,$4,$5,$6)`

	_,err:=s.db.Query(
		query,
        c.Name,
		c.Task,
		c.Organization,
		c.Issued_to,
		c.Issued_at,
		c.Lane,
	)
	if err!=nil{
       return err
	}
	return nil
}
func (s *PostgressStore) GetUserbyUserName(username string) (*User,error){
	rows,err:=s.db.Query("select *from Users where user_name = $1",username)
	if err!=nil{
		return nil,err
	}
	for rows.Next(){
		return Scanintouser(rows)
	}
   return nil,fmt.Errorf("failed to Scan user %s",err)
	
}

func (s *PostgressStore) GetOrganizationByName(name string) (*Organization,error){
	rows,err:=s.db.Query("select * from Organizations where name = $1",name)
	if err!=nil{
		return nil,err
	}
	for rows.Next(){
		return ScanintoOrganization(rows)
	}
   return nil,fmt.Errorf("failed to Scan org %s",err)

}
func (s *PostgressStore) AddMenber(o *Organization,name string) error{
	query:=`UPDATE organizations
	SET team = array_append(team,$1)
	where name=$2`
	_,err:=s.db.Query(
		query,
		name,
		o.Name,
	)
	if err!=nil{
		return err
	}
    return nil
}
func (s *PostgressStore) RemoveMember(o *Organization,name string) error{
	query:=`UPDATE organizations
	SET team = array_remove(team,$1)
	where name=$2`
	_,err:=s.db.Query(
		query,
		name,
		o.Name,
	)
	if err!=nil{
		return err
	}
	return nil
}
func (s *PostgressStore) GetOrganizations(name string) ([]*Organization,error){
	rows,err:=s.db.Query(`SELECT * from organizations
	where $1 = ANY("team"::text[]) or creater = $1`,name)
	if err!=nil{
		return nil,err
	}
    orgs:=[]*Organization{

	}
	for rows.Next(){
      org,err:=ScanintoOrganization(rows)
	  if err!=nil{
		return nil,err
	  }
	  orgs = append(orgs, org)
	}
	return orgs,nil
}
func (s *PostgressStore) GetCards(org string) ([]*Card,error){
	rows,err:=s.db.Query(`SELECT * from cards 
	where organization = $1`,org)
	if err!=nil{
	  return nil,err
	}
	cards:=[]*Card{}
	for rows.Next(){
	  card,err := ScanintoCards(rows)
	  if err!=nil{
		return  nil,err
	  }
	  cards=append(cards,card)
	}
    return cards,nil
}
func (s *PostgressStore) GetCard(id int) (*Card,error){
	rows,err:=s.db.Query(`SELECT * from cards 
	where id = $1`,id)
	if err!=nil{
	  return nil,err
	}
	for rows.Next(){
		return ScanintoCards(rows)
	}
	return nil,fmt.Errorf("Failed to get card")
}
func (s *PostgressStore) UpdateCard(c *Card) error{
	query:=`UPDATE cards
	SET state = $1
	where id=$2`
	_,err:=s.db.Query(
		query,
		c.Lane,
		c.ID,
	)
	if err!=nil{
		return err
	}
	return nil
}
func (s *PostgressStore) DeleteCard(id int) error{
	query:=`DELETE * from cards
	where id=$1`
	_,err:=s.db.Query(
		query,
		id,
	)
	if err!=nil{
		return err
	}
	return nil
}
func (s *PostgressStore) CreateIssue(i *Issue) error{
	query:=`insert into issues
	(content,created_by,organization,created_at)
	values($1,$2,$3,$4)`
	_,err:=s.db.Query(
		query,
		i.Problem,
		i.Created_by,
		i.Organization,
		i.Created_at,
	)
	if err!=nil{
		return err
	}
	return nil
}
func (s *PostgressStore)GetIssuses(org string) ([]*Issue,error){
	rows,err:=s.db.Query(`SELECT * from issues
	where organization = $1`,org)
	if err!=nil{
		return nil,err
	}
	issues:= []*Issue{}
	for rows.Next(){
		issue,err:=ScanintoIssue(rows)
		if err!=nil{
			return nil,err
		}
		issues=append(issues, issue)
	}
	return issues,nil
}
func (s *PostgressStore)GetIssuebyId(id int)(*Issue,[]*DispComment,error){
	row1,err:=s.db.Query(`SELECT * from issues
	where id = $1`,id)
	if err!=nil{
		return nil,nil,err
	}
	 var issue *Issue
	for row1.Next(){
		iss,err:=ScanintoIssue(row1)
		if err!=nil{
			return nil,nil,err
		}
		issue = iss
	}
	row2,err:=s.db.Query(`SELECT * from comments
	where post_id=$1`,id)
	if err!=nil{
		return nil,nil,err
	}
	commnets:= []*DispComment{} 
	for row2.Next(){
		comm,err:=ScanintoComment(row2)
		if err!=nil{
			return nil,nil,err
		}
		commnets=append(commnets,comm)
	}
	return issue,commnets,nil
}
func (s *PostgressStore) CreateComment(c *Comment) error{
	query:=`insert into comments
	(content,created_by,post_id,created_at)
	values($1,$2,$3,$4)`
	_,err:=s.db.Query(
		query,
		c.Content,
		c.CreatedBy,
		c.PostId,
		c.Created_at,
	)
	if err!=nil{
		return err
	}
	return nil
}

func Scanintouser(rows *sql.Rows) (*User, error){
	user:=new(User)
	err:=rows.Scan(&user.ID,&user.Name,&user.UserName,&user.Password,&user.Email,&user.Is_Premium,&user.Created_at)
	if err!=nil{
		return nil,err
	}
	return user,err
}
func ScanintoOrganization(rows *sql.Rows) (*Organization,error){
	org:=new(Organization)
	err:=rows.Scan(&org.ID,&org.Name,&org.CreateadAt,&org.Creator,pq.Array(&org.Team))
	if err!=nil{
		return nil,err
	}
	return org,err
}
func ScanintoCards(rows *sql.Rows) (*Card,error){
	card:=new(Card)
	err:=rows.Scan(&card.ID,&card.Name,&card.Task,&card.Organization,&card.Issued_to,&card.Issued_at,&card.Lane)
	if err!=nil{
		return nil,err
	}
    return card,nil
}
func ScanintoIssue(rows *sql.Rows) (*Issue,error){
	issue:=new(Issue)
	err:=rows.Scan(&issue.ID,&issue.Problem,&issue.Created_by,&issue.Organization,&issue.Created_at)
	if err!=nil{
		return nil,err
	}
	return issue,err
}
func ScanintoComment(rows *sql.Rows) (*DispComment,error){
	comment:=new(Comment)
	err:=rows.Scan(&comment.ID,&comment.Content,&comment.CreatedBy,&comment.PostId,&comment.Created_at)
	if err!=nil{
		return nil,err
	}
	dispcomment:=&DispComment{
		ID: comment.ID,
		Content: comment.Content,
		CreatedBy: comment.CreatedBy,
		Created_at: comment.Created_at,
	}
	return dispcomment,nil
}