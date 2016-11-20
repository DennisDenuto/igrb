package diskstore_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"os"
	"github.com/DennisDenuto/igrb/data/diskstore"
)

func TestDiskstore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Diskstore Suite")
}


var _ = BeforeSuite(func() {

	CleanUpDataDir()
})

var _ = AfterSuite(func() {
	CleanUpDataDir()
})

func CleanUpDataDir() {
	if _, err := os.Stat(diskstore.DataDir); os.IsExist(err) {
		err = os.Remove(diskstore.DataDir)
		Expect(err).ToNot(HaveOccurred())
	}
}