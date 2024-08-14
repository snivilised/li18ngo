# 🌐 i18n Assistance

This document aims to augment that which appears in the [goi18n project](https://github.com/nicksnyder/go-i18n/) and provide information that can help the client use the i18n functionality defined there and its integration into this template project. The translations process is quite laborious, so this project tries alleviate this process by providing helper tasks which will be documented here also.

💥 Warning: when using goi18n and defining content that needs to be extracted, it is vitally important that any code files with extractable content imports directly from go-i18n ie

> import ("github.com/nicksnyder/go-i18n/v2/i18n")

Now this may seem obvious, but you should not fall into the trap of writing your own layer on top of ___go-i18n___; do not do something like this ...:

```go
type Message = i18n.Message
```

... thinking that the new type defined in another package of your own making then defining messages with this new type. This will not work, because now presumably, you would be importing the new type, eg:

```go
import (
  "github.com/snivilised/extendio/i18n"
)

func (td PavementGraffitiReportTemplData) Message() *i18n.Message {
  return &i18n.Message{
    ID:          "pavement-graffiti-report.graffico.unit-test",
    Description: "Report of graffiti found on a pavement",
    Other:       "Found graffiti on pavement; primary colour: '{{.Primary}}'",
  }
}
```

In this example, the `i18n.Message` type refers to the alias not the original type defined in ___go-i18n___, but when you go to extract the content using the `go-i18n extract` command, this will silently extract nothing. This is because the extract command is specifically looking for "github.com/nicksnyder/go-i18n/v2/i18n" in the import statement.

This does also have other implications. When using ___go-i18n___, do not define your own local package of the name `i18n` because that will cause a conflict with ___go-i18n___ which is defined using this package name. I have taken to using an alternative name of `locale` so that no conflicts occur. However, this is not mandatory because you can just use a different import alias, but this is just adding friction that is best avoided, particularly when authoring a library package to be consumed by third parties.

## 📁 Directory Structure

The local directory structure is as follows:

- ___default___: contains the translation file created by the __newt__ (new translation task). Actually, this task creates an __active__ file (`active.en-GB.json`) in the `locale/out` folder, the result  of which needs to be manually copied into the __active__ file in the `default` folder.

- ___deploy___: contains all the translations files that are intended to be deployed with the application. There will be one per supported language and by default this template project includes a translation file for __en-US__ (`li18ngo.active.en-US.json`)

## ⚙️ Translation Workflow

### ✨ New Translations

__goi18n__ instructs the user to manually create an empty translation message file that they want to add (eg `translate.en-US.json`). This is taken care of by the __newt__ task. Then the requirement is to run the goi18n merge \<active\> command (goi18n merge `li18ngo.active.en-US.json` `li18ngo.translate.en-US.json`). This has been wrapped up into the __merge__ task and the result is that the translation file `li18ngo.translation.en-US.json` is populated with the messages to be translated. So the sequence goes:

- run __newt__ task: (generates default language file `./locale/out/l10n/active.en-GB.json` and empty `./locale/out/l10n/li18ngo.translation.en-US.json` file). This task can be run from the root folder, __goi18n__ will recursively search the directory tree for files with translate-able content, ie files with template definitions (___i18n.Message___)
- run __merge__ task: derives a translation file for the requested language __en-US__ using 2 files as inputs: source active file (`./locale/out/active.en-GB.json`) and the empty __en-US__ translate file (`./locale/out/l10n/li18ngo.translation.en-US.json`), both of which were generated in the previous step.
- hand the translate file to your translator for them to translate
- rename the translate file to the active equivalent (`li18ngo.translation.en-US.json`); ie rename `li18ngo.translation.en-US.json` to `li18ngo.active.en-US.json`. If `li18ngo.translation.en-US.json` exists as an empty file, then delete it as its not required. Save this into the __deploy__ folder. This file will be deployed with the application.

### 🧩 Update Existing Translations (⚠️ not finalised)

__goi18n__ instructs the user to run the following steps:

1. Run `goi18n extract` to update `active.en.toml` with the new messages.
2. Run `goi18n merge active.*.toml` to generate updated `translate.*.toml` files.
3. Translate all the messages in the `translate.*.toml` files.
4. Run `goi18n merge active.*.toml translate.*.toml` to merge the translated messages into the active message files.

___The above description is way too vague and ambiguous. It is for this reason that this process has not been finalised. The intention will be to upgrade the instructions as time goes by and experience is gained.___

However, in this template, the user can execute the following steps:

- run task __update__: this will re-extract messages into `active.en-GB.json` and then runs merge to create an updated translate file.
- hand the translate file to your translator for them to translate
- as before, this translated file should be used to update the active file `li18ngo.active.en-US.json` inside the __deploy__ folder.

## 🎓 Task Reference

❗This is a work in progress ...

### 💤 extract

Scans the code base for messages and extracts them into __out/active.en-GB.json__ (The name of the file can't be controlled, it is derived from the language specified). No need to call this task directly.

### 💠 newt

Invokes the `extract` task to extract messages from code. After running this task, the __translate.en-US.json__ file is empty ready to be used as one of the inputs to the merge task (_why do we need an empty input file for the merge? Perhaps the input file instructs the language tag to goi18n_).

Inputs:

- source code

Outputs:

- ./locale/out/active.en-GB.json (messages extracted from code, without hashes)
- ./locale/out/en-US/translate.en-US.json (empty)

### 💠 merge

Inputs:

- ./locale/out/active.en-GB.json
- ./locale/out/en-US/li18ngo.translate-en-US.json

Outputs:

- ./locale/out/active.en-US.json
- ./locale/out/translate.en-US.json
