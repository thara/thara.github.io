---
title: WatchKitでBackground中にHaptic feedbackする
published: '2018-03-07'
tags: [ios]
---

ちょっとApple Watch　のアプリで作ってみたいものがあって、そのとき調べたことのメモ。

やりたいことはタイトルの通りなのだが、通常はバックグラウンドまたは非アクティブ時に Haptic feedback（触覚フィードバック）を行うことはできない。

例外として、ワークアウト中、つまり、`HKWorkoutSession` オブジェクトを `HKHealthStore#start` に渡してから `HKHealthStore#end` を呼ぶまでのみ、 `WKInterfaceDevice.current().play` が意図通りに動作する。

以下は、一定間隔で Haptic feedback を行うサンプルコード。
Watch上のボタンと、 `toggleSession` が connect されている想定。

```swift
import WatchKit
import Foundation
import HealthKit

class InterfaceController: WKInterfaceController, HKWorkoutSessionDelegate {

    let healthStore = HKHealthStore()
    var currentWorkoutSession: HKWorkoutSession?

    var isRunning: Bool?

    weak var timer: Timer?
    var intervalSec: TimeInterval = 300.0

    // awake, willActivateなどの定義は省略

    @IBAction func toggleSession() {
        if isRunning == true {
            isRunning = false
            
            guard let session = currentWorkoutSession else { return }
            healthStore.end(session)
            self.timer?.invalidate()
        } else {
            isRunning = true
            
            let conf = HKWorkoutConfiguration()
            conf.activityType = .other
            
            do {
                let session = try HKWorkoutSession(configuration: conf)
                session.delegate = self
                currentWorkoutSession = session
                healthStore.start(session)
                
                self.timer?.invalidate()
                self.timer = Timer.scheduledTimer(
                    timeInterval: self.intervalSec,
                    target: self, selector: #selector(InterfaceController.play),
                    userInfo: nil, repeats: true)
            } catch let e as NSError {
                fatalError("*** Unable to create the workout session: \(e.localizedDescription) ***")
            }
        }
        
        WKInterfaceDevice.current().play(.start)
    }

    @objc func play() {
        WKInterfaceDevice.current().play(.success)
        NSLog("Played")
    }

    // workoutSessionなどの定義も省略
```

このとき忘れてはいけないのは、いくつかのプロジェクトの設定が必要だということ。

- HealthKit.framework を Link Binary with Libraries に追加する
- WatchKit Extention projectの entitlements ファイルの `HealthKit` Key の Value を `YES` に設定する
- WatchKit Extention projectの Info.plist 内の `Required background modes (Watch)` Key の Value に `Workout Processing` を設定する

（一番最後の Info.plist の更新を忘れて、上記サンプルコードの Timer が動作しなくてハマった）

---

やりたいことは実現はできたのだけれど、単に Background 中に Haptic feedback したいだけなのに、ワークアウトセッションにしなければいけないのが難点。

上記のサンプルコードのような簡単な実装であれば、1日中ワークアウトセッション起動しっぱなしでもバッテリの減りが通常時と比べて極端に多くなる、ということはないのだが、大して運動もしてないのにアクティビティアプリに運動中だと認識されてしまう。

実際に App Store に審査依頼を出してはいないのだけれど、おそらくリジェクトされるだろう。運動してないのにアクティビティアプリのアチーブメントを達成してしまうのは自分としても不本意だし、誰もが嫌がると思う。

結局、現在のWatchKitだと自分が作りたいものが作れなさそう、ということでこのアイディアはお蔵入りに。

うーん、残念。
