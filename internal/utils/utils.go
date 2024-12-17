package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func GetUUIDFromJsonString(input map[string]interface{}, idname string) (uuid.UUID, error) {	
	idstr, ok := input[idname].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("key '%s' it's not a string", idname)
	}
	id, err := uuid.Parse(idstr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error converting '%s' to UUID: %v", idstr, err)
	}

	return id, nil
}

func StringAsAPointer(s string) *string {
	return &s
}