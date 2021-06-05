package model

type LogObject struct {
	ID   string `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
}

type LogRequest struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Container string `json:"container,omitempty"`
	Follow    bool   `json:"follow,omitempty"`
	Previous  bool   `json:"previous,omitempty"`
	TailLines int64  `json:"taillines,omitempty"`
	Stop      bool   `json:"stop,omitempty"`
}

type LogRequests map[string]LogRequest
