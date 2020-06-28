package store

import (
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

type Store struct {
	RootDirectory string
	Archives      []Archive
}

func FindAll(rootDirectory string) Store {
	files, err := ioutil.ReadDir(rootDirectory)
	if err != nil {
		log.Fatalf("Error while reading the backup archive %s:%s", rootDirectory, err)
	}

	ret := Store{
		RootDirectory: "",
		Archives:      make([]Archive, 0),
	}

	for _, element := range files {
		if strings.HasPrefix(element.Name(), ".") {
			continue
		}

		a := Archive{
			FQP: rootDirectory + "/" + element.Name(),
		}
		ret.Archives = append(ret.Archives, a)
	}

	sort.Slice(ret.Archives, func(i, j int) bool {
		return ret.Archives[i].FQP > ret.Archives[j].FQP
	})

	return ret
}
