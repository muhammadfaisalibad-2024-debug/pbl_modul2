package helper

import (
	"api-students/app/model"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidCursor = errors.New("invalid cursor")

func EncodeCursor(t time.Time, id int) string {
	raw := strconv.FormatInt(t.UTC().UnixNano(), 10) + "|" + strconv.Itoa(id)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}
func DecodeCursor(encoded string) (model.Cursor, error) {
	b, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		return model.Cursor{}, ErrInvalidCursor
	}
	p := strings.Split(string(b), "|")
	if len(p) != 2 {
		return model.Cursor{}, ErrInvalidCursor
	}
	n, e1 := strconv.ParseInt(p[0], 10, 64)
	id, e2 := strconv.Atoi(p[1])
	if e1 != nil || e2 != nil || id < 1 {
		return model.Cursor{}, ErrInvalidCursor
	}
	return model.Cursor{CreatedAt: time.Unix(0, n).UTC(), ID: id}, nil
}
