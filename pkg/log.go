package pkg

import (
	"encoding/json"
	"fmt"
)

type LogEntry struct {
	Level string `json:"level"`
	Time  string `json:"time"`
	Msg   string `json:"msg"`
	Raw   string `json:"-"`
}

func UnmarshalLogEntry(data []byte) (LogEntry, error) {
	om := NewOrderedMap()

	if err := json.Unmarshal(data, &om); err != nil {
		return LogEntry{}, fmt.Errorf("new ordered map from json %q: %w", string(data), err)
	}

	it := om.EntriesIter()
	v := LogEntry{Raw: string(data)}

	for {
		entry, ok := it()
		if !ok {
			break
		}

		switch entry.Key {
		case "level":
			v.Level = fmt.Sprintf("%v", entry.Value)
		case "time":
			v.Time = fmt.Sprintf("%v", entry.Value)
		case "msg":
			v.Msg = fmt.Sprintf("%v", entry.Value)
		}
	}

	return v, nil
}
