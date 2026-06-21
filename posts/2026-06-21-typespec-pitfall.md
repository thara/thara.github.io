---
title: typespecの落とし穴
date: '2026-06-21'
published: '2026-06-21'
---

[TypeSpec](https://typespec.io/) を書いてて、OpenAPIスキーマの出力内容が意図通りになっていなくてハマったのでメモ。

TypeSpecでこういうAPI定義した時に、

types.tsp

```
model Widget {
  id: string;
  weight: int32;
  color: "red" | "blue";
}
```

main.tsp
```
import "@typespec/http";
import "./types.tsp";

using Http;
@service(#{ title: "Widget Service" })
namespace DemoService;

@route("/widgets")
interface Widget {
  @get read(@path id: string): Widget;
}
```

`tsp compile .` で生成されるOpenAPIスキーマで `GET /widgets/{id}` のレスポンスが無になる。

```yaml
openapi: 3.1.0
info:
  title: Widget Service
  version: 0.0.0
tags: []
paths:
  /widgets/{id}:
    get:
      operationId: Widget_read
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: The request has succeeded.
          content:
            application/json:
              schema: {}      # <- 無
components: {}
```

原因は main.tsp で、以下の修正で直せる。

```diff
 @route("/widgets")
-interface Widget {
+interface Widgets {
   @get read(@path id: string): Widget;
 }
```

正しく生成されたOpenAPIスキーマ

```yaml
openapi: 3.1.0
info:
  title: Widget Service
  version: 0.0.0
tags: []
paths:
  /widgets/{id}:
    get:
      operationId: Widgets_read
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: The request has succeeded.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Widget'   # <- 無じゃない!
(以下略)
```

いや、気づけよ、という話ではあるんだけど。

model定義とinterface定義が同じファイルの場合、 `tsp compile .` がエラーになるので容易に検知できる。

こういうmain.tspだと、

```
import "@typespec/http";

using Http;

@service(#{ title: "Widget Service" })
namespace DemoService;

model Widget {
  id: string;
  weight: int32;
  color: "red" | "blue";
}

@route("/widgets")
interface Widget {
  @get read(@path id: string): Widget;
}
```

こういうふうにエラーが出る。

```
$ tsp compile .
TypeSpec compiler v1.13.0

× Compiling
Diagnostics were reported during compilation:

main.tsp:8:7 - error duplicate-symbol: Duplicate name: "Widget"
> 8 | model Widget {
    |       ^^^^^^
main.tsp:15:11 - error duplicate-symbol: Duplicate name: "Widget"
> 15 | interface Widget {
     |           ^^^^^^

Found 2 errors.
```

modelとinterfaceの名前空間は同じなので、ファイルを分けるとそのチェックが働かなくなるみたい。

これ起因で、interface定義によっては↓みたいな謎の警告が出たりもする。
(どういう時に出るかまだ踏み込めてない)

```
main.tsp:15:19 - warning @typespec/http/metadata-ignored: path property will be ignored as it is inside of a @body property. Use @bodyRoot instead if wanting to mix.
> 15 |   @get read(@path id: string): Widget | Error;
     |                   ^^
```

これを避けるために、interface定義とmodel定義を同じファイルにしよう、というのは無理筋というか、ファイル構成に対する重い制約になってしまう。namespace分けるという手もあるが、API定義に反映されてしまうしな...

というわけで、一旦はinterface名は複数形にしておく、というルールにしておこうかな...   
もう一度ハマったら、[linter](https://typespec.io/docs/extending-typespec/linters/) 作るかも。
