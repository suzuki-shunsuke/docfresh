# template

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
