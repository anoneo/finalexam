package customerhandler

type Todo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// func todosHandler(w http.ResponseWriter, req *http.Request) {
// 	if req.Method == "DELETE" {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}
// 	w.Write([]byte("hello get todos"))
// }
