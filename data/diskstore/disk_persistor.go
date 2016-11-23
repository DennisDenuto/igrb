package diskstore

import (
	"github.com/peterbourgon/diskv"
	"fmt"
	"os"
	"encoding/json"
	"github.com/pkg/errors"
	logger "github.com/Sirupsen/logrus"
)

var DataDir string = fmt.Sprintf("%s/igrb-data", os.TempDir())

type DiskPersistor struct {
	diskv *diskv.Diskv
}

func NewDiskPersistor() DiskPersistor {
	flatTransform := func(s string) []string {
		return []string{}
	}

	logger.Debugf("Init db in %s", DataDir)
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

func (disk DiskPersistor) ListKeys() ([]string, error) {
	var allKeys []string

	keys := disk.diskv.Keys(nil)
	for value := range keys {
		allKeys = append(allKeys, value)
	}

	return allKeys, nil
}