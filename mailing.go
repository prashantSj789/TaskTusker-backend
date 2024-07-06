package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendMail(token string,email string) {

// Choose auth method and set it up

auth := smtp.PlainAuth("", "singhprashant79072@gmail.com", "woyd mwmr qorq cxif", "smtp.gmail.com")

// Here we do it all: connect to our server, set up a message and send it

to := []string{email}

msg := []byte("To: "+email+ "\r\n" +

"Subject: Verification  mail from TaskTusker\r\n" +

"\r\n" +

"http://192.168.1.7:8080/register?token="+token +"\r\n")

err := smtp.SendMail("smtp.gmail.com:587", auth, "singhprashant79072@gmail.com", to, msg)
fmt.Printf("mail sent")
if err != nil {

log.Fatal(err)

}

}

  func SendInviteMail(token, from, email string) error{

	// Choose auth method and set it up
	
	auth := smtp.PlainAuth("", "singhprashant79072@gmail.com", "woyd mwmr qorq cxif", "smtp.gmail.com")
	
	// Here we do it all: connect to our server, set up a message and send it
	
	to := []string{email}
	
	msg := []byte("To: "+email+ "\r\n" +
	
	"Subject: Invitation to join an organization TaskTusker"+from+"\r\n" +
	
	"\r\n" +
	
	"http://192.168.1.7:8080/invite?token="+token +"\r\n")
	
	err := smtp.SendMail("smtp.gmail.com:587", auth, "singhprashant79072@gmail.com", to, msg)

	if err != nil {
	
	return fmt.Errorf("%s",err)
	
	}
	fmt.Printf("inivite mail sent")
	return nil
	}
	



