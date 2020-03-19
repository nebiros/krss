package output

import "encoding/json"

type UserOutput struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

func (u UserOutput) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserOutput) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}

	return nil
}
