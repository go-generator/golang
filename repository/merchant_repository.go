package repository

import (
	"log"

	. "../models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/reactivex/rxgo/observable"
)

type MerchantRepository struct {
	Server   string
	Database string
}

var merchanDb *mgo.Database

const (
	MERCHANT_COLLECTION = "merchant"
)

// Establish a connection to database
func (m *MerchantRepository) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	merchanDb = session.DB(m.Database)
}

// Find list of movies
func (m *MerchantRepository) FindAllObserver() observable.Observable {
	/*return observable.DefaultObservable.Map(func(i interface{}) interface{} {
		return m.FindAll
	}).Map(func(i interface{}) interface{} {
		return Merchant{}
	})*/
	res, err := m.FindAll()
	if err != nil {
		return observable.Just(err)
	} else {
		return observable.Just(res)
	}
}

// Find list of movies
func (m *MerchantRepository) FindAll() ([]Merchant, error) {
	var movies []Merchant
	err := merchanDb.C(MERCHANT_COLLECTION).Find(bson.M{}).All(&movies)
	return movies, err
}

// Find a movie by its id
func (m *MerchantRepository) FindById(id string) (Merchant, error) {
	var movie Merchant
	err := merchanDb.C(MERCHANT_COLLECTION).FindId(bson.ObjectIdHex(id)).One(&movie)
	return movie, err
}

// Insert a movie into database
func (m *MerchantRepository) Insert(movie Merchant) error {
	err := merchanDb.C(MERCHANT_COLLECTION).Insert(&movie)
	return err
}

// Delete an existing movie
func (m *MerchantRepository) Delete(movie Merchant) error {
	err := merchanDb.C(MERCHANT_COLLECTION).Remove(&movie)
	return err
}

// Update an existing movie
func (m *MerchantRepository) Update(movie Merchant) error {
	err := merchanDb.C(MERCHANT_COLLECTION).UpdateId(movie.ID, &movie)
	return err
}
