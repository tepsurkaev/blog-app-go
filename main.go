package main

// nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run ./main.go

import (
	"net/http"
	"fmt"
	"encoding/json"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Blog struct {
	ID 		string 	`json:"id"`
	Title 	string 	`json:"title"`
	Content string 	`json:"content"`
}

var blogs = []Blog{
	{ID: "1", Title: "Blue Train", Content: "John Coltrane"},
	{ID: "2", Title: "Red Train", Content: "John Deax"},
	{ID: "3", Title: "Green Train", Content: "Blaine DeBears"},
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hi!"))
	})

	r.Route("/blogs", func(r chi.Router) {
		r.Get("/", getAllBlogs)
		r.Get("/{blogID}", getBlogById)
	})

	http.ListenAndServe(":8080", r)
}

func getAllBlogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	value, err := json.Marshal(blogs)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(value)
}

func getBlogById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	blogID := chi.URLParam(r, "blogID")
	blogIDIdx := slices.IndexFunc(blogs, func(blog Blog) bool {
		return blog.ID == blogID
	})
	if blogIDIdx == -1 {
		fmt.Println("There is no blog with provided id")
		return
	}

	blogByID := blogs[blogIDIdx]
	value, err := json.Marshal(blogByID)

	if err != nil {
		fmt.Println(err)
	}

	w.Write(value)
}

func createBlog(w http.ResponseWriter, r *http.Request) {}

func deleteBlog(w http.ResponseWriter, r *http.Request) {}
