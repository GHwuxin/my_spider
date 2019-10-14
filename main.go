package main

import (
	"mypractice/spider/app/controllers"

	"net/http"
)

func main() {

	controllers.NewTasksController("./config/tasks", 1).Run()
	// controllers.NewNormalController("D:/download/yyyyMMdd", "2019-10-09 00:00:00", "2019-10-10 23:59:00").Run()
	mux := http.NewServeMux()
	http.ListenAndServe(":6666", mux)
}
