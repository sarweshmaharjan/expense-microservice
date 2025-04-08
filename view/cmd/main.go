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
	<title>Monthly Income Division</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 0;
			padding: 0;
			background-color: #f4f4f9;
			color: #333;
		}
		header {
			background-color: #4CAF50;
			color: white;
			padding: 1rem;
			text-align: center;
		}
		.container {
			padding: 2rem;
		}
		table {
			width: 100%;
			border-collapse: collapse;
			margin: 1rem 0;
		}
		table, th, td {
			border: 1px solid #ddd;
		}
		th, td {
			padding: 0.75rem;
			text-align: left;
		}
		th {
			background-color: #4CAF50;
			color: white;
		}
		.button {
			display: inline-block;
			padding: 0.5rem 1rem;
			margin-top: 1rem;
			background-color: #4CAF50;
			color: white;
			text-decoration: none;
			border-radius: 5px;
			cursor: pointer;
		}
		.button:hover {
			background-color: #45a049;
		}
	</style>
</head>
<body>
	<header>
		<h1>Monthly Income Division Per Expenses List</h1>
	</header>
	<div class="container">
		<p>Salary Currency: <span id="salaryCurrency">Loading...</span></p>
		<p>Current Salary: <span id="currentSalary">Loading...</span></p>
		<p>Cap Income Limit: <span id="capIncomeLimit">Loading...</span></p>
		<table>
			<thead>
				<tr>
					<th>Expense</th>
					<th>Expected Amount</th>
					<th>Is Fixed</th>
					<th>Min</th>
					<th>Max</th>
					<th>Type</th>
					<th>Active</th>
				</tr>
			</thead>
			<tbody id="expensesTable">
				<tr>
					<td colspan="7">Loading...</td>
				</tr>
			</tbody>
		</table>
		<button href="#" class="button" id="callBtn">View Income Division</button>
	</div>

	<script>
		async function fetchData() {
			try {
				const response = await fetch("http://localhost:8002");
				const data = await response.json();
				console.log(data);
				// Set salary details
				document.getElementById("salaryCurrency").textContent = data.SalaryCurrency;
				document.getElementById("currentSalary").textContent = data.CurrentSalary;
				document.getElementById("capIncomeLimit").textContent = data.CapIncomeLimit;

				// Populate expenses table
				const expensesTable = document.getElementById("expensesTable");
				expensesTable.innerHTML = ""; // Clear existing rows
				data.Expenses.forEach(expense => {
					const row = document.createElement("tr");
					const expenseCell = document.createElement("td");
					const amountCell = document.createElement("td");
					const isFixedCell = document.createElement("td");
					const minCell = document.createElement("td");
					const maxCell = document.createElement("td");
					const typeCell = document.createElement("td");
					const activeCell = document.createElement("td");

					expenseCell.textContent = expense.Name;
					amountCell.textContent = expense.ExpectedAmount || "N/A";
					isFixedCell.textContent = expense.IsFixed ? "Yes" : "No";
					minCell.textContent = expense.Min;
					maxCell.textContent = expense.Max;
					typeCell.textContent = expense.Type;
					activeCell.textContent = expense.Active ? "Yes" : "No";

					row.appendChild(expenseCell);
					row.appendChild(amountCell);
					row.appendChild(isFixedCell);
					row.appendChild(minCell);
					row.appendChild(maxCell);
					row.appendChild(typeCell);
					row.appendChild(activeCell);
					expensesTable.appendChild(row);
				});
			} catch (error) {
				console.error("Error fetching data:", error);
			}
		}

		document.getElementById("callBtn").addEventListener("click", () => {
			window.open("http://localhost:8000", "_blank");
		});

		// Fetch data on page load
		fetchData();
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