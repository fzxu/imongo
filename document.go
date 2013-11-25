package main

import (
	"log"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Document struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Path        string        `bson:"path"`
	ContentType string        `bson:"content_type"`
	Binary      bson.Binary   `bson:"binary"`
}

func (d *Document) Collcetion(s *mgo.Session) *mgo.Collection {
	return s.DB(Configuration.DBName).C(Configuration.Collection)
}

func (d *Document) Save(s *mgo.Session) error {
	coll := d.Collcetion(s)

	if !bson.IsObjectIdHex(d.Id.Hex()) {
		d.Id = bson.NewObjectId()
	}

	_, err := coll.Upsert(bson.M{"_id": d.Id}, d)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}
