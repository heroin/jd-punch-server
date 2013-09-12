package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	SERVER = "JD-PUNCH-SERVER"
	PORT   = 12324
)

var (
	TASK_LAST_MODIFY_DATE time.Time
	TASK_DATA             = &Task{}
)

type Task struct {
	Start  bool     `json:"start"`
	Users  []*User  `json:"users,omitempty"`
	Cancel []string `json:"cancel,omitempty"`
}

type User struct {
	Id       int64  `json:"id,omitempty"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
	Start    bool   `json:"start"`
	Trigger  int64  `json:"trigger"`
}

type Context struct {
	data interface{}
}

func (context *Context) view(write io.Writer, html string) {
	cursor := template.Must(template.ParseFiles(fmt.Sprintf("views/%s.html", html)))
	cursor.Execute(write, context.data)
}

func index(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	requestPath := request.URL.Path
	if requestPath != "/" {
		n := len(requestPath)
		http.ServeFile(out, request, fmt.Sprintf("static/%s", requestPath[1:n]))
	} else {
		app := Context{}
		app.view(out, "index")
	}
}

func task(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	out.Header().Set("Last-Modifyed", TASK_LAST_MODIFY_DATE.Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	data, err_json := json.Marshal(TASK_DATA)
	if err_json != nil {
	} else {
		fmt.Fprintf(out, string(data))
	}
}

func user_add(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	request.ParseForm()
	username := request.FormValue("username")
	password := request.FormValue("password")
	start, _ := strconv.ParseBool(request.FormValue("start"))
	trigger, _ := strconv.ParseInt(request.FormValue("trigger"), 10, 64)
	if strings.TrimSpace(username) != "" && strings.TrimSpace(password) != "" {
		TASK_DATA.Users = append(TASK_DATA.Users, &User{
			UserName: username,
			PassWord: password,
			Trigger:  trigger,
			Start:    start,
		})
		TASK_LAST_MODIFY_DATE = time.Now()
		fmt.Fprintf(out, "success")
	} else {
		fmt.Fprintf(out, "error")
	}
}

func user_del(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	request.ParseForm()
	username := request.FormValue("username")
	trigger, _ := strconv.ParseInt(request.FormValue("trigger"), 10, 64)
	if strings.TrimSpace(username) != "" {
		count := 0
		for _, value := range TASK_DATA.Users {
			if fmt.Sprintf("%s-%d", username, trigger) == fmt.Sprintf("%s-%d", value.UserName, value.Trigger) {
				count++
			}
		}
		if count > 0 {
			for t := 0; t < count; t++ {
				for i := range TASK_DATA.Users {
					if fmt.Sprintf("%s-%d", username, trigger) == fmt.Sprintf("%s-%d", TASK_DATA.Users[i].UserName, TASK_DATA.Users[i].Trigger) {
						TASK_DATA.Users = append(TASK_DATA.Users[:i], TASK_DATA.Users[i+1:]...)
						break
					}
				}
			}
		}
		TASK_LAST_MODIFY_DATE = time.Now()
		fmt.Fprintf(out, "success")
	} else {
		fmt.Fprintf(out, "error")
	}
}

func cancel_add(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	request.ParseForm()
	username := request.FormValue("username")
	trigger, _ := strconv.ParseInt(request.FormValue("trigger"), 10, 64)
	if strings.TrimSpace(username) != "" {
		TASK_DATA.Cancel = append(TASK_DATA.Cancel, fmt.Sprintf("%s-%d", username, trigger))
		TASK_LAST_MODIFY_DATE = time.Now()
		fmt.Fprintf(out, "success")
	} else {
		fmt.Fprintf(out, "error")
	}
}

func cancel_del(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	request.ParseForm()
	username := request.FormValue("username")
	trigger, _ := strconv.ParseInt(request.FormValue("trigger"), 10, 64)
	if strings.TrimSpace(username) != "" {
		count := 0
		for _, value := range TASK_DATA.Cancel {
			if fmt.Sprintf("%s-%d", username, trigger) == value {
				count++
			}
		}
		if count > 0 {
			for t := 0; t < count; t++ {
				for i := range TASK_DATA.Cancel {
					if fmt.Sprintf("%s-%d", username, trigger) == TASK_DATA.Cancel[i] {
						TASK_DATA.Cancel = append(TASK_DATA.Cancel[:i], TASK_DATA.Cancel[i+1:]...)
						break
					}
				}
			}
		}
		TASK_LAST_MODIFY_DATE = time.Now()
		fmt.Fprintf(out, "success")
	} else {
		fmt.Fprintf(out, "error")
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	TASK_DATA.Start = true
	http.HandleFunc("/", index)
	http.HandleFunc("/task.json", task)
	http.HandleFunc("/user/add", user_add)
	http.HandleFunc("/user/del", user_del)
	http.HandleFunc("/cancel/add", cancel_add)
	http.HandleFunc("/cancel/del", cancel_del)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
