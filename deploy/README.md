# AutoRotate-VPN/deploy

## Usage

```go
package main

import (
  "log"
  "github.com/Th3-S1lenc3/AutoRotate/deploy"
)

func main() {

  d := NewDeploy()

  // var apiToken string
  // var debug bool
  // var configFile string
  // var dryRun bool

  // Initialize
  err := d.Init(apiToken, debug)
  if err != nil {
    log.Fatal(err)
  }

  // Load Config From File
  instanceOptions, err := d.LoadConfigFromFile(configFile)
  if err != nil {
    log.Fatal(err)
  }

  // Or Manually Construct
  instanceOptions := deploy.InstanceCreateOptions{
    Region: "us-central",
  	Type: "g6-nanode-1",
  	Label: "test",
  	RootPass: "someSecurePassword",
  	StackScriptID: 123456,
  	StackScriptData map[string]string{
      "username": "demonstration",
    },
  	Image: "linode/debian11",
  	Tags: []string{
      "demonstration",
    },
  }

  res, err := d.Deploy(instanceOptions, dryRun)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(res)
}
```

### JSON Config File
```json
{
  "region": "us-central",
  "type": "g6-nanode-1",
  "label": "test",
  "root_pass": "someSecurePassword",
  "stackscript_id": "123456",
  "stackscript_data": {
    "username": "demonstration"
  },
  "image": "linode/debian11",
  "tags": [
    "demonstration"
  ]
}
```
