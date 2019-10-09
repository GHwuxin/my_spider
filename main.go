package main

import (
	"mypractice/spider/app/controllers"
	"net/http"
)

func main() {

	controllers.NewTasksController("./config/tasks", 1).Run()

	// controllers.NewNormalController("./download/yyyyMMdd", "2019-05-15 00:00:00", "2019-10-01 23:59:00").Run()

	mux := http.NewServeMux()
	http.ListenAndServe(":6666", mux)
}
