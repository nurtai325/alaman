package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nurtai325/alaman/internal/auth"
	"github.com/nurtai325/alaman/internal/config"
	"github.com/nurtai325/alaman/internal/db"
	"github.com/nurtai325/alaman/internal/db/repository"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if 0 == 0 {
		password, _ := bcrypt.GenerateFromPassword([]byte("adminadmin"), bcrypt.DefaultCost)
		fmt.Println(string(password))
	}
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	newDB, err := db.New(conf)
	if err != nil {
		panic(err)
	}
	queries := repository.New(newDB)
	name := "Tolymbek Nurtai"
	// phone := "+77052505839"
	password := "admin_test"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if !validPhone(phone) {
	// 	panic("invalid phone")
	// }
	if !validPassword(password) {
		panic("invalid password")
	}
	phoneNumbers := []string{
		"+77052505839", "+77051234501", "+77051234502", "+77051234503", "+77051234504", "+77051234505",
		"+77051234506", "+77051234507", "+77051234508", "+77051234509", "+77051234510",
		"+77051234511", "+77051234512", "+77051234513", "+77051234514", "+77051234515",
		"+77051234516", "+77051234517", "+77051234518", "+77051234519", "+77051234520",
		"+77051234521", "+77051234522", "+77051234523", "+77051234524", "+77051234525",
		"+77051234526", "+77051234527", "+77051234528", "+77051234529", "+77051234530",
		"+77051234531", "+77051234532", "+77051234533", "+77051234534", "+77051234535",
		"+77051234536", "+77051234537", "+77051234538", "+77051234539", "+77051234540",
		"+77051234541", "+77051234542", "+77051234543", "+77051234544", "+77051234545",
		"+77051234546", "+77051234547", "+77051234548", "+77051234549", "+77051234550",
	}
	for i, phone := range phoneNumbers {
		if i != 0 {
			break
		}
		_, err := queries.InsertUser(context.Background(), repository.InsertUserParams{
			Name:     name,
			Phone:    phone,
			Password: string(hashed),
			Role:     string(auth.AdminRole),
		})
		if err != nil {
			panic(err)
		}
	}
}

func validPhone(phone string) bool {
	if phone == "" {
		return false
	} else if len(phone) != 12 {
		return false
	} else if rune(phone[0]) != '+' {
		return false
	}
	phone = phone[1:]
	for _, r := range phone {
		if r <= 47 || r >= 58 {
			return false
		}
	}
	return true
}

func validPassword(password string) bool {
	if password == "" || len(password) > 72 {
		return false
	} else if len(password) < 8 {
		return false
	}
	return true
}
