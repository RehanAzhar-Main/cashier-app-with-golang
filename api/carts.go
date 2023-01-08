package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (api *API) AddCart(w http.ResponseWriter, r *http.Request) {

	var respErr model.ErrorResponse

	// Get username context to struct model.Cart. ???
	username := fmt.Sprintf("%s", r.Context().Value("username")) // TODO: replace this
	r.ParseForm()

	// Check r.Form with key product, if not found then return response code 400 and message "Request Product Not Found".
	// TODO: answer here
	var product string = r.FormValue("product")
	if product == "" {

		respErr.Error = "Request Product Not Found"
		// encode json
		jsonInBytes, err := json.Marshal(respErr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(400)
		w.Write(jsonInBytes)
		return
	}
	var totalPrice float64

	var list []model.Product
	for _, formList := range r.Form {
		for _, v := range formList {
			item := strings.Split(v, ",")
			p, _ := strconv.ParseFloat(item[2], 64)
			q, _ := strconv.ParseFloat(item[3], 64)
			total := p * q
			list = append(list, model.Product{
				Id:       item[0],
				Name:     item[1],
				Price:    item[2],
				Quantity: item[3],
				Total:    total,
			})

			totalPrice += total
		}
	}
	w.WriteHeader(200)

	// Add data field Name, Cart and TotalPrice with struct model.Cart.
	cart := model.Cart{
		Name:       username,
		Cart:       list,
		TotalPrice: totalPrice,
	} // TODO: replace this

	_, err := api.cartsRepo.CartUserExist(cart.Name)
	if err != nil {
		api.cartsRepo.AddCart(cart)
	} else {
		api.cartsRepo.UpdateCart(cart)
	}
	api.dashboardView(w, r)

}
