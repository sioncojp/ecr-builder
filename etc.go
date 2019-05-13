package main

import "fmt"

// NewDockerImageTags ... 必要なDockerのイメージ:tag名をスライスに入れる
// 0: envタグ
// 1: commitHashタグ
// 2: prodの場合はlatestタグ
func NewDockerImageTags() {
	ImageTags = map[string]string{}

	i := fmt.Sprintf(EcrRepositoryFormat, Aws.AccountID, Aws.Region, ImageName)

	ImageTags["env"] = fmt.Sprintf("%s:%s", i, Env)
	ImageTags["commitHash"] = fmt.Sprintf("%s:%s", i, CommitHash)
	if Env == "prod" {
		ImageTags["latest"] = fmt.Sprintf("%s:%s", i, "latest")
	}
}

// containsSlice ... sliceの中に指定されたsliceが含まれてるか
func isContainsSlice(s []string, ss []string) bool {
	for _, v := range s {
		for _, vv := range ss {
			if v == vv {
				return true
			}
		}
	}
	return false
}
