package main

import (
	"errors"
	"log"

	"github.com/urfave/cli"
)

// FlagSet ... flagを設定する
func FlagSet() *cli.App {
	app := cli.NewApp()
	app.Name = "ecr-builder"
	app.Description = "docker build/push to ecr"
	app.Usage = "ecr-builder -n hoge --profile fuga -e stg -a 11111111 -t 1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name, n",
			Usage: "image name",
		},
		cli.StringFlag{
			Name:  "env, e",
			Usage: "stg / prod",
		},
		cli.StringFlag{
			Name:  "region, r",
			Usage: "default ap-northeast-1",
			Value: "ap-northeast-1",
		},
		cli.StringFlag{
			Name:  "profile, p",
			Usage: "AWS Profile",
		},
		cli.IntFlag{
			Name:  "account, a",
			Usage: "aws account id",
		},
		cli.StringSliceFlag{
			Name:  "tags, t",
			Usage: "[-t hoge -t fuga] save tags. prod / stg are normally not deleted.",
			Value: &cli.StringSlice{},
		},
	}
	return app
}

// SetFlagToVariable ... 設定したflagをvariableにsetする
func SetFlagToVariable(ctx *cli.Context) error {
	ImageName = ctx.String("name")
	Env = ctx.String("env")
	Aws.AccountID = ctx.String("account")
	Aws.Profile = ctx.String("profile")
	Aws.Region = ctx.String("region")
	SaveTags = append(SaveTags, ctx.StringSlice("tags")...)

	if ImageName == "" {
		return errors.New("required -n, --name option")
	}

	if Env == "" {
		return errors.New("required -e, --env option")
	}

	if Aws.Profile == "" {
		return errors.New("required -p, --profile option")
	}

	if Aws.AccountID == "" {
		return errors.New("required -a, --account option")
	}

	log.Printf("save ECR tags: %s\n", SaveTags)
	return nil
}
