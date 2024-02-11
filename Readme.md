## Proton-Assignment CLI Program: Markdown to HTML Blog Converter

This Go-based command-line interface (CLI) program converts markdown files into an HTML blog format. Each markdown file is rendered as an individual blog post within the generated HTML output.

### Usage

```bash
generate <input> <output> <title> [additional-options]
```

### Example usage

```bash
gen-blog generate --input inputFolder --output outputFolder --title "The Proton Blog"
```

```bash
gen-blog generate --input inputFolder --output outputFolder --title "The Proton Blog" --posts-per-page 2 --s --m
```

### Options

- `--input`: A required argument representing the valid input folder path containing markdown files.
- `--output`: A required argument representing the valid output folder path where HTML files will be generated.
- `--title`: A required string argument representing the title for the generated blog.

### Additional Options

- `--posts-per-page <int>`: Optional. Specifies the number of posts per single page. Default is 0.
- `--separator, --s`: Optional. Disables the separator at the end of each blog post. Default is false.
- `--multithread, --m`: Optional. Enables multithreaded output. Default is false.

### General Notes:

- Code was developed on an M1/ARM MacOS system.
- Code was tested on both MacOS and Windows.
