package managers

import (
	"GIG/app/utilities/managers"
	"GIG/tests/app/test_values"
	"github.com/lsflk/gig-sdk/models"
)

/*
TestThatNewEntityTitleIsWithinLifetimeOfExistingEntityReturnsFalseIfEntityIsTerminated
New entity title is within lifetime of existing entity returns false if entity is terminated
*/
func (t *TestManagers) TestThatNewEntityTitleIsWithinLifetimeOfExistingEntityReturnsFalseIfEntityIsTerminated() {
	TestValue := managers.EntityManager{}.NewEntityTitleIsWithinLifetimeOfExistingEntity(models.Attribute{}, models.Attribute{}, true)

	t.AssertEqual(TestValue, false)
}

/*
TestThatNewEntityTitleIsWithinLifetimeOfExistingEntityReturnsTrueIfNewTitleDateIsAfterLastTitleDate
New entity title is within lifetime of existing entity returns true if new title date is after last title date
*/
func (t *TestManagers) TestThatNewEntityTitleIsWithinLifetimeOfExistingEntityReturnsTrueIfNewTitleDateIsAfterLastTitleDate() {

	lastTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj)
	newTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj2)

	TestValue := managers.EntityManager{}.NewEntityTitleIsWithinLifetimeOfExistingEntity(newTitleAttribute, lastTitleAttribute, false)

	t.AssertEqual(newTitleAttribute.GetValue().GetDate().After(lastTitleAttribute.GetValue().GetDate()), true)
	t.AssertEqual(TestValue, true)
}

/*
TestThatNewEntityTitleIsWithinLifetimeOfExistingEntityReturnsFalseIfNewTitleDateIsBeforeLastTitleDate
New entity title is within lifetime of existing entity returns false if new title date is before last title date
*/
func (t *TestManagers) TestThatNewEntityTitleIsWithinLifetimeOfExistingEntityReturnsFalseIfNewTitleDateIsBeforeLastTitleDate() {

	lastTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj2)
	newTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj)

	TestValue := managers.EntityManager{}.NewEntityTitleIsWithinLifetimeOfExistingEntity(newTitleAttribute, lastTitleAttribute, false)

	t.AssertEqual(newTitleAttribute.GetValue().GetDate().Before(lastTitleAttribute.GetValue().GetDate()), true)
	t.AssertEqual(TestValue, false)
}

/*
TestThatNewEntityTitleIsWithinLifetimeOfExistingEntityReturnsFalseIfNewTitleDateIsZero
New entity title is within lifetime of existing entity returns false if new title date zero
*/
func (t *TestManagers) TestThatNewEntityTitleIsWithinLifetimeOfExistingEntityReturnsFalseIfNewTitleDateIsZero() {

	lastTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj2)
	newTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj0)

	TestValue := managers.EntityManager{}.NewEntityTitleIsWithinLifetimeOfExistingEntity(newTitleAttribute, lastTitleAttribute, false)

	t.AssertEqual(newTitleAttribute.GetValue().GetDate().IsZero(), true)
	t.AssertEqual(TestValue, false)
}

/*
TestThatNewEntityTitleIsWithinLifetimeOfExistingEntityReturnsFalseIfNewTitleDateEqualsLastTitleDate
New entity title is within lifetime of existing entity returns false if new title date equals last title date
*/
func (t *TestManagers) TestThatNewEntityTitleIsWithinLifetimeOfExistingEntityReturnsFalseIfNewTitleDateEqualsLastTitleDate() {

	lastTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj2)
	newTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj2)

	TestValue := managers.EntityManager{}.NewEntityTitleIsWithinLifetimeOfExistingEntity(newTitleAttribute, lastTitleAttribute, false)

	t.AssertEqual(newTitleAttribute.GetValue().GetDate(), lastTitleAttribute.GetValue().GetDate())
	t.AssertEqual(TestValue, false)
}

/*
TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsTrueIfExistingEntityIsNotTerminatedAndNewEntitySourceDateIsAfterExistingEntitySourceDate
New entity is within lifetime of existing entity returns
true if existing entity is not terminated and new entity source date is after existing entity source date
*/
func (t *TestManagers) TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsTrueIfExistingEntityIsNotTerminatedAndNewEntitySourceDateIsAfterExistingEntitySourceDate() {

	lastTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj)

	testEntity := *new(models.Entity).SetTitle(test_values.TestValueObj2).SetSourceDate(test_values.TestValueObj2.GetDate())

	TestValue := managers.EntityManager{}.NewEntityIsWithinLifeTimeOfExistingEntity(testEntity, lastTitleAttribute, false)

	t.AssertEqual(testEntity.IsTerminated(), false)
	t.AssertEqual(testEntity.GetSourceDate().After(lastTitleAttribute.GetValue().Date), true)
	t.AssertEqual(TestValue, true)
}

/*
TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsTrueIfExistingEntityIsNotTerminatedAndNewEntitySourceDateEqualsExistingEntitySourceDate
New entity is within lifetime of existing entity returns
true if existing entity is not terminated and new entity source date equals existing entity source date
*/
func (t *TestManagers) TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsTrueIfExistingEntityIsNotTerminatedAndNewEntitySourceDateEqualsExistingEntitySourceDate() {

	lastTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj).SetValue(test_values.TestValueObj2).SetValue(test_values.TestValueObj3)

	testEntity := *new(models.Entity).SetTitle(test_values.TestValueObj).SetSourceDate(test_values.TestValueObj.GetDate())

	TestValue := managers.EntityManager{}.NewEntityIsWithinLifeTimeOfExistingEntity(testEntity, lastTitleAttribute, false)

	t.AssertEqual(testEntity.IsTerminated(), false)
	t.AssertEqual(testEntity.GetSourceDate().Equal(lastTitleAttribute.GetValueByDate(test_values.TestValueObj.GetDate()).Date), true)
	t.AssertEqual(TestValue, true)
}

/*
TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsTrueIfExistingEntityIsTerminatedAndNewEntitySourceDateIsWithinEntityLifetime
New entity is within lifetime of existing entity returns
true if existing entity is terminated but new entity source date is between existing entity lifetime
*/
func (t *TestManagers) TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsTrueIfExistingEntityIsTerminatedAndNewEntitySourceDateIsWithinEntityLifetime() {

	lastTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj).SetValue(test_values.TestValueObj2).SetValue(test_values.TestValueObj3)

	testEntity := *new(models.Entity).SetTitle(test_values.TestValueObj).SetSourceDate(test_values.TestValueObj2.GetDate())

	TestValue := managers.EntityManager{}.NewEntityIsWithinLifeTimeOfExistingEntity(testEntity, lastTitleAttribute, true)

	t.AssertEqual(testEntity.GetSourceDate().After(lastTitleAttribute.GetValues()[0].GetDate()), true)
	t.AssertEqual(testEntity.GetSourceDate().Before(lastTitleAttribute.GetValue().GetDate()), true)
	t.AssertEqual(TestValue, true)
}

/*
TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsTrueIfExistingEntityIsTerminatedButNewEntitySourceDateEqualsExistingEntitySourceDate
New entity is within lifetime of existing entity returns
true if existing entity is terminated but new entity source date equals existing entity source date
*/
func (t *TestManagers) TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsTrueIfExistingEntityIsTerminatedButNewEntitySourceDateEqualsExistingEntitySourceDate() {

	lastTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj).SetValue(test_values.TestValueObj2).SetValue(test_values.TestValueObj3)

	testEntity := *new(models.Entity).SetTitle(test_values.TestValueObj).SetSourceDate(test_values.TestValueObj.GetDate())

	TestValue := managers.EntityManager{}.NewEntityIsWithinLifeTimeOfExistingEntity(testEntity, lastTitleAttribute, true)

	t.AssertEqual(testEntity.GetSourceDate().Equal(lastTitleAttribute.GetValues()[0].GetDate()), true)
	t.AssertEqual(TestValue, true)
}

/*
TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsFalseIfExistingEntityIsTerminatedAndNewEntitySourceDateAfterExistingEntityTerminationDate
New entity is within lifetime of existing entity returns
false if existing entity is terminated and new entity source date is after existing entity termination date
*/
func (t *TestManagers) TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsFalseIfExistingEntityIsTerminatedAndNewEntitySourceDateAfterExistingEntityTerminationDate() {

	lastTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj).SetValue(test_values.TestValueObj2)

	testEntity := *new(models.Entity).SetTitle(test_values.TestValueObj3).SetSourceDate(test_values.TestValueObj3.GetDate())

	TestValue := managers.EntityManager{}.NewEntityIsWithinLifeTimeOfExistingEntity(testEntity, lastTitleAttribute, true)

	t.AssertEqual(testEntity.GetSourceDate().After(lastTitleAttribute.GetValue().GetDate()), true)
	t.AssertEqual(TestValue, false)
}

/*
TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsFalseIfNewEntitySourceDateIsBeforeExistingEntitySourceDate
New entity is within lifetime of existing entity returns
false if new entity source date is before existing entity source date
*/
func (t *TestManagers) TestThatNewEntityIsWithinLifetimeOfExistingEntityReturnsFalseIfNewEntitySourceDateIsBeforeExistingEntitySourceDate() {

	lastTitleAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj3).SetValue(test_values.TestValueObj2)
	newAttribute := *new(models.Attribute).SetValue(test_values.TestValueObj)

	testEntity := *new(models.Entity).SetTitle(test_values.TestValueObj).SetSourceDate(test_values.TestValueObj.GetDate())

	TestValue := managers.EntityManager{}.NewEntityIsWithinLifeTimeOfExistingEntity(testEntity, lastTitleAttribute, true)

	t.AssertEqual(newAttribute.GetValue().GetDate().Before(lastTitleAttribute.GetValue().GetDate()), true)
	t.AssertEqual(testEntity.GetSourceDate().Before(lastTitleAttribute.GetValues()[0].GetDate()), true)
	t.AssertEqual(TestValue, false)
}
