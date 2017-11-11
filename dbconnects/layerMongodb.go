package dbconnect

import (
	mgo "gopkg.in/mgo.v2"
)

/*Get Collect and copy new database session*/
func WithCollection(collection string, s func(*mgo.Collection) error) error {
	session := GetMongoSessionCopy()
	defer session.Close()
	c := session.DB(LoadConfigMongoDB.Name).C(collection)
	return s(c)
}

/*Function search for multiple Collection*/
func SearchData(collectionName string, q interface{}, skip int, limit int) (searchResults []interface{}, searchErr error) {
	query := func(c *mgo.Collection) error {
		fn := c.Find(q).Skip(skip).Limit(limit).All(&searchResults)
		if limit < 0 {
			fn = c.Find(q).Skip(skip).All(&searchResults)
		}
		return fn
	}
	search := func() error {
		return WithCollection(collectionName, query)
	}
	err := search()
	if err != nil {
		searchErr = err
	}
	return
}

/*Insert To Collection*/
func InserToCollection(collectionName string, dataInsert interface{}) error {
	err := WithCollection(collectionName, func(c *mgo.Collection) error {
		return c.Insert(&dataInsert)
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdateOneInCollection(collectionName string, condition interface{}, dataSet interface{}) error {

	ErrQuery := WithCollection(collectionName, func(c *mgo.Collection) error {
		return c.Update(condition, dataSet)
	})
	if ErrQuery != nil {
		return ErrQuery
	}
	return nil
}

/*Count collection with query*/
func CountRowsInCollection(collectionName string, condition interface{}) int {
	var Rows int
	ErrQuery := WithCollection(collectionName, func(c *mgo.Collection) error {
		rows, errx := c.Find(condition).Count()
		if errx != nil {
			return errx
		}
		Rows = rows
		return nil
	})
	if ErrQuery != nil {
		return Rows
	}
	return Rows
}
