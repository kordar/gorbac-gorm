module github.com/kordar/gorbac-gorm

go 1.16

replace (
	github.com/kordar/gorbac => ../../github.com/gorbac
)

require (
	github.com/kordar/gorbac v1.0.7
	gorm.io/gorm v1.25.10
)
