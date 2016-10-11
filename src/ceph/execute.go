package ceph

import (
  _ "fmt"
  "log"
  "os/exec"
  "encoding/json"
  "strings"
)

func runCephCommand(command string, out interface{}) error {

    cmd := strings.Split(command, " ")

    command_out, err := exec.Command(cmd[0], cmd[1:]...).Output()
    if err != nil { log.Fatal(err) }

    if out == nil {
      return err
    }

    err = json.Unmarshal(command_out, &out)
    if err != nil { log.Fatal(err) }

    // replace Fatal above with proper err handling
    return err
}
