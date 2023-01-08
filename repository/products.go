package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
)

type ProductRepository struct {
	db db.DB
}

func NewProductRepository(db db.DB) ProductRepository {
	return ProductRepository{db}
}

func (u *ProductRepository) ReadProducts() ([]model.Product, error) {
	// // read file JSON
	// var jsonData, err = ioutil.ReadFile("user.json")
	// if err != nil {
	// 	panic(err)
	// }

	//load data
	jsonData, err := u.db.Load("products")
	if err != nil {
		return []model.Product{}, err
	}

	// decode JSON ke struct
	var listProduct []model.Product
	err = json.Unmarshal(jsonData, &listProduct)
	if err != nil {
		panic(err)
	}

	return listProduct, err // TODO: replace this
}

func (u *ProductRepository) ResetProducts() error {
	err := u.db.Reset("products", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}
