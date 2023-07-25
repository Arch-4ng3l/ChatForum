package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Arch-4ng3l/ChatForum/types"
	"github.com/Arch-4ng3l/ChatForum/util"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgresDB() *Postgres {
	connStr := "user=moritz dbname=chat password=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil
	}

	if err := db.Ping(); err != nil {
		return nil
	}

	psql := &Postgres{
		db: db,
	}
	if err := psql.Init(); err != nil {
		return nil
	}
	fmt.Println("Table Created")
	return psql
}

func (psql *Postgres) Init() error {

	query := `CREATE TABLE IF NOT EXISTS users(
		id serial PRIMARY KEY,
		username TEXT UNIQUE,
		email TEXT UNIQUE,
		password TEXT,
		created_at TIMESTAMP
	)`

	query2 := `CREATE TABLE IF NOT EXISTS messages(
		msgid serial PRIMARY KEY,
		sender TEXT,
		receiver TEXT,
		created_at TIMESTAMP,
		content TEXT, 
		edited INT
	)`
	fmt.Println("Creating Table")

	_, err := psql.db.Exec(query)

	_, err = psql.db.Exec(query2)

	return err
}

func (psql *Postgres) CreateNewUser(req *types.SignUpRequest) error {
	query := ` INSERT INTO users (username, email, password, created_at) VALUES($1, $2, $3, $4)`

	passwd := util.CreateHash(req.Password)

	res, err := psql.db.Exec(query, req.Name, req.Email, passwd, time.Now())

	if err != nil {
		return err
	}

	n, _ := res.RowsAffected()

	if n != 1 {
		return fmt.Errorf("User Already Exists")
	}

	return nil
}

func (psql *Postgres) GetUserPassword(req *types.LoginRequest) (string, error) {
	query := `SELECT password FROM users WHERE username=$1`

	row, err := psql.db.Query(query, req.Name)

	if err != nil {
		return "", err
	}
	passwd := ""
	row.Next()
	row.Scan(&passwd)
	row.Close()
	return passwd, nil
}

func (psql *Postgres) CreateNewMessage(req *types.CreateMessageRequest) error {

	query := `
	INSERT INTO messages (sender, receiver, content, created_at)
	SELECT $1::TEXT, $2::TEXT, $3::TEXT, $4
	WHERE EXISTS (SELECT 1 FROM users WHERE username = $1)
		AND EXISTS (SELECT 1 FROM users WHERE username = $2)
	`

	res, err := psql.db.Exec(query, req.Name, req.Receiver, req.Message, time.Now())

	n, _ := res.RowsAffected()
	fmt.Println(n)
	return err
}
func (psql *Postgres) UpdateMessage(req *types.UpdateMessageRequest) error {

	query := `
		UPDATE messages SET content=$1, edited=1 WHERE msgid=$2	
	`

	_, err := psql.db.Exec(query, req.NewContent, req.MessageID)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (psql *Postgres) GetMessages(req *types.GetMessagesRequest) ([]*types.Message, error) {
	var msgs []*types.Message

	query := `
		SELECT msgid, content, created_at FROM messages
		WHERE (sender=$1 AND receiver=$2) OR (sender=$2 AND receiver=$1)
	`

	rows, err := psql.db.Query(query, req.Sender, req.Receiver)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		msg := &types.Message{}
		rows.Scan(&msg.ID, &msg.Message, &msg.CreatedAt)
		msgs = append(msgs, msg)
	}

	rows.Close()
	for _, n := range msgs {
		fmt.Println(n.Message)
	}
	return msgs, nil
}
