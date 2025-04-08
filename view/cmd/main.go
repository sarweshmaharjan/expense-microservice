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
	<title>Call Backend</title>
</head>
<body>
	<h1>Frontend (Port 3000)</h1>
	<button id="callBtn">Call Backend</button>

	<script>
		document.getElementById("callBtn").addEventListener("click", async () => {
			try {
				const res = await fetch("http://hub:8000/", { method: "GET" });
				if (!res.ok) throw new Error("Network response was not ok");

				const data = await res.text();
				alert("Response from backend: " + data);
			} catch (err) {
				alert("Error calling backend: " + err.message);
			}
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