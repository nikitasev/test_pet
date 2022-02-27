package entity

type User struct {
	Id        int64  `db:"id"`
	Name      string `db:"name"`
	IsDeleted int8   `db:"is_deleted"`
}
