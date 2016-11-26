package commands_test

import (
	. "github.com/DennisDenuto/igrb/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/concourse/atc"
	"github.com/DennisDenuto/igrb/ui/bitbar"
	"github.com/DennisDenuto/igrb/data/diskstore"
	"github.com/DennisDenuto/igrb/multicast"
	"os"
	"time"
)

var _ = Describe("Status", func() {

	BeforeEach(func() {
		err := os.RemoveAll(diskstore.DataDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("AddMenuItemsToPainter", func() {
		Context("red builds", func() {
			var failedBuilds map[string][]atc.Build = make(map[string][]atc.Build)

			BeforeEach(func() {
				failedBuilds["pipeline-1"] = []atc.Build{
					{ID: 1, JobName: "job-1", PipelineName: "pipeline-1"},
					{ID: 2, JobName: "job-2", PipelineName: "pipeline-1"},
					{ID: 3, JobName: "job-3", PipelineName: "pipeline-1"},
				}
			})
			It("Should add all red builds to painter", func() {
				painter := &bitbar.Painter{}

				AddMenuItemsToPainter("url", failedBuilds, painter)

				Expect(painter.MainItems).To(HaveLen(3))
			})

			Context("Some red builds have been 'ignored'", func() {

				BeforeEach(func() {
					ignoredBuild := multicast.DevLookingIntoBuild{
						DevName      : "",
						PipelineName : "pipeline-1",
						JobName      : "job-1",
						Ignore       : true,
						ID           : 1,
					}

					diskstore.NewDiskPersistor().Save(ignoredBuild.Key(), ignoredBuild)
				})

				It("Should not add 'ignored' red builds", func() {
					painter := &bitbar.Painter{}

					AddMenuItemsToPainter("url", failedBuilds, painter)

					Expect(painter.MainItems).To(HaveLen(2))
				})

			})
		})
	})

	Describe("SaveFailedBuildsSummary", func() {
		var failedBuilds map[string][]atc.Build = make(map[string][]atc.Build)

		BeforeEach(func() {
			failedBuilds["pipeline-1"] = []atc.Build{
				{ID: 1, JobName: "job-1", PipelineName: "pipeline-1"},
			}
			failedBuilds["pipeline-2"] = []atc.Build{
				{ID: 1, JobName: "job-1", PipelineName: "pipeline-2"},
			}
		})

		Context("failed builds returned", func() {
			It("Should persist to disk", func() {
				err := SaveFailedBuildsSummary("url", failedBuilds)
				Expect(err).ToNot(HaveOccurred())

				summary := &FailedBuildsSummary{}
				err = diskstore.NewDiskPersistor().ReadAndUnmarshal("summary", summary)
				Expect(err).ToNot(HaveOccurred())

				Expect(summary.URL).To(Equal("url"))
				Expect(summary.FailedBuilds).To(HaveLen(2))
				Expect(summary.CreatedAt).To(BeTemporally("~", time.Now(), 5 * time.Second))
			})
		})

	})
})
