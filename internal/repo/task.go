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

type taskRepo struct {
	conn *pgxpool.Pool
}

func newTaskRepo(conn *pgxpool.Pool) *taskRepo {
	return &taskRepo{conn}
}

func (r *taskRepo) Insert(ctx context.Context, userID string, task entity.Task) error {
	query := "INSERT INTO tasks (title, description, detail, submission, task_type, deadline, class_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	taskID := ""
	err := r.conn.QueryRow(ctx, query,
		task.Title, task.Description, task.Detail, task.Submission, task.TaskType, task.Deadline, task.ClassID).Scan(&taskID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	query = "SELECT user_id FROM user_classes WHERE class_id = $1"
	rowClasses, err := r.conn.Query(ctx, query, task.ClassID)
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

	userTasks := make([]entity.UserTask, 0, 10)
	for rowClasses.Next() {
		task := entity.UserTask{}
		if err := rowClasses.Scan(&task.UserID); err != nil {
			return err
		}
		task.TaskID = taskID
		task.Status = false
		userTasks = append(userTasks, task)
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

func (r *taskRepo) MarkTaskDone(ctx context.Context, userID string, taskID string) error {
	q2 := `UPDATE user_tasks SET status = true WHERE user_id = $1 AND task_id = $2`

	_, err := r.conn.Exec(ctx, q2, userID, taskID)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		return err
	}

	return nil
}

func (r *taskRepo) GetMyTasks(ctx context.Context, sub, classID string) ([]entity.TaskWithStatus, error) {
	rowsTask, err := r.conn.Query(ctx, `SELECT t.id, t.title, t.description, t.task_type, t.deadline, ut.status FROM tasks t JOIN user_tasks ut ON t.id = ut.task_id WHERE t.class_id = $1 AND ut.user_id = $2`, classID, sub)
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

	tasks := make([]entity.TaskWithStatus, 0, 10)
	for rowsTask.Next() {
		task := entity.TaskWithStatus{}
		if err := rowsTask.Scan(&task.ID, &task.Title, &task.Description, &task.TaskType, &task.Deadline, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepo) GetTask(ctx context.Context, sub, taskID string) (entity.TaskWithStatus, error) {
	task := entity.TaskWithStatus{}
	err := r.conn.QueryRow(ctx, `SELECT t.id, t.title, t.description, t.task_type, t.deadline, t.detail, t.submission, t.class_id, ut.status FROM tasks t JOIN user_tasks ut ON ut.task_id = t.id WHERE t.id = $1 AND ut.task_id = $1 AND ut.user_id = $2`,
		taskID, sub).Scan(&task.ID, &task.Title, &task.Description, &task.TaskType, &task.Deadline, &task.Detail, &task.Submission, &task.ClassID, &task.Status)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return task, ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return task, ierr.ErrNotFound
			}
		}
		return task, err
	}

	return task, nil
}
