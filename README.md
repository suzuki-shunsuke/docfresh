# docfresh

docfresh is a CLI making document with command and code snippet maintainable, reusable, and testable.
It prevents document from being outdated.

## Status

This is still alpha.
Many features are not implemented yet, and the API is subject to change.
Please don't use this in production.

## Features

- Execute external commands and embed their output into document
- Test external commands in document
- Unify a template file and generated file, which improves the maintainability of document

Note that docfresh is intended to update markdown files.
Other markup language isn't supported.

<!--
- Fetch document from local and remote files and embed them into document
-->

## Getting Started

1. [Install docfresh](#install)

```sh
: Check version
docfresh -v
```

2. Checkout the repository

```sh
git clone https://github.com/suzuki-shunsuke/docfresh
cd examples
```

Please see [date.md](examples/date.md).
In this document, the result of `date` command is embedded.

```sh
cat date.md
```

Please run `docfresh run date.md` to update date.md.

```sh
docfresh run date.md
```

Then the datetime is updated.

## Motivation

Keeping documentation accurate is not easy.
Commands and code in documentation can quickly become outdated, which discourages readers.
When execution results are included in the documentation, it is tedious to manually rerun commands and update the results each time something changes.

With docfresh, commands in documentation can be executed automatically, and their results can be embedded directly into the documentation.
It also helps you quickly detect when commands start failing.

By running docfresh in CI, you can automate documentation updates and validation.

Another key feature of docfresh is that templates and generated files are unified, making documentation easier to maintain.

When templates and generated files are separate, it creates the question of where and how to manage the template files:

- Using a special extension like `.tpl`: syntax highlighting in editors may no longer work properly.
- Adding a suffix like `-tpl` to filenames: it becomes harder to navigate as the number of files increases, and static site generators may include template files unintentionally.
- Separating them into different directories: it becomes less obvious that templates exist and where they are located.

There are also practical issues.
Even if you delete a template file, the generated file may remain.
Since the editable file and the generated file are separate, it is also harder to edit while previewing the final output.

With docfresh, templates and generated files are unified, so these problems do not occur.

## Install

```sh
go install github.com/szksh-lab/docfresh/cmd/docfresh@latest
```

## Security

docfresh may execute arbitrary external commands defined in templates. Therefore, it is important to take appropriate security precautions.
Running docfresh on untrusted templates can be dangerous. It is recommended to execute docfresh in an isolated environment such as a container. Secrets should not be provided unless absolutely necessary.
Support for executing commands inside containers is also being considered for future releases.

## Template Syntax

In docfresh, instructions are embedded into Markdown using HTML comments, as shown below:

```md
<!-- docfresh begin
command:
  command: npm test
-->

The result will be embedded here.

<!-- docfresh end -->
```

Instructions are written in YAML format inside `<!-- docfresh begin -->`, and the execution result is embedded between `<!-- docfresh begin -->` and `<!-- docfresh end -->`.
Since HTML comments are not rendered in Markdown, they don't affect the view of the documentation.
Because this mechanism relies on HTML comments, docfresh is designed specifically for Markdown and does not support other document formats.
Each directive must start with `<!-- docfresh begin` and must have a corresponding closing `<!-- docfresh end -->`.

## Command and File Processing Order

docfresh executes all file processing and commands sequentially.
Commands within the same file are executed from top to bottom. If a command fails, the file will not be updated.

Support for parallel processing across multiple files may be added in the future.

## YAML Syntax In Begin Comment

```md
<!-- docfresh begin
command:
  command: npm test
  shell:
    - bash
    - "-c"
-->
```

- command.command: External Command
- command.shell: The list of shell command executing command. By default, `["bash", "-c"]`
- file.path: The relative path from the current file to the loaded file

### Run Command

```md
<!-- docfresh begin
command:
  command: npm test
-->
```

### Change Shell

```md
<!-- docfresh begin
command:
  command: echo hello
  shell:
    - zsh
    - "-c"
-->
```

### Read File

```md
<!-- docfresh begin
file:
  path: foo.md
-->
```
