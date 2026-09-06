---
title: Hagging Faceのモデルを使って推論
date: '2026-09-06'
published: '2026-09-06'
---

元同僚がHagging Faceのモデル使ってアプリ作ってると聞いて、
そういえばHagging Faceからモデル取ってきて推論するのやってないな、と思ったのでやってみた。

今更感が半端ないんだけど。

どうせなので、量子化手法とバックエンドを切り替えたベンチマークも取ってみる。

## Hagging Faceからモデル取得

Hagging FaceのCLIツール [hf](https://huggingface.co/docs/huggingface_hub/en/guides/cli) を使って適当なモデルを取得する。

```bash
$ hf download Qwen/Qwen3-0.6B --local-dir models/Qwen3-0.6B
```

ダウンロードしてきたモデルの中身を見てみる

```
$ ll -h models/Qwen3-0.6B
Permissions Size User  Group Date Modified    Name
drwxr-xr-x     - thara staff 2026-09-05 19:41 .cache
.rw-r--r--  1.6k thara staff 2026-09-05 19:41 .gitattributes
.rw-r--r--   11k thara staff 2026-09-05 19:41 LICENSE
.rw-r--r--   14k thara staff 2026-09-05 19:41 README.md
.rw-r--r--   726 thara staff 2026-09-05 19:41 config.json
.rw-r--r--   239 thara staff 2026-09-05 19:41 generation_config.json
.rw-r--r--  1.7M thara staff 2026-09-05 19:41 merges.txt
.rw-r--r--  1.5G thara staff 2026-09-05 19:43 model.safetensors
.rw-r--r--   11M thara staff 2026-09-05 19:41 tokenizer.json
.rw-r--r--  9.7k thara staff 2026-09-05 19:41 tokenizer_config.json
.rw-r--r--  2.8M thara staff 2026-09-05 19:41 vocab.json

$ du -sh models/Qwen3-0.6B
1.4G    models/Qwen3-0.6B

$ ll -h models/Qwen3-0.6B/*.safetensors
Permissions Size User  Group Date Modified    Name
.rw-r--r--  1.5G thara staff 2026-09-05 19:43 models/Qwen3-0.6B/model.safetensors

$ cat models/Qwen3-0.6B/config.json
{
  "architectures": [
    "Qwen3ForCausalLM"
  ],
  "attention_bias": false,
  "attention_dropout": 0.0,
  "bos_token_id": 151643,
  "eos_token_id": 151645,
  "head_dim": 128,
  "hidden_act": "silu",
  "hidden_size": 1024,
  "initializer_range": 0.02,
  "intermediate_size": 3072,
  "max_position_embeddings": 40960,
  "max_window_layers": 28,
  "model_type": "qwen3",
  "num_attention_heads": 16,
  "num_hidden_layers": 28,
  "num_key_value_heads": 8,
  "rms_norm_eps": 1e-06,
  "rope_scaling": null,
  "rope_theta": 1000000,
  "sliding_window": null,
  "tie_word_embeddings": true,
  "torch_dtype": "bfloat16",
  "transformers_version": "4.51.0",
  "use_cache": true,
  "use_sliding_window": false,
  "vocab_size": 151936
}
```

config.jsonの値とかを参照すればパラメータ数とか算出できるはず。


## GGUFに変換

モデル -> GGUF -> 量子化 -> 推論という工程を見たいので、自前でやってみる。

llama.cppをインストールしてセットアップ。

```
$ git clone https://github.com/ggml-org/llama.cpp
$ cd llama.cpp

$ python3 -m venv .venv
$ source .venv/bin/activate

$ pip install -r requirements.txt
```

GGUFに変換。
(ここのf16はtensorの出力データ型の指定）

```
$ python convert_hf_to_gguf.py \
  ../models/Qwen3-0.6B \
  --outfile ../models/qwen3-f16.gguf \
  --outtype f16
...
INFO:hf-to-gguf:Model successfully exported to ../models/qwen3-f16.gguf

$ ll -h ../models/*.gguf
Permissions Size User  Group Date Modified    Name
.rw-r--r--  1.5G thara staff 2026-09-05 19:53 ../models/qwen3-f16.gguf
```

## 量子化

llama.cppに量子化コマンドがあるのでビルドする。
(別プロセスのvimがカクつくほどめっちゃ重い)

```
$ cmake -B build
$ cmake --build build --config Release -j
```

Q4_K_M で量子化。

```
$ ./build/bin/llama-quantize \
  ../models/qwen3-f16.gguf \
  ../models/qwen3-q4_k_m.gguf \
  Q4_K_M
...
llama_model_quantize_impl: model size  =  1433.75 MiB (16.00 BPW)
llama_model_quantize_impl: quant size  =   456.11 MiB (5.09 BPW)

llama_quantize: quantize time =  6560.77 ms
llama_quantize:    total time =  6560.77 ms

$ ls -lh ../models/
Permissions Size User  Date Modified     Name                                 Permissions Size User  Date Modified     Name
drwxr-xr-x     - thara 2026-09-05 16:01  .cache                               .rw-r--r--  1.5G thara 2026-09-05 19:53  qwen3-f16.gguf
drwxr-xr-x     - thara 2026-09-05 19:43  Qwen3-0.6B                           .rw-r--r--  484M thara 2026-09-05 19:59  qwen3-q4_k_m.gguf
.rw-r--r--  986M thara 2026-09-05 16:03  Qwen2.5-1.5B-Instruct-Q4_K_M.gguf
```

### llama.cppでの推論

試しに量子化したモデルで推論してみる。

CPU使用:

```
$ ./build/bin/llama-cli \
  -m ../models/qwen3-q4_k_m.gguf \
  -p "TCPとUDPの違いを説明して"
...
[ Prompt: 378.5 t/s | Generation: 140.7 t/s ]
```

GPU使用:

```
$ ./build/bin/llama-cli \
  -m ../models/qwen3-q4_k_m.gguf \
  -ngl 99 \
  -p "TCPとUDPの違いを説明して"
...
> [ Prompt: 361.2 t/s | Generation: 142.6 t/s ]
```

PromptはGPUの方が並列化されているので速いが、Generationの方は逐次処理なので影響がないことが分かる。

### パラメータの実態

パラメータの中身は実際どんなもんだろな、というのが気になったので見てみる。

以下のようなPythonスクリプトを作る。

```
from safetensors import safe_open

path = "models/Qwen3-0.6B/model.safetensors"

total = 0

with safe_open(path, framework="pt") as f:
    for name in f.keys():
        tensor = f.get_tensor(name)
        n = tensor.numel()
        total += n
        print(f"{name:60} {str(tuple(tensor.shape)):20} {n:,}")

print()
print(f"Total: {total:,}")
print(f"Total: {total / 1e9:.3f}B")
```

実在するtensorを全部列挙できる。

```
$ python3 ./list_params.py
lm_head.weight                                               (151936, 1024)       155,582,464
model.embed_tokens.weight                                    (151936, 1024)       155,582,464
model.layers.0.input_layernorm.weight                        (1024,)              1,024
model.layers.0.mlp.down_proj.weight                          (1024, 3072)         3,145,728
model.layers.0.mlp.gate_proj.weight                          (3072, 1024)         3,145,728
model.layers.0.mlp.up_proj.weight                            (3072, 1024)         3,145,728
model.layers.0.post_attention_layernorm.weight               (1024,)              1,024
model.layers.0.self_attn.k_norm.weight                       (128,)               128
model.layers.0.self_attn.k_proj.weight                       (1024, 1024)         1,048,576
model.layers.0.self_attn.o_proj.weight                       (1024, 2048)         2,097,152
model.layers.0.self_attn.q_norm.weight                       (128,)               128
model.layers.0.self_attn.q_proj.weight                       (2048, 1024)         2,097,152
model.layers.0.self_attn.v_proj.weight                       (1024, 1024)         1,048,576
model.layers.1.input_layernorm.weight                        (1024,)              1,024
model.layers.1.mlp.down_proj.weight                          (1024, 3072)         3,145,728
model.layers.1.mlp.gate_proj.weight                          (3072, 1024)         3,145,728
model.layers.1.mlp.up_proj.weight                            (3072, 1024)         3,145,728
model.layers.1.post_attention_layernorm.weight               (1024,)              1,024
model.layers.1.self_attn.k_norm.weight                       (128,)               128
model.layers.1.self_attn.k_proj.weight                       (1024, 1024)         1,048,576
model.layers.1.self_attn.o_proj.weight                       (1024, 2048)         2,097,152
model.layers.1.self_attn.q_norm.weight                       (128,)               128
model.layers.1.self_attn.q_proj.weight                       (2048, 1024)         2,097,152
model.layers.1.self_attn.v_proj.weight                       (1024, 1024)         1,048,576
model.layers.10.input_layernorm.weight                       (1024,)              1,024
model.layers.10.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.10.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.10.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.10.post_attention_layernorm.weight              (1024,)              1,024
model.layers.10.self_attn.k_norm.weight                      (128,)               128
model.layers.10.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.10.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.10.self_attn.q_norm.weight                      (128,)               128
model.layers.10.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.10.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.11.input_layernorm.weight                       (1024,)              1,024
model.layers.11.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.11.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.11.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.11.post_attention_layernorm.weight              (1024,)              1,024
model.layers.11.self_attn.k_norm.weight                      (128,)               128
model.layers.11.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.11.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.11.self_attn.q_norm.weight                      (128,)               128
model.layers.11.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.11.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.12.input_layernorm.weight                       (1024,)              1,024
model.layers.12.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.12.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.12.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.12.post_attention_layernorm.weight              (1024,)              1,024
model.layers.12.self_attn.k_norm.weight                      (128,)               128
model.layers.12.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.12.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.12.self_attn.q_norm.weight                      (128,)               128
model.layers.12.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.12.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.13.input_layernorm.weight                       (1024,)              1,024
model.layers.13.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.13.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.13.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.13.post_attention_layernorm.weight              (1024,)              1,024
model.layers.13.self_attn.k_norm.weight                      (128,)               128
model.layers.13.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.13.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.13.self_attn.q_norm.weight                      (128,)               128
model.layers.13.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.13.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.14.input_layernorm.weight                       (1024,)              1,024
model.layers.14.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.14.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.14.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.14.post_attention_layernorm.weight              (1024,)              1,024
model.layers.14.self_attn.k_norm.weight                      (128,)               128
model.layers.14.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.14.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.14.self_attn.q_norm.weight                      (128,)               128
model.layers.14.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.14.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.15.input_layernorm.weight                       (1024,)              1,024
model.layers.15.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.15.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.15.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.15.post_attention_layernorm.weight              (1024,)              1,024
model.layers.15.self_attn.k_norm.weight                      (128,)               128
model.layers.15.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.15.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.15.self_attn.q_norm.weight                      (128,)               128
model.layers.15.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.15.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.16.input_layernorm.weight                       (1024,)              1,024
model.layers.16.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.16.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.16.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.16.post_attention_layernorm.weight              (1024,)              1,024
model.layers.16.self_attn.k_norm.weight                      (128,)               128
model.layers.16.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.16.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.16.self_attn.q_norm.weight                      (128,)               128
model.layers.16.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.16.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.17.input_layernorm.weight                       (1024,)              1,024
model.layers.17.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.17.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.17.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.17.post_attention_layernorm.weight              (1024,)              1,024
model.layers.17.self_attn.k_norm.weight                      (128,)               128
model.layers.17.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.17.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.17.self_attn.q_norm.weight                      (128,)               128
model.layers.17.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.17.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.18.input_layernorm.weight                       (1024,)              1,024
model.layers.18.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.18.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.18.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.18.post_attention_layernorm.weight              (1024,)              1,024
model.layers.18.self_attn.k_norm.weight                      (128,)               128
model.layers.18.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.18.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.18.self_attn.q_norm.weight                      (128,)               128
model.layers.18.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.18.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.19.input_layernorm.weight                       (1024,)              1,024
model.layers.19.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.19.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.19.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.19.post_attention_layernorm.weight              (1024,)              1,024
model.layers.19.self_attn.k_norm.weight                      (128,)               128
model.layers.19.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.19.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.19.self_attn.q_norm.weight                      (128,)               128
model.layers.19.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.19.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.2.input_layernorm.weight                        (1024,)              1,024
model.layers.2.mlp.down_proj.weight                          (1024, 3072)         3,145,728
model.layers.2.mlp.gate_proj.weight                          (3072, 1024)         3,145,728
model.layers.2.mlp.up_proj.weight                            (3072, 1024)         3,145,728
model.layers.2.post_attention_layernorm.weight               (1024,)              1,024
model.layers.2.self_attn.k_norm.weight                       (128,)               128
model.layers.2.self_attn.k_proj.weight                       (1024, 1024)         1,048,576
model.layers.2.self_attn.o_proj.weight                       (1024, 2048)         2,097,152
model.layers.2.self_attn.q_norm.weight                       (128,)               128
model.layers.2.self_attn.q_proj.weight                       (2048, 1024)         2,097,152
model.layers.2.self_attn.v_proj.weight                       (1024, 1024)         1,048,576
model.layers.20.input_layernorm.weight                       (1024,)              1,024
model.layers.20.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.20.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.20.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.20.post_attention_layernorm.weight              (1024,)              1,024
model.layers.20.self_attn.k_norm.weight                      (128,)               128
model.layers.20.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.20.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.20.self_attn.q_norm.weight                      (128,)               128
model.layers.20.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.20.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.21.input_layernorm.weight                       (1024,)              1,024
model.layers.21.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.21.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.21.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.21.post_attention_layernorm.weight              (1024,)              1,024
model.layers.21.self_attn.k_norm.weight                      (128,)               128
model.layers.21.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.21.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.21.self_attn.q_norm.weight                      (128,)               128
model.layers.21.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.21.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.22.input_layernorm.weight                       (1024,)              1,024
model.layers.22.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.22.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.22.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.22.post_attention_layernorm.weight              (1024,)              1,024
model.layers.22.self_attn.k_norm.weight                      (128,)               128
model.layers.22.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.22.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.22.self_attn.q_norm.weight                      (128,)               128
model.layers.22.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.22.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.23.input_layernorm.weight                       (1024,)              1,024
model.layers.23.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.23.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.23.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.23.post_attention_layernorm.weight              (1024,)              1,024
model.layers.23.self_attn.k_norm.weight                      (128,)               128
model.layers.23.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.23.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.23.self_attn.q_norm.weight                      (128,)               128
model.layers.23.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.23.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.24.input_layernorm.weight                       (1024,)              1,024
model.layers.24.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.24.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.24.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.24.post_attention_layernorm.weight              (1024,)              1,024
model.layers.24.self_attn.k_norm.weight                      (128,)               128
model.layers.24.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.24.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.24.self_attn.q_norm.weight                      (128,)               128
model.layers.24.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.24.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.25.input_layernorm.weight                       (1024,)              1,024
model.layers.25.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.25.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.25.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.25.post_attention_layernorm.weight              (1024,)              1,024
model.layers.25.self_attn.k_norm.weight                      (128,)               128
model.layers.25.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.25.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.25.self_attn.q_norm.weight                      (128,)               128
model.layers.25.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.25.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.26.input_layernorm.weight                       (1024,)              1,024
model.layers.26.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.26.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.26.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.26.post_attention_layernorm.weight              (1024,)              1,024
model.layers.26.self_attn.k_norm.weight                      (128,)               128
model.layers.26.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.26.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.26.self_attn.q_norm.weight                      (128,)               128
model.layers.26.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.26.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.27.input_layernorm.weight                       (1024,)              1,024
model.layers.27.mlp.down_proj.weight                         (1024, 3072)         3,145,728
model.layers.27.mlp.gate_proj.weight                         (3072, 1024)         3,145,728
model.layers.27.mlp.up_proj.weight                           (3072, 1024)         3,145,728
model.layers.27.post_attention_layernorm.weight              (1024,)              1,024
model.layers.27.self_attn.k_norm.weight                      (128,)               128
model.layers.27.self_attn.k_proj.weight                      (1024, 1024)         1,048,576
model.layers.27.self_attn.o_proj.weight                      (1024, 2048)         2,097,152
model.layers.27.self_attn.q_norm.weight                      (128,)               128
model.layers.27.self_attn.q_proj.weight                      (2048, 1024)         2,097,152
model.layers.27.self_attn.v_proj.weight                      (1024, 1024)         1,048,576
model.layers.3.input_layernorm.weight                        (1024,)              1,024
model.layers.3.mlp.down_proj.weight                          (1024, 3072)         3,145,728
model.layers.3.mlp.gate_proj.weight                          (3072, 1024)         3,145,728
model.layers.3.mlp.up_proj.weight                            (3072, 1024)         3,145,728
model.layers.3.post_attention_layernorm.weight               (1024,)              1,024
model.layers.3.self_attn.k_norm.weight                       (128,)               128
model.layers.3.self_attn.k_proj.weight                       (1024, 1024)         1,048,576
model.layers.3.self_attn.o_proj.weight                       (1024, 2048)         2,097,152
model.layers.3.self_attn.q_norm.weight                       (128,)               128
model.layers.3.self_attn.q_proj.weight                       (2048, 1024)         2,097,152
model.layers.3.self_attn.v_proj.weight                       (1024, 1024)         1,048,576
model.layers.4.input_layernorm.weight                        (1024,)              1,024
model.layers.4.mlp.down_proj.weight                          (1024, 3072)         3,145,728
model.layers.4.mlp.gate_proj.weight                          (3072, 1024)         3,145,728
model.layers.4.mlp.up_proj.weight                            (3072, 1024)         3,145,728
model.layers.4.post_attention_layernorm.weight               (1024,)              1,024
model.layers.4.self_attn.k_norm.weight                       (128,)               128
model.layers.4.self_attn.k_proj.weight                       (1024, 1024)         1,048,576
model.layers.4.self_attn.o_proj.weight                       (1024, 2048)         2,097,152
model.layers.4.self_attn.q_norm.weight                       (128,)               128
model.layers.4.self_attn.q_proj.weight                       (2048, 1024)         2,097,152
model.layers.4.self_attn.v_proj.weight                       (1024, 1024)         1,048,576
model.layers.5.input_layernorm.weight                        (1024,)              1,024
model.layers.5.mlp.down_proj.weight                          (1024, 3072)         3,145,728
model.layers.5.mlp.gate_proj.weight                          (3072, 1024)         3,145,728
model.layers.5.mlp.up_proj.weight                            (3072, 1024)         3,145,728
model.layers.5.post_attention_layernorm.weight               (1024,)              1,024
model.layers.5.self_attn.k_norm.weight                       (128,)               128
model.layers.5.self_attn.k_proj.weight                       (1024, 1024)         1,048,576
model.layers.5.self_attn.o_proj.weight                       (1024, 2048)         2,097,152
model.layers.5.self_attn.q_norm.weight                       (128,)               128
model.layers.5.self_attn.q_proj.weight                       (2048, 1024)         2,097,152
model.layers.5.self_attn.v_proj.weight                       (1024, 1024)         1,048,576
model.layers.6.input_layernorm.weight                        (1024,)              1,024
model.layers.6.mlp.down_proj.weight                          (1024, 3072)         3,145,728
model.layers.6.mlp.gate_proj.weight                          (3072, 1024)         3,145,728
model.layers.6.mlp.up_proj.weight                            (3072, 1024)         3,145,728
model.layers.6.post_attention_layernorm.weight               (1024,)              1,024
model.layers.6.self_attn.k_norm.weight                       (128,)               128
model.layers.6.self_attn.k_proj.weight                       (1024, 1024)         1,048,576
model.layers.6.self_attn.o_proj.weight                       (1024, 2048)         2,097,152
model.layers.6.self_attn.q_norm.weight                       (128,)               128
model.layers.6.self_attn.q_proj.weight                       (2048, 1024)         2,097,152
model.layers.6.self_attn.v_proj.weight                       (1024, 1024)         1,048,576
model.layers.7.input_layernorm.weight                        (1024,)              1,024
model.layers.7.mlp.down_proj.weight                          (1024, 3072)         3,145,728
model.layers.7.mlp.gate_proj.weight                          (3072, 1024)         3,145,728
model.layers.7.mlp.up_proj.weight                            (3072, 1024)         3,145,728
model.layers.7.post_attention_layernorm.weight               (1024,)              1,024
model.layers.7.self_attn.k_norm.weight                       (128,)               128
model.layers.7.self_attn.k_proj.weight                       (1024, 1024)         1,048,576
model.layers.7.self_attn.o_proj.weight                       (1024, 2048)         2,097,152
model.layers.7.self_attn.q_norm.weight                       (128,)               128
model.layers.7.self_attn.q_proj.weight                       (2048, 1024)         2,097,152
model.layers.7.self_attn.v_proj.weight                       (1024, 1024)         1,048,576
model.layers.8.input_layernorm.weight                        (1024,)              1,024
model.layers.8.mlp.down_proj.weight                          (1024, 3072)         3,145,728
model.layers.8.mlp.gate_proj.weight                          (3072, 1024)         3,145,728
model.layers.8.mlp.up_proj.weight                            (3072, 1024)         3,145,728
model.layers.8.post_attention_layernorm.weight               (1024,)              1,024
model.layers.8.self_attn.k_norm.weight                       (128,)               128
model.layers.8.self_attn.k_proj.weight                       (1024, 1024)         1,048,576
model.layers.8.self_attn.o_proj.weight                       (1024, 2048)         2,097,152
model.layers.8.self_attn.q_norm.weight                       (128,)               128
model.layers.8.self_attn.q_proj.weight                       (2048, 1024)         2,097,152
model.layers.8.self_attn.v_proj.weight                       (1024, 1024)         1,048,576
model.layers.9.input_layernorm.weight                        (1024,)              1,024
model.layers.9.mlp.down_proj.weight                          (1024, 3072)         3,145,728
model.layers.9.mlp.gate_proj.weight                          (3072, 1024)         3,145,728
model.layers.9.mlp.up_proj.weight                            (3072, 1024)         3,145,728
model.layers.9.post_attention_layernorm.weight               (1024,)              1,024
model.layers.9.self_attn.k_norm.weight                       (128,)               128
model.layers.9.self_attn.k_proj.weight                       (1024, 1024)         1,048,576
model.layers.9.self_attn.o_proj.weight                       (1024, 2048)         2,097,152
model.layers.9.self_attn.q_norm.weight                       (128,)               128
model.layers.9.self_attn.q_proj.weight                       (2048, 1024)         2,097,152
model.layers.9.self_attn.v_proj.weight                       (1024, 1024)         1,048,576
model.norm.weight                                            (1024,)              1,024

Total: 751,632,384
Total: 0.752B
```

なるほどわからん。

## ベンチマーク

使ったモデルに対して、異なる量子化手法を使って量子化する。

今回は、F16, Q8_0, Q4_K_M を使う。

```
#!/usr/bin/env bash
set -euo pipefail

MODEL_DIR="${1:-../models/Qwen3-0.6B}"
OUT_DIR="${2:-../models/gguf}"

MODEL_NAME="$(basename "$MODEL_DIR")"

LLAMA_DIR="$(cd "$(dirname "$0")" && pwd)"
CONVERT="$LLAMA_DIR/convert_hf_to_gguf.py"
QUANTIZE="$LLAMA_DIR/build/bin/llama-quantize"

mkdir -p "$OUT_DIR"

F16="$OUT_DIR/${MODEL_NAME}-F16.gguf"
Q8="$OUT_DIR/${MODEL_NAME}-Q8_0.gguf"
Q4="$OUT_DIR/${MODEL_NAME}-Q4_K_M.gguf"

echo "==> Converting Hugging Face model to F16 GGUF"
python "$CONVERT" "$MODEL_DIR" --outfile "$F16" --outtype f16

echo
echo "==> Quantizing to Q8_0"
"$QUANTIZE" "$F16" "$Q8" Q8_0

echo
echo "==> Quantizing to Q4_K_M"
"$QUANTIZE" "$F16" "$Q4" Q4_K_M

echo
echo "==> Done"
ls -lh "$F16" "$Q8" "$Q4"
```

```
./build-models.sh ../models/Qwen3-0.6B ../models/gguf
...
-rw-r--r--  1 thara  staff   1.4G Sep  6 16:17 ../models/gguf/Qwen3-0.6B-F16.gguf
-rw-r--r--  1 thara  staff   462M Sep  6 16:17 ../models/gguf/Qwen3-0.6B-Q4_K_M.gguf
-rw-r--r--  1 thara  staff   767M Sep  6 16:17 ../models/gguf/Qwen3-0.6B-Q8_0.gguf
```


できたモデルを使って、それぞれCPU/GPUバックエンドでの性能を見てみる。

```
#!/usr/bin/env bash
set -euo pipefail

MODEL_DIR="${1:-../models/gguf}"
LLAMA_BENCH="${LLAMA_BENCH:-./build/bin/llama-bench}"

shopt -s nullglob
models=("$MODEL_DIR"/*.gguf)

run_bench() {
  local model="$1"
  local backend="$2"

  local result

  case "$backend" in
    CPU)
      result="$("$LLAMA_BENCH" -m "$model" -p 512 -n 128 -r 5 -dev none -ngl 0 -o json 2>/dev/null)"
      ;;

    Metal)
      result="$("$LLAMA_BENCH" -m "$model" -p 512 -n 128 -r 5 -ngl 99 -o json 2>/dev/null)"
      ;;
    *)
      echo "unknown backend: $backend" >&2
      exit 1
      ;;
  esac

  local prompt
  local generation

  prompt="$(
    jq -r '.[] | select(.n_prompt == 512 and .n_gen == 0) | .avg_ts' <<< "$result"
  )"

  generation="$(
    jq -r '.[] | select(.n_prompt == 0 and .n_gen == 128) | .avg_ts' <<< "$result"
  )"

  printf "%-32s %-8s %12.1f t/s %14.1f t/s\n" "$(basename "$model")" "$backend" "$prompt" "$generation"
}

printf "%-32s %-8s %16s %18s\n" "MODEL" "BACKEND" "PROMPT" "GENERATION"
printf "%-32s %-8s %16s %18s\n" "--------------------------------" "--------" "----------------" "------------------"

for model in "${models[@]}"; do
  run_bench "$model" "CPU"
  run_bench "$model" "Metal"
done
```

```
[17:18:41] ❯❯❯ ./benchmark-all.sh ../models/gguf
MODEL                            SIZE       BACKEND  PROMPT           GENERATION
-------------------------------- ---------- -------- ---------------- ------------------
Qwen3-0.6B-F16.gguf               1439.4 MB CPU             631.8 t/s           51.4 t/s
Qwen3-0.6B-F16.gguf               1439.4 MB Metal          2647.4 t/s           70.0 t/s
Qwen3-0.6B-Q4_K_M.gguf             461.8 MB CPU             554.2 t/s          133.7 t/s
Qwen3-0.6B-Q4_K_M.gguf             461.8 MB Metal          2233.7 t/s          162.6 t/s
Qwen3-0.6B-Q8_0.gguf               767.5 MB CPU             690.7 t/s           92.4 t/s
Qwen3-0.6B-Q8_0.gguf               767.5 MB Metal          2452.5 t/s          111.5 t/s
```

ファイルサイズは F16 -> Q8で約53%、F16 -> Q4_K_Mで約32%まで縮んでる。
が、Prompt処理速度はGPUバックエンド(Macbook AirなのでMetal)のF16が一番早い。
一方、Generation速度はQ4_K_Mが圧倒的に速い。

量子化したら計算量が減るんだから全部速くなる、というわけではないことがわかる。
量子化は「データが小さい」メリットがある一方、圧縮表現を扱うための追加処理もあるため。

Prompt処理はスループットの影響が大きく、Generationはメモリ帯域の影響が大きい。
(なので、Generationには量子化が効く)

---

大体、実体験としては、こんな感じか。

次は、onnxあたりを試してみるかな。
