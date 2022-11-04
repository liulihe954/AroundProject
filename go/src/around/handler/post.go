package handler

import (
	"encoding/json"
	"net/http"
	"fmt"
	"path/filepath"
	"around/model"
	"around/service"
	jwt "github.com/form3tech-oss/jwt-go"
    "github.com/gorilla/mux"
	"github.com/pborman/uuid"
)


var (
	mediaTypes = map[string]string{
		".jpeg": "image",
		".jpg": "image",
		".gif": "image",
		".png": "image",
		".mov": "video",
		".mp4": "video",
		".avi": "video",
		".flv": "video",
		".wmv": "video",
	}
)


func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// parse from body of request to get a json object
	fmt.Println("Received one post request")

	user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

	p := model.Post{
		Id: uuid.New(),
		User: username.(string),
		Message: r.FormValue("message"),
	}

	file, header, err := r.FormFile("media_file")
	
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v \n", err)
		return 
	}

	suffix := filepath.Ext(header.Filename)
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}

	err = service.SavePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to backend %v \n", err)
		return
	}
	fmt.Println("Post is saved successfully.")


	// decoder := json.NewDecoder(r.Body)
	// var p model.Post

	// if err := decoder.Decode(&p); err != nil {
	// 	panic(err)
	// }
	// fmt.Fprintf(w, "Post received: %s\n", p.Message)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request from search")
	w.Header().Set("Content-Type", "application/json")

	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	// fmt.Printf("search by user received: %v\n", user)
	// fmt.Printf("search by keywords received: %v\n", keywords)

	var posts []model.Post
	var err error
	if user != "" {
		posts, err = service.SearchPostsByUser(user)
	} else {
		posts, err = service.SearchPostsByKeywords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		fmt.Printf("Failed to read post form backend %v. \n", err)
		return
	}

	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v. \n", err)
	}
	w.Write(js)
}


func deleteHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for delete")

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"].(string)
    id := mux.Vars(r)["id"]

    if err := service.DeletePost(id, username); err != nil {
        http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to delete post from backend %v\n", err)
        return
    }
    fmt.Println("Post is deleted successfully")
}
