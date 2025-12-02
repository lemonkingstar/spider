package pdebug

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/lemonkingstar/spider/pkg/iserver"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func system(w http.ResponseWriter, r *http.Request) {
	loads, err := load.Avg()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	memory, err := mem.VirtualMemory()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse(w, map[string]any{
		"Args":     strings.Join(os.Args, " "),
		"Load1m":   loads.Load1,
		"Load5m":   loads.Load5,
		"Load15m":  loads.Load15,
		"MemUsed":  memory.Used,
		"MemTotal": memory.Total,
		"Uptime":   time.Now().Unix(),
	})
}

func version(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, map[string]any{
		"app":            iserver.GetAppName(),
		"version":        iserver.GetVersion(),
		"build branch":   iserver.GetBuildBranch(),
		"build commit":   iserver.GetBuildCommit(),
		"build time":     iserver.GetBuildTime(),
		"golang version": runtime.Version(),
	})
}

func jsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
