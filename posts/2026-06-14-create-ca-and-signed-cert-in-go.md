---
title: GoでCAと署名付き証明書を作る
date: '2026-06-14'
published: '2026-06-14'
---

前やったような気もするけど、思い出しも兼ねて。

全体的なフローは、

1. ルート証明書作成
2. ルート証明書からサーバー証明書作成
3. ルート証明書からクライアント証明書作成
4. httptest.Serverにサーバー証明書設定
5. http.Clientにクライアント証明書設定
6. http.Clientからリクエスト

って感じの流れで動作確認する。

## 1. 証明書作成

subjectは適当

```go
	subject := pkix.Name{
		Country:            []string{"JP"},
		Organization:       []string{"My Organization"},
		OrganizationalUnit: []string{"My Unit"},
		Locality:           []string{"My City"},
		Province:           []string{""},
		StreetAddress:      []string{""},
		PostalCode:         []string{"100-0004"},
		CommonName:         "localhost",
	}
```

ルート証明書をこんな感じで作成する。
今回は中間証明書は作らない。

```
	serialNumber, err := randomSerialNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %v", err)
	}

	now := time.Now()
	const rootCAValidYears = 10

	caCert := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             now.Add(-5 * time.Minute),
		NotAfter:              now.AddDate(rootCAValidYears, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		// Intermediate CAなし
		MaxPathLen:     0,
		MaxPathLenZero: true,
	}
```

`x509.Certificate` は単なるパラメータオブジェクト的な存在（たぶん）。

実際の証明書作成は `x509.CreateCertificate` を使う。

```go
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate CA private key: %v", err)
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create CA certificate: %v", err)
	}
```

後続処理でPEM形式で扱うので、変換しておく。

```go
	caPEM := new(bytes.Buffer)
	if err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}); err != nil {
		return nil, fmt.Errorf("failed to encode CA certificate to PEM: %v", err)
	}
	caPrivKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	}); err != nil {
		return nil, fmt.Errorf("failed to encode CA private key to PEM: %v", err)
	}
```

## 2. ルート証明書からサーバー証明書作成

まずは `x509.Certificate` を作る。
subjectは1のをそのまま使ってる想定。

```
	serialNumber, err := randomSerialNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %v", err)
	}

	now := time.Now()

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		// 動作確認用
		DNSNames:    []string{"localhost"},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},

		NotBefore: now.Add(-5 * time.Minute),
		NotAfter:  now.AddDate(1, 0, 0),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
	}
```

DNSNames, IPAddressesは動作確認用にローカルにHTTPサーバーを立てるので。
`NotBefore` はちょっとだけ過去のにしておかないと、クライアント/サーバー間で時刻ズレがあった時に作成直後の証明書が使えないことがある。まぁ今回の動作確認程度では問題にならないと思うけど。

KeyUsage, ExtKeyUsageもサーバー証明書としての設定。

そのあとは、おんなじ感じでx509.CreateCertificateを使って証明書を作成する。

クライアント証明書も同じように作成するので、↓のような関数を用意した。
今回は`parent` にルート証明書を指定する想定。

```go
type Certificate struct {
	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey

	PEM           *bytes.Buffer
	PrivateKeyPEM *bytes.Buffer
}

func generateCertificate(parent *Certificate, cert *x509.Certificate) (*Certificate, error) {
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate certificate private key: %v", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, parent.Certificate, &certPrivKey.PublicKey, parent.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %v", err)
	}

	certPEM := new(bytes.Buffer)
	if err := pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return nil, fmt.Errorf("failed to encode certificate to PEM: %v", err)
	}

	certPrivKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	}); err != nil {
		return nil, fmt.Errorf("failed to encode certificate private key to PEM: %v", err)
	}

	return &Certificate{Certificate: cert, PrivateKey: certPrivKey, PEM: certPEM, PrivateKeyPEM: certPrivKeyPEM}, nil
}
```

## 3. ルート証明書からクライアント証明書作成

サーバー証明書とほぼ同じ。

```go
	serialNumber, err := randomSerialNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %v", err)
	}

	now := time.Now()
	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    now.Add(-5 * time.Minute),
		NotAfter:     now.AddDate(1, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
	}
```

クライアント証明書の用途的に`KeyUsage` に `x509.KeyUsageKeyEncipherment` はいらんやろ...

## 4. httptest.Serverにサーバー証明書設定

こっからが動作確認のためのフェーズ。

あらかじめ、tls.Configを作っておく。

```go
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certs.CA.PEM.Bytes())

	serverTLSCert, err := tls.X509KeyPair(certs.Server.PEM.Bytes(), certs.Server.PrivateKeyPEM.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to create server TLS certificate: %v", err)
	}
	serverTLSConf := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
		// for mTLS
		ClientCAs:  certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	clientTLSCert, err := tls.X509KeyPair(certs.Client.PEM.Bytes(), certs.Client.PrivateKeyPEM.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to create client TLS certificate: %v", err)
	}
	clientTLSConf := &tls.Config{
		RootCAs: certPool,
		// for mTLS
		Certificates: []tls.Certificate{clientTLSCert},
	}
```

serverTLSConfのClientCAs, ClientAuth、clientTLSConfのCertificatesはmTLSのため。不要であればフィールドを設定しなくていい。

で、httptest.Serverを適当に起動する。

```go
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "success!")
	}))
	server.TLS = serverTLSConf
	server.StartTLS()
	defer server.Close()
```

## 5. http.Clientにクライアント証明書設定

4で作ったのをそのまま使うだけ。

```go
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: clientTLSConf,
		},
	}
```

## 6. http.Clientからリクエスト

```go
	resp, err := client.Get(server.URL)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	body := strings.TrimSpace(string(b[:]))
	fmt.Printf("Response: %s\n", body)
```

こんな感じかな。
ソースコードは [thara-playground/go-ca-playground](https://github.com/thara-playground/go-ca-playground) に置いてある。(自分が理解しやすいようにリファクタリング済み)

普段ゴリゴリに使うパッケージじゃないから、ちょこちょこ忘れてた。

こういうPKIの一連の流れをシュッと確認できるのは、Goのいいところかもしれない。
