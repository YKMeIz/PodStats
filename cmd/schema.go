package main

type message struct {
	Type    string
	Content interface{}
}

type systemReport struct {
	Cpu     float64
	Memory  float64
	Service float64
	Status  string
}

type containerReport struct {
	Name         string
	Cpu          float64
	Memory       float64
	NetworkIO    string
	BlockIO      string
	Created      string
	StartedAt    string
	RestartCount int32
}

type eventReport struct {
	ID     string
	Name   string
	Action string
	Time   int64
}

type statusReport struct {
	Name        string
	Description string
	Status      string
}
