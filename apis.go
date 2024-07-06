package main

import (

	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"log"

	"net/http"
	"os"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	listenAddr string
	store      storage
}
type ApiError struct {
	Error string
}
type apiFunc func(http.ResponseWriter, *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func makemuxhandlefunc(f apiFunc) http.HandlerFunc {
	
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
func NewApiServer(listenAddr string, store storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}
func (s *APIServer) Run() {
	router := mux.NewRouter()
	c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"}, // Allow specific origin
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
    })
    router.HandleFunc("/register",makemuxhandlefunc(s.handleregisteruser))
	router.HandleFunc("/login",makemuxhandlefunc(s.handleloginUser))
	router.HandleFunc("/createorganization",makemuxhandlefunc(s.handlecreateorg))
	router.HandleFunc("/invite",makemuxhandlefunc(s.handleinvite))
	router.HandleFunc("/remove",makemuxhandlefunc(s.handleremovemember))
	router.HandleFunc("/exit",makemuxhandlefunc(s.handlexitorganization))
	router.HandleFunc("/getallorg",makemuxhandlefunc(s.handlegetallorganizations))
	router.HandleFunc("/issue_card",makemuxhandlefunc(s.handlecreatecards))
	router.HandleFunc("/getcard/{id}",makemuxhandlefunc(s.handlegetcardbyId))
	router.HandleFunc("/getallcard/{org}",makemuxhandlefunc(s.handlegetallcards))
	router.HandleFunc("/forwardcard/{id}",makemuxhandlefunc(s.handleforwardcard))
	router.HandleFunc("/moveback/{id}",makemuxhandlefunc(s.handlemovebackcard))
	router.HandleFunc("/deletecard",makemuxhandlefunc(s.handledeletecard))
	router.HandleFunc("/createissue",makemuxhandlefunc(s.handlecreateissue))
	router.HandleFunc("/issue/{org}",makemuxhandlefunc(s.handlegetallissues))
	router.HandleFunc("/getissue/{id}",makemuxhandlefunc(s.handlegetissuebyid))
	router.HandleFunc("/comment/{id}",makemuxhandlefunc(s.handlecomment))
    // Use the CORS middleware with the router
    handler := c.Handler(router)
	log.Println("JSON Api Running on port:", s.listenAddr)
	http.ListenAndServe(s.listenAddr,handler)
}
func (s *APIServer) handleregisteruser( w http.ResponseWriter, r *http.Request) error{
   if r.Method!="POST"&&r.Method!="GET"{
	return fmt.Errorf("Method not allowed %s",r.Method)
   }
   if r.Method=="GET"{
	err,req:=validateUserCreateToken(w,r)
	if err!=nil{
		return fmt.Errorf("token can't be read ")
	}
	user,err:= NewUser(req.Name,req.UserName,req.Password,req.Email)
	if err!=nil{
		return err
	}
	err=s.store.CreateUser(user)
	if err!=nil{
		return err
	}
	return WriteJSON(w,http.StatusOK,user)
   }
  if r.Method=="POST" {
   req:= new(CreateuserRequest)
   if err := json.NewDecoder(r.Body).Decode(req); err != nil {
	   return err
   }
   _,err:= s.store.GetUserbyUserName(req.UserName)
   if err==nil{
	return fmt.Errorf("UserName already exists")
   }

   token,err:=CreateUsercreateJWT(req)	
   msg:= createuserMessage{
   Message: "User verification mail sent",
   }
   SendMail(token,req.Email)
   if err!=nil{
	return fmt.Errorf("Error occured while generating the token")
}
   return WriteJSON(w,http.StatusOK,msg)
}
return nil
}
func (s *APIServer) handleloginUser(w http.ResponseWriter, r *http.Request) error{
	if (r.Method!="POST"){
		return fmt.Errorf("Method not allowed :%s",r.Method)
	}
	req:=new(UserLoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	user,err:=s.store.GetUserbyUserName(req.UserName)
	if err!=nil{
		return fmt.Errorf("failed to Fetch the User")
	}
    if user.ValidatePassword(req.Password)==false{
		return fmt.Errorf("Invalid Password or UserName entered!!")
	}
	token,err:=CreateJWT(req)
	msg:=&LoginMessage{
		Token: token,
		Message: "User Verified Successfully",
	}
	return WriteJSON(w,http.StatusOK,msg)
}
func (s *APIServer) handlecreateorg(w http.ResponseWriter,r *http.Request) error{
	if r.Method!="POST"{
		return fmt.Errorf("Method not allowed :%s",r.Method)
	}
	err,user:=s.ValidateJWT(w,r)
	if err!=nil{
	   return fmt.Errorf("%s",err)
	}
	req:=new(CreateOrganizationRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	org,err:=NewOrganization(req.Name,user.UserName,)
	if err!=nil{
		return fmt.Errorf("%s",err)
	}
	err=s.store.CreateOrganization(org)
	if err!=nil{
		return fmt.Errorf("%s",err)
	}
	return WriteJSON(w,http.StatusOK,org)
}
func (s *APIServer) handleinvite(w http.ResponseWriter,r *http.Request) error {
	if r.Method!="POST"&&r.Method!="GET"{
		return fmt.Errorf("Method not allowed :%s",r.Method)
	}
	if r.Method=="GET"{
        err,org,user:=s.ValidateInviteJWT(w,r)
		if err!=nil{
			return err
		}
		err = s.store.AddMenber(org,user)
		if err!=nil{
			return err
		}
		return WriteJSON(w,http.StatusOK,org)
	}
	if r.Method=="POST"{
		req:=new(InviteRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			return err
		}
		err,user:=s.ValidateJWT(w,r)
		if err!=nil{
			return fmt.Errorf("%s",err)
		}	
        for i:=0;i<len(req.Team);i++{
			us,err:=s.store.GetUserbyUserName(req.Team[i])
			if err!=nil{
				return err
			}
			token,err:=CreateInviteToken(req.Organization,req.Team[i])
			if err!=nil{
				return err
			}
			err=SendInviteMail(token,user.Name,us.Email)
			if err!=nil{
				return err;
			}
		}
		return WriteJSON(w,http.StatusOK,"invitation sent to all the members")
	}
 return nil
}
func (s *APIServer) handleremovemember(w http.ResponseWriter,r *http.Request) error {
	if r.Method!="POST"{
		return fmt.Errorf("Method not allowed:%s",r.Method)
	}
	req:=new(DeleteMemberRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	org,err:= s.store.GetOrganizationByName(req.Organization)
	if err!=nil{
		return err
	}
	err,user:=s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	fmt.Println(user.Name)
    if org.Creator!=user.UserName{
		return fmt.Errorf("Action Prohibitted!!")
	}
	err =s.store.RemoveMember(org,req.Member)
	if err!=nil{
		return err
	}
   return WriteJSON(w,http.StatusOK,org)
}
func (s *APIServer) handlexitorganization(w http.ResponseWriter,r *http.Request)error{
	if r.Method!="POST"{
		return fmt.Errorf("Method not allowed:%s",r.Method)
	}
	req:=new(CreateOrganizationRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	err,user:= s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	org,err:=s.store.GetOrganizationByName(req.Name)
	if err!=nil{
		return err
	}
	err=s.store.RemoveMember(org,user.UserName)
	return WriteJSON(w,http.StatusOK,"you are not a member now")
}
func (s *APIServer) handlegetallorganizations(w http.ResponseWriter,r *http.Request) error{
	if r.Method!="GET"{
		return fmt.Errorf("Method not allowed:%s",r.Method)
	}
	err,user:=s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	org,err:=s.store.GetOrganizations(user.UserName)
	if err!=nil{
		return err
	}
	return WriteJSON(w,http.StatusOK,org)
}
func (s *APIServer) handlecreatecards(w http.ResponseWriter,r *http.Request) error{
	if r.Method!="POST"{
		return fmt.Errorf("Method not allowed:%s",r.Method)
	}
	req:=new(CreateCardRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	err,user:=s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	org,err:=s.store.GetOrganizationByName(req.Organization)
	if user.UserName!=org.Creator{
		return fmt.Errorf("Action Unauthorized !!")
	}
	cd:=NewCard(req.Name,req.Task,req.Organization,req.Issued_to)
	err= s.store.CreateCard(cd)
	if err !=nil{
		return err
	}
	return WriteJSON(w,http.StatusOK,cd)
}
func (s *APIServer) handlegetallcards(w http.ResponseWriter,r *http.Request) error{
	if r.Method!="GET"{
		return fmt.Errorf("Method not allowed:%s",r.Method)
	}
	og:=mux.Vars(r)["org"]
	cds,err:=s.store.GetCards(og)
	if err!=nil{
		return err
	}
	resp:=GetAllCardResponse{
		Cards: cds,
	}
	return WriteJSON(w,http.StatusOK,resp)
}
func (s *APIServer) handlegetcardbyId(w http.ResponseWriter,r *http.Request) error{
	if r.Method!="GET"{
		return fmt.Errorf("Method not allowed:%s",r.Method)
	}
	idstr:=mux.Vars(r)["id"]
	id,err:=strconv.Atoi(idstr)
	if err!=nil{
		return err;
	}
	card,er:=s.store.GetCard(id)
	if er!=nil{
		return er
	}
	return WriteJSON(w,http.StatusOK,card)
}
func (s *APIServer) handleforwardcard(w http.ResponseWriter,r *http.Request) error{
	if r.Method!="PUT"{
		return fmt.Errorf("Method not allowed:%s",r.Method)
	}
	err,user:=s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	idstr:=mux.Vars(r)["id"]
	id,err:=strconv.Atoi(idstr)
	if err!=nil{
		return err;
	}
	card,er:=s.store.GetCard(id)
	if er!=nil{
		return er
	}
	org,err:=s.store.GetOrganizationByName(card.Organization)
	if err!=nil{
		return err
	}
	if user.UserName!=card.Issued_to && user.UserName!=org.Creator {
        return fmt.Errorf("Action Unauthorized")
	}
	l:=GetLane(card.Lane)
	if(l==4){
		return fmt.Errorf("Card allready at final state")
	}
	card.Lane= (l+1).String()
	err = s.store.UpdateCard(card)
	if err!=nil{
		return nil
	}
	return WriteJSON(w,http.StatusOK,card)
}
func (s *APIServer) handlemovebackcard(w http.ResponseWriter, r *http.Request) error{
	if r.Method!="PUT"{
		return fmt.Errorf("Method not allowed:%s",r.Method)
	}
	err,user:=s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	idstr:=mux.Vars(r)["id"]
	id,err:=strconv.Atoi(idstr)
	if err!=nil{
		return err;
	}
	card,er:=s.store.GetCard(id)
	if er!=nil{
		return er
	}
	org,err:=s.store.GetOrganizationByName(card.Organization)
	if err!=nil{
		return err
	}
	if user.UserName!=card.Issued_to && user.UserName!=org.Creator {
        return fmt.Errorf("Action Unauthorized")
	}
	l:=GetLane(card.Lane)
	if(l==0){
		return fmt.Errorf("Card allready at initial state")
	}
	card.Lane= (l-1).String()
	err = s.store.UpdateCard(card)
	if err!=nil{
		return nil
	}
	return WriteJSON(w,http.StatusOK,card)	
}
func (s *APIServer) handledeletecard(w http.ResponseWriter,r *http.Request)error{
	if r.Method!="DELETE"{
		return fmt.Errorf("Method not allowed:%s",r.Method)	
	}
	err,user:=s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	idstr:=mux.Vars(r)["id"]
	id,err:=strconv.Atoi(idstr)
	if err!=nil{
		return err;
	}
	card,er:=s.store.GetCard(id)
	if er!=nil{
		return er
	}
	org,err:=s.store.GetOrganizationByName(card.Organization)
	if err!=nil{
		return err
	}
	if user.UserName!=org.Creator {
       return fmt.Errorf("Action Unauthorized")
	}
	err =s.store.DeleteCard(id)
	if err!=nil{
		return err
	}
	return WriteJSON(w,http.StatusOK,"Card Deleted SuccesFully")
}
func (s *APIServer) handlecreateissue(w http.ResponseWriter,r *http.Request) error{
    if r.Method!="POST"{
		return fmt.Errorf("Method not allowed: %s",r.Method)
	}
    err,user:=s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	req:=new(CreateIssueRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	org,err:=s.store.GetOrganizationByName(req.Organization)
	if err!=nil{
		return err
	}
	if user.UserName!=org.Creator && !contains(org.Team,user.UserName){
       return fmt.Errorf("Action UnAuthorized")
	}
	issue,err:=NewIssue(user.UserName,req.Problem,req.Organization)
	if err!=nil {
		return err
	}
	err = s.store.CreateIssue(issue)
	if err!=nil{
		return err
	}
    return  WriteJSON(w,http.StatusOK,issue)
}
func (s *APIServer) handlegetallissues(w http.ResponseWriter, r *http.Request) error{
	if r.Method!="GET"{
		return fmt.Errorf("Method not allowed: %s",r.Method)
	}
	err,user:=s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	orgstr:=mux.Vars(r)["org"]

	org,err:=s.store.GetOrganizationByName(orgstr)
	if err!=nil{
		return err
	}
	if user.UserName!=org.Creator && !contains(org.Team,user.UserName){
		return fmt.Errorf("Action UnAuthorized")
	 }
	iss,err:=s.store.GetIssuses(orgstr)
	if err!=nil{
		return err
	}
	isssue:=getallissue{
		Issues: iss,
	}
	return WriteJSON(w,http.StatusOK,isssue)
}
func (s *APIServer) handlegetissuebyid(w http.ResponseWriter, r *http.Request) error{
	if r.Method!="GET"{
		return fmt.Errorf("Method not allowed: %s",r.Method)
	}
	err,_:=s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	idstr:=mux.Vars(r)["id"]
	id,err:=strconv.Atoi(idstr)
	if err!=nil{
		return err;
	}
	is,com,err:=s.store.GetIssuebyId(id)
	if err!=nil{
		return err
	}
	resp:=getissue{
     Issue: is,
	 Comments: com,
	}
	return WriteJSON(w,http.StatusOK,resp)
}
func (s *APIServer) handlecomment(w http.ResponseWriter, r *http.Request) error{
	if r.Method!="POST"{
		return fmt.Errorf("Metod not allowed: %s",r.Method)
	}
	err,user:=s.ValidateJWT(w,r)
	if err!=nil{
		return err
	}
	idstr:=mux.Vars(r)["id"]
	id,err:=strconv.Atoi(idstr)
	if err!=nil{
		return err;
	}
	req:=new(CommentRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	comt,err:=NewComment(user.UserName,req.Content,id)
	if err!=nil{
		return err
	}
	err=s.store.CreateComment(comt)
	if err!=nil{
		return err
	}
	return WriteJSON(w,http.StatusOK,comt)
}

func CreateUsercreateJWT(req *CreateuserRequest,) (string,error){
	claims:=jwt.MapClaims{
		"expiresAt":  jwt.NewNumericDate(time.Now().Local().Add(time.Minute * 5)),
		"name":  req.Name,
		"userName": req.UserName,
		"email": req.Email,
		"password": req.Password,

	}
	secret := os.Getenv("SECRET")
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString([]byte(secret))

}
func validateUserCreateToken(w http.ResponseWriter, r *http.Request,) (error,*CreateuserRequest) {
	if r.URL.Query().Get("token") == "" {
		fmt.Fprintf(w, "can not find token in header")
		return fmt.Errorf( "can not find token in header %s",w),nil
	}

	token,_  := jwt.Parse(r.URL.Query().Get("token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing:",)
		}
		return os.Getenv("SECRET"), nil
	})


	if token == nil {
		fmt.Fprintf(w, "invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("couldn't parse claims"),nil
	}

	exp := claims["expiresAt"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return fmt.Errorf("token expired"),nil
	}

	req:= &CreateuserRequest{
		Name: claims["name"].(string),
		UserName: claims["userName"].(string),
		Password: claims["password"].(string),
		Email: claims["email"].(string),
	}
	return nil,req
}
func CreateJWT(req *UserLoginRequest) (string,error){
	claims:=jwt.MapClaims{
		"expiresAt":  jwt.NewNumericDate(time.Now().Local().Add(time.Minute * 5)),
		"userName": req.UserName,
	}
	secret := os.Getenv("SECRET")
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString([]byte(secret))

}
func CreateInviteToken(org,user string) (string,error){
	claims:=jwt.MapClaims{
		"expiresAt":  jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 72)),
		"userName": user,
		"organization":org,
	}
	secret := os.Getenv("SECRET")
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString([]byte(secret))
}
func (s *APIServer)ValidateJWT(w http.ResponseWriter,r *http.Request) (error,*User){
    if r.Header["Token"]==nil{
        return fmt.Errorf("can't find token in the headers "),nil
	}
	token,_  := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing:",)
		}
		return os.Getenv("SECRET"), nil
	})
	if token == nil {
		 return fmt.Errorf( "invalid token"),nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Fprintf(w, "couldn't parse claims")
		return errors.New("Token error"),nil
	}
	exp := claims["expiresAt"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return errors.New("Token expired"),nil
	}
    username:=claims["userName"].(string)
	user,err:=s.store.GetUserbyUserName(username)
	if err!=nil{
		return fmt.Errorf("%s",err),nil
	}
	return nil,user
}
func (s *APIServer)ValidateInviteJWT(w http.ResponseWriter,r *http.Request) (error,*Organization,string){
	if r.URL.Query().Get("token") == "" {
		fmt.Fprintf(w, "can not find token in header")
		return fmt.Errorf( "can not find token in header %s",w),nil,""
	}

	token,_  := jwt.Parse(r.URL.Query().Get("token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing:",)
		}
		return os.Getenv("SECRET"), nil
	})
	if token == nil {
		 return fmt.Errorf( "invalid token"),nil,""
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Fprintf(w, "couldn't parse claims")
		return errors.New("Token error"),nil,""
	}
	exp := claims["expiresAt"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return errors.New("Token expired"),nil,""
	}
    username:=claims["userName"].(string)
	org:=claims["organization"].(string)
	fmt.Printf("userName:%s,Organization:%s",username,org)
	organization,err:=s.store.GetOrganizationByName(org)
	if err!=nil{
		return fmt.Errorf("%s",err),nil,""
	}
   return nil,organization,username
}





func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}