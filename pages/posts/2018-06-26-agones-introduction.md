---
title: ゲームサーバー管理サービス Agones
tags: [game-dev, agones]
---

Google と Ubisoft が協働で開発している [Agones](https://github.com/GoogleCloudPlatform/agones) が 2018/3/14 に公開された。

- [GoogleがUbisoftと協働でオープンソースのゲームサーバーホスティングシステムAgonesをローンチ](https://jp.techcrunch.com/2018/03/14/2018-03-13-google-partners-with-ubisoft-to-launch-agones-an-open-source-game-server-hosting-system/)
- [Agones ―― Kubernetes 上でのゲーム サーバー構築をサポートするオープンソース プロジェクトが始動](https://cloudplatform-jp.googleblog.com/2018/04/introducing-Agones-open-source-multiplayer-dedicated-game-server-hosting-built-on-Kubernetes.html)

ゲームサーバーやオンラインゲーム業界に携わったことがない人から見ると、どのような問題を解決するのかわかりにくいかもしれない。

自分もまだしっかり触ったわけでもプロダクションに導入したわけでもないが、現在の理解の範囲内で簡単に説明すると、
Agones は Kubernetes(以下k8s) 上のアプリケーションのホスティング/スケーリング管理を行うサービス。

アプリケーション開発者は、SDK を組み込んだ Docker image と設定ファイルを作成し、
それを用いて Agones が管理する k8s にデプロイすることで、簡単にスケールアウト可能なクラスタを構成することができる。

ステートフルなゲームサーバーを運用しているサービスであれば何らかの方法でこのようなゲームサーバー管理サービスを構築していると思われるが、
Agones の大きな利点として、k8s 上に構築されているため minikube を用いてローカル環境でも簡単にゲームサーバーを稼働させることができる点が挙げられる。

minikube を用いれば macOS 上でも動作させられるはずなので、実際に試してみた。

## Agones サンプルを macOS 上で試す

[公式のQuickStart](https://github.com/GoogleCloudPlatform/agones/blob/291825b6f837982adb2b0198110a4e74cf9d5f09/docs/create_gameserver.md) の焼き直しになるが、
以下のような手順で macOS 上で Agones の動作確認ができる。

### 事前準備

事前準備として、Docker for Mac および minikube のインストールが必要。
minikube のインストールは、バージョン管理のため、[asdf](https://github.com/asdf-vm/asdf) 経由で行うことをおすすめする。

```bash
$ brew install asdf
$ asdf plugin-add minikube
$ asdf global minikube 0.25.2
```

### Agones の取得

2018/6/26 現在の GitHub 上の Agones リポジトリには、SDKやサンプルなどがすべて含まれている。

```bash
$ git clone https://github.com/GoogleCloudPlatform/agones.git
$ cd agones
```

### k8s 準備

```bash
$ minikube profile agones
$ minikube start --kubernetes-version v1.9.4 --vm-driver virtualbox \
  --extra-config=apiserver.Admission.PluginNames=NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,DefaultTolerationSeconds,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota \
  --extra-config=apiserver.Authorization.Mode=RBAC
$ kubectl create clusterrolebinding cluster-admin-binding \
  --clusterrole=cluster-admin --serviceaccount=kube-system:default
```

### Agones インストール

以下を実行することで、 k8s クラスタ上に Agones がインストールされる。

```bash
$ kubectl apply -f https://github.com/GoogleCloudPlatform/agones/raw/release-0.2.0/install/yaml/install.yaml
$ kubectl describe --namespace agones-system pods  # 確認
```

### Agones サンプルの実行

`eval $(minikube docker-env) ` を忘れてドハマリしたので注意。

```bash
$ eval $(minikube docker-env)
$ kubectl create -f examples/cpp-simple/gameserver.yaml
gameserver.stable.agones.dev "cpp-simple-s776l" created
```

各アプリケーションの標準出力は、`kubectl logs` で確認できる。

```bash
$ kubectl get pods
NAME                     READY     STATUS    RESTARTS   AGE
cpp-simple-s776l-j5894   2/2       Running   0          42s
$ kubectl logs cpp-simple-s776l-j5894 -c cpp-simple
```

また、 `kubectl describe pods` で、Agones が出力するメッセージを確認できる。

```bash
$ kubectl describe pods
Events:
  Type     Reason                 Age                From               Message
  ----     ------                 ----               ----               -------
  Normal   Scheduled              1m                 default-scheduler  Successfully assigned cpp-simple-s776l-j5894 to agones
  Normal   SuccessfulMountVolume  1m                 kubelet, agones    MountVolume.SetUp succeeded for volume "agones-sdk-token-bzhdp"
  Normal   Pulling                1m                 kubelet, agones    pulling image "gcr.io/agones-images/cpp-simple-server:0.1"
  Normal   Pulled                 1m                 kubelet, agones    Successfully pulled image "gcr.io/agones-images/cpp-simple-server:0.1"
  Normal   Created                1m                 kubelet, agones    Created container
  Normal   Started                1m                 kubelet, agones    Started container
```

## まとめ

Agones の簡単な紹介と、macOS + minikube でのサンプルの実行を試した。

個人的には、今までクローズドであったステートフルなゲームサーバーの実装が Agones 上に再構築されることでオープンソースとなり、
より洗練された設計/実装になったり、オンラインゲーム開発の障壁の低減につながることを期待したい。
また、k8s 上に構築可能なサービスとして良い例 だと思うので、ゲーム業界以外でも何かの参考になるかもしれない。

ちなみに、Agones にはすでに [Rust SDK のサポート追加という形でコントリビュートした](https://github.com/GoogleCloudPlatform/agones/pull/230) 。

またの機会にSDKの追加方法や、ローカルでのゲームサーバー動作確認方法を解説する予定。

