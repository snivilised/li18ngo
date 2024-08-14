package nfs_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo/internal/helpers"
	"github.com/snivilised/li18ngo/internal/nfs"
)

var _ = Describe("EnsurePathAt", Ordered, func() {
	var (
		mocks *ResolveMocks
		mfs   *mkDirAllMapFS
	)

	BeforeEach(func() {
		mocks = &ResolveMocks{
			HomeFunc: func() (string, error) {
				return filepath.Join(string(filepath.Separator), "home", "prodigy"), nil
			},
			AbsFunc: func(_ string) (string, error) {
				return "", errors.New("not required for these tests")
			},
		}

		mfs = &mkDirAllMapFS{
			mapFS: fstest.MapFS{
				filepath.Join("home", "prodigy"): &fstest.MapFile{
					Mode: os.ModeDir,
				},
			},
		}
	})

	DescribeTable("with mapFS",
		func(entry *ensureTE) {
			home, _ := mocks.HomeFunc()
			location := TrimRoot(filepath.Join(home, entry.relative))

			if entry.directory {
				location += string(filepath.Separator)
			}

			actual, err := nfs.EnsurePathAt(location, "default-test.log", perm, mfs)
			directory, _ := filepath.Split(actual)
			directory = filepath.Clean(directory)
			expected := TrimRoot(helpers.Path(home, entry.expected))

			Expect(err).Error().To(BeNil())
			Expect(actual).To(Equal(expected))
			Expect(AsDirectory(TrimRoot(directory))).To(ExistInFS(mfs))
		},
		func(entry *ensureTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &ensureTE{
			given:    "path is file",
			should:   "create parent directory and return specified file path",
			relative: filepath.Join("logs", "test.log"),
			expected: "logs/test.log",
		}),

		Entry(nil, &ensureTE{
			given:     "path is directory",
			should:    "create parent directory and return default file path",
			relative:  "logs/",
			directory: true,
			expected:  "logs/default-test.log",
		}),
	)
})
