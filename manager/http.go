package manager

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type HttpConfig struct {
	Port     uint
	External bool
}

func (m *Manager) ServeHTTP(conf HttpConfig) {
	if conf.Port == 0 {
		return
	}

	m.logger.Infoln("Serving api on port:", conf.Port)
	server := http.NewServeMux()
	server.HandleFunc("/", m.rootHandler)
	server.HandleFunc("/shares", m.sharesHandler)
	server.HandleFunc("/hashrates", m.hashratesHandler)
	server.HandleFunc("/hashrates/total", m.totalHashratesHandler)

	if conf.External {
		go http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(conf.Port)), server)
	} else {
		go http.ListenAndServe("127.0.0.1:"+strconv.Itoa(int(conf.Port)), server)
	}
}

func (m *Manager) rootHandler(w http.ResponseWriter, r *http.Request) {
	raw := []byte(`
/shares => all mined shares' status
/hashrates => realtime hashrates
`)
	_, _ = w.Write(raw)
}

func (m *Manager) sharesHandler(w http.ResponseWriter, r *http.Request) {
	shares := map[string]uint64{
		"accept": m.client.AcceptNum + m.dClient.AcceptNum,
		"reject": m.client.RejectNum + m.dClient.RejectNum,
	}

	shares["total"] = shares["accept"] + shares["reject"]
	raw, _ := json.Marshal(shares)
	_, _ = w.Write(raw)
}

func (m *Manager) hashratesHandler(w http.ResponseWriter, r *http.Request) {
	hrs := m.client.Rx.GetWorkerHashrates()
	dhrs := m.dClient.Rx.GetWorkerHashrates()
	for k, v := range dhrs {
		if v != 0 {
			hrs[k] += v
		}
	}

	raw, _ := json.Marshal(hrs)
	_, _ = w.Write(raw)
}

func (m *Manager) totalHashratesHandler(w http.ResponseWriter, r *http.Request) {
	thr := map[string]float64{
		"total": 0,
	}
	hrs := m.client.Rx.GetWorkerHashrates()
	dhrs := m.dClient.Rx.GetWorkerHashrates()
	for _, v := range hrs {
		thr["total"] += v
	}

	for _, v := range dhrs {
		thr["total"] += v
	}

	raw, _ := json.Marshal(thr)
	_, _ = w.Write(raw)
}
