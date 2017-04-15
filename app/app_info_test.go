package app

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(" Info ", func() {

	var info *Info

	BeforeEach(func() {
		info = NewAppInfo()
	})

	It("should comprise the current app version by default", func() {
		Expect(info.Version).To(Equal(VERSION))
	})

	It("should comprise the current system by default", func() {
		Expect(info.System).ToNot(BeNil())
		Expect(*info.System).To(Equal(*NewSystem()))
	})
})
