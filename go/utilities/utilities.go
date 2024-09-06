package utilities

import (
	"encoding/json"
	"reflect"
)

// PatchJson takes two JSON strings and merges the second JSON string into the first JSON string
func PatchJSON(originalJSON, patchJSON string) string {
	var originalJSONUnMarshalled map[string]interface{}
	originalJSONUnmarshalErr := json.Unmarshal([]byte(originalJSON), &originalJSONUnMarshalled)
	if originalJSONUnmarshalErr != nil {
		panic("Error unmarshalling original JSON, please verify the JSON input.")
	}
	var patchJSONUnMarshalled map[string]interface{}
	patchJSONUnmarshalErr := json.Unmarshal([]byte(patchJSON), &patchJSONUnMarshalled)
	if patchJSONUnmarshalErr != nil {
		panic("Error unmarshalling patch JSON, please verify the JSON input.")
	}
	for key, _ := range originalJSONUnMarshalled {
		mergeJson(key, originalJSONUnMarshalled[key], patchJSONUnMarshalled[key], patchJSONUnMarshalled)
	}
	mergedJSON, mergedJSONErr := json.Marshal(patchJSONUnMarshalled)
	if mergedJSONErr != nil {
		panic("Error marshalling merged JSON")
	}
	return string(mergedJSON)
}

func mergeJson(key string, originalMap interface{}, patchedMap interface{}, mergedMap map[string]interface{}) {
	if !reflect.DeepEqual(originalMap, patchedMap) {
		switch originalMap.(type) {
		// nested map
		case map[string]interface{}:
			// nothing to patch
			if patchedMap == nil {
				mergedMap[key] = originalMap
			} else if _, ok := patchedMap.(map[string]interface{}); ok { // check if the patched map value is also a map[string]interface and recurse
				// cast interfaces to map[string]interface{}
				originalMap := originalMap.(map[string]interface{})
				patchedMap := patchedMap.(map[string]interface{})
				for newKey, _ := range originalMap {
					mergeJson(newKey, originalMap[newKey], patchedMap[newKey], mergedMap[key].(map[string]interface{}))
				}
			}
		default:
			if patchedMap == nil {
				mergedMap[key] = originalMap
			}
		}
	}
}
