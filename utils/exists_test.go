package utils_test

import (
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo/internal/helpers"

	"github.com/snivilised/li18ngo/utils"
)

func path(parent, relative string) string {
	segments := strings.Split(relative, "/")
	return filepath.Join(append([]string{parent}, segments...)...)
}

var _ = Describe("Exists Utils", Ordered, func() {
	var repo string

	BeforeAll(func() {
		repo = helpers.Repo("")
		Expect(utils.FolderExists(repo)).To(BeTrue())
	})

	DescribeTable("Exists",
		func(_, relative string, expected bool, msg string) {
			path := path(repo, relative)

			GinkgoWriter.Printf("---> 🔰 FULL-PATH: '%v'\n", path)
			Expect(utils.Exists(path)).To(Equal(expected), msg)
		},

		func(message, _ string, _ bool, _ string) string {
			return fmt.Sprintf("🥣 message: '%v'", message)
		},
		Entry(nil, "folder exists", "/", true, "failed: root path should exist"),
		Entry(nil, "file exists", "README.md", true, "failed: README.md path should exist"),
		Entry(nil, "does not exist", "foo-bar", false, "failed: foo-bar path should not exist"),
	)

	DescribeTable("FolderExists",
		func(_, relative string, expected bool, msg string) {
			path := path(repo, relative)
			GinkgoWriter.Printf("---> 🔰 FULL-PATH: '%v'\n", path)

			Expect(utils.FolderExists(path)).To(Equal(expected), msg)
		},
		func(message, _ string, _ bool, _ string) string {
			return fmt.Sprintf("🍤 message: '%v'", message)
		},
		Entry(nil, "folder exists", "/", true, "failed: root folder should exist"),
		Entry(nil, "exists as file", "README.md", false, "failed: README.md file should exist"),
		Entry(nil, "folder does not exist", "foo-bar", false, "failed: foo-bar folder should not exist"),
	)

	DescribeTable("FileExists",
		func(_, relative string, expected bool, msg string) {
			path := path(repo, relative)
			GinkgoWriter.Printf("---> 🔰 FULL-PATH: '%v'\n", path)

			Expect(utils.FileExists(path)).To(Equal(expected), msg)
		},
		func(message, _ string, _ bool, _ string) string {
			return fmt.Sprintf("🍤 message: '%v'", message)
		},
		Entry(nil, "file exists", "README.md", true, "failed: root file should exist"),
		Entry(nil, "file does not exist", "foo-bar", false, "failed: foo-bar file should not exist"),
		Entry(nil, "does not exist as file", "Test", false, "failed: Test file should not exist"),
	)
})
