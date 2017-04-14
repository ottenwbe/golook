package app

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(" AppInfo ", func() {

	It("should comprise the current app version by default", func() {
		info := NewAppInfo()
		Expect(info.Version).To(Equal(VERSION))
	})

	It("should comprise the current system by default", func() {
		info := NewAppInfo()
		Expect(info.System).ToNot(BeNil())
	})
})
