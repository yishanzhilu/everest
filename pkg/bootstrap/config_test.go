package bootstrap

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	When("init with right file name", func() {
		It("should work", func() {
			Ω(func() {
				initConfig("viper")
			}).ShouldNot(Panic())
		})
	})
	When("init with wrong file name", func() {
		It("should work", func() {
			Ω(func() {
				initConfig("wrong-name")
			}).Should(Panic())
		})
	})
})
