package envoy

import (
	"encoding/json"
	"strings"
)

type Server struct {
	i *Instance
	endpoints ServerEndpoints
}

type ServerEndpoints struct {
	info string
	memory string
}
func (i *Instance) Server() *Server {
	return &Server{
		i: i,
		endpoints: ServerEndpoints{
			info: i.Address + "/server_info",
			memory: i.Address + "/memory",
		},
	}
}

func (l *Server) Info() (ServerInfoJson, error) {
	rawJson, err := l.i.nsenter.Curl("-s", l.endpoints.info)
	if err != nil { return ServerInfoJson{}, err }

	jsonData := []byte(strings.Trim(rawJson, " "))

	var response ServerInfoJson
	err = json.Unmarshal(jsonData, &response)
	if err != nil { return ServerInfoJson{}, err }

	return response, nil
}

type ServerInfoJson struct {
	Version string 				`json:"version"`
	State string 				`json:"state"`
	HotRestartVersion string	`json:"hot_restart_version"`
	ServerCommandLineOptions ServerCommandLineOptionsJson `json:"command_line_options"`
	Node ServerNode				`json:"node"`
	UptimeCurrentEpoch string   `json:"uptime_current_epoch"`
	UptimeAllEpochs string		`json:"uptime_all_epochs"`
}

type ServerCommandLineOptionsJson struct {
	Concurrency int             `json:"concurrency"`
	ConfigPath string			`json:"config_path"`
	LogLevel string				`json:"log_level"`
	ComponentLogLevel string    `json:"component_log_level"`
	LogFormat string            `json:"log_format"`
	LogFormatEscaped bool       `json:"log_format_escaped"`
	Mode string					`json:"mode"`
	DrainStrategy string        `json:"drain_strategy"`
	DrainTime string            `json:"drain_time"`
	ParentShutdownTime string   `json:"parent_shutdown_time"`
}

type ServerNode struct {
	ID string				`json:"id"`
	Cluster string			`json:"cluster"`
	Metadata ServerMetadata `json:"metadata"`
	UserAgentName string	`json:"user_agent_name"`
}

type ServerMetadata struct {
	Namespace string    `json:"namespace"`
	EnvoyVersion string `json:"envoy_version"`
}

// /memory

func (l *Server) Memory() (ServerMemoryJson, error) {
	rawJson, err := l.i.nsenter.Curl("-s", l.endpoints.memory)
	if err != nil { return ServerMemoryJson{}, err }

	jsonData := []byte(strings.Trim(rawJson, " "))

	var response ServerMemoryJson
	err = json.Unmarshal(jsonData, &response)
	if err != nil { return ServerMemoryJson{}, err }

	return response, nil
}

type ServerMemoryJson struct {
	Allocated int 			`json:"allocated,string"`
	HeapSize int			`json:"heap_size,string"`
	PageHeapUnmapped int 	`json:"pageheap_unmapped,string"`
	PageHeapFree int 		`json:"pageheap_free,string"`
	TotalThreadCache int 	`json:"total_thread_cache,string"`
	TotalPhysicalBytes int 	`json:"total_physical_bytes,string"`
}