package repositories

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestClients(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Test Suite")
}
