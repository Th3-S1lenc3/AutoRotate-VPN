package main

import (
  "log"
  "fmt"
  "flag"

  "github.com/Th3-S1lenc3/AutoRotate-VPN/deploy"
)

func main() {
  d := deploy.NewDeploy()

  debug := flag.Bool("d", false, "Enable Debug Mode")
  apiToken := flag.String("apiToken", "", "Linode API Token")
  configFile := flag.String("f", "", "Path to JSON Config File")
  dryRun := flag.Bool("dryRun", false, "Dry Run")

  flag.Parse()

  err := d.Init(*apiToken, *debug)
  if err != nil {
    log.Fatal(err)
  }

  instanceOptions, err := d.LoadConfigFromFile(*configFile)
  if err != nil {
    log.Fatal(err)
  }

  res, err := d.Deploy(instanceOptions, *dryRun)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(res)
}
