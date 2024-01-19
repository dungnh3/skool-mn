package utils

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

func GenerateId() string {
	raw := int64(uuid.New().ID())
	base36 := strconv.FormatInt(raw, 36)
	if len(base36) < 7 {
		padding := "0000000"
		base36 = padding[:7-len(base36)] + base36
	}
	now := time.Now().Format("20060102")
	return fmt.Sprintf("%v.%v", now, strings.ToUpper(base36))
}
