package store

import (
	"container-angel/angel"
	"strconv"
	"time"
)

type Archive struct {
	FQP string // full qualified path including filename
}

func NewArchive(conf angel.Configuration) Archive {
	return Archive{
		FQP: conf.AngelDirectory + "/" + strconv.FormatInt(time.Now().UnixNano(), 10),
	}
}
