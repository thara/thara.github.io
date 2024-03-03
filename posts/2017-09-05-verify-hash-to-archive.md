---
title: アーカイブのハッシュ値検証方法
published: '2017-09-05'
---

今まで使ったことのあるものだけ、雑にまとめておく。

## MD5やSHA256など

### ハッシュ値のみ

ダウンロードサイトなどでアーカイブのハッシュ値だけが記載されているパターン。

アーカイブをダウンロードしたディレクトリで、以下のコマンドを実行。

```
# 例) MD5
$ echo "${ハッシュ値}  ${アーカイブファイル名}" | md5sum -c - \
```

### ハッシュ値検証ファイル

ダウンロードサイトなどで複数のアーカイブ形式のハッシュ値を一つのファイルにまとめて記載しているパターン。

アーカイブをダウンロードしたディレクトリで、以下のコマンドを実行。

```
# 例) SHA256 (sha256sumに --ignore-missing が実装されている場合)
$ sha256sum -c {各アーカイブごとのハッシュ値が記載されているファイルのパス} --ignore-missing

# 例) SHA256 (sha256sumに --ignore-missing が実装されていない場合)
$ sha256sum -c {各アーカイブごとのハッシュ値が記載されているファイルのパス} 2>&1 | grep OK
```


### 例1) MySQL Connector/C++ source install in Docker image

```
RUN set -ex; \
  key='A4A9406876FCBD3C456770C88C718D3B5072E1F5'; \
  gpg --keyserver ha.pool.sks-keyservers.net --recv-keys "$key"

RUN mkdir -p /usr/src/mysql-connector-cpp \
  && curl -fSL https://dev.mysql.com/get/Downloads/Connector-C++/mysql-connector-c++-1.1.9.tar.gz -o mysql-connector-cpp.tar.gz \
  && echo "f262bef7e70178f95ceb72a71f0915f7  mysql-connector-cpp.tar.gz" | md5sum -c - \
  && curl -fSL https://dev.mysql.com/downloads/gpg/?file=mysql-connector-c%2B%2B-1.1.9.tar.gz -o mysql-connector-cpp.tar.gz.asc \
  && gpg --batch --verify mysql-connector-cpp.tar.gz.asc mysql-connector-cpp.tar.gz \
  && tar xf mysql-connector-cpp.tar.gz -C /usr/src/mysql-connector-cpp --strip-components=1 \
  && rm -f mysql-connector-cpp.* \
  && cd /usr/src/mysql-connector-cpp \
  && cmake . && make clean && make -j$(nproc) && make install
```


### 例2) CMake source install in Docker image

```
ENV CMAKE_MINOR_VERSION=v3.7 \
    CMAKE_FULL_VERSION=3.7.1

RUN mkdir -p /usr/src/cmake \
  && curl -fSLO https://cmake.org/files/${CMAKE_MINOR_VERSION}/cmake-${CMAKE_FULL_VERSION}.tar.gz \
  && curl -fSLO https://cmake.org/files/${CMAKE_MINOR_VERSION}/cmake-${CMAKE_FULL_VERSION}-SHA-256.txt \
  && sha256sum -c cmake-${CMAKE_FULL_VERSION}-SHA-256.txt 2>&1 | grep OK \
  && tar xf cmake-${CMAKE_FULL_VERSION}.tar.gz -C /usr/src/cmake --strip-components=1 \
  && rm -f cmake-${CMAKE_FULL_VERSION}.* \
  && cd /usr/src/cmake \
  && ./configure && make -j$(nproc) && make install
```
