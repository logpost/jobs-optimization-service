package models

// JobExpected use for destruct job from object
type JobExpected struct {
	Job		Job		`json:"job"`
}

// GetterExpected use for parse only field required
type GetterExpected struct { 
	Getter	[]JobExpected	`json:"getter"`
}