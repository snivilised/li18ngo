# 🌐 li18ngo: ___i18n helper___

[![A B](https://img.shields.io/badge/branching-commonflow-informational?style=flat)](https://commonflow.org)
[![A B](https://img.shields.io/badge/merge-rebase-informational?style=flat)](https://git-scm.com/book/en/v2/Git-Branching-Rebasing)
[![A B](https://img.shields.io/badge/branch%20history-linear-blue?style=flat)](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/managing-a-branch-protection-rule)
[![Go Reference](https://pkg.go.dev/badge/github.com/snivilised/li18ngo.svg)](https://pkg.go.dev/github.com/snivilised/li18ngo)
[![Go report](https://goreportcard.com/badge/github.com/snivilised/li18ngo)](https://goreportcard.com/report/github.com/snivilised/li18ngo)
[![Coverage Status](https://coveralls.io/repos/github/snivilised/li18ngo/badge.svg?branch=master)](https://coveralls.io/github/snivilised/li18ngo?branch=master&kill_cache=1)
[![Li18ngo Continuous Integration](https://github.com/snivilised/li18ngo/actions/workflows/ci-workflow.yml/badge.svg)](https://github.com/snivilised/li18ngo/actions/workflows/ci-workflow.yml)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![A B](https://img.shields.io/badge/commit-conventional-commits?style=flat)](https://www.conventionalcommits.org/)

<!-- MD013/Line Length -->
<!-- MarkDownLint-disable MD013 -->

<!-- MD014/commands-show-output: Dollar signs used before commands without showing output mark down lint -->
<!-- MarkDownLint-disable MD014 -->

<!-- MD033/no-inline-html: Inline HTML -->
<!-- MarkDownLint-disable MD033 -->

<!-- MD040/fenced-code-language: Fenced code blocks should have a language specified -->
<!-- MarkDownLint-disable MD040 -->

<!-- MD028/no-blanks-blockquote: Blank line inside blockquote -->
<!-- MarkDownLint-disable MD028 -->

<p align="left">
  <a href="https://go.dev"><img src="resources/images/go-logo-light-blue.png" width="50" alt="go dev" /></a>
</p>

## 🔰 Introduction

This project contains helpers to aid in the development of libraries and programs that require i18n support. It doesn't perform i18n itself, that is delegated to [go-i18n](https://github.com/nicksnyder/go-i18n). Rather it aims to provide functionality that makes using go-i18n easier. For example, implementing localised error messages can be a bit tedious and error prone so included within this module is a cli app, ___lingo___, that can generate all error related code (PS, this has not been implemented yet, so is still in the works).

## 📚 Usage

## 🎀 Features

<p align="left">
  <a href="https://onsi.github.io/ginkgo/"><img src="https://onsi.github.io/ginkgo/images/ginkgo.png" width="100" alt="ginkgo" /></a>
  <a href="https://onsi.github.io/gomega/"><img src="https://onsi.github.io/gomega/images/gomega.png" width="100" alt="gomega" /></a>
</p>

+ unit testing with [Ginkgo](https://onsi.github.io/ginkgo/)/[Gomega](https://onsi.github.io/gomega/)
+ i18n with [go-i18n](https://github.com/nicksnyder/go-i18n)
+ linting configuration and pre-commit hooks, (see: [linting-golang](https://freshman.tech/linting-golang/)).

### 🌐 l10n Translations

This template has been setup to support localisation. The default language is `en-GB` with support for `en-US`. There is a translation file for `en-US` defined as __src/i18n/deploy/arcadia.active.en-US.json__. This is the initial translation for `en-US` that should be deployed with the app.

Make sure that the go-i18n package has been installed so that it can be invoked as cli, see [go-i18n](https://github.com/nicksnyder/go-i18n) for installation instructions.

To maintain localisation of the application, the user must take care to implement all steps to ensure translate-ability of all user facing messages. Whenever there is a need to add/change user facing messages including error messages, to maintain this state, the user must:

+ define template struct (__xxxTemplData__) in __src/i18n/messages.go__ and corresponding __Message()__ method. All messages are defined here in the same location, simplifying the message extraction process as all extractable strings occur at the same place. Please see [go-i18n](https://github.com/nicksnyder/go-i18n) for all translation/pluralisation options and other regional sensitive content.

For more detailed workflow instructions relating to i18n, please see [i18n README](./resources/doc/i18n-README.md)
