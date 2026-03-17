# GoRBAC Gorm 实现

基于 [GoRBAC](https://github.com/kordar/gorbac) 的 Gorm 数据持久化实现，提供 RBAC 角色/权限/规则/分配 等实体的 MySQL 存储与查询能力。可与核心库无缝配合，通过实现的 AuthRepository 接口对接服务层。

---

## 安装

```bash
go get github.com/kordar/gorbac-gorm
```

依赖：
- Go 1.16+
- gorm.io/gorm（自行选择驱动，如 gorm.io/driver/mysql）
- github.com/kordar/gorbac

---

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/kordar/gorbac"
    gorbac_gorm "github.com/kordar/gorbac-gorm"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    // 连接数据库（按需替换 DSN）
    dsn := "user:pass@tcp(127.0.0.1:3306)/rbac?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    // 如需自定义表名，在迁移前设置
    // gorbac.SetTableName("item", "auth_item")
    // gorbac.SetTableName("item-child", "auth_item_child")
    // gorbac.SetTableName("assignment", "auth_assignment")
    // gorbac.SetTableName("rule", "auth_rule")

    // 自动迁移表结构
    _ = db.AutoMigrate(
        &gorbac_gorm.AuthItem{},
        &gorbac_gorm.AuthItemChild{},
        &gorbac_gorm.AuthAssignment{},
        &gorbac_gorm.AuthRule{},
    )

    // 创建仓库与服务
    repo := gorbac_gorm.NewSqlRbac(db)
    service := gorbac.NewRbacService(repo, true)

    // 创建角色与权限
    service.AddRole("admin", "管理员角色", "")
    service.AddPermission("view_dashboard", "查看仪表盘", "")
    _ = service.AssignChildren("admin", "view_dashboard")

    // 分配用户
    userId := 1001
    service.GetAuthManager().Assign(service.GetAuthManager().GetRole("admin"), userId)

    // 权限检查
    ctx := context.Background()
    allowed := service.GetAuthManager().CheckAccess(ctx, userId, "view_dashboard")
    fmt.Println("是否允许：", allowed)
}
```

---

## 表结构

默认表名由核心库决定，可通过 `gorbac.SetTableName` 覆盖：
- item: `auth_item`
- item-child: `auth_item_child`
- assignment: `auth_assignment`
- rule: `auth_rule`

仓库中附带了示例 SQL 文件，亦可选择使用 Gorm 自动迁移。

---

## 已实现接口

SqlRbac 实现了核心库的 AuthRepository 接口，涵盖：
- 角色/权限的增删改查与批量删除
- 规则的增删改查
- 父子关系维护（AddItemChild/RemoveChild/FindChildren 等）
- 用户与角色/权限的分配、批量分配与查询
- 基于用户的角色与权限查询

---

## 注意事项

- Rule 的 ExecuteName 会被持久化到 `auth_rule.execute_name`，用于在 gorbac 的 CheckAccess 中通过执行器做动态校验。
- 事务方法（如 RemoveItem/RemoveRule/RemoveAll）会在任意一步失败时返回错误，避免 silent failure。

---

## License

MIT License
