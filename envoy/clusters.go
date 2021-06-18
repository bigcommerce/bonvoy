package envoy

import (
	"github.com/spf13/cast"
	"strings"
)

type Clusters struct {
	i *Instance
	endpoints ClustersEndpoints
}

type ClustersEndpoints struct {
	list string
}

func (i *Instance) Clusters() *Clusters {
	return &Clusters{
		i: i,
		endpoints: ClustersEndpoints{
			list: i.Address + "/clusters",
		},
	}
}


type ClusterStatistics struct {
	Host string `json:"host"`
	Outlier ClusterOutlierStatistics `json:"outlier"`
	DefaultPriority ClusterPriorityStatistics `json:"default_priority"`
	HighPriority ClusterPriorityStatistics `json:"high_priority"`
	AddedViaApi bool `json:"added_via_api"`
	Instances map[string]ClusterInstance `json:"instances"`
}
type ClusterInstance struct {
	Hostname string `json:"hostname"`
	Connections ClusterConnectionStatistics `json:"connections"`
	Requests ClusterRequestStatistics `json:"requests"`
	HealthFlags string `json:"health_flags"`
	Weight int `json:"weight"`
	Region string `json:"region"`
	Zone string `json:"zone"`
	SubZone string `json:"sub_zone"`
	Canary bool `json:"canary"`
	Priority int `json:"priority"`
	SuccessRate string `json:"success_rate"`
	LocalOriginSuccessRate string `json:"local_origin_success_rate"`
}
type ClusterOutlierStatistics struct {
	SuccessRateAverage string `json:"success_rate_average"`
	SuccessRateEjectionThreshold string `json:"success_rate_ejection_threshold"`
	LocalOriginSuccessRateAverage string `json:"local_origin_success_rate_average"`
	LocalOriginSuccessRateEjectionThreshold string `json:"local_origin_success_rate_ejection_threshold"`
}
type ClusterPriorityStatistics struct {
	MaxConnections int `json:"max_connections"`
	MaxPendingRequests int `json:"max_pending_requests"`
	MaxRequests int `json:"max_requests"`
	MaxRetries int `json:"max_retries"`
}
type ClusterConnectionStatistics struct {
	Active int `json:"active"`
	Failed int `json:"failed"`
	Total int `json:"total"`
}
type ClusterRequestStatistics struct {
	Active int `json:"active"`
	Error int `json:"error"`
	Success int `json:"success"`
	Timeout int `json:"timeout"`
	Total int `json:"total"`
}

func (c *Clusters) GetStatistics(specific string) (map[string]ClusterStatistics, error) {
	var clusters = make(map[string]ClusterStatistics)

	raw, err := c.i.nsenter.Curl("-s", c.endpoints.list)
	if err != nil { return clusters, err }

	data := strings.Split(raw, "\n")
	for _, str := range data {
		parts := strings.Split(str, "::")
		if len(parts) < 2 { continue }

		var cs ClusterStatistics
		if _, ok := clusters[parts[0]]; ok {
			cs = clusters[parts[0]]
		} else {
			cs = ClusterStatistics{
				Outlier:         ClusterOutlierStatistics{},
				DefaultPriority: ClusterPriorityStatistics{},
				HighPriority:    ClusterPriorityStatistics{},
				Instances:       make(map[string]ClusterInstance),
			}
		}
		cs.Host = parts[0]

		// If specific cluster specified, filter others out
		if specific != "" && cs.GetConsulName() != specific {
			continue
		}

		if parts[1] == "outlier" {
			cs.Outlier.Deserialize(parts[2], parts[3])
		} else if parts[1] == "default_priority" {
			cs.DefaultPriority.Deserialize(parts[2], parts[3])
		} else if parts[1] == "high_priority" {
			cs.HighPriority.Deserialize(parts[2], parts[3])
		} else if parts[1] == "added_via_api" {
			cs.AddedViaApi = parts[2] == "true"
		} else if strings.Contains(parts[1], ":") { // then this is an instance addr
			var i ClusterInstance
			if _, ok := cs.Instances[parts[1]]; ok {
				i = cs.Instances[parts[1]]
			} else {
				i = ClusterInstance{
					Hostname: parts[1],
					Connections: ClusterConnectionStatistics{},
					Requests: ClusterRequestStatistics{},
				}
			}
			i.Deserialize(parts[2], parts[3])
			cs.Instances[parts[1]] = i
		} else {
			continue
		}
		clusters[parts[0]] = cs
	}

	return clusters, nil
}

func (s *ClusterOutlierStatistics) Deserialize(fieldName string, value string) {
	switch fieldName {
	case "success_rate_average":
		s.SuccessRateAverage = value
	case "success_rate_ejection_threshold":
		s.SuccessRateEjectionThreshold = value
	case "local_origin_success_rate_average":
		s.LocalOriginSuccessRateAverage = value
	case "local_origin_success_rate_ejection_threshold":
		s.LocalOriginSuccessRateEjectionThreshold = value
	}
}

func (s *ClusterPriorityStatistics) Deserialize(fieldName string, value string) {
	switch fieldName {
	case "max_connections":
		s.MaxConnections = cast.ToInt(value)
	case "max_pending_requests":
		s.MaxPendingRequests = cast.ToInt(value)
	case "max_requests":
		s.MaxRequests = cast.ToInt(value)
	case "max_retries":
		s.MaxRetries = cast.ToInt(value)
	}
}

func (s *ClusterInstance) Deserialize(fieldName string, value string) {
	switch fieldName {
	case "cx_active":
		s.Connections.Active = cast.ToInt(value)
	case "cx_connect_fail":
		s.Connections.Failed = cast.ToInt(value)
	case "cx_total":
		s.Connections.Total = cast.ToInt(value)
	case "rq_active":
		s.Requests.Active = cast.ToInt(value)
	case "rq_error":
		s.Requests.Error = cast.ToInt(value)
	case "rq_success":
		s.Requests.Success = cast.ToInt(value)
	case "rq_total":
		s.Requests.Total = cast.ToInt(value)
	case "health_flags":
		s.HealthFlags = value
	case "weight":
		s.Weight = cast.ToInt(value)
	case "region":
		s.Region = value
	case "zone":
		s.Zone = value
	case "sub_zone":
		s.SubZone = value
	case "canary":
		s.Canary = value == "true"
	case "priority":
		s.Priority = cast.ToInt(value)
	case "success_rate":
		s.SuccessRate = value
	case "local_origin_success_rate":
		s.LocalOriginSuccessRate = value
	}
}

func (s *ClusterStatistics) GetConsulName() string {
	return strings.Split(s.Host, ".")[0]
}