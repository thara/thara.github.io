---
layout: post
title: Agones動作確認 2019-03現在
description: Agonesの最新をローカルで動かそうと思ったが予想以上にハマってしまったので、2019-03現段階でのサンプルを動かすまでをメモっておく。
tags: [agones, kubernetes]
---

久しぶりにAgonesを調べてみたら、[Webサイト](https://agones.dev/site/) ができていた。

Agones本体も、ゲームサーバーの情報が得られるようになっていたり、Node.jsのSDKをサポートしたりと色々変わっていたので、Agonesの最新をローカルで動かそうと思ったが、予想以上にハマってしまったので、
[以前のAgones紹介記事](../agones-introduction) と少しかぶってしまうが、2019-03現段階でのサンプルを動かすまでをメモっておく。

## Agonesのインストール方法

Agonesのインストール方法には、 Google Kubernetes Engine (GKE)を使うか、Minikubeを使うか、Amazon Web Services EKSを使うかの大きく3通りがあるが、自分は動作確認がしたいだけなので、Minikubeを選択する。
その他2つは、まだ試してない。

MinikubeではVMドライバを選択できるが、公式の[Installation](https://agones.dev/site/docs/installation/#setting-up-a-minikube-cluster)どおり、VirtualBoxを使ったほうが変にハマらなくて済むと思う。   
このとき、VirtualBoxは最新のバージョンに更新しておいた方がよい。
自分はVirtualBoxのバージョンが古くて、後述のセットアップコマンドが正常に動作しなかったので、VirtualBoxのバージョンを更新して `6.0.4 r128413` にした。

そして、もう一つ決めておいた方がよいのが、`kubectl apply -f` に渡すinstall.yamlをどれにするか。
これは、Agonesの動作確認をどのような用途で行いたいかによる。

単に、Agonesというサービスを手元で動作確認したいだけであれば、公式ドキュメントのとおり、GitHubリポジトリ上においてある `install/yaml/install.yaml` を使えばよい。

自分のようにsimple-udpなどのサンプルだけではなく、SDK周りのサンプルを触ってみたい場合は、Agonesリポジトリをgit cloneした上で、その作業ツリー上の `install/yaml/install.yaml` を使う。

このとき、masterではなく、リリースブランチを使用した方がよい。   
masterのinstall.yml内のテンプレートに用いるimageが、Agonesが使用するDockerレジストリ(`gcr.io/agones-images`)上にまだデプロイされていないことがあるからだ。
[開発者向けREADME](https://github.com/GoogleCloudPlatform/agones/blob/master/build/README.md)通りにやれば自分でイメージをビルドしてAgonesを立ち上げることもできるが、
かなりの遠回りになるので、Agones自体の開発にコントリビュートしたいのでなければ避けておいた方が無難だ。


```bash
$ git clone https://github.com/GoogleCloudPlatform/agones.git && cd agones && git checkout release-0.9.0-rc

$ minikube profile agones
-   minikube profile was successfully set to agones

$ minikube start --kubernetes-version v1.11.0 --vm-driver virtualbox \
		--extra-config=apiserver.authorization-mode=RBAC
o   minikube v0.35.0 on darwin (amd64)
>   Creating virtualbox VM (CPUs=2, Memory=2048MB, Disk=20000MB) ...
-   "agones" IP address is 192.168.99.102
-   Configuring Docker as the container runtime ...
-   Preparing Kubernetes environment ...
    - apiserver.authorization-mode=RBAC
-   Pulling images required by Kubernetes v1.11.0 ...
-   Launching Kubernetes v1.11.0 using kubeadm ...
:   Waiting for pods: apiserver proxy etcd scheduler controller addon-manager dns
-   Configuring cluster permissions ...
-   Verifying component health .....
+   kubectl is now configured to use "agones"
=   Done! Thank you for using minikube!

$ kubectl create clusterrolebinding cluster-admin-binding \
  --clusterrole=cluster-admin --serviceaccount=kube-system:default
clusterrolebinding.rbac.authorization.k8s.io "cluster-admin-binding" created

$ kubectl create namespace agones-system
namespace "agones-system" created

$ kubectl apply -f install/yaml/install.yaml
serviceaccount "agones-controller" created
clusterrole.rbac.authorization.k8s.io "agones-controller" created
clusterrolebinding.rbac.authorization.k8s.io "agones-controller-access" created
serviceaccount "agones-sdk" created
clusterrole.rbac.authorization.k8s.io "agones-sdk" created
rolebinding.rbac.authorization.k8s.io "agones-sdk-access" created
customresourcedefinition.apiextensions.k8s.io "fleets.stable.agones.dev" created
customresourcedefinition.apiextensions.k8s.io "fleetallocations.stable.agones.dev" created
customresourcedefinition.apiextensions.k8s.io "fleetautoscalers.stable.agones.dev" created
customresourcedefinition.apiextensions.k8s.io "gameservers.stable.agones.dev" created
customresourcedefinition.apiextensions.k8s.io "gameserverallocations.stable.agones.dev" created
customresourcedefinition.apiextensions.k8s.io "gameserversets.stable.agones.dev" created
service "agones-controller-service" created
deployment.apps "agones-controller" created
deployment.apps "agones-ping" created
service "agones-ping-http-service" created
service "agones-ping-udp-service" created
priorityclass.scheduling.k8s.io "agones-system" created
validatingwebhookconfiguration.admissionregistration.k8s.io "agones-validation-webhook" created
mutatingwebhookconfiguration.admissionregistration.k8s.io "agones-mutation-webhook" created
secret "agones-manual-cert" created
```

## simple-udpサンプルの動作確認

Agonesをインストールしたら、[Quickstart](https://agones.dev/site/docs/getting-started/create-gameserver/) を参考に動作確認をする。
が、ドキュメントでは先述の「Agonesというサービスを手元で動作確認したい」というケースを前提としているため、自分のように手元の `install/yaml/install.yaml` でAgonesをインストールした場合は `kubectl create -f` に作業ツリーの `gameserver.yml` を指定するように注意する。

誤って異なるソースのinstall.ymlとgameserver.ymlを用いた場合、Agones内の設定ファイルのフォーマットやプロトコルに差異が出て、不可解なエラーになることがある。

```bash
$ eval $(minikube docker-env)
$ kubectl create -f examples/simple-udp/gameserver.yaml
gameserver.stable.agones.dev "simple-udp-fkhqj" created
```

simple-udpサンプルの動作確認では最終的に `nc` コマンドでUDPパケットを送って、ACKが返ってくることを確認する。   
公式ドキュメントでは `kubectl get gs` で接続アドレスとポートが表示されるはずなのだが、自分の環境では表示されなかった。

```bash
$ kubectl get gs
NAME               AGE
simple-udp-fkhqj   1h
```

接続アドレスとポートを調べるのに、`kubectl get` でピンポイントで探そうとして、結構な時間ハマってしまった。   
愚直に `kubectl describe all` の中を探せば見つかる。

```bash
$ kubectl describe all
Name:           simple-udp-fkhqj
Namespace:      default
Node:           minikube/10.0.2.15
Start Time:     Sat, 30 Mar 2019 23:04:20 +0900
Labels:         stable.agones.dev/gameserver=simple-udp-fkhqj
                stable.agones.dev/role=gameserver
Annotations:    cluster-autoscaler.kubernetes.io/safe-to-evict=false
                stable.agones.dev/container=simple-udp
                stable.agones.dev/sdk-version=0.9.0-rc
Status:         Running
IP:             172.17.0.7
Controlled By:  GameServer/simple-udp-fkhqj
Containers:
  simple-udp:
    Container ID:   docker://a20277e8c04cb92e5b3885824f9d6b98b36488312e92ed5bf664cb4bd3653917
    Image:          gcr.io/agones-images/udp-server:0.7
    Image ID:       docker-pullable://gcr.io/agones-images/udp-server@sha256:324f398fcee52edd0dee847496b350f3717e69536a1d70ae6a22b6fd8aab8bf0
    Port:           7654/UDP
    Host Port:      7383/UDP     <<<<<<<<< ここ
    ...
  agones-gameserver-sidecar:
    ...
Conditions:
  Type              Status
  Initialized       True 
  Ready             True 
  ContainersReady   True 
  PodScheduled      True 
Volumes:
  ...
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  39m   default-scheduler  Successfully assigned default/simple-udp-fkhqj to minikube
  Normal  Pulling    39m   kubelet, minikube  pulling image "gcr.io/agones-images/udp-server:0.7"
  Normal  Pulled     39m   kubelet, minikube  Successfully pulled image "gcr.io/agones-images/udp-server:0.7"
  Normal  Created    39m   kubelet, minikube  Created container
  Normal  Started    39m   kubelet, minikube  Started container
  Normal  Pulled     39m   kubelet, minikube  Container image "gcr.io/agones-images/agones-sdk:0.9.0-rc" already present on machine
  Normal  Created    39m   kubelet, minikube  Created container
  Normal  Started    39m   kubelet, minikube  Started container


Name:              kubernetes
Namespace:         default
Labels:            component=apiserver
                   provider=kubernetes
Annotations:       <none>
Selector:          <none>
Type:              ClusterIP
IP:                10.96.0.1
Port:              https  443/TCP
TargetPort:        8443/TCP
Endpoints:         192.168.99.102:8443         <<<<<<<<< ここ
Session Affinity:  None
Events:            <none>
```

上で「ここ」と示した `Endpoints` のIP部と `Host Port` のポートを使用する。

```bash
$ nc -u 192.168.99.102 7383
hello
ACK: hello
test
ACK: test
GAMESERVER
NAME: simple-udp-fkhqj
EXIT
```


ゲームサーバ側のログを見ながらだと、挙動がわかりやすい。



```bash
$ kubectl logs simple-udp-fkhqj -c simple-udp -f
2019/03/30 14:04:25 Starting UDP server, listening on port 7654
2019/03/30 14:04:25 Creating SDK instance
2019/03/30 14:04:25 Starting Health Ping
2019/03/30 14:04:25 Marking this server as ready
2019/03/30 14:40:06 Received packet from 192.168.99.1:49977: hello
2019/03/30 14:40:49 Received packet from 192.168.99.1:49977: test
2019/03/30 14:41:32 Received packet from 192.168.99.1:49977: GAMESERVER
2019/03/30 14:41:32 GameServer: {"object_meta":{"name":"simple-udp-fkhqj","namespace":"default","uid":"b70af2f8-52f4-11e9-adb9-080027031e5a","resource_version":"4697","generation":1,"creation_timestamp":1553954660,"annotations":{"stable.agones.dev/sdk-version":"0.9.0-rc"},"labels":{"stable.agones.dev/sdk-timestamp":"1553956874"}},"spec":{"health":{"PeriodSeconds":5,"FailureThreshold":3,"InitialDelaySeconds":5}},"status":{"state":"Ready","address":"10.0.2.15","ports":[{"name":"default","port":7383}]}}
2019/03/30 15:31:17 Received packet from 192.168.99.1:49977: EXIT
2019/03/30 15:31:17 Received EXIT command. Exiting.
```

## まとめ

- Agonesの動作確認は、自分が何をしたいかによって適切なインストール方法を選ぶ必要がある
    - Agonesが何者かを手っ取り早くしりたいのであれば、ドキュメントどおりでOK
    - SDK周りのサンプルまで確認したい場合は、リリースブランチを使う
- 自分のようなk8s弱者は、ゲームサーバ動作確認用の接続先アドレスとポートを `kubectl describe all` から探す
    - もっといい方法があるとは思うが、自分にはまだ圧倒的にk8s力が足りない。
