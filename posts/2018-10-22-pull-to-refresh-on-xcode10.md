---
title: UITableViewControllerでPull-to-Refresh Xcode10版
published: '2018-10-22'
tags: [ios]
---

少しググったらわかるけど微妙にバージョン違いが多かったのでメモ。   
XCode 10というのはただ確認したのがこのバージョンだっただけで、厳密には違うかもしれない。

```swift
class MyViewController: UITableViewController {

    override func viewDidLoad() {
        refreshControl = UIRefreshControl()
        refreshControl?.addTarget(self, action: #selector(refresh(sender:)), for: .valueChanged)
    }

    @objc func refresh(sender: UIRefreshControl) {
        // 更新処理
        sender.endRefreshing()  // アクティビティインジケータを消す
    }
```


## 参考

- [UITableViewController - UIKit / Apple Developer Documentation](https://developer.apple.com/documentation/uikit/uitableviewcontroller)
- [UIRefreshControl - UIKit / Apple Developer Documentation](https://developer.apple.com/documentation/uikit/uirefreshcontrol)
