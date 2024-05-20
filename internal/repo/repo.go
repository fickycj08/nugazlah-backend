package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool

	User   *userRepo
	Class  *classRepo
	Friend *friendRepo
	Post   *postRepo
	Task   *taskRepo
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	repo := Repo{}
	repo.conn = conn

	repo.User = newUserRepo(conn)
	repo.Class = newClassRepo(conn)
	repo.Friend = newFriendRepo(conn)
	repo.Post = newPostRepo(conn)
	repo.Task = newTaskRepo(conn)

	return &repo
}
