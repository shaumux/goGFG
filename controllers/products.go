package controllers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/olivere/elastic/v7"
	"goGFG/storage"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Product struct {
	storage.Product
}

func (p *Product) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

func (p *Product) Bind(r *http.Request) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return nil
}

type Products struct {
	Products []*storage.Product
}

func (p *Products) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

func (p *Products) MarshalJSON() ([]byte, error) {
	if p.Products != nil {
		return json.Marshal(p.Products)
	} else {
		products := make([]*storage.Product, 0)
		return json.Marshal(products)
	}
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	sizeStr := r.URL.Query().Get("size")
	size := new(int)
	var err error
	if sizeStr != "" {
		*size, err = strconv.Atoi(sizeStr)
		if err != nil {
			render.Render(w, r, ErrBadRequest(errors.New("Invalid size")))
			return
		}
	} else {
		size = nil
	}
	pageStr := r.URL.Query().Get("page")
	page := new(int)
	if pageStr != "" {
		*page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			render.Render(w, r, ErrBadRequest(errors.New("Invalid page")))
			return
		}
	} else {
		page = nil
	}

	search := r.URL.Query().Get("q")
	filter := r.URL.Query().Get("filter")

	search = strings.TrimSpace(search)
	filter = strings.TrimSpace(filter)
	searchQuery := &elastic.MultiMatchQuery{}
	filterQuery := &elastic.TermQuery{}
	if search != "" {
		searchQuery = storage.SearchProduct([]string{"title", "brand"}, search)
	} else {
		searchQuery = nil
	}
	if filter != "" {
		filterParams := strings.Split(strings.TrimSpace(filter), ":")
		filterQuery = storage.FilterProduct(filterParams[0], filterParams[1])
	} else {
		filterQuery = nil
	}

	sortField := r.URL.Query().Get("sort")
	searchResults, err := storage.GetProducts(sortField, page, size, filterQuery, searchQuery)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	var ptyp storage.Product
	searchedProds := &Products{}
	products := searchedProds.Products
	for _, item := range searchResults.Each(reflect.TypeOf(ptyp)) {
		if p, ok := item.(storage.Product); ok {
			products = append(products, &p)
		}
	}
	searchedProds.Products = products
	render.Render(w, r, searchedProds)
}

func IndexProduct(w http.ResponseWriter, r *http.Request) {
	product := &Product{}
	if err := render.Bind(r, product); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	if err := product.IndexProduct(); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	render.Render(w, r, product)
}
