package diskstore

import (
	"github.com/peterbourgon/diskv"
	"fmt"
	"os"
	"encoding/json"
	"github.com/pkg/errors"

)


var DataDir string = fmt.Sprintf("%sigrb-data", os.TempDir())

type DiskPersistor struct {
	diskv *diskv.Diskv
}

func NewDiskPersistor() DiskPersistor {
	flatTransform := func(s string) []string {
		return []string{}
	}

	diskV := diskv.New(diskv.Options{
		BasePath:     DataDir,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	return DiskPersistor{
		diskv: diskV,
	}
}

func (disk DiskPersistor) Save(key string, val interface{}) error {
	jsonVal, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal into json")
	}

	return disk.diskv.Write(key, jsonVal)
}

func (disk DiskPersistor) Read(key string) ([]byte, error) {
	return disk.diskv.Read(key)
}

func (disk DiskPersistor) ReadAndUnmarshal(key string, val interface{}) (error) {
	data, err := disk.diskv.Read(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, val)
}

