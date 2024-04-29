package gorbac_gorm

import (
	"fmt"
	"github.com/kordar/gorbac/db"
	"gorm.io/gorm"
)

type SqlRbac struct {
	db *gorm.DB
}

func NewSqlRbac(db *gorm.DB) *SqlRbac {
	return &SqlRbac{db: db}
}

func (rbac *SqlRbac) AddItem(authItem db.AuthItem) error {
	if authItem.RuleName == "" {
		return rbac.db.Omit("rule_name").Create(&authItem).Error
	}
	return rbac.db.Create(&authItem).Error
}

func (rbac *SqlRbac) GetItem(name string) (*db.AuthItem, error) {
	item := db.AuthItem{}
	err := rbac.db.Where("name = ?", name).First(&item).Error
	return &item, err
}

func (rbac *SqlRbac) GetItems(t int32) ([]*db.AuthItem, error) {
	var items []*db.AuthItem
	err := rbac.db.Where("type = ?", t).Find(&items).Error
	return items, err
}

func (rbac *SqlRbac) FindAllItems() ([]*db.AuthItem, error) {
	var items []*db.AuthItem
	err := rbac.db.Find(&items).Error
	return items, err
}

func (rbac *SqlRbac) AddRule(rule db.AuthRule) error {
	return rbac.db.Create(&rule).Error
}

func (rbac *SqlRbac) GetRule(name string) (*db.AuthRule, error) {
	var rule db.AuthRule
	err := rbac.db.Where("name = ?", name).First(&rule).Error
	return &rule, err
}

func (rbac *SqlRbac) GetRules() ([]*db.AuthRule, error) {
	var rules []*db.AuthRule
	err := rbac.db.Find(&rules).Error
	return rules, err
}

func (rbac *SqlRbac) RemoveItem(name string) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		itemChild := db.AuthItemChild{}
		tx.Where("parent = ? or child = ?", name, name).Delete(&itemChild)
		assignment := db.AuthAssignment{}
		tx.Where("item_name = ?", name).Delete(&assignment)
		item := db.AuthItem{}
		tx.Where("name = ?", name).Delete(&item)
		return nil
	})
}

func (rbac *SqlRbac) RemoveRule(ruleName string) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		item := db.AuthItem{}
		tx.Model(&item).Where("rule_name = ?", ruleName).Update("rule_name", nil)
		rule := db.AuthRule{}
		tx.Where("name = ?", ruleName).Delete(&rule)
		return nil
	})
}

func (rbac *SqlRbac) UpdateItem(itemName string, updateItem db.AuthItem) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		if itemName != updateItem.Name {
			child := db.AuthItemChild{}
			assignment := db.AuthAssignment{}
			tx.Model(&child).Where("parent = ?", itemName).Update("parent", updateItem.Name)
			tx.Model(&child).Where("child = ?", itemName).Update("child", updateItem.Name)
			tx.Model(&assignment).Where("item_name = ?", itemName).Update("item_name", updateItem.Name)
		}
		authItem := db.AuthItem{}
		return tx.Model(&authItem).Where("name = ?", itemName).Omit("create_at").Updates(&updateItem).Error
	})
}

func (rbac *SqlRbac) UpdateRule(ruleName string, updateRule db.AuthRule) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		if ruleName != updateRule.Name {
			item := db.AuthItem{}
			tx.Model(&item).Where("rule_name = ?", ruleName).Update("rule_name", updateRule.Name)
		}
		rule := db.AuthRule{}
		return tx.Model(&rule).Where("name = ?", ruleName).Omit("create_at").Updates(&updateRule).Error
	})
}

// FindRolesByUser 通过会员id获取关联的所有角色
func (rbac *SqlRbac) FindRolesByUser(userId interface{}) ([]*db.AuthItem, error) {
	assignment := db.AuthAssignment{}
	var items []*db.AuthItem
	err := rbac.db.Model(&assignment).
		Joins(fmt.Sprintf("inner join %s on %s.item_name = %s.name", db.GetTableName("item"), db.GetTableName("assignment"), db.GetTableName("item"))).
		Where(fmt.Sprintf("%s.user_id = ? and %s.`type` = 1", db.GetTableName("assignment"), db.GetTableName("item")), userId).
		Find(&items).Error
	return items, err
}

func (rbac *SqlRbac) FindChildrenList() ([]*db.AuthItemChild, error) {
	var children []*db.AuthItemChild
	err := rbac.db.Find(&children).Error
	return children, err
}

func (rbac *SqlRbac) FindChildrenFormChild(child string) ([]*db.AuthItemChild, error) {
	var children []*db.AuthItemChild
	err := rbac.db.Where("child = ?", child).Find(&children).Error
	return children, err
}

func (rbac *SqlRbac) GetItemList(t int32, names []string) ([]*db.AuthItem, error) {
	var items []*db.AuthItem
	if len(names) > 0 {
		err := rbac.db.Where("type = ? and name in ?", t, names).Find(&items).Error
		return items, err
	} else {
		err := rbac.db.Where("type = ?", t).Find(&items).Error
		return items, err
	}
}

func (rbac *SqlRbac) FindPermissionsByUser(userId interface{}) ([]*db.AuthItem, error) {
	assignment := db.AuthAssignment{}
	var items []*db.AuthItem
	err := rbac.db.Model(&assignment).
		Joins(fmt.Sprintf("inner join %s on %s.item_name = %s.name", db.GetTableName("item"), db.GetTableName("assignment"), db.GetTableName("item"))).
		Where(fmt.Sprintf("%s.user_id = ? and %s.type = 2", db.GetTableName("assignment"), db.GetTableName("item")), userId).
		Find(&items).Error
	return items, err
}

func (rbac *SqlRbac) FindAssignmentByUser(userId interface{}) ([]*db.AuthAssignment, error) {
	var assignments []*db.AuthAssignment
	err := rbac.db.Where("user_id = ?", userId).Find(&assignments).Error
	return assignments, err
}

func (rbac *SqlRbac) AddItemChild(itemChild db.AuthItemChild) error {
	return rbac.db.Create(&itemChild).Error
}

func (rbac *SqlRbac) RemoveChild(parent string, child string) error {
	var itemChild db.AuthItemChild
	return rbac.db.Where("parent = ? and child = ?", parent, child).Delete(&itemChild).Error
}

func (rbac *SqlRbac) RemoveChildren(parent string) error {
	var itemChild db.AuthItemChild
	return rbac.db.Where("parent = ?", parent).Delete(&itemChild).Error
}

func (rbac *SqlRbac) HasChild(parent string, child string) bool {
	var itemChild db.AuthItemChild
	first := rbac.db.Model(&itemChild).Where("parent = ? and child = ?", parent, child).First(&itemChild)
	return first.Error == nil
}

func (rbac *SqlRbac) FindChildren(name string) ([]*db.AuthItem, error) {
	var items []*db.AuthItem
	item := db.AuthItem{}
	err := rbac.db.Model(&item).
		Joins(fmt.Sprintf("inner join %s on %s.name = %s.child", db.GetTableName("item-child"), db.GetTableName("item"), db.GetTableName("item-child"))).
		Where(fmt.Sprintf("%s.parent = ?", db.GetTableName("item-child")), name).Error
	return items, err
}

func (rbac *SqlRbac) Assign(assignment db.AuthAssignment) error {
	return rbac.db.Create(&assignment).Error
}

func (rbac *SqlRbac) RemoveAssignment(userId interface{}, name string) error {
	var assignment db.AuthAssignment
	return rbac.db.Where("user_id = ? and item_name = ?", userId, name).Delete(&assignment).Error
}

func (rbac *SqlRbac) RemoveAllAssignmentByUser(userId interface{}) error {
	var assignment db.AuthAssignment
	return rbac.db.Where("user_id = ?", userId).Delete(&assignment).Error
}

func (rbac *SqlRbac) RemoveAllAssignments() error {
	var assignment db.AuthAssignment
	return rbac.db.Delete(&assignment).Error
}

func (rbac *SqlRbac) GetAssignment(userId interface{}, name string) (*db.AuthAssignment, error) {
	var assignments *db.AuthAssignment
	err := rbac.db.Where("user_id = ? and item_name = ?", userId, name).First(assignments).Error
	return assignments, err
}

func (rbac *SqlRbac) GetAssignmentByItems(name string) ([]*db.AuthAssignment, error) {
	var assignments []*db.AuthAssignment
	err := rbac.db.Where("item_name = ?", name).Find(&assignments).Error
	return assignments, err
}

func (rbac *SqlRbac) GetAssignments(userId interface{}) ([]*db.AuthAssignment, error) {
	var assignments []*db.AuthAssignment
	err := rbac.db.Where("user_id = ?", userId).Find(&assignments).Error
	return assignments, err
}

func (rbac *SqlRbac) GetAllAssignment() ([]*db.AuthAssignment, error) {
	var assignments []*db.AuthAssignment
	err := rbac.db.Find(&assignments).Error
	return assignments, err
}

func (rbac *SqlRbac) RemoveAll() error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		var assignment db.AuthAssignment
		tx.Delete(&assignment)
		var item db.AuthItem
		tx.Delete(&item)
		var rule db.AuthRule
		tx.Delete(&rule)
		return nil
	})
}

func (rbac *SqlRbac) RemoveChildByNames(key string, names []string) error {
	if names != nil && len(names) > 0 {
		var itemChild db.AuthItemChild
		return rbac.db.Where(key+" in (?)", names).Delete(&itemChild).Error
	}
	return nil
}

func (rbac *SqlRbac) RemoveAssignmentByName(names []string) error {
	if names != nil && len(names) > 0 {
		var assignments db.AuthAssignment
		return rbac.db.Where("item_name in (?)", names).Delete(&assignments).Error
	}
	return nil
}

func (rbac *SqlRbac) RemoveItemByType(t int32) error {
	var item db.AuthItem
	return rbac.db.Where("type = ?", t).Delete(&item).Error
}

func (rbac *SqlRbac) RemoveAllRules() error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		var item db.AuthItem
		tx.Model(&item).Update("rule_name", nil)
		var rule db.AuthRule
		tx.Delete(&rule)
		return nil
	})
}
