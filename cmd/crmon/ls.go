package main

import (
  "context"
  "fmt"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "github.com/urfave/cli"
)

var lsCommand = cli.Command{
  Name:   "ls",
  Usage:  "list all docker containers",
  Action: ls,
}

func ls(c *cli.Context) error {
  ctx := context.Background()

  cli, err := client.NewEnvClient()
  if err != nil {
    return err
  }
  ver, err := cli.ServerVersion(ctx)
  if err != nil {
    return err
  }
  fmt.Println(cli.ClientVersion(), ver.Version, ver.APIVersion)

  containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
    All: true,
  })
  if err != nil {
    return err
  }

  for _, ctn := range containers {
    fmt.Printf("%s %s %#v\n", ctn.ID[:10], ctn.Image, ctn.State)
    if ctn.State == "running" {
      continue
    }
    // delete container
    if err = cli.ContainerRemove(ctx, ctn.ID, types.ContainerRemoveOptions{}); err != nil {
      fmt.Println(err)
      return err
    }
  }
  return nil
}
