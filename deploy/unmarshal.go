package deploy

import (
  "encoding/json"
  "strconv"

  linode "github.com/linode/linodego"
)

type StringInt int

type InstanceCreateOptions struct {
  Region          string                    `json:"region"`
	Type            string                    `json:"type"`
	Label           string                    `json:"label,omitempty"`
	Group           string                    `json:"group,omitempty"`
	RootPass        string                    `json:"root_pass,omitempty"`
	AuthorizedKeys  []string                  `json:"authorized_keys,omitempty"`
	AuthorizedUsers []string                  `json:"authorized_users,omitempty"`
	StackScriptID   StringInt                 `json:"stackscript_id,omitempty"`
	StackScriptData map[string]string         `json:"stackscript_data,omitempty"`
	BackupID        StringInt                 `json:"backup_id,omitempty"`
	Image           string                    `json:"image,omitempty"`
	Interfaces      []linode.InstanceConfigInterface `json:"interfaces,omitempty"`
	BackupsEnabled  bool                      `json:"backups_enabled,omitempty,string"`
	PrivateIP       bool                      `json:"private_ip,omitempty,string"`
	Tags            []string                  `json:"tags,omitempty"`

	// Creation fields that need to be set explicitly false, "", or 0 use pointers
	SwapSize *int  `json:"swap_size,omitempty"`
	Booted   *bool `json:"booted,omitempty"`
}

func (st *StringInt) UnmarshalJSON(b []byte) error {
  var item interface{}

  err := json.Unmarshal(b, &item)
  if err != nil {
    return err
  }

  switch v := item.(type) {
  case int:
    *st = StringInt(v)
  case float64:
    *st = StringInt(int(v))
  case string:
    i, err := strconv.Atoi(v)
    if err != nil {
      return err
    }
    *st = StringInt(i)
  }

  return nil
}
