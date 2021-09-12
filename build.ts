import * as path from "https://deno.land/std/path/mod.ts";
import * as datetime from 'https://deno.land/std/datetime/mod.ts';
import * as flags from "https://deno.land/std/flags/mod.ts";
import { ensureDirSync, walkSync, copySync } from "https://deno.land/std/fs/mod.ts";

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

  const templates = await getTemplatePathMap(templateDir);

  // posts
  const posts = await buildPosts(postsSrc, postsDst, config);
  // other pages
  await buildPages({ posts: posts, ...config }, templates);
  // assets
  copySync("css", path.join(dstRoot, "css"), { overwrite: true });
  copySync("images", path.join(dstRoot, "images"), { overwrite: true });
}

function md2html(srcPath: string) {
  const c = decoder.decode(Deno.readFileSync(srcPath));
  return Marked.parse(c);
}

function parsePageName(name: string): {date: string | null, pageName: string} {
  const date = name.substring(0, 10);
  try {
    datetime.parse(date, "yyyy-MM-dd");
    return { date, pageName: name.substring(11) };
  } catch {
    return { date: null, pageName: name};
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
  const entries = Deno.readDirSync(postsDir);
  const posts = Array.from(entries)
    .filter((e) => e.isFile)
    .map((e) => {
      const p = path.join(postsDir, e.name);
      const { meta, content } = md2html(p);
      const { title, date: pageDate } = meta;
      const { date } = parsePageName(e.name)!;
      const timestamp = pageDate ?? date;
      const { name: dstName, path: dstPath } = buildPagePath(p, dstDir);
      buildPage(dstPath, content, "templates/post.ejs", {
          pageTitle: title,
          pageCreatedDate: date,
          ...config
      });
      return { timestamp, title, date, path: path.join(postsDir, dstName)};
    });
  return (await Promise.all(posts)).sort(({timestamp: a}, {timestamp: b}) => a < b ? 1 : -1);
}

type TemplatePathMap = Map<string, string>;

async function getTemplatePathMap(templateDir: string) {
  const entries = Deno.readDirSync(templateDir);
  const templates = Array.from(entries)
    .filter((e) => !e.isSymlink && !e.isDirectory && path.extname(e.name) == ".ejs")
    .map((e) => {
      const p = path.join(templateDir, e.name);
      const { name } = path.parse(p);
      return [name, p] as [string, string];
    });
  return new Map(templates) as TemplatePathMap;
}

async function buildPages(config: Config, templates: TemplatePathMap) {
  const entries = Array.from(walkSync(srcRoot));
  const p = entries
    .filter((e) => !e.isSymlink && !e.isDirectory && path.extname(e.name) == ".md" && !e.path.startsWith("posts/"))
    .map((e) => {
        const { name } = path.parse(e.path);
        const template = templates.get(name) ?? defaultTemplatePath;
        const { meta, content } = md2html(e.path);
        const { title, path: metaPath } = meta;
        const dstPath = metaPath ? path.join(dstRoot, metaPath) : buildPagePath(e.path, dstRoot).path;
        buildPage(dstPath, content, template, {
          pageTitle: title,
          ...config
        });
    });
  await Promise.all(p);
}

if (import.meta.main) {
  await main();
}
