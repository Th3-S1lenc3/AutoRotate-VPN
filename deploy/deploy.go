package deploy

import (
  "context"
  "fmt"
  "encoding/json"
  "net/http"
  "os"
  "io/ioutil"

  linode "github.com/linode/linodego"
  "golang.org/x/oauth2"
  "github.com/imdario/mergo"
)

type Deploy struct {
  linodeClient linode.Client
  debug bool
}

func NewDeploy() *Deploy {
  return &Deploy{}
}

func (d *Deploy) Deploy(options InstanceCreateOptions, dryRun bool) (string, error) {
  instanceOptions, err := d.convertToLinode(options)
  if err != nil {
    return "", err
  }

  err = d.validateInstanceOptions(instanceOptions)
  if err != nil {
    return "", err
  }

  stackscriptData := instanceOptions.StackScriptData

  stackscriptUDFs, err := d.getStackscriptUDFs(instanceOptions.StackScriptID)
  if err != nil {
    return "", err
  }

  udfData := d.parseStackscriptUDFs(*stackscriptUDFs)

  err = mergo.Merge(&udfData, stackscriptData, mergo.WithOverride)
  if err != nil {
    return "", err
  }

  instanceOptions.StackScriptData = udfData

  if dryRun == false {
    res, err := d.linodeClient.CreateInstance(context.Background(), instanceOptions)
    if err != nil {
      return "", err
    }

    json, err := json.MarshalIndent(res, "", "  ")
    if err != nil {
      return "", nil
    }

    return string(json), nil
  }

  return fmt.Sprintf("%v", instanceOptions), nil
}

func (d *Deploy) getStackscriptUDFs(id int) (*[]linode.StackscriptUDF, error) {
  stackscript, err := d.linodeClient.GetStackscript(context.Background(), id)
  if err != nil {
    return &[]linode.StackscriptUDF{}, err
  }

  udfs := stackscript.UserDefinedFields

  return udfs, nil
}

func (d *Deploy) parseStackscriptUDFs(stackscriptUDFs []linode.StackscriptUDF) map[string]string {
  udfs := make(map[string]string)

  for i := 0; i < len(stackscriptUDFs); i++ {
    udf := stackscriptUDFs[i]

    if udf.Default != "" {
      udfs[udf.Name] = udf.Default
    }
  }

  return udfs
}

func (d *Deploy) validateInstanceOptions(options linode.InstanceCreateOptions) error {
  if options.Region == "" {
    return fmt.Errorf("Invalid Options. Region Required.")
  }

  if options.Type == "" {
    return fmt.Errorf("Invalid Options. Type Required.")
  }

  regionValid := false
  typeVaild := false
  imageValid := false

  regions, err := d.linodeClient.ListRegions(context.Background(), nil)
  if err != nil {
    return err
  }

  for i := 0; i < len(regions); i++ {
    if regions[i].ID == options.Region {
      regionValid = true
    }
  }

  types, err := d.linodeClient.ListTypes(context.Background(), nil)
  if err != nil {
    return err
  }

  for i := 0; i < len(types); i++ {
    if types[i].ID == options.Type {
      typeVaild = true
    }
  }

  images, err := d.linodeClient.ListImages(context.Background(), nil)
  if err != nil {
    return err
  }

  for i := 0; i < len(images); i++ {
    if images[i].ID == options.Image {
      imageValid = true
    }
  }


  if regionValid == false {
    return fmt.Errorf("Invalid Options. Invalid Region: %s", options.Region)
  }

  if typeVaild == false {
    return fmt.Errorf("Invalid Options. Invalid Type: %s", options.Type)
  }

  if imageValid == false {
    return fmt.Errorf("Invalid Options. Invalid Image: %s", options.Image)
  }

  return nil
}

func (d *Deploy) convertToLinode(instanceOptions InstanceCreateOptions) (linode.InstanceCreateOptions, error) {
  data, err := json.Marshal(instanceOptions)
  if err != nil {
    return linode.InstanceCreateOptions{}, err
  }

  linodeInstanceOptions := linode.InstanceCreateOptions{}
  err = json.Unmarshal(data, &linodeInstanceOptions)
  if err != nil {
    return linode.InstanceCreateOptions{}, err
  }

  return linodeInstanceOptions, nil
}

func (d *Deploy) LoadConfigFromFile(configFile string) (InstanceCreateOptions, error) {
  if configFile == "" {
    return InstanceCreateOptions{}, fmt.Errorf("Config File Required")
  }
	_, err := os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		return InstanceCreateOptions{}, fmt.Errorf("Cannot find file: %s", configFile)
	}

  data, err := ioutil.ReadFile(configFile)
  if err != nil {
    return InstanceCreateOptions{}, err
  }

  instanceOptions := InstanceCreateOptions{}
  err = json.Unmarshal(data, &instanceOptions)
  if err != nil {
    return InstanceCreateOptions{}, err
  }

  return instanceOptions, nil
}

func (d *Deploy) Init(apiToken string, debug bool) error {
  if apiToken == "" {
    return fmt.Errorf("Invalid API Token.")
  }

  tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})

  oauth2Client := &http.Client{
    Transport: &oauth2.Transport{
      Source: tokenSource,
    },
  }

  d.linodeClient = linode.NewClient(oauth2Client)
  d.linodeClient.SetDebug(debug)

  d.debug = debug

  return nil
}
