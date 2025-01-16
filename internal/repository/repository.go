package repository

import (
	"context"
	"database/sql"
	"example/internal/domain"
	"example/pkg/storage/mysql"
	"go.uber.org/zap"
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
	query := `select id, username, email from user`

	var users []domain.User
	err := ur.db.SelectContext(ctx, &users, query)
	if err != nil {
		ur.logger.Error("UserRepo.GetAllUsers select error", zap.Error(err))
		return nil, err
	}
	return users, nil
}

func (ur *UserRepo) GetUser(ctx context.Context, ID int) (domain.User, error) {
	ur.logger.Info("UserRepo.GetUser", zap.Int("id", ID))
	query := `select id, username, email from user where id = ?`

	var user domain.User
	err := ur.db.SelectContext(ctx, &user, query, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ur.logger.Warn("UserRepo.GetUser: no user found", zap.Int("ID", ID))
			return domain.User{}, nil
		}
		ur.logger.Error("UserRepo.GetUser query error", zap.Error(err))
		return domain.User{}, err
	}
	return user, nil
}

func (ur *UserRepo) CreateUser(ctx context.Context, user domain.User) (int, error) {
	ur.logger.Info("UserRepo.CreateUser")

	query := `insert into user (id, username, email) values (?, ?)`

	res, err := ur.db.ExecContext(ctx, query, user.ID, user.Username, user.Email)

	if err != nil {
		ur.logger.Error("UserRepo.CreateUser exec error", zap.Error(err))
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		ur.logger.Error("UserRepo.CreateUser LastInsertId error", zap.Error(err))
		return 0, err
	}

	ur.logger.Info("UserRepo.CreateUser: user created successfully", zap.Int64("ID", id))

	return int(id), nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ur.logger.Info("UserRepo.UpdateUser", zap.Int("id", user.ID))

	query := "UPDATE user SET"
	params := []interface{}{}

	if user.Username != nil {
		query += " Username = ?,"
		params = append(params, *user.Username)
	}

	if user.Email != nil {
		query += " Email = ?,"
		params = append(params, *user.Email)
	}

	// Remove trailing comma and add WHERE clause
	query = query[:len(query)-1] + " WHERE ID = ?"
	params = append(params, user.ID)

	res, err := ur.db.ExecContext(ctx, query, user.Username, user.Email, user.ID)
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

	ur.logger.Info("UserRepo.UpdateUser updated user successfully", zap.Int64("ID", int64(user.ID)))
	return user, nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, ID int) error {
	ur.logger.Info("UserRepo.DeleteUser", zap.Int("id", ID))

	query := "delete from user where id = ?"

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

	ur.logger.Info("UserRepo.DeleteUser: user deleted successfully", zap.Int("ID", ID))

	return nil
}
