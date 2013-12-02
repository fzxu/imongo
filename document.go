package main

import (
	"log"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Document struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Path        string        `bson:"path"`
	CreatedAt   time.Time     `bson:"created_at"`
	Binary      []byte        `bson:"binary"`
	ContentType string        `bson:"content_type,omitempty"`
}

func (d Document) Collection(s *mgo.Session) *mgo.Collection {
	return s.DB(Configuration.DBName).C(Configuration.Collection)
}

func (d *Document) Save(s *mgo.Session) error {
	coll := d.Collection(s)

	if !bson.IsObjectIdHex(d.Id.Hex()) {
		d.Id = bson.NewObjectId()
		d.CreatedAt = time.Now()
	}

	_, err := coll.Upsert(bson.M{"_id": d.Id}, d)
	if err != nil {
		log.Panicln(err)
		return err
	}
	return nil
}

func (d Document) Find(s *mgo.Session, name string, path string) (*Document, error) {
	result := new(Document)
	coll := d.Collection(s)

	query := coll.Find(bson.M{"name": name, "path": path})
	err := query.One(result)

	return result, err
}
