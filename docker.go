package main

import (
	"fmt"
	"log"
	"os/exec"
)

// DockerBuild ... カレントのDockerfileを元にenvタグを付けてbuildする
func DockerBuild() error {
	if out, err := exec.Command(
		"docker",
		"build",
		"--pull",
		"-t",
		ImageTags["env"],
		".",
	).CombinedOutput(); err != nil {
		return fmt.Errorf("docker build: %s: %s", err, string(out))
	}

	fmt.Println()

	log.Printf("docker build done: %s\n", ImageTags["env"])

	fmt.Println()

	return nil
}

// DockerSetTag ... 必要なtagを付与する
// Env=stg : stg, commitHash
// Env=prod: prod, commitHash, latest
func DockerSetTag() error {
	if AlreadyBuildImage != "" {
		if out, err := exec.Command(
			"docker",
			"tag",
			AlreadyBuildImage,
			ImageTags["env"],
		).CombinedOutput(); err != nil {
			return fmt.Errorf("docker tag: %s: %s", err, string(out))
		}
		log.Printf("docker SetTag from %s done: %s\n", AlreadyBuildImage, ImageTags["commitHash"])
	}

	if out, err := exec.Command(
		"docker",
		"tag",
		ImageTags["env"],
		ImageTags["commitHash"],
	).CombinedOutput(); err != nil {
		return fmt.Errorf("docker tag: %s: %s", err, string(out))
	}
	log.Printf("docker SetTag from %s done: %s\n", ImageTags["env"], ImageTags["commitHash"])

	// For prod, tagging with latest
	if _, ok := ImageTags["latest"]; ok {
		if out, err := exec.Command(
			"docker",
			"tag",
			ImageTags["env"],
			ImageTags["latest"],
		).CombinedOutput(); err != nil {
			return fmt.Errorf("docker tag: %s: %s", err, string(out))
		}
		log.Printf("docker SetTag done: %s\n", ImageTags["latest"])
	}

	fmt.Println()

	return nil
}

// DockerRun ... Docker処理を実行
func DockerRun() error {
	if AlreadyBuildImage == "" {
		if err := DockerBuild(); err != nil {
			return err
		}
	}

	if err := DockerSetTag(); err != nil {
		return err
	}

	return nil
}
