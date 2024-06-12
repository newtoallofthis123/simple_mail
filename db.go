package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DbInstance struct {
	db *sql.DB
}

func NewDbInstance(connString string) *DbInstance {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	return &DbInstance{db: db}
}

func (d *DbInstance) Prep() error {
	query := `
	CREATE TABLE IF NOT EXISTS emails (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL,
		subject TEXT NOT NULL,
		body TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

type EmailEntry struct {
	Id        int    `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Body      string `json:"body,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type CreateEmailEntry struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (d *DbInstance) InsertEmail(e *CreateEmailEntry) error {
	query := `
	INSERT INTO emails (email, subject, body)
	VALUES ($1, $2, $3);
	`

	_, err := d.db.Exec(query, e.Email, e.Subject, e.Body)
	if err != nil {
		return err
	}

	return nil
}

func (d *DbInstance) GetEmails() ([]EmailEntry, error) {
	query := `
	SELECT id, email, subject, body, created_at
	FROM emails;
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	emails := []EmailEntry{}
	for rows.Next() {
		var email EmailEntry
		err := rows.Scan(&email.Id, &email.Email, &email.Subject, &email.Body, &email.CreatedAt)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

func (d *DbInstance) GetEmail(id string) (*EmailEntry, error) {
	query := `
	SELECT id, email, subject, body, created_at
	FROM emails
	WHERE id = $1;
	`

	var email EmailEntry
	err := d.db.QueryRow(query, id).Scan(&email.Id, &email.Email, &email.Subject, &email.Body, &email.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &email, nil
}

func (d *DbInstance) GetByEmail(email string) ([]EmailEntry, error) {
	query := `
	SELECT id, email, subject, body, created_at
	FROM emails
	WHERE email = $1;
	`

	rows, err := d.db.Query(query, email)
	if err != nil {
		return nil, err
	}

	emails := []EmailEntry{}
	for rows.Next() {
		var email EmailEntry
		err := rows.Scan(&email.Id, &email.Email, &email.Subject, &email.Body, &email.CreatedAt)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}
