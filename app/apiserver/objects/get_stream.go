package objects

import (
	"fmt"
	"goss/app/apiserver/locate"
	"goss/pkg/objectstream"
	"io"
)

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}

