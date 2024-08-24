package gorbac_gorm

import (
	"fmt"
	"github.com/kordar/gorbac"
	"gorm.io/gorm"
)

// 基于Gorm的rbac实现类

type SqlRbac struct {
	db *gorm.DB
}

func NewSqlRbac(db *gorm.DB) *SqlRbac {
	return &SqlRbac{db: db}
}

func (rbac *SqlRbac) AddItem(item gorbac.Item) error {
	authItem := ToAuthItem(item)
	if item.GetRuleName() == "" {
		return rbac.db.Omit("rule_name").Create(&authItem).Error
	}
	return rbac.db.Create(&authItem).Error
}

func (rbac *SqlRbac) GetItem(name string) (gorbac.Item, error) {
	authItem := AuthItem{}
	tx := rbac.db.Where("name = ?", name).First(&authItem)
	if err := tx.Error; err == nil {
		item := ToItem(authItem)
		return item, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) GetItems(t int32) ([]gorbac.Item, error) {
	var authItems []AuthItem
	tx := rbac.db.Where("type = ?", t).Find(&authItems)
	if err := tx.Error; err == nil {
		items := ToItems(authItems)
		return items, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) FindAllItems() ([]gorbac.Item, error) {
	var authItems []AuthItem
	tx := rbac.db.Find(&authItems)
	if err := tx.Error; err == nil {
		items := ToItems(authItems)
		return items, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) RemoveItem(name string) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		authItemChild := AuthItemChild{}
		tx.Where("parent = ? or child = ?", name, name).Delete(&authItemChild)
		authAssignment := AuthAssignment{}
		tx.Where("item_name = ?", name).Delete(&authAssignment)
		authItem := AuthItem{}
		tx.Where("name = ?", name).Delete(&authItem)
		return nil
	})
}

func (rbac *SqlRbac) UpdateItem(itemName string, updateItem gorbac.Item) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		if itemName != updateItem.GetName() {
			authItemChild := AuthItemChild{}
			authAssignment := AuthAssignment{}
			tx.Model(&authItemChild).Where("parent = ?", itemName).Update("parent", updateItem.GetName())
			tx.Model(&authItemChild).Where("child = ?", itemName).Update("child", updateItem.GetName())
			tx.Model(&authAssignment).Where("item_name = ?", itemName).Update("item_name", updateItem.GetName())
		}
		authItem := ToAuthItem(updateItem)
		return tx.Model(&authItem).Where("name = ?", itemName).Omit("create_time").Updates(&authItem).Error
	})
}

func (rbac *SqlRbac) AddRule(rule gorbac.Rule) error {
	authRule := ToAuthRule(rule)
	return rbac.db.Create(&authRule).Error
}

func (rbac *SqlRbac) GetRule(name string) (*gorbac.Rule, error) {
	var authRule AuthRule
	tx := rbac.db.Where("name = ?", name).First(&authRule)
	if err := tx.Error; err == nil {
		rule := ToRule(authRule)
		return rule, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) GetRules() ([]*gorbac.Rule, error) {
	var authRules []AuthRule
	tx := rbac.db.Find(&authRules)
	if err := tx.Error; err == nil {
		rules := ToRules(authRules)
		return rules, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) RemoveRule(ruleName string) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		authItem := AuthItem{}
		tx.Model(&authItem).Where("rule_name = ?", ruleName).Update("rule_name", nil)
		authRule := AuthRule{}
		tx.Where("name = ?", ruleName).Delete(&authRule)
		return nil
	})
}

func (rbac *SqlRbac) UpdateRule(ruleName string, updateRule gorbac.Rule) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		if ruleName != updateRule.Name {
			authItem := AuthItem{}
			tx.Model(&authItem).Where("rule_name = ?", ruleName).Update("rule_name", updateRule.Name)
		}
		authRule := AuthRule{}
		return tx.Model(&authRule).Where("name = ?", ruleName).Omit("create_time").Updates(&updateRule).Error
	})
}

// FindRolesByUser 通过会员id获取关联的所有角色
func (rbac *SqlRbac) FindRolesByUser(userId interface{}) ([]gorbac.Item, error) {
	authAssignment := AuthAssignment{}
	var authItems []AuthItem
	tx := rbac.db.Model(&authAssignment).
		Joins(fmt.Sprintf("inner join %s on %s.item_name = %s.name", gorbac.GetTableName("item"), gorbac.GetTableName("assignment"), gorbac.GetTableName("item"))).
		Where(fmt.Sprintf("%s.user_id = ? and %s.`type` = 1", gorbac.GetTableName("assignment"), gorbac.GetTableName("item")), userId).
		Find(&authItems)
	if err := tx.Error; err == nil {
		items := ToItems(authItems)
		return items, nil
	} else {
		return nil, err
	}

}

func (rbac *SqlRbac) FindChildrenList() ([]*gorbac.ItemChild, error) {
	var children []AuthItemChild
	tx := rbac.db.Find(&children)
	if err := tx.Error; err == nil {
		itemChildren := ToItemChildren(children)
		return itemChildren, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) FindChildrenFormChild(child string) ([]*gorbac.ItemChild, error) {
	var children []AuthItemChild
	tx := rbac.db.Where("child = ?", child).Find(&children)
	if err := tx.Error; err == nil {
		itemChildren := ToItemChildren(children)
		return itemChildren, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) GetItemList(t int32, names []string) ([]gorbac.Item, error) {
	var authItems []AuthItem
	if len(names) > 0 {
		tx := rbac.db.Where("type = ? and name in ?", t, names).Find(&authItems)
		if err := tx.Error; err == nil {
			items := ToItems(authItems)
			return items, nil
		} else {
			return nil, err
		}
	} else {
		tx := rbac.db.Where("type = ?", t).Find(&authItems)
		if err := tx.Error; err == nil {
			items := ToItems(authItems)
			return items, nil
		} else {
			return nil, err
		}
	}
}

func (rbac *SqlRbac) FindPermissionsByUser(userId interface{}) ([]gorbac.Item, error) {
	authAssignment := AuthAssignment{}
	var authItems []AuthItem
	tx := rbac.db.Model(&authAssignment).
		Joins(fmt.Sprintf("inner join %s on %s.item_name = %s.name", gorbac.GetTableName("item"), gorbac.GetTableName("assignment"), gorbac.GetTableName("item"))).
		Where(fmt.Sprintf("%s.user_id = ? and %s.type = %d", gorbac.GetTableName("assignment"), gorbac.GetTableName("item"), gorbac.PermissionType), userId).
		Find(&authItems)
	if err := tx.Error; err == nil {
		items := ToItems(authItems)
		return items, err
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) FindAssignmentByUser(userId interface{}) ([]*gorbac.Assignment, error) {
	var authAssignments []AuthAssignment
	tx := rbac.db.Where("user_id = ?", userId).Find(&authAssignments)
	if err := tx.Error; err == nil {
		assignments := ToAssignments(authAssignments)
		return assignments, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) AddItemChild(itemChild gorbac.ItemChild) error {
	authItemChild := ToAuthItemChild(itemChild.Parent, itemChild.Child)
	return rbac.db.Create(&authItemChild).Error
}

func (rbac *SqlRbac) RemoveChild(parent string, child string) error {
	var itemChild AuthItemChild
	return rbac.db.Where("parent = ? and child = ?", parent, child).Delete(&itemChild).Error
}

func (rbac *SqlRbac) RemoveChildren(parent string) error {
	var itemChild AuthItemChild
	return rbac.db.Where("parent = ?", parent).Delete(&itemChild).Error
}

func (rbac *SqlRbac) HasChild(parent string, child string) bool {
	var itemChild AuthItemChild
	first := rbac.db.Model(&itemChild).Where("parent = ? and child = ?", parent, child).First(&itemChild)
	return first.Error == nil
}

func (rbac *SqlRbac) FindChildren(name string) ([]gorbac.Item, error) {
	var authItems []AuthItem
	tx := rbac.db.Model(&AuthItem{}).
		Joins(fmt.Sprintf("inner join %s on %s.name = %s.child", gorbac.GetTableName("item-child"), gorbac.GetTableName("item"), gorbac.GetTableName("item-child"))).
		Where(fmt.Sprintf("%s.parent = ?", gorbac.GetTableName("item-child")), name)
	if err := tx.Error; err == nil {
		items := ToItems(authItems)
		return items, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) Assign(assignment gorbac.Assignment) error {
	authAssignment := ToAuthAssignment(assignment)
	return rbac.db.Create(&authAssignment).Error
}

func (rbac *SqlRbac) RemoveAssignment(userId interface{}, name string) error {
	var assignment AuthAssignment
	return rbac.db.Where("user_id = ? and item_name = ?", userId, name).Delete(&assignment).Error
}

func (rbac *SqlRbac) RemoveAllAssignmentByUser(userId interface{}) error {
	var assignment AuthAssignment
	return rbac.db.Where("user_id = ?", userId).Delete(&assignment).Error
}

func (rbac *SqlRbac) RemoveAllAssignments() error {
	var assignment AuthAssignment
	return rbac.db.Delete(&assignment).Error
}

func (rbac *SqlRbac) GetAssignment(userId interface{}, name string) (*gorbac.Assignment, error) {
	var authAssignment AuthAssignment
	tx := rbac.db.Where("user_id = ? and item_name = ?", userId, name).First(&authAssignment)
	if err := tx.Error; err == nil {
		assignment := ToAssignment(authAssignment)
		return assignment, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) GetAssignmentByItems(name string) ([]*gorbac.Assignment, error) {
	var authAssignments []AuthAssignment
	tx := rbac.db.Where("item_name = ?", name).Find(&authAssignments)
	if err := tx.Error; err == nil {
		assignments := ToAssignments(authAssignments)
		return assignments, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) GetAssignments(userId interface{}) ([]*gorbac.Assignment, error) {
	var authAssignments []AuthAssignment
	tx := rbac.db.Where("user_id = ?", userId).Find(&authAssignments)
	if err := tx.Error; err == nil {
		assignments := ToAssignments(authAssignments)
		return assignments, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) GetAllAssignment() ([]*gorbac.Assignment, error) {
	var authAssignments []AuthAssignment
	tx := rbac.db.Find(&authAssignments)
	if err := tx.Error; err == nil {
		assignments := ToAssignments(authAssignments)
		return assignments, nil
	} else {
		return nil, err
	}
}

func (rbac *SqlRbac) RemoveAll() error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		var authAssignment AuthAssignment
		tx.Delete(&authAssignment)
		var authItem AuthItem
		tx.Delete(&authItem)
		var authRule AuthRule
		tx.Delete(&authRule)
		return nil
	})
}

func (rbac *SqlRbac) RemoveChildByNames(key string, names []string) error {
	if names != nil && len(names) > 0 {
		var authItemChild AuthItemChild
		return rbac.db.Where(key+" in (?)", names).Delete(&authItemChild).Error
	}
	return nil
}

func (rbac *SqlRbac) RemoveAssignmentByName(names []string) error {
	if names != nil && len(names) > 0 {
		var authAssignments AuthAssignment
		return rbac.db.Where("item_name in ?", names).Delete(&authAssignments).Error
	}
	return nil
}

func (rbac *SqlRbac) RemoveItemByType(t int32) error {
	var authItem AuthItem
	return rbac.db.Where("type = ?", t).Delete(&authItem).Error
}

func (rbac *SqlRbac) RemoveAllRules() error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		var authItem AuthItem
		tx.Model(&authItem).Update("rule_name", nil)
		var authRule AuthRule
		tx.Delete(&authRule)
		return nil
	})
}
