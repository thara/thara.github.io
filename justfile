build:
  go build .

sitegen baseURL:
  go run . --base {{baseURL}}

serve port *FLAGS:
  go run . --serve --port {{port}} --base http://localhost:{{port}} {{FLAGS}}

post:
  read -p "post title: " title; (vim pages/posts/`date +%Y-%m-%d`-$title.md)
