package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

// EcrRepositoryFormat ... ECRのイメージ名の書式
const EcrRepositoryFormat = "%s.dkr.ecr.%s.amazonaws.com/%s"

var (
	// ImageName ... イメージ名
	ImageName string

	// Env ... prod /stg
	Env string

	// Aws ... aws用config data
	Aws awsConfig

	// CommitHash ... カレントディレクトリのgit commit hash
	CommitHash string

	// ImageTags ... env, commitHash, （latest）のイメージ名を格納
	ImageTags map[string]string

	// SaveTags ... ECR lifecycleで削除を保護するタグを格納
	SaveTags = []string{"prod", "stg"}
)

type awsConfig struct {
	AccountID string
	Profile   string
	Region    string
}

// run ...flagを元に処理を実行する
func run(ctx *cli.Context) error {
	if err := SetFlagToVariable(ctx); err != nil {
		return err
	}

	// 必要な"image:tag"を初期化
	NewDockerImageTags()

	if err := DockerRun(); err != nil {
		return err
	}

	if err := EcrRun(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := GitGetCommitHash(); err != nil {
		log.Fatal(err)
	}

	app := FlagSet()
	app.Action = run
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
