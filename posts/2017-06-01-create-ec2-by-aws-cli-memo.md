---
title: aws-cliでEC2インスタンスを作成するときのメモ
published: '2017-06-01'
---

毎回ググってる気がするのでメモ。

## 準備

#### 使用するツール

[aws-cli](https://github.com/aws/aws-cli)

#### インストール

```
$ pip install awscli
```

#### プロファイル設定

```
$ aws configure --profile ${プロファイル名｝
```

アクセスキーID、シークレットアクセスキー、デフォルトリージョン、出力フォーマットを指定。
出力フォーマットはだいたい `json` にしてる。 `jq` コマンドで使うので。


## EC2インスタンス作成

### キーペア作成
（必要であれば）

``` 
$ aws --profile=${プロファイル名｝ ec2 create-key-pair \
    --key-name ${新規キー名} \
    --query 'KeyMaterial' \
    --output text \
    > $HOME/.ssh/${保存する秘密鍵のファイル名} \
    && chmod 400 $HOME/.ssh/${保存する秘密鍵のファイル名}
```

### インスタンス作成

*設定するSecurityGroupとSubnetを事前に調べておく*
```
$ aws --profile=${プロファイル名｝ ec2 describe-security-groups | jq -c ".SecurityGroups[] | [.GroupName,.GroupId]"
$ aws --profile=${プロファイル名｝ ec2 describe-subnets | jq -c ".Subnets[] | [.Tags[0].Value,.SubnetId]"
```

```
aws --profile=${プロファイル名｝ ec2 run-instances \
	--image-id ami-923d12f5 \
	--instance-type t2.micro \
	--key-name ${新規キー名} \
	--security-group-ids ${SecurityGroupのGroupID} \
	--subnet-id ${SubnetのSubnetID} \
	--monitoring Enabled=true
```

インスタンスIDを控えておく。

### インスタンスに名前を付ける

```
$ aws --profile=${プロファイル名} ec2 create-tags --resources ${インスタンスID} --tags Key=Name,Value="${インスタンス名}"
```
