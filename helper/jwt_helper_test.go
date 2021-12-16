package helper_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	helper "personal-blog/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JwtHelper", func() {
	jsonBytes, _ := json.Marshal("body")
	fakeRequest, _ := http.NewRequest(http.MethodPost, "deneme.com", bytes.NewReader(jsonBytes))
	Context("When calling CreateToken", func() {
		token, _ := helper.CreateToken("deneme")
		fmt.Println(token)
		It("should not be a nil", func() {
			Expect(token).ToNot(BeNil())
		})
	})

	Context("When calling ValidateToken", func() {
		token, _ := helper.CreateToken("deneme")
		header := "Bearer " + token
		fakeRequest.Header.Add("Authorization", header)
		validateToken := helper.TokenValid(fakeRequest)
		It("should be nil", func() {
			Expect(validateToken).To(BeNil())
		})
	})
})
