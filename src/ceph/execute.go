package ceph

import (
  "log"
  "os/exec"
  "encoding/json"
)

// check out if nil, if it is don't unmarshal, just return
func runCephCommand(command []string, out interface{}) error {

    command_out, err := exec.Command(command[0], command[1:]...).Output()
    if err != nil { log.Fatal(err) }

    if out == nil {
      return err
    }

    err = json.Unmarshal(command_out, &out)
    if err != nil { log.Fatal(err) }

    // replace Fatal above with proper err handling
    return err
}
