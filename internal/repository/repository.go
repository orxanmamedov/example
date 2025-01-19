package repository

import (
	"context"
	"database/sql"
	"example/internal/domain"
	"example/pkg/storage/mysql"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"strings"
)

type UserRepo struct {
	db     mysql.MySQL
	logger *zap.Logger
}

func NewUserRepo(db mysql.MySQL, logger *zap.Logger) *UserRepo {
	return &UserRepo{db: db, logger: logger}
}
func (ur *UserRepo) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	ur.logger.Info("UserRepo.GetAllUsers")
	const query = `select id, username, email from user`

	var users []domain.User
	err := ur.db.SelectContext(ctx, &users, query)
	if err != nil {
		ur.logger.Error("UserRepo.GetAllUsers select error", zap.Error(err))
		return nil, err
	}

	if len(users) == 0 {
		ur.logger.Error("UserRepo.GetAllUsers no users found", zap.Error(sql.ErrNoRows))
		return nil, sql.ErrNoRows
	}
	return users, nil
}

func (ur *UserRepo) GetUser(ctx context.Context, ID uuid.UUID) (domain.User, error) {
	ur.logger.Info("UserRepo.GetUser", zap.String("id", ID.String()))
	const query = `select id, username, email from user where id = ?`

	var user domain.User
	err := ur.db.SelectContext(ctx, &user, query, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ur.logger.Warn("UserRepo.GetUser: no user found", zap.String("ID", ID.String()))
			return domain.User{}, nil
		}
		ur.logger.Error("UserRepo.GetUser query error", zap.Error(err))
		return domain.User{}, err
	}
	return user, nil
}

func (ur *UserRepo) CreateUser(ctx context.Context, user domain.User) (uuid.UUID, error) {
	ur.logger.Info("UserRepo.CreateUser")

	id, err := uuid.NewV7()
	if err != nil {
		ur.logger.Error("UserRepo.CreateUser UUID generation error", zap.Error(err))
		return uuid.Nil, err
	}

	user.ID = id

	const query = `INSERT INTO user (id, username, email) VALUES (?, ?, ?)`

	_, err = ur.db.ExecContext(ctx, query, user)

	if err != nil {
		ur.logger.Error("UserRepo.CreateUser exec error", zap.Error(err))
		return uuid.Nil, err
	}

	ur.logger.Info("UserRepo.CreateUser: user created successfully", zap.String("ID", id.String()))

	return id, nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ur.logger.Info("UserRepo.UpdateUser", zap.String("id", user.ID.String()))

	query := `UPDATE user SET`
	params := []interface{}{}

	if user.Name != "" {
		query += `Name = ?,`
		params = append(params, user.Name)
	}

	if user.Email != "" {
		query += " Email = ?,"
		params = append(params, user.Email)
	}

	// Remove trailing comma and add WHERE clause
	query = strings.TrimSuffix(query, ",") + " WHERE id = ?"
	params = append(params, user.ID)

	res, err := ur.db.ExecContext(ctx, query, params...)
	if err != nil {
		ur.logger.Error("UserRepo.UpdateUser exec error", zap.Error(err))
		return domain.User{}, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		ur.logger.Error("UserRepo.UpdateUser RowsAffected error", zap.Error(err))
		return domain.User{}, err
	}
	if affected == 0 {
		ur.logger.Error("UserRepo.UpdateUser no rows affected error", zap.Error(err))
		return domain.User{}, err
	}

	ur.logger.Info("UserRepo.UpdateUser updated user successfully", zap.String("ID", user.ID.String()))
	return user, nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, ID uuid.UUID) error {
	ur.logger.Info("UserRepo.DeleteUser", zap.String("id", ID.String()))

	const query = `DELETE FROM user WHERE id = ?`

	res, err := ur.db.ExecContext(ctx, query, ID)

	if err != nil {
		ur.logger.Error("UserRepo.DeleteUser exec error", zap.Error(err))
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		ur.logger.Error("UserRepo.DeleteUser RowsAffected error", zap.Error(err))
		return err
	}

	if affected == 0 {
		ur.logger.Error("UserRepo.DeleteUser no rows affected error", zap.Error(err))
		return err
	}

	ur.logger.Info("UserRepo.DeleteUser: user deleted successfully", zap.String("ID", ID.String()))

	return nil
}
