package mongodb

import (
	"GIG/app/constants/database"
	"GIG/app/databases/mongodb"
	"GIG/app/repositories/interfaces"

	"github.com/lsflk/gig-sdk/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type NormalizedNameRepository struct {
	interfaces.NormalizedNameRepositoryInterface
}

func (n NormalizedNameRepository) newNormalizedNameCollection() *mongodb.Collection {
	return mongodb.NewCollectionSession(database.NormalizedNameCollection)
}

// AddNormalizedName insert a new NormalizedName into database and returns
// last inserted normalized_name on success.
func (n NormalizedNameRepository) AddNormalizedName(m models.NormalizedName) (normalizedName models.NormalizedName, err error) {
	c := n.newNormalizedNameCollection()
	defer c.Close()
	m = m.NewNormalizedName()
	err = c.Collection.Insert(m)
	if mgo.IsDup(err) {
		err = nil
	}
	return m, err
}

// GetNormalizedNames Get all NormalizedNames from database and returns
// list of NormalizedName on success
func (n NormalizedNameRepository) GetNormalizedNames(searchString string, limit int) ([]models.NormalizedName, error) {
	var (
		normalizedNames []models.NormalizedName
		err             error
		resultQuery     *mgo.Query
	)

	query := bson.M{}
	c := n.newNormalizedNameCollection()
	defer c.Close()

	if searchString != "" {
		query = bson.M{
			"$text": bson.M{"$search": searchString},
		}
	}

	resultQuery = c.Collection.Find(query).Select(bson.M{
		"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score")

	err = resultQuery.Limit(limit).All(&normalizedNames)

	return normalizedNames, err
}

// GetNormalizedName Get a NormalizedName from database and returns
// a NormalizedName on success
func (n NormalizedNameRepository) GetNormalizedName(id string) (models.NormalizedName, error) {
	var (
		normalizedName models.NormalizedName
		err            error
	)

	c := n.newNormalizedNameCollection()
	defer c.Close()

	err = c.Collection.Find(bson.M{"_id": id}).One(&normalizedName)
	return normalizedName, err
}

/*
GetNormalizedNameBy Get an Entity from database and returns
a models.Entity on success
*/
func (n NormalizedNameRepository) GetNormalizedNameBy(attribute string, value string) (models.NormalizedName, error) {
	var (
		normalizedName models.NormalizedName
		err            error
	)

	c := n.newNormalizedNameCollection()
	defer c.Close()

	err = c.Collection.Find(bson.M{attribute: value}).One(&normalizedName)
	return normalizedName, err
}
