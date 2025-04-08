package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	tmpl := template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Expenses Division</title>
</head>
<body>
	<h1>Expenses Division</h1>
	<button id="callBtn">View Expense Report</button>

	<script>
		document.getElementById("callBtn").addEventListener("click", () => {
			window.open("http://localhost:8000/", "_blank");
		});
	</script>
</body>
</html>
	`))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	})

	log.Println("Frontend running at http://view:3000/")
	if err := http.ListenAndServe("view:3000", nil); err != nil {
		log.Fatal(err)
	}
}