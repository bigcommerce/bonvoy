package envoy

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Server struct {
	i *Instance
	endpoints ServerEndpoints
}

type ServerEndpoints struct {
	info string
}
func (i *Instance) Server() *Server {
	return &Server{
		i: i,
		endpoints: ServerEndpoints{
			info: i.Address + "/server_info",
		},
	}
}

func (l *Server) Info() ServerInfoJson {
	rawJson := l.i.nsenter.Curl("-s", l.endpoints.info)
	jsonData := []byte(strings.Trim(rawJson, " "))

	var response ServerInfoJson
	err := json.Unmarshal(jsonData, &response)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return response
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