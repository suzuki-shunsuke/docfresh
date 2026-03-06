{{if .Command -}}
```sh
{{trimSuffix "\n" .Command}}
```
{{- else if .EmbedScript -}}
```{{.ScriptLanguage}}
{{trimSuffix "\n" .Content}}
```
{{- else -}}
```sh
{{join " " .Shell}} {{trimSuffix "\n" .Script}}
```
{{- end}}

```
{{trimSuffix "\n" .CombinedOutput}}
```
