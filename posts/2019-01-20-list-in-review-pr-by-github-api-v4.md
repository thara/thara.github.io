---
layout: post
title: GitHub API v4でレビューPRを一覧にするスクリプトを書いた
description: GitHub API v4 + graphql-clientによるレビューPRを一覧にするスクリプトの公開
tags: [ruby, github, graphql]
---

レビュー対象のPRをサクッとCLIで確認したいなーと思い、 せっかくならGitHub API v4を使おうと思ってGraphQLクライアントを探したら、GitHub自体が [github/graphql-client](https://github.com/github/graphql-client) というGraphQLクライアントを公開していたので、それを使ってみた。

実際のクエリを投げる前に、以下のようにHTTPクライアントの初期化やスキーマの読み込みが必要。

```ruby
require 'graphql/client'
require 'graphql/client/http'

TOKEN = ENV['GITHUB_ACCESS_TOKEN']
ENDPOINT = 'https://api.github.com/graphql'

module GitHub
  HTTP = GraphQL::Client::HTTP.new(ENDPOINT) do
    def headers(context)
      {'Authorization': "bearer #{TOKEN}" }
    end
  end

  Schema = GraphQL::Client.load_schema(HTTP)

  Client = GraphQL::Client.new(schema: Schema, execute: HTTP)
end
```

ここでは、 `REVIEW` というラベルがつけられたPRをリポジトリごとに一覧する。

```ruby
org = ARGV[0]

InReviewPullRequestQuery = GitHub::Client.parse <<-"GRAPHQL"
  query($org:String!) {
      organization(login: $org) {
      repositories(first: 100) {
        nodes {
          pullRequests(labels: "REVIEW", states: [OPEN], first: 100, orderBy: {field: CREATED_AT, direction: DESC}) {
            edges {
              node {
                number
                title
                url
                repository {
                  nameWithOwner
                }
              }
            }
          }
        }
      }
    }
  }
GRAPHQL

result = GitHub::Client.query(InReviewPullRequestQuery, variables: {org: org})

prs = result.data.organization.repositories.nodes.flat_map(&:pull_requests).flat_map(&:edges).flat_map(&:node)

repos = prs.group_by{|pr| pr.repository.name_with_owner}
repos.each do |repo_name, prs|
  puts "- #{repo_name}"
  prs.each do |pr|
    puts "  - [#{pr.title} ##{pr.number}](#{pr.url})"
  end
end
```

GraphQLのクエリを文字列で書いているところは、 [GraphQL API Explorer](https://developer.github.com/v4/explorer/) で実際にクエリを投げて検証できる。
Variableとかも対応しているし、補完も効くので、とても使いやすい。

GraphQLは、cURLで手軽に動作確認できないだろうと今まで手をつけてなかったけど、 API Explorerみたいなのがあると簡単に動作確認できてDX良いな、と思った。

API Explorerで書いたクエリを共有するような機能があると、もっと捗りそう。
