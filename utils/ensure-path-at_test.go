package utils_test

import (
	"errors"
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo/internal/helpers"
	"github.com/snivilised/li18ngo/matchers"
	"github.com/snivilised/li18ngo/storage"
	"github.com/snivilised/li18ngo/utils"
)

type ensureTE struct {
	given     string
	should    string
	relative  string
	expected  string
	directory bool
}

const perm = 0o766

var _ = Describe("EnsurePathAt", Ordered, func() {
	var (
		vfs   storage.VirtualFS
		mocks *utils.ResolveMocks
	)

	BeforeEach(func() {
		mocks = &utils.ResolveMocks{
			HomeFunc: func() (string, error) {
				return filepath.Join(string(filepath.Separator), "home", "prodigy"), nil
			},
			AbsFunc: func(_ string) (string, error) {
				return "", errors.New("not required for these tests")
			},
		}

		vfs = storage.UseMemFS()
	})

	DescribeTable("with vfs",
		func(entry *ensureTE) {
			home, _ := mocks.HomeFunc()
			location := filepath.Join(home, entry.relative)
			if entry.directory {
				location += string(filepath.Separator)
			}

			actual, err := utils.EnsurePathAt(location, "default-test.log", perm, vfs)
			directory, _ := filepath.Split(actual)
			expected := helpers.Path(home, entry.expected)

			Expect(err).Error().To(BeNil())
			Expect(actual).To(Equal(expected))
			Expect(matchers.AsDirectory(directory)).To(matchers.ExistInFS(vfs))
		},
		func(entry *ensureTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &ensureTE{
			given:    "path with file",
			should:   "create parent directory and return specified file path",
			relative: "logs/test.log",
			expected: "logs/test.log",
		}),

		Entry(nil, &ensureTE{
			given:     "path with file",
			should:    "create parent directory and return default file path",
			relative:  "logs/",
			directory: true,
			expected:  "logs/default-test.log",
		}),
	)
})
