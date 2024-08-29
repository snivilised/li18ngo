package sorter

import (
	"encoding/json"
	"sort"
)

// When we need to use a second command, we'll need to use cobra to enable
// sub-commands. Just copy arcadia into this code. Can't reuse arcadia or
// cobrass as that would cause cyclic dependency chain, because they require
// this module.

type MessageEntry struct {
	Description string `json:"description"`
	Other       string `json:"other"`
}

type HashedMessageEntry struct {
	MessageEntry
	Hash string `json:"hash"`
}

func Apply[T MessageEntry | HashedMessageEntry](data []byte) ([]byte, error) {
	var messages map[string]T
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(messages))
	for k := range messages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make(map[string]T)
	for _, k := range keys {
		sorted[k] = messages[k]
	}

	return json.MarshalIndent(sorted, "", "  ")
}
