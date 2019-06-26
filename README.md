# ecr-builder

```bash
### help
$ ecr-builder help
NAME:
   ecr-builder - ecr-builder -n hoge --profile fuga -e stg -a 11111111 -t 1.0

USAGE:
   ecr-builder [global options] command [command options] [arguments...]

VERSION:
   0.0.0

DESCRIPTION:
   docker build/push to ecr

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --name value, -n value        image name
   --env value, -e value         ex. stg / prod
   --region value, -r value      default ap-northeast-1 (default: "ap-northeast-1")
   --profile value, -p value     AWS Profile
   --skip-build value, -s value  if skip-build, please specify the docker image already build
   --account value, -a value     aws account id (default: 0)
   --tags value, -t value        [-t hoge -t fuga] save tags. prod / stg are normally not deleted.
   --help, -h                    show help
   --version, -v                 print the version


### build & push: test image
$ ecr-builder --name test --profile stg --account 11111111 --env stg -t 1.0.0
2019/05/13 21:54:29 save ECR tags: [prod stg 1.0.0]


2019/05/13 21:54:32 docker build done: 11111111.dkr.ecr.ap-northeast-1.amazonaws.com/test:prod

2019/05/13 21:54:32 docker SetTag done: 11111111.dkr.ecr.ap-northeast-1.amazonaws.com/test:xxxxxxxxxxxxxxxxxxxxx

2019/05/13 21:54:33 ecr Login done

2019/05/13 21:54:34 docker push done: 11111111.dkr.ecr.ap-northeast-1.amazonaws.com/test:prod
2019/05/13 21:54:35 docker push done: 11111111.dkr.ecr.ap-northeast-1.amazonaws.com/test:xxxxxxxxxxxxxxxxxxx
```


## overview

ecr-builder will do the below

1. build from current Dockerfile
2. set tags
    - stg, commitHash
    - prod, commitHash, latest
3. ECR login
4. ECR lifecycle(remove except prod, stg and -t options tags.)
5. ECR push


## Development

```
make help
make lint
make build
```
