package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SiteStatus struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}

var (
	sites     = make(map[string]string)
	sitesLock sync.Mutex
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{}
)

func main() {
	go monitorSites()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add", addSiteHandler)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
		<html>
		<head>
			<script>
				let socket = new WebSocket("ws://" + location.host + "/ws");
				socket.onmessage = function(event) {
					let list = document.getElementById("status-list");
					list.innerHTML = "";
					let data = JSON.parse(event.data);
					data.forEach(item => {
						let li = document.createElement("li");
						li.textContent = item.url + " - " + item.status;
						list.appendChild(li);
					});
				}
			</script>
		</head>
		<body>
			<h2>Pingador de Sites</h2>
			<form method="POST" action="/add">
				URL: <input name="url" />
				<input type="submit" value="Adicionar">
			</form>
			<h3>Status:</h3>
			<ul id="status-list"></ul>
		</body>
		</html>
	`
	t, _ := template.New("webpage").Parse(tmpl)
	t.Execute(w, nil)
}

func addSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		url := r.FormValue("url")
		if url != "" {
			sitesLock.Lock()
			sites[url] = "Aguardando verificação..."
			sitesLock.Unlock()
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func monitorSites() {
	for {
		sitesLock.Lock()
		for url := range sites {
			go func(u string) {
				status := "OK"
				resp, err := http.Get(u)
				if err != nil || resp.StatusCode != 200 {
					status = "FORA DO AR"
				}

				sitesLock.Lock()
				sites[u] = status
				sitesLock.Unlock()

				broadcastSites()
			}(url)
		}
		sitesLock.Unlock()

		time.Sleep(10 * time.Second)
	}
}

func broadcastSites() {
	sitesLock.Lock()
	var list []SiteStatus
	for url, status := range sites {
		list = append(list, SiteStatus{URL: url, Status: status})
	}
	sitesLock.Unlock()

	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		err := client.WriteJSON(list)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao fazer upgrade:", err)
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	sitesLock.Lock()
	var list []SiteStatus
	for url, status := range sites {
		list = append(list, SiteStatus{URL: url, Status: status})
	}
	sitesLock.Unlock()
	conn.WriteJSON(list)
}
