package dbconnect

import (
	mgo "gopkg.in/mgo.v2"
)

/*WithCollection get Collect and copy new database session*/
func WithCollection(collection string, s func(*mgo.Collection) error) error {
	session := GetMongoSessionCopy()
	defer session.Close()
	c := session.DB(LoadConfigMongoDB.Name).C(collection)
	return s(c)
}

/*SearchDataInCollection search for multiple Collection*/
func SearchDataInCollection(collectionName string, q interface{}, skip, limit int) (searchResults []interface{}, searchErr error) {
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

/*GetOneDataInCollection get one data*/
func GetOneDataInCollection(collectionName string, q interface{}) (searchResults interface{}, searchErr error) {

	query := func(c *mgo.Collection) error {
		fn := c.Find(q).One(&searchResults)
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

/*InserToCollection Insert data to Collection*/
func InserToCollection(collectionName string, dataInsert interface{}) error {
	err := WithCollection(collectionName, func(c *mgo.Collection) error {
		return c.Insert(&dataInsert)
	})
	if err != nil {
		return err
	}
	return nil
}

/*UpdateOneInCollection update data with condition*/
func UpdateOneInCollection(collectionName string, condition interface{}, dataSet interface{}) error {

	ErrQuery := WithCollection(collectionName, func(c *mgo.Collection) error {
		return c.Update(condition, dataSet)
	})
	if ErrQuery != nil {
		return ErrQuery
	}
	return nil
}

/*CountRowsInCollection count collection with query*/
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
