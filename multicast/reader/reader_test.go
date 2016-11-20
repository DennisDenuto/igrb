package reader_test

import (
	. "github.com/DennisDenuto/igrb/multicast/reader"

	. "github.com/onsi/ginkgo"
	"net"
	"github.com/DennisDenuto/igrb/multicast/types"
	"encoding/json"
	. "github.com/onsi/gomega"
	"github.com/DennisDenuto/igrb/data/diskstore"
	"time"
)

func SendMulticast(mGroup string, payload string) error {
	addr, err := net.ResolveUDPAddr("udp", mGroup)
	if err != nil {
		return err
	}
	c, err := net.DialUDP("udp", nil, addr)
	c.Write([]byte(payload))

	return nil
}

var _ = Describe("Reader", func() {

	Describe("ServeMulticastUDP", func() {
		Context("A build is taken by a dev", func() {
			var finish chan bool

			BeforeEach(func() {
				finish = make(chan bool, 1)

				go func() {
					NewServeMulticastUDP("224.0.0.1:9999").ServeMulticastUDP(finish)
				}()

			})

			AfterEach(func() {
				finish <- true
			})

			It("Should correctly parse and update the store with the info", func() {
				req := types.DevLookingIntoBuild{
					DevName: "dev-name",
					PipelineName: "pipeline",
					JobName: "job-name",
					BuildId: "build-id",
				}

				reqJson, err := json.Marshal(req)
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() types.DevLookingIntoBuild {
					SendMulticast("224.0.0.1:9999", string(reqJson))
					resp := types.DevLookingIntoBuild{}
					diskstore.NewDiskPersistor().ReadAndUnmarshal("pipeline_job_build-id", &resp)
					return resp
				}, 5 * time.Second, 1 * time.Second).Should(Equal(req))

			})

			It("Should handle being sent an unknown request", func() {
				type UnknownRequest struct {
					A string
				}
				req := UnknownRequest{
					A: "foo",
				}

				reqJson, err := json.Marshal(req)
				Expect(err).ToNot(HaveOccurred())

				allKeys, _ := diskstore.NewDiskPersistor().ListKeys()
				SendMulticast("224.0.0.1:9999", string(reqJson))

				Expect(diskstore.NewDiskPersistor().ListKeys()).To(HaveLen(len(allKeys)))
			})
		})
	})
})
