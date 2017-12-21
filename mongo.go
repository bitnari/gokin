package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

var (
	ErrNoAccount        = errors.New("no account")
	ErrAccountExist     = errors.New("account already exists")
	ErrIncorrectPassword = errors.New("incorrect password")
)

type Account struct {
	Id          string
	Hash        string
	Credit		int
	Gold        int
	Score       map[string]int
}

type MongoConnection struct {
	session     *mgo.Session

	account     *mgo.Collection
	score       *mgo.Collection
}

func (m *MongoConnection) Init(host, db string) (err error) {
	m.session, err = mgo.Dial(host)

	if err != nil {
		return
	}

	m.account = m.session.DB(db).C("account")
	m.score = m.session.DB(db).C("score")

	return nil
}

func (m *MongoConnection) AddAccount(id, password string, defaultCredit int, defaultGold int) (err error) {
	count, err := m.account.Find(bson.M{"id": id}).Count()
	if err != nil {
		return
	}
	if count > 0 {
		return ErrAccountExist
	}

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	m.account.Insert(&Account {id, string(hash), defaultCredit, defaultGold, map[string]int{}})
	return
}

func (m *MongoConnection) VerifyAccount(id, password string) (err error) {
	var account Account
	err = m.account.Find(bson.M{"id": id}).One(&account)
	if err != nil {
		if err == mgo.ErrNotFound {
			return ErrNoAccount
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrIncorrectPassword
		}
	}

	return
}

func (m *MongoConnection) SubtractGold(id string, gold int) (err error) {
	var account Account
	account, err = m.GetAccount(id)
	if err != nil {
		return
	}

	err = m.account.Update(bson.M{"id": account.Id}, bson.M{"$set": bson.M{"gold": account.Gold - gold}})
	return
}

func (m *MongoConnection) SetScore(account Account, gameId string, score int) (err error) {
	account.Score[gameId] = score

	err = m.account.Update(bson.M{"id": account.Id}, bson.M{"$set": bson.M{"score": account.Score}})
	return
}

func (m *MongoConnection) SetCredit(id string, gold, credit int) (err error) {
	err = m.account.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"gold": gold, "credit": credit}})
	return
}

func (m *MongoConnection) GetAccount(id string) (account Account, err error) {
	err = m.account.Find(bson.M{"id": id}).One(&account)
	if err == mgo.ErrNotFound {
		return account, ErrNoAccount
	}

	return
}
