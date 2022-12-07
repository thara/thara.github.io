---
title: aws-cliでSecurityGroupを作成するときのメモ
---

これも毎回ググってる気がするのでメモ。

```sh
$ GROUP_ID=`aws --profile=${プロファイル名｝ ec2 create-security-group \
    --group-name ${新規セキュリティグループ名} \
    --description ${説明} --vpc-id ${VPCのID} | jq -r ".GroupId"`

```

全開になっている22ポートを塞ぐ・・・が、そんなルールがないと怒られた。    
昔はデフォで22空いてたような気がしたけど、気のせい・・・？

```sh
$ aws --profile=${プロファイル名｝ ec2 revoke-security-group-ingress \
    --group-id $GROUP_ID --protocol tcp --port 22 --cidr "0.0.0.0/0"
```


一例）　`192.168.1.0/24` からの22ポートと80ポートへのアクセスを許可

```sh
$ aws --profile=${プロファイル名｝ ec2 authorize-security-group-ingress --group-id $GROUP_ID \
    --protocol tcp --port 22 --cidr "192.168.1.0/24"
$ aws --profile=${プロファイル名｝ ec2 authorize-security-group-ingress --group-id $GROUP_ID \
    --protocol tcp --port 80 --cidr "192.168.1.0/24"
```
