package mongodb

import (
	"GIG/app/constants/database"
	"GIG/app/databases/mongodb"
	"GIG/app/repositories/interfaces"
	"time"

	"github.com/lsflk/gig-sdk/models"
)

type StatRepository struct {
	interfaces.StatRepositoryInterface
}

func (e StatRepository) newStatCollection() *mongodb.Collection {
	return mongodb.NewCollectionSession(database.StatCollection)
}

/*
AddStat insert a new Stat into database and returns
last inserted stat on success.
*/
func (e StatRepository) AddStat(stat models.EntityStats) (models.EntityStats, error) {
	c := e.newStatCollection()
	defer c.Close()
	stat.CreatedAt = time.Now()
	return stat, c.Collection.Insert(stat)
}

/*
GetLastStat Get a Last Stat from database and returns
a models. Stat on success
*/
func (e StatRepository) GetLastStat() (models.EntityStats, error) {
	var (
		stat models.EntityStats
		err  error
	)

	c := e.newStatCollection()
	defer c.Close()

	err = c.Collection.Find(nil).Sort("-created_at").One(&stat)
	return stat, err
}
