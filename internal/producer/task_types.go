package producer

type TaskTypes int

const (
	Invalid TaskTypes = iota
	VMDeployTask
	VMInfoTask
)
