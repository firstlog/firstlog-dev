package output

import (
	"firstlog-docker/output"
	"log"
	"testing"
)

func TestNewEs(t *testing.T) {
	toEs, err := NewEs([]string{"http://localhost:9200"},[]string{"www"},"2","1","7")
	if err != nil {
		log.Println(err)
		return
	}
	toEs.ToEs(output.Storage)
}
