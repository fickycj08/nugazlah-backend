package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vandenbill/nugazlah-backend/internal/entity"
	"github.com/vandenbill/nugazlah-backend/internal/ierr"
)

type classRepo struct {
	conn *pgxpool.Pool
}

func newClassRepo(conn *pgxpool.Pool) *classRepo {
	return &classRepo{conn}
}

func (r *classRepo) Insert(ctx context.Context, class entity.Class) error {
	query := "INSERT INTO classes (name, lecturer, description, icon, code, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	classID := ""
	err := r.conn.QueryRow(ctx, query,
		class.Name, class.Lecturer, class.Description, class.Icon, class.Code, class.UserID).Scan(&classID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	query = "INSERT INTO user_classes (class_id, user_id) VALUES ($1, $2)"

	_, err = r.conn.Exec(ctx, query, classID, class.UserID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (r *classRepo) JoinClass(ctx context.Context, userID string, classCode string) error {
	q2 := `SELECT id FROM classes WHERE code = $1`

	classID := ""
	err := r.conn.QueryRow(ctx, q2, classCode).Scan(&classID)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		return err
	}

	query := "INSERT INTO user_classes (user_id, class_id) VALUES ($1, $2)"

	_, err = r.conn.Exec(ctx, query, userID, classID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	rowsTasks, err := r.conn.Query(ctx, `SELECT id FROM tasks WHERE class_id = $1`, classID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	taskIDs := make([]string, 0)
	for rowsTasks.Next() {
		taskID := ""
		if err := rowsTasks.Scan(&taskID); err != nil {
			return err
		}
		taskIDs = append(taskIDs, taskID)
	}

	userTasks := make([]entity.UserTask, 0)
	for _, v := range taskIDs {
		uTask := entity.UserTask{}
		uTask.UserID = userID
		uTask.TaskID = v
		uTask.Status = false
		userTasks = append(userTasks, uTask)
	}

	queryTask := "INSERT INTO user_tasks (user_id, task_id, status) VALUES "
	values := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for i, uTask := range userTasks {
		placeholder := fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
		placeholders = append(placeholders, placeholder)
		values = append(values, uTask.UserID, uTask.TaskID, uTask.Status)
	}

	queryTask += strings.Join(placeholders, ",")
	_, err = r.conn.Exec(ctx, queryTask, values...)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	return nil
}

func (r *classRepo) GetMyClasses(ctx context.Context, userID string) ([]entity.Class, error) {
	rowsClass, err := r.conn.Query(ctx, `SELECT class_id FROM user_classes WHERE user_id = $1`, userID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return nil, ierr.ErrNotFound
			}
		}
		return nil, err
	}

	classIDs := make([]string, 0)
	for rowsClass.Next() {
		classID := ""
		if err := rowsClass.Scan(&classID); err != nil {
			return nil, err
		}
		classIDs = append(classIDs, classID)
	}

	rows, err := r.conn.Query(ctx, `SELECT id, name, lecturer, description, icon, code, user_id FROM classes WHERE id = ANY($1)`,
		classIDs)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return nil, ierr.ErrNotFound
			}
		}
		return nil, err
	}

	classes := make([]entity.Class, 0, 10)
	for rows.Next() {
		class := entity.Class{}
		if err := rows.Scan(&class.ID, &class.Name, &class.Lecturer, &class.Description, &class.Icon, &class.Code, &class.UserID); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}

	return classes, nil
}

func (r *classRepo) GetClass(ctx context.Context, classID string) (entity.Class, error) {
	class := entity.Class{}
	err := r.conn.QueryRow(ctx, `SELECT id, name, lecturer, description, icon, code, user_id FROM classes WHERE id = ANY($1)`,
		classID).Scan(&class.ID, &class.Name, &class.Lecturer, &class.Description, &class.Icon, &class.Code,
		&class.UserID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return class, ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return class, ierr.ErrNotFound
			}
		}
		return class, err
	}

	return class, nil
}

func (r *classRepo) IsAlreadyJoin(ctx context.Context, userID string, classCode string) (bool, error) {
	q2 := `SELECT id FROM classes WHERE code = $1`

	classID := ""
	err := r.conn.QueryRow(ctx, q2, classCode).Scan(&classID)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, ierr.ErrNotFound
		}
		return false, err
	}

	q := `SELECT 1 FROM user_classes WHERE class_id = $1 AND user_id = $2`

	v := 0
	err = r.conn.QueryRow(ctx, q, classID, userID).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
