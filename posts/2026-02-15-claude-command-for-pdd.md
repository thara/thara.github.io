---
title: Protocol Driven DevelopmentのためのClaude Code slash command
date: '2026-02-15'
published: '2026-02-15'
---

作った。実際に使ってみてどうかは、これから検証する。

Claude codeのslash commandはskillに統合されたっぽいのだけれど、勝手に読み込んでほしくはないので description は指定しない。
また、その意図を自分が忘れないように `.claude/commands` に配置している。


```markdown
---
disable-model-invocation: true
---

Proceed with Protocol-Driven Development.

## What is Protocol-Driven Development

A development methodology where you define protocols to understand the problem, then implement those protocols to solve it.

Defining a protocol is itself an act of deepening problem understanding. By defining protocols and getting them to compile, the structure of the problem becomes clear and implementation details emerge.

"Protocol" here is not limited to Swift's language feature. It refers to any "contract" that makes problem structure and responsibilities explicit:
- Type definitions: interface, type, trait
- Data structures: struct, enum
- Tests: unit tests, property-based tests

## Role Division

**Human responsibilities (AI must not take over):**
- Define protocols that represent the problem structure
- Write tests for constraints that the compiler cannot express

**AI responsibilities:**
- Implement the protocols that humans have defined
- Run the feedback loop until tests pass
- Do not deviate from the intent of the protocols

## Workflow

1. Check if protocols are defined
2. If not, prompt the human to define them or clarify the structure through dialogue
3. Implement using protocols and tests as guides
4. Verify with compiler and tests

## Do Not

- Define protocols without human input (this is the human's responsibility)
- Modify human-defined structures without permission
- Consider implementation complete without tests
```

とりあえず、これでターミナルエミュレータでも作ってみようかな。
