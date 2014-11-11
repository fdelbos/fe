package rw_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestRw(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rw Suite")
}
