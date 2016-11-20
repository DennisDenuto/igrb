package red_test

import (
	. "github.com/DennisDenuto/igrb/builds/red"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/concourse/fly/rc/rcfakes"
	"github.com/concourse/go-concourse/concourse/concoursefakes"
	"github.com/concourse/atc"
	"errors"
)

var _ = Describe("Builds", func() {

	Describe("ListBuilds", func() {

		var pipelineName string
		var build FailedBuildFetcher
		var target *rcfakes.FakeTarget
		var team *concoursefakes.FakeTeam

		BeforeEach(func() {
			target = &rcfakes.FakeTarget{}
			team = &concoursefakes.FakeTeam{}
			target.TeamReturns(team)

			build = FailedBuildFetcher{
				Target: target,
			}
		})

		Context("A Single pipeline with a Single Job with a Red Build", func() {
			BeforeEach(func() {
				pipelineName = "bosh"
				config := atc.Config{}
				config.Jobs = []atc.JobConfig{
					{Name: "job-1"},
				}

				team.PipelineConfigReturns(config, "", "", false, nil)
				team.JobReturns(atc.Job{FinishedBuild: &atc.Build{Status: "failed"}}, true, nil)
			})

			It("Should return the correct build", func() {
				redBuilds, err := build.Fetch(pipelineName)
				Expect(err).ToNot(HaveOccurred())

				Expect(target.TeamCallCount()).To(Equal(1))

				Expect(team.PipelineConfigCallCount()).To(Equal(1))
				pipelineName := team.PipelineConfigArgsForCall(0)
				Expect(pipelineName).To(Equal(pipelineName))

				Expect(team.JobCallCount()).To(Equal(1))
				pipelineName, jobName := team.JobArgsForCall(0)
				Expect(pipelineName).To(Equal(pipelineName))
				Expect(jobName).To(Equal("job-1"))

				Expect(redBuilds).To(HaveLen(1))
				Expect(redBuilds[0].Status).To(Equal("failed"))
			})
		})

		Context("A Single pipeline with multiple Jobs all being Red builds", func() {
			BeforeEach(func() {
				pipelineName = "bosh"
				config := atc.Config{}
				config.Jobs = []atc.JobConfig{
					{Name: "job-1"},
					{Name: "job-2"},
				}

				team.PipelineConfigReturns(config, "", "", false, nil)
				team.JobReturns(atc.Job{FinishedBuild: &atc.Build{Status: "failed"}}, true, nil)
			})

			It("Should return the correct build", func() {
				redBuilds, err := build.Fetch(pipelineName)
				Expect(err).ToNot(HaveOccurred())

				Expect(redBuilds).To(HaveLen(2))
				Expect(redBuilds[0].Status).To(Equal("failed"))
				Expect(redBuilds[1].Status).To(Equal("failed"))
			})
		})

		Context("A Single pipeline with multiple Jobs with some Red and Green builds", func() {
			BeforeEach(func() {
				pipelineName = "bosh"
				config := atc.Config{}
				config.Jobs = []atc.JobConfig{
					{Name: "job-1"},
					{Name: "job-2"},
					{Name: "job-3"},
				}

				team.PipelineConfigReturns(config, "", "", false, nil)
				team.JobStub = func(pipelineName, jobName string) (atc.Job, bool, error) {
					switch jobName {
					case "job-1": return atc.Job{FinishedBuild: &atc.Build{Status: "failed"}}, true, nil
					case "job-2": return atc.Job{FinishedBuild: &atc.Build{Status: "succeeded"}}, true, nil
					case "job-3": return atc.Job{FinishedBuild: &atc.Build{Status: "errored"}}, true, nil
					default: return atc.Job{FinishedBuild: &atc.Build{Status: "failed"}}, true, nil
					}
				}
			})

			It("Should return the correct build", func() {
				redBuilds, err := build.Fetch(pipelineName)
				Expect(err).ToNot(HaveOccurred())

				Expect(redBuilds).To(HaveLen(2))
				Expect(redBuilds[0].Status).To(Equal("failed"))
				Expect(redBuilds[1].Status).To(Equal("errored"))
			})
		})

		Context("concourse errors out", func() {
			BeforeEach(func() {
				pipelineName = "bosh"
			})

			It("should return correct error if pipeline fails", func() {
				team.PipelineConfigReturns(atc.Config{}, "", "", false, errors.New("some error"))
				_, err := build.Fetch(pipelineName)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to get pipeline config: some error"))
			})

			It("should return correct error if getting a job fails", func() {
				config := atc.Config{}
				config.Jobs = []atc.JobConfig{
					{Name: "job-1"},
					{Name: "job-2"},
				}

				team.PipelineConfigReturns(config, "", "", false, nil)
				team.JobReturns(atc.Job{}, true, errors.New("some strange error"))

				_, err := build.Fetch(pipelineName)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to get pipeline config: some strange error"))
			})
		})
	})
})
