package locale

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// 🧊 Root Cmd Short Description

// 🧊 Widget Cmd Short Description

// WidgetCmdShortDescTemplData
type WidgetCmdShortDescTemplData struct {
	li18ngoTemplData
}

func (td WidgetCmdShortDescTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "widget-command-short-description",
		Description: "short description for the widget command",
		Other:       "A brief description of widget command",
	}
}

// 🧊 Widget Cmd Long Description

// WidgetCmdLongDescTemplData
type WidgetCmdLongDescTemplData struct {
	li18ngoTemplData
}

func (td WidgetCmdLongDescTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "widget-command-long-description",
		Description: "long description for the widget command",
		Other: `A longer description that spans multiple lines and likely contains
		examples and usage of using your application.`,
	}
}
