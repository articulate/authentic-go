package authentic_go_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAuthenticGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthenticGo Suite")
}
