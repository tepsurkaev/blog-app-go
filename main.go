package main

// nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run ./main.go

import (
	"net/http"
	"fmt"
	"encoding/json"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gofrs/uuid/v5"
)

type Blog struct {
	ID 		uuid.UUID 	`json:"id"`
	Title 	string 		`json:"title"`
	Content string		`json:"content"`
}

type Blogs struct {
	Blogs []Blog `json:"blogs"`
}

var blogs = Blogs{[]Blog{}}

func (blogs *Blogs) FindBlogByID(blogID string) int {
	blogIDToUUID, _ := uuid.FromString(blogID)
	blogIDIdx := slices.IndexFunc(blogs.Blogs, func(blog Blog) bool {
		return blog.ID == blogIDToUUID
	})
	return blogIDIdx
}

func (blogs *Blogs) AddNewBlog(blog Blog, uid uuid.UUID) []Blog {
	blog.ID = uid
    blogs.Blogs = append(blogs.Blogs, blog)
    return blogs.Blogs
}

func (blogs *Blogs) RemoveBlogByID(blogID string) []Blog {
	blogIDIdx := blogs.FindBlogByID(blogID)
	if blogIDIdx == -1 {
		fmt.Println("There is no blog with provided id!")
	}
    blogs.Blogs = append(blogs.Blogs[:blogIDIdx], blogs.Blogs[blogIDIdx+1:]...)
    return blogs.Blogs
}

func main() {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
	    AllowedOrigins:   []string{"http://*"},
	    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	    ExposedHeaders:   []string{"Link"},
	    AllowCredentials: false,
	    MaxAge:           300,
	}))
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Get("/", func(response http.ResponseWriter, _ *http.Request) {
		response.Write([]byte("Hi!"))
	})

	router.Route("/blogs", func(r chi.Router) {
		r.Get("/", getAllBlogs)
		r.Get("/{blogID}", getBlogByID)
		r.Post("/", createBlog)
		r.Delete("/{blogID}", deleteBlogByID)
	})

	http.ListenAndServe(":8080", router)
}

func getAllBlogs(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	value, err := json.Marshal(blogs)
	if err != nil {
		fmt.Println(err)
		return
	}
	response.Write(value)
}

func getBlogByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	blogID := chi.URLParam(request, "blogID")
	blogIDIdx := blogs.FindBlogByID(blogID)
	if blogIDIdx == -1 {
		http.Error(response, http.StatusText(404), 404)
		return
	}

	blogByID := blogs.Blogs[blogIDIdx]
	value, err := json.Marshal(blogByID)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Write(value)
}

func createBlog(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var newBlog Blog
	uid, _ := uuid.NewV4()

	err := json.NewDecoder(request.Body).Decode(&newBlog)
    if err != nil {
        response.WriteHeader(http.StatusBadRequest)
        response.Write([]byte(err.Error()))
        return
    }

	blogs.AddNewBlog(newBlog, uid)
}

func deleteBlogByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	blogID := chi.URLParam(request, "blogID")
	blogs.RemoveBlogByID(blogID)
}
