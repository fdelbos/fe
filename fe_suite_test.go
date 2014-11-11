package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestFe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fe Suite")
}
