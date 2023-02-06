package main

import (
	"su-ning-yi-gou/api"
	"su-ning-yi-gou/dao"
)

func main() {
	dao.InitDB()
	api.InitRouter()

}
