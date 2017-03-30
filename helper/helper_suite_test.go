package helper

import (
. "github.com/onsi/ginkgo"
. "github.com/onsi/gomega"
"testing"
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helper Suite")
}


