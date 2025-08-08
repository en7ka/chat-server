package main

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbDSN = "host=localhost port=5433 dbname=data user=data-user password=note-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	//создаем пул соединений
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	//делаем запрос на вставку
	builderInsert := sq.Insert("data").
		PlaceholderFormat(sq.Dollar).
		Columns("username", "email").
		Values(gofakeit.Name(), gofakeit.Email()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
	}

	var noteId int64

	err = pool.QueryRow(ctx, query, args...).Scan(&noteId)
	if err != nil {
		log.Fatalf("failed to query row: %v", err)
	}

	log.Printf("note id: %d", noteId)

	//делаем запрос на выборку
	builderSelect := sq.Select("id", "username", "email", "created_at").
		From("data").
		PlaceholderFormat(sq.Dollar).
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to query rows: %v", err)
	}

	var id int
	var username, email string
	var createdAt time.Time

	for rows.Next() {
		err = rows.Scan(&id, &username, &email, &createdAt)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}
		log.Printf("id: %d, username: %s, email: %s, created_at: %s", id, username, email, createdAt)
	}

	//делаем записи на обновление данных
	builderUpdate := sq.Update("data").
		PlaceholderFormat(sq.Dollar).
		Set("username", gofakeit.Name()).
		Set("email", gofakeit.Email()).
		Set("created_at", time.Now()).
		Where(sq.Eq{"id": noteId})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
	}
	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update rows: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	////делаем запрос на удаление
	//builderDelete := sq.Delete("data").
	//	PlaceholderFormat(sq.Dollar).
	//	Where(sq.Eq{"id": noteId})
	//
	//query, args, err = builderDelete.ToSql()
	//if err != nil {
	//	log.Fatalf("failed to generate query: %v", err)
	//}
	//res, err = pool.Exec(ctx, query, args...)
	//if err != nil {
	//	log.Fatalf("failed to delete rows: %v", err)
	//}
	//
	//log.Printf("deleted %d rows", res.RowsAffected())
}
