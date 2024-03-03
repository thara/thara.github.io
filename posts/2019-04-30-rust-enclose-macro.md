---
title: Rustのクロージャ内で外側のスコープのimmutableな変数をmutableな変数に同名で束縛したいときにつかうマクロ
description: Rustのクロージャ内で外側のスコープのimmutableな変数をmutableな変数に同名で束縛したいときにつかうマクロ
published: '2019-04-30'
tags: [rust]
---

クロージャ内で外側のスコープの変数をmutableな変数に同名で束縛したいときに、外側の変数の不変に保つために、
以下のようなマクロを使っている。



```rust
macro_rules! enclose {
    ( ($( $x:ident ),*) $y:expr ) => {
        {
            $(let mut $x = $x.clone();)*
            $y
        }
    };
}
```

以下のような感じで書ける。

```rust
let sdk = agones::Sdk::new().map_err(|_| "Could not connect to the sidecar. Exiting!")?;

let _health = thread::spawn(enclose! {(sdk) move || {
    loop {
        match sdk.health() {
            (s, Ok(_)) => {
                println!("Health ping sent");
                sdk = s;
            },
            (s, Err(e)) => {
                println!("Health ping failed : {:?}", e);
                sdk = s;
            }
        }
        thread::sleep(Duration::from_secs(2));
    }
}});
```

`thread::spawn` に渡しているクロージャ内では、sdk変数を特定の条件で再束縛しているため、sdk変数はmutableである必要がある。   
マクロなしでmutableなsdk変数を作ろうとcloneしてしまうと、シャドウイングによって外側のスコープのimmutableなsdk変数を後続の処理で使うことができなくなってしまう。   

上記の `enclose` マクロは、このシャドウイングを防止する。健全なマクロバンザイ。


なお、自分はマクロ素人 & ちょっと前のバージョンでこのマクロを使ったので、Rust2018ではもっといい方法があるかもしれない。
あしからず。
