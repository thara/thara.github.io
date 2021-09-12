---
layout: post
title: Swift製ファミコンエミュレータの進捗、そしてCライブラリへのSwiftバインディング開発の一例
tags:
  - swift
  - advent-calendar
---

Misocaの [thara](https://twitter.com/zetta1985) です。   
この記事は [Misoca+弥生 Advent Calendar 2019 - Qiita](https://qiita.com/advent-calendar/2019/misoca-yayoi) 23日目の記事です。
遅くなってすみません...

3ヶ月前にiOSDC Japan 2019にて [Swiftでつくるファミコンエミュレータのススメ](https://fortee.jp/iosdc-japan-2019/proposal/92904657-beda-46fe-8ecb-b27c75ee0f16) というタイトルでLTをしてきました。

<iframe width="560" height="315" src="https://www.youtube.com/embed/9pBPF77XQX0?start=1361" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

LTでは[NES研究室で配布されているHallo, WorldのサンプルROM](http://hp.vector.co.jp/authors/VA042397/nes/sample.html) が動作したことをお伝えしました。
あれから数ヶ月を経て、どれぐらいまで動くようになったのか、イベント会場でお話した「自分が子どもの頃にプレイしたゲームを自作エミュレータ上でプレイして、子どもの頃にクリアできなかったゲームをクリアしたい」という目標までどれぐらい近づいたのかをお伝えしたいと思います。


## 進捗

10月の半ばには、スーパーマリオブラザーズはプレイできるようになりました。

<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">自作エミュレータでちゃんとマリオをプレイできるようになった（音はまだ無い） <a href="https://twitter.com/hashtag/nesdev?src=hash&amp;ref_src=twsrc%5Etfw">#nesdev</a> <a href="https://t.co/apuUo851BA">https://t.co/apuUo851BA</a> <a href="https://t.co/EYLbe3SqoF">pic.twitter.com/EYLbe3SqoF</a></p>&mdash; thara (@zetta1985) <a href="https://twitter.com/zetta1985/status/1183747534109343749?ref_src=twsrc%5Etfw">October 14, 2019</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

↑のgifでは60FPS出ていないんですが、 これは録画ソフトがCPUリソースを持っていってしまってその影響でモッサリしているだけなので、ちゃんとエミュレータ単体では60FPS出てます。
(色々試したんだけれど、Webブラウザ上でぬるぬる動く動画撮れなかった...)

ここまで動いたら、あと残すは「音」です。

## APU

ファミコンではAPU（Audio Processing Unit）と呼ばれるマイクロプロセッサで、プログラムから動的に設定された値を元に矩形波や三角波などを生成し、音を生成しています。

そもそもオーディオプログラミングの知識が皆無だった自分は，[nesdev.comのAPU関連のwiki](http://wiki.nesdev.com/w/index.php/APU) を読み漁り、
理解をコードに落とし込むかのように分周器やタイマー、エンベロープ・ジェネレータなどを細かく実装していきました。

そもそも用語の意味さえ分からなかった自分は、ここでだいぶ時間を食いました。[^1]

さらに、実際に音を再生するためには何らかのクロスプラットフォーム対応されたライブラリを使う必要があります。
Webブラウザ上で動作しているものは[Web Audio API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Audio_API)を使用しており、自分はクロスプラットフォームのネイティブで動作するエミュレータを実装しているので参考になりません。

C言語やC++製のファミコンエミュレータのコードを見てみるとSDL2のAudioを使っているものが多かったのですが [libsoundio-sharpとPInvokeGeneratorについて - ものがたり](http://atsushieno.hatenablog.com/entry/2017/12/07/041145) を参考に、ソースコードが読みやすく、エラー関連のドキュメントがしっかり書かれた[libsoundio](https://github.com/andrewrk/libsoundio)にチャレンジしてみました。

当然のようにSwiftバインディングが存在しないので、Swiftから扱いやすくするためのバインディングを書くことにしました <- イマココ

## Swiftバインディングを少しずつ書いていく

言語バインディングと言うと難しい印象ですが、C言語製のライブラリをSwiftから扱うのは非常にかんたんです。   

ただ単に使うぐらいであれば、Package.swiftの修正、modulemapという設定ファイル、ほんの数行のヘッダーファイルを用意するだけで、Swiftからそのライブラリの構造体や関数を直接使うことができます。

以下は [libsoundioのREADME中のサンプルコード](https://github.com/andrewrk/libsoundio#synopsis) とほぼ同じ処理をSwiftで書いたものです。

<script src="https://gist.github.com/thara/f1a4c725d8ef743c543a4e808b06db52.js?file=main.swift"></script>

たしかにSwiftからC言語の構造体や関数を直接使えていますが、このままだとさすがにSwiftのコードとしては読みづらいものがあります。
`UnsafeMutablePointer` は普段はめったに使わないものですし、エラーコードが`CInt`として扱われていて扱いづらく、そもそも関数の命名規則がSwiftのそれとは違います。

そこで、このmain.swiftをSwiftらしくすることを当面の目的として、徐々にSwiftバインディグのコードを書いていくことにします。

### エラーコードをSwiftでのErrorとして扱う

まず目につくのが、このようなエラーコード判定です。

```swift
    var err: CInt = soundio_connect(soundio)
    if 0 < err {
        fatalError("error connecting: \(soundioError(err))")
    }
```

Swiftではエラー処理の方法に、Result, Error(いわゆる検査例外), fatalErrorなどといくつか方法がありますが、ライブラリ内でfatalErrorを発生させるのは不適切であり、さらに[libsoundioのエラーコードの定義](http://libsound.io/doc-2.0.0/soundio_8h.html#a9aed679ba44aaa7d9bcbe2fa1b1156c5)を見るとアプリケーションがその後に復帰して継続可能なエラーの類はなさそうなので、Errorとして表現する方法を考えてみることにします。

```swift
public struct SoundIOError: Error {
   public let message: String

   init(errorCode: CInt) {
       self.message = String(cString: soundio_strerror(errorCode))
   }
}
```

[soundio_strerror](https://github.com/andrewrk/libsoundio/blob/b810bf2e9c4afc822c4843322cd08f7b36668109/src/soundio.c#L76-L96) は与えられたerrorCodeが既知のエラーコードでない場合には`(invalid error)`を返してくれるので、安心してこう書けます。


そして、`CInt`がエラーコードである場合にエラーをthrowする拡張メソッドを追加します。

```swift
//TODO Make internal
public extension CInt {

    @inline(__always)
    func ensureSuccess() throws {
        if 0 < self {
            throw SoundIOError(errorCode: self)
        }
    }
}
```

ここではpublicとしていますが、これはmain.swiftを配置しているデモ用のビルドターゲットから一時的に呼び出すためのもので、最終的なライブラリの実装ではinternalとしてライブラリAPIとしては公開しないようにします。

これで、main.swiftのエラーコードを扱っている箇所は、以下のように書けるようになりました。

```diff
func main() {
+func main() throws {
     // 略

-     var err: CInt = soundio_connect(soundio)
-     if 0 < err {
-         fatalError("error connecting: \(soundioError(err))")
-     }
+     try soundio_connect(soundio).ensureSuccess()
```

### Optional<UnsafeMutablePointer> がnilの場合のエラーを表現

次に気になるのが、各ポインタのnil（C言語でいうNULLポインタ）チェックの箇所です。

このままOptionalとして扱ってもよいのですが、コードをより説明的にするために、これもErrorとして表現してみます。

```swift
//TODO Make internal
 public extension Optional {
     @inline(__always)
     func ensureAllocatedMemory() throws -> Wrapped {
         guard let value = self else {
             throw SoundIOError(message: "Ouf of memory: \(Wrapped.self)")
         }
         return value
     }
 }
```

こちらも後に非公開にするAPIとしておきます。

これで、以下のようにnilチェックをErrorとして扱えるようになります。

```diff
-    guard let soundio: UnsafeMutablePointer<SoundIo> = soundio_create() else {
-        fatalError("out of memory")
-    }
+    let soundio = try soundio_create().ensureAllocatedMemory()
```

### 構造体をクラスで表現する

libsoundioで公開されている構造体は、その構造体のポインタを扱う関数とセットで扱う、つまりC言語とはいえ、オブジェクト指向で設計されたAPIです。
よって、それらの構造体をSwiftではクラスで表現することを試みます。

まず、構造体へのポインタをラップするクラスを定義します。

```swift
import CSoundIO

public class SoundIO {
    private let internalPointer: UnsafeMutablePointer<CSoundIO.SoundIo>

    public init() throws {
         self.internalPointer = try soundio_create().ensureAllocatedMemory()
    }

    deinit {
        soundio_destroy(internalPointer)
    }
}
```

`CSoundIO.SoundIo` が`libsoundio` が公開している構造体です。
ここで定義した`SoundIO`クラスはインスタンス生成時に`soundio_create`を呼び出して`CSoundIO.SoundIo`のポインタを取得・保持します。
そして、`deinit` によって自身のインスタンスが破棄されるタイミングで適切に `soundio_destroy` を呼び出し、保持している`CSoundIO.SoundIo`のポインタを破棄します。

ちなみに、deinitはクラスでしか使えないのでSwiftの構造体として`SoundIO`を定義することはできません。

さらに、他の関数呼び出しも`SoundIO`のメソッドでラップします。

```swift
    public func connect() throws {
        try soundio_connect(self.internalPointer).ensureSuccess()
    }

    public func flushEvents() {
        soundio_flush_events(self.internalPointer)
    }
```

これで、以下のようなSwiftらしいコードが書けるようになります。

```diff
-    let soundio = try soundio_create().ensureAllocatedMemory()
-    try soundio_connect(soundio).ensureSuccess()
-    soundio_flush_events(soundio)
+    let soundio = try SoundIO()
+    try soundio.connect()
+    soundio.flushEvents()
```

### 型エイリアスで説明的にする

`CSoundIO.SoundIo`構造体を扱う関数をメソッドにしていくうえで、以下のコードを書きました。

```swift
    public func defaultOutputDeviceIndex() throws -> CInt {
        let index = soundio_default_output_device_index(self.internalPointer)
        guard 0 <= index else {
            throw SoundIOError(message: "No output device found")
        }
        return index
    }

    //TODO Wrap SoundIODevice
    public func getOutputDevice(at index: CInt) throws -> UnsafeMutablePointer<CSoundIO.SoundIoDevice> {
        guard let device = soundio_get_output_device(self.internalPointer, index) else {
            throw SoundIOError(message: "invalid parameter value")
        }
        return device
    }
```

`soundio_default_output_device_index` の戻り値は、 `soundio_get_output_device` の引数として使っています。
`CInt`のままだとそれがわからないので、型エイリアスを使って、よりメソッド間の関係がわかりやすいようにします。

```swift
public typealias DeviceIndex = Int

public class SoundIO {
    // 略
    public func defaultOutputDeviceIndex() throws -> DeviceIndex {
        let index = soundio_default_output_device_index(self.internalPointer)
        guard 0 <= index else {
            throw SoundIOError(message: "No output device found")
        }
        return DeviceIndex(index)
    }

    //TODO Wrap SoundIODevice
    public func getOutputDevice(at index: DeviceIndex) throws -> UnsafeMutablePointer<CSoundIO.SoundIoDevice> {
        guard let device = soundio_get_output_device(self.internalPointer, index) else {
            throw SoundIOError(message: "invalid parameter value")
        }
        return device
    }
}
```

これで、ライブラリを扱うユーザーが誤った値を`getOutputDevice(at:)`に渡す可能性が減ります。 [^2]

### 生のポインタを扱う余地を残す

`CSoundIO.SoundIo`構造体を扱う関数を全て`SoundIO`クラスのメソッドとして定義したいところですが、このバインディングライブラリの保守の手間や元のライブラリの変更への追従を考えると、ライブラリユーザーにはある程度の柔軟性を持たせた方が良いことがあります。

例えば、ライブラリユーザーは最新バージョンのライブラリで追加された関数を使いたいが、バインディングライブラリがその関数をサポートしておらず、かつ内部のポインタへのアクセスを完全にライブラリ内に閉じている場合、ライブラリユーザーはそのバインディングライブラリの使用を諦めるかforkするしか手がなくなってしまいます。

よって、以下のような限定的なスコープで内部ポインタにアクセス可能な手段を提供します。

```swift
    public func withInternalPointer(_ unsafeTask: (_ pointer: UnsafeMutablePointer<CSoundIO.SoundIo>) throws -> Void) throws {
        try unsafeTask(self.internalPointer)
    }
```

このメソッドは以下のように使用できます。


```swift
         try soundio.withInternalPointer {
             soundio_wait_events($0)
         }
```

もちろん、バインディングライブラリの意図しない変更をポインタに対して行われる可能性があるので、ライブラリユーザーに注意を促すドキュメントを残しておくと良さそうです。

## SoundIOクラスはほぼ完成

先程のmain.swift内で使っている`CSoundIO.SoundIo`とその関数を`SoundIO`クラスとして扱えるようになると 以下のように大部分をSwiftらしく書けるようになります。

```swift
func main() throws {
    let soundio = try SoundIO()
    try soundio.connect()
    soundio.flushEvents()

    let outputDeviceIndex = try soundio.defaultOutputDeviceIndex()
    let device = try soundio.getOutputDevice(at: outputDeviceIndex)
    let deviceName = String(cString: device.pointee.name)
    print("Output device: \(deviceName)")

    let outstream = try soundio_outstream_create(device).ensureAllocatedMemory()
    outstream.pointee.format = SoundIoFormatFloat32LE;
    outstream.pointee.write_callback = writeCallback;

    try soundio_outstream_open(outstream).ensureSuccess()
    try outstream.pointee.layout_error.ensureSuccess()
    try soundio_outstream_start(outstream).ensureSuccess()

    while true {
        soundio.waitEvents()
    }

    soundio_outstream_destroy(outstream)
    soundio_device_unref(device)
}
```

あとは`CSoundIO.SoundIoDevice`構造体や`CSoundIO.SoundIoOutStream`構造体などが残っていますが、今まで見てきたように徐々にラッパーを書いていくことで、バインディングライブラリを完成に近づけていけると思います。

## Swift as an system programming language

ファミコンエミュレータ実装においてAPUを実装中であること、そして、そのAPU実装のためにlibsoundioのバインディングライブラリを徐々に開発していることをご紹介しました。

Swiftは、Linux上でも動作するにも限らずiOSやmacOS上での開発言語としてしか注目を浴びていませんが、C言語との相互運用性やバインディングライブラリの書きやすさ、コンパイラによる最適化など、Rustまでは行かないものの、Go言語と同じぐらいのレイヤーまでは十分に対応できるポテンシャルを秘めていると自分は考えています。


ファミコンエミュレータの実装、そしてCライブラリ向けのバインディングライブラリの実装を通じて、Swiftの可能性を広げていきたいです。

---


[Misoca+弥生 Advent Calendar 2019 - Qiita](https://qiita.com/advent-calendar/2019/misoca-yayoi)、   
次回24日は [yayoi_ueshima](https://qiita.com/yayoi_ueshima)さんの「仮想通貨のやり取りを体感してみよう！」です。   
お楽しみに！

---

[^1]: Courseraで機械学習系のコースを受講してたのが一番の要因の気もする...
[^2]: Swiftの型エイリアスはコンパイル時には削除されるので、あくまでもドキュメンテーションとライブラリユーザーへの説明のためのものです
