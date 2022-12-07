---
title: Githubのアカウント名を変えた
---

Githubのアカウント名を `tomochikahara` から `thara` に変えた。

最初にGithubのアカウントを取った時に `thara` がすでに取られていて諦めていたが、 そのアカウントがGithubの [Name Squatting Policy](https://help.github.com/articles/name-squatting-policy/) に違反していることがわかったので、思い切ってサポートにメールしたらあっさり開放してくれた。

送ったメールは下記のとおり。英語は適当。

```
Subject: Request to remove "thara" account
Hi,

My name is Tomochika Hara. I have a GitHub account “tomochikahara” (https://github.com/tomochikahara), but I want to use the account ”thara” that is the same as my handle name.

I see the profile of "thara" (https://github.com/thara), and I found the account ”thara” is inactive.
Can you  remove “thara” account to use it?

Thanks.
```

これに対する返信が以下。

```
Hi there, 
You are in luck — we have classified the thara account as inactive and released the username for you to claim, as per our Name Squatting Policy:
https://help.github.com/articles/name-squatting-policy

Be quick, as the username is now publicly available. Glad to help!

Cheers
```

送信から返信までの時間、わずか30分。   
これには流石にびっくりした。送信した時間は深夜0:32 で、返信に気づいたのは翌朝。   
あくまで当該Githubアカウントが解放されただけなので、誰かに先を越されていないかヒヤヒヤした...


Githubのアカウント名を変えるのはGithubのサイトから簡単にできるが、そのあとにGithub Pagesやらローカルリポジトリやらに対応が必要だった。

### Github Pages

Github Pagesは、Github上のリモートリポジトリ名を `{アカウント名}.github.io` にすることでWebサイトを公開できる仕組みだが、Githubアカウント名を変更した場合は、当然このリポジトリ名も変更しなければならない。

自分の場合はさらに独自ドメインでhttps化していたため、Github Pagesにフォワードしていたnginxのフォワード先も変更する必要があった。


### リモートリポジトリURLの変更

開発に使っているマシン内のローカルリポジトリのリモートリポジトリURLも変更しなければならない。
自分は `$HOME/src/github.com/{アカウント名}/{リポジトリ名}` というディレクトリ構造を愛用しているため、当然ディレクトリ名も変える必要があった。   

ディレクトリ名変更後、以下のスクリプトでリモートリポジトリの向き先を変更した。

```
# macOSのみを想定

cd $HOME/src/github.com/thara

for f in *; do
    cd "$f"
    NEW_URL=`git remote -v | grep fetch | sed -e 's/tomochikahara/thara/' | sed s/$'\t'/$' '/g | cut -d ' ' -f 2`
    if echo $NEW_URL | grep -v "fatal"
    then
      git remote set-url origin "$NEW_URL"
    fi
    cd ..
done
```

### symlinkの更新

前述のように自分は `$HOME/src/github.com/{アカウント名}/{リポジトリ名}` というディレクトリ構造を愛用しているが、リポジトリには `dotfiles` が含まれており、そのリポジトリ内に `.zshrc` や `.vim` などホームディレクトリからsymlinkを貼っているものがあった。


ディレクトリ名を変更したことにより、当然そのsymlinkが壊れたので、それも修正した。

```
# macOSのみを想定

cd $HOME

for f in `find . -maxdepth 1 -type l -print`; do
    o=`readlink "$f"`
    new=`echo $o | sed -e 's/tomochikahara/thara/'`
    ln -s -f $new $f
done
```


## Conclusion

今まで地味にコンプレックスだった長めのGithubアカウント名が解消されて、清々しい気分。   
Golang用のライブラリもこれで安心して公開できる。
さっさとやっておけばよかった...
