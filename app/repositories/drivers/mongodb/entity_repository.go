package mongodb

import (
	"GIG/app/constants/database"
	"GIG/app/databases/mongodb"
	"GIG/app/repositories/constants"
	"GIG/app/repositories/interfaces"
	"log"
	"time"

	"github.com/lsflk/gig-sdk/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type EntityRepository struct {
	interfaces.EntityRepositoryInterface
}

func (e EntityRepository) newEntityCollection() *mongodb.Collection {
	return mongodb.NewCollectionSession(database.EntityCollection)
}

/*
AddEntity insert a new Entity into database and returns
last inserted entity on success.
*/
func (e EntityRepository) AddEntity(entity models.Entity) (models.Entity, error) {
	c := e.newEntityCollection()
	defer c.Close()
	return entity, c.Collection.Insert(entity)
}

func (e EntityRepository) GetEntityByPreviousTitle(title string, date time.Time) (models.Entity, error) {
	var (
		entity models.Entity
		err    error
	)

	query := bson.M{
		"attributes.titles.values.value_string": title,
		"attributes.titles.values.date":         bson.M{"$lt": date.Add(time.Duration(1) * time.Second)},
	}

	c := e.newEntityCollection()
	defer c.Close()

	err = c.Collection.Find(query).Sort("-attributes.titles.values.date", "-attributes.titles.values.updated_at").One(&entity)
	return entity, err
}

/*
GetRelatedEntities Get all Entities where a given title is linked from
list of models.Entity on success
*/
func (e EntityRepository) GetRelatedEntities(entity models.Entity, limit int, offset int) ([]models.Entity, error) {
	var (
		entities []models.Entity
		err      error
	)

	query := bson.M{}
	c := e.newEntityCollection()
	defer c.Close()

	entityTitle := entity.GetTitle()
	if entityTitle != "" {
		query = bson.M{"links.title": bson.M{"$in": append(entity.GetLinkTitles(), entity.GetTitle())}}
	}
	log.Println(query)
	err = c.Collection.Find(query).Sort(constants.UpdatedAtDecending).Skip(offset).Limit(limit).All(&entities)

	for _, item := range entities {
		log.Println(item.GetTitle())
	}
	return entities, err
}

/*
GetEntities Get all Entities from database and returns
list of models.Entity on success
*/
func (e EntityRepository) GetEntities(search string, categories []string, limit int, offset int) ([]models.Entity, error) {
	var (
		entities    []models.Entity
		err         error
		resultQuery *mgo.Query
	)

	query := bson.M{}
	c := e.newEntityCollection()
	defer c.Close()

	if search != "" {
		query = bson.M{
			"$text": bson.M{"$search": search},
			//"attributes": bson.M{"$exists": true, "$not": bson.M{"$size": 0}},
		}
	}

	if categories != nil && len(categories) != 0 {
		query["categories"] = bson.M{"$all": categories}
	}

	// sort by search score for text indexed search, otherwise sort by latest first in category
	if search == "" {
		resultQuery = c.Collection.Find(query).Sort("-source_date")
	} else {
		resultQuery = c.Collection.Find(query).Select(bson.M{
			"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score")
	}

	err = resultQuery.Skip(offset).Limit(limit).All(&entities)

	return entities, err
}

/*
GetEntity Get an Entity from database and returns
a models. Entity on success
*/
func (e EntityRepository) GetEntity(id string) (models.Entity, error) {
	var (
		entity models.Entity
		err    error
	)

	c := e.newEntityCollection()
	defer c.Close()

	err = c.Collection.Find(bson.M{"_id": id}).One(&entity)
	return entity, err
}

/*
GetEntityBy Get an Entity from database and returns
Search entity by metadata e.g. title, source, created at
Note that only a single entity is returned
Only use this function when from a guaranteed unique attribute and value. for e.g. title is unique
a models.Entity on success
*/
func (e EntityRepository) GetEntityBy(attribute string, value string) (models.Entity, error) {
	var (
		entity models.Entity
		err    error
	)

	c := e.newEntityCollection()
	defer c.Close()
	err = c.Collection.Find(bson.M{attribute: value}).Sort(constants.UpdatedAtDecending).One(&entity)
	return entity, err
}

/*
UpdateEntity update a Entity into database and returns
last nil on success.
*/
func (e EntityRepository) UpdateEntity(entity models.Entity) error {
	c := e.newEntityCollection()
	defer c.Close()

	err := c.Collection.Update(bson.M{
		"_id": entity.GetId(),
	}, bson.M{
		"$set": entity,
	})
	return err
}

/*
DeleteEntity Delete Entity from database and returns
last nil on success.
*/
func (e EntityRepository) DeleteEntity(entity models.Entity) error {
	c := e.newEntityCollection()
	defer c.Close()

	err := c.Collection.Remove(bson.M{"_id": entity.GetId()})
	return err
}

/*
GetStats Get entity states from the DB
*/
func (e EntityRepository) GetStats() (models.EntityStats, error) {
	var (
		entityStats models.EntityStats
		err         error
	)
	entityStats.CreatedAt = time.Now()

	c := e.newEntityCollection()
	defer c.Close()

	// Get total number of entities
	entityStats.EntityCount, err = c.Collection.Find(nil).Count()
	var linkCount []map[string]interface{}

	//Get category wise count
	categoryCountPipeline := []bson.M{
		{constants.UnwindAttribute: constants.CategoryAttribute},
		{constants.GroupAttribute: bson.M{
			"_id":            constants.CategoryAttribute,
			"category_count": bson.M{"$sum": 1}}},
		{constants.SortAttribute: bson.M{"category_count": -1}},
	}
	err = c.Collection.Pipe(categoryCountPipeline).All(&entityStats.CategoryWiseCount)

	//Get category group wise count
	categoryGroupCountPipeline := []bson.M{
		{constants.UnwindAttribute: constants.CategoryAttribute},
		{constants.SortAttribute: bson.M{"categories": 1}},
		{constants.GroupAttribute: bson.M{"_id": "$_id", "sortedCategories": bson.M{"$push": constants.CategoryAttribute}}},
		{
			constants.GroupAttribute: bson.M{
				"_id":            "$sortedCategories",
				"category_count": bson.M{"$sum": 1}}},
		{constants.SortAttribute: bson.M{"category_count": -1}},
	}
	err = c.Collection.Pipe(categoryGroupCountPipeline).All(&entityStats.CategoryGroupWiseCount)

	// Get total number of relations
	linkSumPipeline := []bson.M{{
		constants.GroupAttribute: bson.M{
			"_id":      "$link_sum",
			"link_sum": bson.M{"$sum": bson.M{"$size": "$links"}}}}}

	err = c.Collection.Pipe(linkSumPipeline).All(&linkCount)
	entityStats.RelationCount, _ = linkCount[0]["link_sum"].(int)

	return entityStats, err
}

/*
GetGraph - Get the entity relations summary for graph visualization
*/
func (e EntityRepository) GetGraph() (graph map[string]models.GraphArray, err error) {
	graph = make(map[string]models.GraphArray)
	c := e.newEntityCollection()
	resultQuery := c.Collection.Find(bson.D{}).Select(bson.M{"title": 1, "links": 1, "categories": 1}).Iter()
	if err != nil {
		return
	}
	// iterate through all documents and map to graph array
	var entity models.Entity
	for resultQuery.Next(&entity) {

		var links []string
		for _, link := range entity.Links {
			links = append(links, link.Title)
		}
		if err = resultQuery.Close(); err != nil {
			return
		}
		graph[entity.GetTitle()] = models.GraphArray{Title: entity.GetTitle(), Categories: entity.Categories, Links: links}
	}
	return
}
