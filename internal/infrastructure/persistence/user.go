package persistence

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"test_pet/internal/domain/entity"
	"test_pet/internal/domain/repository"
)

type User struct {
	Db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.User {
	return &User{Db: db}
}

func (p *User) GetList(limit, offset int32) ([]entity.User, error) {
	var (
		err      error
		rows     *sqlx.Rows
		userList []entity.User
		user     entity.User
		query    = `SELECT * FROM "user" WHERE is_deleted = 0`
	)

	if limit != 0 {
		query = fmt.Sprintf(`SELECT * FROM "user" WHERE is_deleted = 0 LIMIT %d OFFSET %d`, limit, offset)
	}

	if err := p.Db.Ping(); err != nil {
		return nil, err
	}

	if rows, err = p.Db.Queryx(query); err != nil {
		if err == sql.ErrNoRows {
			return userList, nil
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.StructScan(&user); err != nil {
			return nil, err
		}

		userList = append(userList, user)
	}

	return userList, nil
}

func (p *User) Add(name string) (int64, error) {
	if err := p.Db.Ping(); err != nil {
		return 0, err
	}

	var lastInsertId int64
	err := p.Db.QueryRow(`INSERT INTO "user"(name) VALUES ($1) RETURNING id`, name).Scan(&lastInsertId)

	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (p *User) DeleteById(userId int64) error {
	if err := p.Db.Ping(); err != nil {
		return err
	}

	res, err := p.Db.NamedExec(`UPDATE "user" SET is_deleted = 1 WHERE id = :id AND is_deleted = 0`,
		map[string]interface{}{
			"id": userId,
		})

	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return repository.MissingClientWithId
	}

	return nil
}
