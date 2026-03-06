# Customize Template

<!-- docfresh begin
command:
  command: echo Hello
template:
  content: |
    ```console
    $ {{.Command}}
    ```
    
    ```
    {{trimSuffix "\n" .CombinedOutput}}
    ```
-->
```console
$ echo Hello
```

```
Hello
```
<!-- docfresh end -->

## template.path

<!-- docfresh begin
command:
  command: echo "read template file"
template:
  path: file/template2.md
-->
```console
$ echo "read template file"
read template file
```
<!-- docfresh end -->

## template variables

<!-- docfresh begin
file:
  path: file/template.md
  template:
    vars:
      name: foo
-->
Hello, foo
<!-- docfresh end -->
