package main

import (
	"b1-taskday7/connection"
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	connection.DataBaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/addProject", addProject).Methods("GET")
	route.HandleFunc("/addProject", addProjectPost).Methods("POST")
	route.HandleFunc("/contactMe", contactMe).Methods("GET")
	route.HandleFunc("/addContactMe", contactMePost).Methods("POST")
	route.HandleFunc("/projectDetail/{index}", projectDetail).Methods("GET")
	route.HandleFunc("/editProject/{index}", editProject).Methods("GET")
	// route.HandleFunc("/update-project/{Id}", submitEdit).Methods("POST")
	route.HandleFunc("/deleteProject/{index}", deleteProject).Methods("GET")

	fmt.Println("server running on port 5000")
	http.ListenAndServe("localhost:5000", route)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("views/home.html")

	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name_project, start_date, end_date, description, duration FROM public.tb_project")

	var result []Project
	for data.Next() {
		var each = Project{}
		err := data.Scan(&each.Id, &each.NameProject, &each.StartDate, &each.EndDate, &each.Description, &each.Duration)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)
	}
	resData := map[string]interface{}{
		"Projects": result,
	}

	tmpl.Execute(w, resData)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("views/add-my-project.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

var dataProject = []Project{}

type Project struct {
	NameProject  string
	Duration     string
	StartDate    string
	EndDate      string
	Description  string
	Technologies string
	Reactjs      string
	Javascript   string
	Golang       string
	Nodejs       string
	Image        string
	Id           int
}

func addProjectPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var nameProject = r.PostForm.Get("input-nameProject")
	var description = r.PostForm.Get("description")
	var startDate = r.PostForm.Get("input-startDate")
	var endDate = r.PostForm.Get("input-endDate")
	var reactjs = r.PostForm.Get("react")
	var javascript = r.PostForm.Get("javascript")
	var golang = r.PostForm.Get("golang")
	var nodejs = r.PostForm.Get("nodejs")
	var image = r.PostForm.Get("input-image")

	timePost, _ := time.Parse("2006-01-02", startDate)
	timeNow, _ := time.Parse("2006-01-02", endDate)
	println(timeNow.String())
	println(timePost.String())

	hours := timeNow.Sub(timePost).Hours()
	days := hours / 24
	weeks := math.Floor(days / 7)
	months := math.Floor(days / 30)
	years := math.Floor(days / 365)

	var duration string

	if years > 0 {
		duration = strconv.FormatFloat(years, 'f', 0, 64) + " Years"
	} else if months > 0 {
		duration = strconv.FormatFloat(months, 'f', 0, 64) + " Months"
	} else if weeks > 0 {
		duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " Weeks"
	} else if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " Days"
	} else if hours > 0 {
		duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hours"
	} else {
		duration = "0 Days"
	}
	println(hours)

	var newPoject = Project{
		NameProject: nameProject,
		Duration:    duration,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		Reactjs:     reactjs,
		Javascript:  javascript,
		Golang:      golang,
		Nodejs:      nodejs,
		Image:       image,
		Id:          len(dataProject),
	}
	fmt.Println(newPoject)

	dataProject = append(dataProject, newPoject)
	// fmt.Println(dataProject)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("views/my-project-detail.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	var ProjectDetail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if index == i {
			ProjectDetail = Project{
				NameProject: data.NameProject,
				Duration:    data.Duration,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Description: data.Description,
				Reactjs:     data.Reactjs,
				Javascript:  data.Javascript,
				Golang:      data.Golang,
				Nodejs:      data.Nodejs,
				Image:       data.Image,
			}
		}
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpl.Execute(w, data)
}

func contactMePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("name : " + r.PostForm.Get("input-name"))
	fmt.Println("email : " + r.PostForm.Get("input-email"))
	fmt.Println("phoneNumber : " + r.PostForm.Get("input-phonenumber"))
	fmt.Println("subject : " + r.PostForm.Get("input-subject"))
	fmt.Println("message : " + r.PostForm.Get("input-yourmessage"))

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/update-my-project.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}
	var projectDetail = Project{}
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	for i, data := range dataProject {
		if index == i {
			projectDetail = Project{
				NameProject: data.NameProject,
				Duration:    data.Duration,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Description: data.Description,
				Reactjs:     data.Reactjs,
				Javascript:  data.Javascript,
				Golang:      data.Golang,
				Nodejs:      data.Nodejs,
				Image:       data.Image,
				Id:          data.Id,
			}

		}
	}
	data := map[string]interface{}{
		"editProject": projectDetail,
	}
	tmpl.Execute(w, data)
}

// func submitEdit(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var nameProject = r.PostForm.Get("input-nameProject")
// 	var description = r.PostForm.Get("description")
// 	var startDate = r.PostForm.Get("input-startDate")
// 	var endDate = r.PostForm.Get("input-endDate")
// 	var reactjs = r.PostForm.Get("react")
// 	var javascript = r.PostForm.Get("javascript")
// 	var golang = r.PostForm.Get("golang")
// 	var nodejs = r.PostForm.Get("nodejs")
// 	var image = r.PostForm.Get("input-image")

// 	timePost, _ := time.Parse("2006-01-02", startDate)
// 	timeNow, _ := time.Parse("2006-01-02", endDate)
// 	println(timeNow.String())
// 	println(timePost.String())

// 	hours := timeNow.Sub(timePost).Hours()
// 	days := hours / 24
// 	weeks := math.Floor(days / 7)
// 	months := math.Floor(days / 30)
// 	years := math.Floor(days / 365)

// 	var duration string

// 	if years > 0 {
// 		duration = strconv.FormatFloat(years, 'f', 0, 64) + " Years"
// 	} else if months > 0 {
// 		duration = strconv.FormatFloat(months, 'f', 0, 64) + " Months"
// 	} else if weeks > 0 {
// 		duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " Weeks"
// 	} else if days > 0 {
// 		duration = strconv.FormatFloat(days, 'f', 0, 64) + " Days"
// 	} else if hours > 0 {
// 		duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hours"
// 	} else {
// 		duration = "0 Days"
// 	}
// 	println(hours)

// 	var newPoject = Project{
// 		NameProject: nameProject,
// 		Duration:    duration,
// 		StartDate:   startDate,
// 		EndDate:     endDate,
// 		Description: description,
// 		Reactjs:     reactjs,
// 		Javascript:  javascript,
// 		Golang:      golang,
// 		Nodejs:      nodejs,
// 		Image:       image,
// 		Id:          len(dataProject),
// 	}
// 	fmt.Println(newPoject)

// 	dataProject = append(dataProject, newPoject)
// 	// fmt.Println(dataProject)

// 	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

// }

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataProject = append(dataProject[:index], dataProject[index+1:]...)
	http.Redirect(w, r, "/home", http.StatusFound)
}

// go run main.go
