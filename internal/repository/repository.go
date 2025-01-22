package repository

import (
	"context"
	"database/sql"
	"example/internal/domain"
	UserUUIS "example/internal/uuid"
	"example/pkg/storage/mysql"
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
	const query = `SELECT uuid, username, email FROM user`

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

func (ur *UserRepo) GetUser(ctx context.Context, UUID UserUUIS.UUID) (*domain.User, error) {
	ur.logger.Info("UserRepo.GetUser", zap.String("uuid", UUID.String()))
	const query = `SELECT uuid, username, email FROM user WHERE uuid = ?`

	var user domain.User
	err := ur.db.GetContext(ctx, &user, query, UUID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			ur.logger.Warn("UserRepo.GetUser: no user found", zap.String("UUID", UUID.String()))
			return nil, nil
		}
		ur.logger.Error("UserRepo.GetUser query error", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepo) CreateUser(ctx context.Context, user domain.User) (UserUUIS.UUID, error) {
	ur.logger.Info("UserRepo.CreateUser")

	uuid := UserUUIS.NewUUID()
	user.UUID = uuid.String()

	const query = `INSERT INTO user (uuid, username, email) VALUES (?, ?, ?)`

	_, err := ur.db.ExecContext(ctx, query, user.UUID, user.Name, user.Email)
	if err != nil {
		ur.logger.Error("UserRepo.CreateUser exec error", zap.Error(err))
		return "", err
	}

	ur.logger.Info("UserRepo.CreateUser: user created successfully", zap.String("UUID", uuid.String()))

	return uuid, nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	ur.logger.Info("UserRepo.UpdateUser", zap.String("uuid", user.UUID))

	query := `UPDATE user SET `
	var args []any

	if user.Name != "" {
		query += "username = ?,"
		args = append(args, user.Name)
	}

	if user.Email != "" {
		query += " email = ?,"
		args = append(args, user.Email)
	}

	// Remove trailing comma and add WHERE clause
	query = strings.TrimSuffix(query, ",") + " WHERE uuid = ?"
	args = append(args, user.UUID)

	res, err := ur.db.ExecContext(ctx, query, args...)
	if err != nil {
		ur.logger.Error("UserRepo.UpdateUser exec error", zap.Error(err))
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		ur.logger.Error("UserRepo.UpdateUser RowsAffected error", zap.Error(err))
		return nil, err
	}
	if affected == 0 {
		ur.logger.Error("UserRepo.UpdateUser no rows affected error", zap.Error(err))
		return nil, sql.ErrNoRows
	}

	ur.logger.Info("UserRepo.UpdateUser updated user successfully", zap.String("UUID", user.UUID))
	return &user, nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, UUID UserUUIS.UUID) error {
	ur.logger.Info("UserRepo.DeleteUser", zap.String("UUID", UUID.String()))

	const query = `DELETE FROM user WHERE uuid = ?`

	res, err := ur.db.ExecContext(ctx, query, UUID.String())
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
		return sql.ErrNoRows
	}

	ur.logger.Info("UserRepo.DeleteUser: user deleted successfully", zap.String("UUID", UUID.String()))

	return nil
}
