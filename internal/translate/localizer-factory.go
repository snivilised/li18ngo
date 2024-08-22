package translate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/snivilised/li18ngo/internal/lo"
	"github.com/snivilised/li18ngo/internal/nfs"
	"golang.org/x/text/language"
)

func createLocalizer(lang *LanguageInfo, sourceID string,
	dirFS nfs.MkDirAllFS,
) (*i18n.Localizer, error) {
	bundle := i18n.NewBundle(lang.Tag)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	if lang.Tag != lang.Default {
		txSource := lang.From.Sources[sourceID]
		path := resolveBundlePath(lang, txSource, dirFS)
		_, err := bundle.LoadMessageFile(path)

		if (err != nil) && (!lang.DefaultIsAcceptable) {
			return nil, NewCouldNotLoadTranslationsNativeError(lang.Tag, path, err)
		}
	}

	supported := lo.Map(lang.Supported, func(t language.Tag, _ int) string {
		return t.String()
	})

	return i18n.NewLocalizer(bundle, supported...), nil
}

func resolveBundlePath(lang *LanguageInfo, txSource TranslationSource,
	dirFS nfs.MkDirAllFS,
) string {
	path := lo.Ternary(txSource.Path != "" && dirFS.DirectoryExists(txSource.Path),
		txSource.Path,
		lang.From.Path,
	)

	directory := lo.TernaryF(path != "" && dirFS.DirectoryExists(path),
		func() string {
			resolved, _ := filepath.Abs(path)
			return resolved
		},
		func() string {
			exe, _ := os.Executable()
			return filepath.Dir(exe)
		},
	)

	filename := lo.TernaryF(txSource.Name == "",
		func() string {
			return fmt.Sprintf("active.%v.json", lang.Tag)
		},
		func() string {
			return fmt.Sprintf("%v.active.%v.json", txSource.Name, lang.Tag)
		},
	)

	return filepath.Join(directory, filename)
}
