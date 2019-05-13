package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// ECR ... Store ECR data.
type ECR struct {
	Session *ecr.ECR
	Auth    EcrAuth
}

// EcrAuth ... Store a auth data of the ECR to login.
type EcrAuth struct {
	Token         string
	User          string
	Pass          string
	ProxyEndpoint string
	ExpiresAt     time.Time
}

// EcrLogin ... ECRにログイン
func (e *ECR) EcrLogin() error {
	token, err := e.GetAuthorizationToken()
	if err != nil {
		return err
	}

	e.Auth, err = NewECRAuth(token)
	if err != nil {
		return err
	}

	if err := ecrLogin(e.Auth); err != nil {
		return err
	}

	log.Println("ecr Login done")

	fmt.Println()

	return nil
}

// ecrLogin ... login ECR repository
func ecrLogin(auth EcrAuth) error {
	out, err := exec.Command(
		"docker",
		"login",
		"-u",
		auth.User,
		"-p",
		auth.Pass,
		auth.ProxyEndpoint,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker login: %s: %s", err, string(out))
	}

	return nil
}

// EcrLifeCycle ... ECRのライフサイクル。削除する
// Env=stg : stgがついてるイメージ以外は削除する
// Env=prod: prodがついてるイメージ以外は削除する
// 両方: -t, --tagsで指定されたタグ以外は削除する
func (e *ECR) EcrLifeCycle() error {
	// Get ECR all Image
	describeInput := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(ImageName),
	}
	images, err := e.Session.DescribeImages(describeInput)
	if err != nil {
		return err
	}

	// Delete用データを初期化
	deleteInput := &ecr.BatchDeleteImageInput{
		RepositoryName: aws.String(ImageName),
	}

	// 削除するtagの抽出
	for _, v := range images.ImageDetails {
		tags := aws.StringValueSlice(v.ImageTags)

		// SaveTagsに入ってるもの以外を抽出
		if !(isContainsSlice(tags, SaveTags)) {
			deleteInput.ImageIds = append(deleteInput.ImageIds,
				&ecr.ImageIdentifier{
					ImageDigest: v.ImageDigest,
				})
		}
	}

	// イメージが何もなければdelete処理を行わない
	if len(deleteInput.ImageIds) == 0 {
		return nil
	}

	// Delete
	if _, err := e.Session.BatchDeleteImage(deleteInput); err != nil {
		return err
	}

	return nil
}

// EcrPush ... ECRにpushする
func (e *ECR) EcrPush() error {
	for _, i := range ImageTags {
		if i != "" {
			if out, err := exec.Command(
				"docker",
				"push",
				i,
			).CombinedOutput(); err != nil {
				return fmt.Errorf("docker push: %s: %s", err, string(out))
			}
			log.Printf("docker push done: %s\n", i)
		}
	}

	return nil
}

// EcrRun ... ECRの処理を実行
func EcrRun() error {
	e := &ECR{}

	// ECR用のSession初期化
	err := e.NewSession()
	if err != nil {
		return err
	}

	if err := e.EcrLogin(); err != nil {
		return err
	}

	// 先にlifecycleを入れることで、前のイメージを残す
	if err := e.EcrLifeCycle(); err != nil {
		return err
	}

	if err := e.EcrPush(); err != nil {
		return err
	}

	return nil
}

// GetAuthorizationToken ... Retrieves a token that is valid for a specified registry for 12 hours for ECR.
func (e *ECR) GetAuthorizationToken() (*ecr.GetAuthorizationTokenOutput, error) {
	input := &ecr.GetAuthorizationTokenInput{
		RegistryIds: GetRegistryIds(),
	}

	result, err := e.Session.GetAuthorizationToken(input)
	if err != nil {
		return nil, fmt.Errorf("authorizing: %s", err)
	}
	return result, nil
}

// GetRegistryIds ... AwsAccountIDをGetしてくる
func GetRegistryIds() []*string {
	var registryIds []*string
	registryIds = append(registryIds, aws.String(Aws.AccountID))

	return registryIds
}

// NewECRAuth ... ECRのAuthrizationデータをstructに入れる
func NewECRAuth(token *ecr.GetAuthorizationTokenOutput) (EcrAuth, error) {
	auth := token.AuthorizationData[0]
	data, err := base64.StdEncoding.DecodeString(*auth.AuthorizationToken)
	if err != nil {
		return EcrAuth{}, fmt.Errorf("decode to base64: %s", err)
	}
	// extract username and password
	t := strings.SplitN(string(data), ":", 2)

	// object to pass to template
	a := EcrAuth{
		Token:         *auth.AuthorizationToken,
		User:          t[0],
		Pass:          t[1],
		ProxyEndpoint: *(auth.ProxyEndpoint),
		ExpiresAt:     *(auth.ExpiresAt),
	}

	return a, nil
}

// NewSession ... ECR用Sessionの初期化
func (e *ECR) NewSession() error {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           Aws.Profile,
	})
	if err != nil {
		return err
	}
	e.Session = ecr.New(sess, aws.NewConfig().WithMaxRetries(10).WithRegion(Aws.Region))
	return nil
}
