package sessions

import (
	"bytes"
	"encoding/json"
)

type serializedSession struct {
	Id     string `json:"id"`
	UserId int    `json:"userid"`
}

func serializeSession(s Session) ([]byte, error) {
	ss := serializedSession{
		Id:     s.ID,
		UserId: s.UserId,
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(ss); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func deserializeSession(b []byte) (Session, error) {
	var ss serializedSession
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&ss); err != nil {
		return Session{}, err
	}

	return Session{
		ID:     ss.Id,
		UserId: ss.UserId,
	}, nil
}
