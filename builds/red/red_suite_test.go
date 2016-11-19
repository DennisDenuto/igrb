package red_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRed(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Red Suite")
}
