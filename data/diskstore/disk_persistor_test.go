package diskstore_test

import (
	. "github.com/DennisDenuto/igrb/data/diskstore"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"io/ioutil"
)

var _ = Describe("DiskPersistor", func() {

	Describe("Save", func() {
		Context("Saving a json record", func() {
			var diskPersistor DiskPersistor
			BeforeEach(func() {
				diskPersistor = NewDiskPersistor()
			})

			It("Should persist to disk", func() {
				val := make(map[string]string)
				val["val1"] = "val2"

				err := diskPersistor.Save("key", val)
				Expect(err).ToNot(HaveOccurred())

				_, err = os.Stat(DataDir)
				Expect(err).ToNot(HaveOccurred())

				content, err := ioutil.ReadFile(DataDir + "/key")
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(Equal(`{"val1":"val2"}`))
			})

			It("Should persist keys with colon in the name", func() {
				val := make(map[string]string)
				val["val1"] = "val2"

				err := diskPersistor.Save("key:key1", val)
				Expect(err).ToNot(HaveOccurred())

				_, err = os.Stat(DataDir)
				Expect(err).ToNot(HaveOccurred())

				content, err := ioutil.ReadFile(DataDir + "/key:key1")
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(Equal(`{"val1":"val2"}`))
			})

			It("Should return error if record cannot be marhalled into json", func() {
				value := func() {}
				err := diskPersistor.Save("key1", value)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Read", func() {
		var diskPersistor DiskPersistor
		BeforeEach(func() {
			diskPersistor = NewDiskPersistor()
			diskPersistor.Save("test", "test-123")
		})

		Context("Reading a record", func() {
			It("Should read successfully", func() {
				val, err := diskPersistor.Read("test")
				Expect(err).ToNot(HaveOccurred())
				Expect(string(val)).To(ContainSubstring("test-123"))
			})
		})
	})

	Describe("ReadAndUnmarshal", func() {
		var diskPersistor DiskPersistor
		type TestStruct struct {
			A string
		}

		BeforeEach(func() {
			diskPersistor = NewDiskPersistor()
			diskPersistor.Save("test", TestStruct{A: "testing123"})

		})

		Context("Reading a record", func() {
			It("Should read successfully", func() {
				val := TestStruct{}
				err := diskPersistor.ReadAndUnmarshal("test", &val)
				Expect(err).ToNot(HaveOccurred())
				Expect(val.A).To(ContainSubstring("testing123"))
			})
		})
	})
})
