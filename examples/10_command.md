# Embed Command Result

## Hello

<!-- docfresh begin
command:
  command: echo hello
-->
```sh
echo hello
```

```
hello
```
<!-- docfresh end -->

## command.dir

<!-- docfresh begin
command:
  command: cat foo.md
  dir: file
-->
```sh
cat foo.md
```

```
This is read from foo.md
```
<!-- docfresh end -->

## Change Shell

<!-- docfresh begin
command:
  command: console.log("hello")
  shell:
    - node
    - "-e"
template:
  content: |
    ```js
    {{.Command}}
    ```
    
    ```
    {{trimSuffix "\n" .CombinedOutput}}
    ```
-->
```js
console.log("hello")
```

```
hello
```
<!-- docfresh end -->

## ignore_fail

By default, `docfresh run` fails if any command fails.
If `.command.ignore_fail` is set to `true`, the command failure will be ignored.

<!-- docfresh begin
command:
  command: |
    echo "failed to install" >&2
    exit 1
  ignore_fail: true
-->
```sh
echo "failed to install" >&2
exit 1
```

```
failed to install
```
<!-- docfresh end -->

## pre_command, post_command

<!-- docfresh begin
pre_command:
  command: echo temporary > temporary.txt
command:
  command: cat temporary.txt
post_command:
  command: rm temporary.txt
-->
```sh
cat temporary.txt
```

```
temporary
```
<!-- docfresh end -->
