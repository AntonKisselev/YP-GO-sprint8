package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, err
	}
	number, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(number), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE number = ?", number)
	var (
		client     int
		status     string
		address    string
		created_at string
	)
	err := row.Scan(&client, &status, &address, &created_at)
	if err != nil {
		return Parcel{}, err
	}
	// заполните объект Parcel данными из таблицы
	p := Parcel{
		Number:    number,
		Client:    client,
		Status:    status,
		Address:   address,
		CreatedAt: created_at,
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT number, status, address, created_at FROM parcel WHERE client = ?", client)
	if err != nil {
		fmt.Println(err)
		return []Parcel{}, err
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		var (
			number     int
			status     string
			address    string
			created_at string
		)
		err := rows.Scan(&number, &status, &address, &created_at)
		if err != nil {
			fmt.Println(err)
			return []Parcel{}, err
		}
		res = append(res, Parcel{
			Number:    number,
			Client:    client,
			Status:    status,
			Address:   address,
			CreatedAt: created_at,
		})
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}
	if parcel.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number),
		)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("parcel already registered")
	}
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}
	if parcel.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("DELETE FROM parcel WHERE number = ?", number)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("parcel already registered")
	}
}
