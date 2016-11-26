package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"os"
	"github.com/DennisDenuto/igrb/data/diskstore"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var _ = BeforeSuite(func() {
	CleanUpDataDir()
})

var _ = AfterSuite(func() {
	CleanUpDataDir()
})

func CleanUpDataDir() {
	if _, err := os.Stat(diskstore.DataDir); err == nil {
		err = os.RemoveAll(diskstore.DataDir)
		Expect(err).ToNot(HaveOccurred())
	}
}