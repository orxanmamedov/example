package repository

import (
	"context"
	"database/sql"
	"example/internal/domain"
	"example/internal/uuid"
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
	const query = `SELECT id, username, email FROM user`

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
	const query = `SELECT id, username, email FROM user WHERE id = ?`

	var user domain.User
	err := ur.db.GetContext(ctx, &user, query, ID.String())
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

	id := uuid.NewUUID()
	user.UUID = id.String()

	const query = `INSERT INTO user (id, username, email) VALUES (?, ?, ?)`

	_, err := ur.db.ExecContext(ctx, query, user.UUID, user.Name, user.Email)
	if err != nil {
		ur.logger.Error("UserRepo.CreateUser exec error", zap.Error(err))
		return "", err
	}

	ur.logger.Info("UserRepo.CreateUser: user created successfully", zap.String("ID", id.String()))

	return id, nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ur.logger.Info("UserRepo.UpdateUser", zap.String("id", user.UUID))

	query := `UPDATE user SET `
	params := []interface{}{}

	if user.Name != "" {
		query += "username = ?,"
		params = append(params, user.Name)
	}

	if user.Email != "" {
		query += " email = ?,"
		params = append(params, user.Email)
	}

	// Remove trailing comma and add WHERE clause
	query = strings.TrimSuffix(query, ",") + " WHERE id = ?"
	params = append(params, user.UUID)

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
		return domain.User{}, sql.ErrNoRows
	}

	ur.logger.Info("UserRepo.UpdateUser updated user successfully", zap.String("ID", user.UUID))
	return user, nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, ID uuid.UUID) error {
	ur.logger.Info("UserRepo.DeleteUser", zap.String("id", ID.String()))

	const query = `DELETE FROM user WHERE id = ?`

	res, err := ur.db.ExecContext(ctx, query, ID.String())
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

	ur.logger.Info("UserRepo.DeleteUser: user deleted successfully", zap.String("ID", ID.String()))

	return nil
}
