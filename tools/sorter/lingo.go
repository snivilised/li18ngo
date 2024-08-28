package sorter

import (
	"encoding/json"
	"sort"
)

const (
	perm = 0o644
)

// When we need to use a second command, we'll need to use cobra to enable
// sub-commands. Just copy arcadia into this code. Can't reuse arcadia or
// cobrass as that would cause cyclic dependency chain, because they require
// this module.

type MessageEntry struct {
	Description string `json:"description"`
	Other       string `json:"other"`
}

// When we add the ability to sort json file with Hashes, we'll need to
// provide an extract flag to indicate which type of file is being sorted, ie:
// - by default, we assume the native file which does not contain hashes, so
// we use MessageEntry.
// - if -t is specified, this indicates that the json file is a translation file
// created by the merge command which contains the hashes, so we use MessageEntryT
// (not yet defined, but will contain the same fields as MessageEntry but with
// an extra string member Hash). This will be tacked onto the update task in the
// same way the sort is invoked as part of the extract task.
func Sort(data []byte) ([]byte, error) {
	var messages map[string]MessageEntry
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(messages))
	for k := range messages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedMessages := make(map[string]MessageEntry)
	for _, k := range keys {
		sortedMessages[k] = messages[k]
	}

	return json.MarshalIndent(sortedMessages, "", "  ")
}
