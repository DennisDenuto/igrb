package reader_test

import (
	. "github.com/DennisDenuto/igrb/multicast/reader"

	. "github.com/onsi/ginkgo"
	"encoding/json"
	. "github.com/onsi/gomega"
	"github.com/DennisDenuto/igrb/data/diskstore"
	"time"
	"github.com/DennisDenuto/igrb/multicast"
	"github.com/DennisDenuto/igrb/multicast/sender"
	"net"
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
			var multicastSender sender.MulticastSender

			BeforeEach(func() {
				multicastSender = sender.NewMultiCastSender("224.0.0.1:9999")

				finish = make(chan bool, 1)

				go func() {
					NewServeMulticastUDP("224.0.0.1:9999").ServeMulticastUDP(finish)
				}()

			})

			AfterEach(func() {
				finish <- true
			})

			It("Should correctly parse and update the store with the info", func() {
				req := multicast.DevLookingIntoBuild{
					DevName: "dev-name",
					PipelineName: "pipeline",
					JobName: "job-name",
					ID: 123,
					CreatedAt: time.Now(),
				}

				Eventually(func() multicast.DevLookingIntoBuild {
					multicastSender.SendMulticast(req)
					resp := multicast.DevLookingIntoBuild{}
					diskstore.NewDiskPersistor().ReadAndUnmarshal(req.Key(), &resp)
					return resp
				}, 5 * time.Second, 1 * time.Second).Should(Equal(req))

			})

			It("Should not update build info if request is stale / old compared to what is stored locally", func() {
				req := multicast.DevLookingIntoBuild{
					DevName: "dev-name1",
					PipelineName: "pipeline1",
					JobName: "job-name1",
					ID: 1234,
					CreatedAt: time.Now(),
				}

				jsonReq, err := json.Marshal(req)
				Expect(err).ToNot(HaveOccurred())
				MsgHandler(nil, len(jsonReq), jsonReq)


				oldRequest := multicast.DevLookingIntoBuild{
					DevName: "dev-name2",
					PipelineName: "pipeline1",
					JobName: "job-name1",
					ID: 1234,
					CreatedAt: time.Now().Add(-1 * time.Minute),
				}

				oldReq, err := json.Marshal(oldRequest)
				Expect(err).ToNot(HaveOccurred())
				MsgHandler(nil, len(oldReq), oldReq)

				resp := multicast.DevLookingIntoBuild{}
				diskstore.NewDiskPersistor().ReadAndUnmarshal(req.Key(), &resp)

				Expect(resp.DevName).To(Equal("dev-name1"))
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
				time.Sleep(2 * time.Second)

				Expect(diskstore.NewDiskPersistor().ListKeys()).To(HaveLen(len(allKeys)))
			})
		})
	})
})
