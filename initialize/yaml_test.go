package initialize

import (
	"fmt"
	"testing"
)

func TestYaml(t *testing.T) {
	svc  := NewInitSvc("../conf/firstlog.yaml")
	fmt.Println(svc.OutputES())
}
