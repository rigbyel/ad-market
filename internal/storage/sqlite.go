package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/rigbyel/ad-market/internal/models"
)

type Storage struct {
	db *sql.DB
}

// creates new instance of Storage
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return &Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return &Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// stops storage work
func (s *Storage) Stop() error {
	return s.db.Close()
}

// saves user in storage
func (s *Storage) SaveUser(u *models.User) (*models.User, error) {
	const op = "storage.sqlite.SaveUser"

	// prepare query
	stmt, err := s.db.Prepare(
		"INSERT INTO users (login, passHash) VALUES ($1, $2)",
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// execute query
	res, err := stmt.Exec(u.Login, u.PassHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// get id of created user
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u.Id = id

	return u, nil
}

// gets user with the given login from storage
func (s *Storage) User(login string) (*models.User, error) {
	const op = "storage.sqlite.User"

	// get user's passhash from database
	row := s.db.QueryRow(
		"SELECT passHash FROM users WHERE login = $1",
		login,
	)

	var passHash []byte
	if err := row.Scan(&passHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// get user id from database
	row = s.db.QueryRow(
		"SELECT id FROM users WHERE login = $1",
		login,
	)

	var id int64
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// create models.User struct with the given login, id and passHash
	user := &models.User{
		Id:       id,
		Login:    login,
		PassHash: passHash,
	}

	return user, nil
}

// saves advert in the storage
func (s *Storage) SaveAd(ad *models.Advert) (*models.Advert, error) {
	const op = "storage.sqlite.SaveAd"

	// prepare query
	stmt, err := s.db.Prepare(
		`INSERT INTO adverts (header, body, imageURL, price, date, authorLogin)
		VAlUES ($1, $2, $3, $4, $5, $6)`,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// execute query
	res, err := stmt.Exec(ad.Header, ad.Body, ad.ImageURL, ad.Price, ad.Date, ad.AuthorLogin)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// get id of created advert
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ad.Id = id

	return ad, nil
}

// gets all adverts within given price range from storage
func (s *Storage) Adverts(minPrice, maxPrice int) (*[]models.Advert, error) {
	const op = "storage.sqlite.Adverts"

	var adverts []models.Advert

	// get adverts from database
	rows, err := s.db.Query(
		"SELECT * FROM adverts WHERE price >= $1 AND price <= $2",
		minPrice,
		maxPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var ad models.Advert

		err = rows.Scan(&ad.Id, &ad.Header, &ad.Body, &ad.ImageURL, &ad.Price, &ad.Date, &ad.AuthorLogin)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		adverts = append(adverts, ad)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &adverts, nil
}
