package crypto_test

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/yishanzhilu/everest/pkg/crypto"
)

var _ = Describe("JWT", func() {
	guard := NewJWTGuard("111", 3600)
	It("should be able to sign for a user with user name", func() {
		token, err := guard.SignToken(1, "user1")
		Ω(err).Should(BeNil())
		userID, err := guard.CheckToken(token)
		Ω(err).Should(BeNil())
		Ω(userID).Should(Equal(uint64(1)))
	})
	It("should be able to handle exp", func() {
		token, err := guard.SignToken(1, "user1")
		Ω(err).Should(BeNil())
		jwt.TimeFunc = func() time.Time {
			return time.Unix(time.Now().Unix()+5000, 0)
		}
		userID, err := guard.CheckToken(token)
		jwt.TimeFunc = time.Now
		Ω(userID).Should(Equal(uint64(1)))
	})
})
