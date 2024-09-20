package managers

import (
	"github.com/lsflk/gig-sdk/enums/ValueType"
	"github.com/lsflk/gig-sdk/models"
)

type EntityManager struct {
}

func (e EntityManager) IsFromVerifiedSource(entity models.Entity) bool {
	return entity.GetSourceSignature() == "trusted"
}

func (e EntityManager) NewEntityTitleIsWithinLifetimeOfExistingEntity(newTitleAttribute models.Attribute, lastTitleAttribute models.Attribute, entityIsTerminated bool) bool {
	return !(entityIsTerminated || newTitleAttribute.GetValue().GetDate().IsZero()) &&
		newTitleAttribute.GetValue().GetDate().After(lastTitleAttribute.GetValue().GetDate())
}
func (e EntityManager) NewEntityIsWithinLifeTimeOfExistingEntity(entity models.Entity, lastTitleAttribute models.Attribute, entityIsTerminated bool) bool {
	return !entityIsTerminated ||
		(entityIsTerminated && !entity.GetSourceDate().IsZero() &&
			entity.GetSourceDate().Before(lastTitleAttribute.GetValue().GetDate())) &&
			!entity.GetSourceDate().Before(lastTitleAttribute.GetValues()[0].GetDate())
}

func (e EntityManager) CheckEntityCompatibility(existingEntity models.Entity, entity models.Entity) (bool, models.Entity) {
	//if an entity doesn't exist
	if existingEntity.GetTitle() != "" {

		lastTitleAttribute, _ := existingEntity.GetAttribute("titles")
		isValidEntity := e.NewEntityIsWithinLifeTimeOfExistingEntity(entity, lastTitleAttribute, existingEntity.IsTerminated())
		existingEntity, _ = e.MergeEntityTitle(existingEntity, entity)

		if isValidEntity {

			if existingEntity.GetSourceDate().IsZero() && e.IsFromVerifiedSource(entity) {
				newValue := models.Value{}
				newValue.SetValueString(entity.GetTitle()).
					SetSource(entity.Source).
					SetDate(entity.GetSourceDate()).
					SetType(ValueType.String)

				existingEntity.SetSourceDate(entity.GetSourceDate()).
					SetTitle(newValue).RemoveCategories([]string{"arbitrary-entities"})
			}

			// merge links
			existingEntity.AddLinks(entity.GetLinks())
			// merge categories
			existingEntity.AddCategories(entity.GetCategories())
			// merge attributes
			for name := range entity.GetAttributes() {
				if name != "new_title" && name != "titles" {
					entityAttribute, _ := entity.GetAttribute(name)
					existingEntity.SetAttribute(name, entityAttribute.GetValue())
				}
			}
			return true, existingEntity
		}
	}
	return false, existingEntity

}
