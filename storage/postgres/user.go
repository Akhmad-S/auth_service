package postgres

import (
	"errors"
	"time"

	ecom "github.com/uacademy/e_commerce/auth_service/proto-gen/e_commerce"
)

func (stg Postgres) CreateUser(id string, input *ecom.CreateUserRequest) error {
	_, err := stg.db.Exec(`INSERT INTO "user" (id, username, password, user_type) VALUES ($1, $2, $3, $4)`, id, input.Username, input.Password, input.UserType)
	if err != nil {
		return err
	}
	return nil
}

func (stg Postgres) GetUserList(offset, limit int, search string) (resp *ecom.GetUserListResponse, err error) {
	resp = &ecom.GetUserListResponse{
		Users: make([]*ecom.User, 0),
	}

	rows, err := stg.db.Queryx(`SELECT
	id,
	username,
	password,
	user_type,
	created_at,
	updated_at
	FROM "user" WHERE deleted_at IS NULL AND (username ILIKE '%' || $1 || '%')
	LIMIT $2
	OFFSET $3
	`, search, limit, offset)

	if err != nil {
		return resp, err
	}
	for rows.Next() {
		u := &ecom.User{}
		var updatedAt *string

		err := rows.Scan(
			&u.Id,
			&u.Username,
			&u.Password,
			&u.UserType,
			&u.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return resp, err
		}

		if updatedAt != nil {
			u.UpdatedAt = *updatedAt
		}

		resp.Users = append(resp.Users, u)
	}

	return resp, err
}

func (stg Postgres) GetUserById(id string) (*ecom.User, error) {
	res := &ecom.User{}
	var updatedAt *string

	err := stg.db.QueryRow(`SELECT id, username, password, user_type, created_at, updated_at FROM "user" WHERE id=$1 AND deleted_at IS NULL`, id).Scan(
		&res.Id, &res.Username, &res.Password, &res.UserType, &res.CreatedAt, &updatedAt,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	return res, nil
}

func (stg Postgres) UpdateUser(input *ecom.UpdateUserRequest) error {
	res, err := stg.db.NamedExec(`UPDATE "user"  SET password=:p, updated_at=now() WHERE deleted_at IS NULL AND id=:id`, map[string]interface{}{
		"id": input.Id,
		"p":  input.Password,
	})
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("user not found")
}

func (stg Postgres) DeleteUser(id string) error {
	res, err := stg.db.Exec(`UPDATE "user"  SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}
	return errors.New("user not found")
}

func (stg Postgres) GetUserByUsername(username string) (*ecom.User, error) {
	res := &ecom.User{}
	var deletedAt *time.Time
	var updatedAt *string
	err := stg.db.QueryRow(`SELECT 
		id,
		username,
		password,
		user_type,
		created_at,
		updated_at,
		deleted_at
    FROM "user" WHERE username = $1`, username).Scan(
		&res.Id,
		&res.Username,
		&res.Password,
		&res.UserType,
		&res.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return res, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	if deletedAt != nil {
		return res, errors.New("user not found")
	}

	return res, err
}
