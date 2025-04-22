package handlers

import (
	"fmt"
	"net/http"
)

func HandleFrontend(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <title>Распределённый калькулятор</title>
  <!-- Bootstrap CSS из CDN -->
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
  <style>
    body { padding-top: 40px; }
    .spinner-border { width: 3rem; height: 3rem; }
  </style>
</head>
<body>
  <div class="container">
    <div class="row justify-content-center">
      <div class="col-md-8">
        <div class="card shadow">
          <div class="card-header bg-primary text-white">
            <h3 class="card-title mb-0">Распределённый калькулятор выражений</h3>
          </div>
          <div class="card-body">
            <form id="calcForm">
              <div class="mb-3">
                <label for="expression" class="form-label">Введите арифметическое выражение</label>
                <input type="text" class="form-control" id="expression" placeholder="2+2*(2+5)*3" required>
              </div>
              <button type="submit" class="btn btn-primary">Вычислить</button>
            </form>
            <hr>
            <div id="result" class="mt-3"></div>
            <div id="status" class="mt-3"></div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Bootstrap JS и зависимости (Popper) -->
  <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
  <script>
    document.getElementById("calcForm").addEventListener("submit", function(e) {
      e.preventDefault();
      var expr = document.getElementById("expression").value;
      document.getElementById("result").innerHTML = "";
      document.getElementById("status").innerHTML = "<div class='d-flex align-items-center'><strong>Отправка запроса...</strong><div class='spinner-border ms-2 text-primary' role='status'><span class='visually-hidden'>Загрузка...</span></div></div>";
      fetch("/api/v1/calculate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ expression: expr })
      })
      .then(res => res.json())
      .then(data => {
        document.getElementById("result").innerHTML = "<div class='alert alert-success'>Задача принята. ID: " + data.id + "</div>";
        pollExpressionStatus(data.id);
      })
      .catch(err => {
        console.error("Ошибка отправки запроса", err);
        document.getElementById("result").innerHTML = "<div class='alert alert-danger'>Ошибка отправки запроса</div>";
        document.getElementById("status").innerHTML = "";
      });
    });

    function pollExpressionStatus(id) {
      var pollInterval = setInterval(function() {
        fetch("/api/v1/expressions/" + id)
          .then(res => res.json())
          .then(data => {
            if(data.expression.status === "completed") {
              document.getElementById("status").innerHTML = "<div class='alert alert-success'>Результат: " + data.expression.result + "</div>";
              clearInterval(pollInterval);
            } else {
              document.getElementById("status").innerHTML = "<div class='alert alert-info'>Статус: " + data.expression.status + "</div>";
            }
          })
          .catch(err => {
            console.error("Ошибка получения статуса", err);
            clearInterval(pollInterval);
            document.getElementById("status").innerHTML = "<div class='alert alert-danger'>Ошибка получения статуса</div>";
          });
      }, 2000);
    }
  </script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}
