package server

import "time"

type Log struct {
	Level      string    `json:"level" bson:"level" validate:"required"`
	Message    string    `json:"message" bson:"message" validate:"required"`
	ResourceId string    `json:"resourceId" bson:"resourceId" validate:"required"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp" validate:"required"`
	TraceId    string    `json:"traceId" bson:"traceId" validate:"required"`
	SpanId     string    `json:"spanId" bson:"spanId" validate:"required"`
	Commit     string    `json:"commit" bson:"commit" validate:"required"`
	Metadata   Metadata  `json:"metadata" bson:"metadata" validate:"required"`
}

type Metadata struct {
	ParentResourceId string `json:"parentResourceId" bson:"parentResourceId"`
}

type Logs struct {
	Logs []Log `json:"logs"`
}
