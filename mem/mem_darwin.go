package mem

import (
	"os/exec"

	"github.com/cloudfoundry/gosigar"
)

func Get() (uint64, error) {
	mem := sigar.Mem{}
	err := mem.Get()
	if err != nil {
		return 0, err
	}
	return mem.Total, nil
}

func ClearBufferCache() error {
	return exec.Command("/usr/sbin/purge").Run()
}
