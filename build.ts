import * as path from "https://deno.land/std/path/mod.ts";
import * as datetime from 'https://deno.land/std/datetime/mod.ts';
import * as flags from "https://deno.land/std/flags/mod.ts";
import { ensureDirSync, walk, copySync } from "https://deno.land/std/fs/mod.ts";

import { Marked } from 'https://deno.land/x/markdown@v2.0.0/mod.ts';
import { renderToString, renderFileToString, Params } from "https://deno.land/x/dejs/mod.ts";

const srcRoot = ".";
const dstRoot = "public";
const templateDir = "templates";
const defaultTemplatePath = path.join(templateDir, "layout.ejs");

type Config = {[k: string]: any};

const decoder = new TextDecoder("utf-8");

interface Args {
  // -b --base-rul
  "base-url"?: string;
  "b"?: string;
}
const args = flags.parse(Deno.args) as Args;

async function main() {
  const templates = await getTemplatePathMap(templateDir);

  const baseUrl = args['base-url'] ?? args.b ?? "https://thara.github.io"

  const config: Config = {
    author: "Tomochika Hara",
    siteTitle: "thara",
    baseUrl: baseUrl,
    year: datetime.format(new Date(), 'yyyy'),
  };

  const postsSrc = path.join(srcRoot, "posts");
  const postsDst = path.join(dstRoot, "posts");

  ensureDirSync(postsSrc);
  ensureDirSync(postsDst);

  const posts = await buildPosts(postsSrc, postsDst, config);

  await buildPages({ posts: posts, ...config }, templates);

  copySync("css", path.join(dstRoot, "css"), { overwrite: true });
}

async function md2html(srcPath: string) {
  const content = decoder.decode(await Deno.readFile(srcPath));
  return Marked.parse(content);
}

function parsePageName(name: string): {date: string | null, pageName: string} {
  const date = name.substring(0, 10);
  try {
    datetime.parse(date, "yyyy-MM-dd");
    return {
      date: date,
      pageName: name.substring(11),
    };
  } catch {
    return {date: null, pageName: name};
  }
}

function buildPagePath(srcPath: string, dstDir: string) {
  const { name } = path.parse(srcPath);
  const { pageName } = parsePageName(name);
  const dstPath = path.format({dir: dstDir, name: pageName, ext: ""});
  return {name: pageName, path: dstPath};
}

async function buildPage(dstPath: string, content: string, templatePath: string, config: Params) {
  const c = await renderFileToString(templatePath, { content: content, ...config });
  await Deno.writeTextFile(dstPath, c);
}

interface Post {
  timestamp: string;
  date: string;
  title: string;
  path: string;
}

async function buildPosts(postsDir: string, dstDir: string, config: Config) {
  const entries = Deno.readDir(postsDir);

  var posts: Post[] = [];
  for await (const e of entries) {
    if (e.isFile) {
      const p = path.join(postsDir, e.name);
      const { meta, content } = await md2html(p);
      const { title, date: pageDate } = meta;

      const { date } = parsePageName(e.name)!;
      if (!date) {
        continue;
      }

      const timestamp = pageDate ?? date;

      const { name: dstName, path: dstPath } = buildPagePath(p, dstDir);

      await buildPage(dstPath, content, "templates/post.ejs", {
          pageTitle: title,
          pageCreatedDate: date!,
          ...config
      });
      posts.push({
        timestamp: timestamp,
        date: date!,
        title: title,
        path: path.join(postsDir, dstName),
      });
    }
  }
  return posts.sort(({timestamp: a}, {timestamp: b}) => a < b ? 1 : -1);
}

type TemplatePathMap = {[k: string]: string};

async function getTemplatePathMap(templateDir: string) {
  var m: TemplatePathMap = {};
  for await (const e of Deno.readDir(templateDir)) {
    if (e.isSymlink || e.isDirectory || path.extname(e.name) != ".ejs") {
      continue;
    }
    const p = path.join(templateDir, e.name);
    const { name } = path.parse(p);
    m[name] = p;
  }
  return m;
}

async function buildPages(config: Config, templates: TemplatePathMap) {
  for await (const e of walk(srcRoot)) {
    if (e.isSymlink || e.isDirectory || path.extname(e.name) != ".md" || e.path.startsWith("posts/")) {
      continue;
    }

    switch (path.extname(e.name)) {
      case ".md":
        const { name } = path.parse(e.path);
        const template = templates[name] ?? defaultTemplatePath;

        const { meta, content } = await md2html(e.path);
        const { title, path: metaPath } = meta;

        const dstPath = (() => {
          if (metaPath) {
            return path.join(dstRoot, metaPath);
          } else {
            const { path: p } = buildPagePath(e.path, dstRoot);
            return p;
          }
        })();

        await buildPage(dstPath, content, template, {
          pageTitle: title,
          ...config
        });
        break;
      default:
        continue;
    }
  }
}

if (import.meta.main) {
  await main();
}
