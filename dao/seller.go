package dao

// InsertProduct 插入一条记录到product表
func InsertProduct(name string, typee string, price float64, number int, url string) error {
	sql := "insert into product(name, type, price, number, image, sales) values(?, ?, ?, ?, ?, 0);"
	_, err := db.Exec(sql, name, typee, price, number, url)

	return err
}
