# ドメインの整理

- 問題: Problem
- コード: Code
- 結果: Result

```
url := ...
code := NewCodeFromFile(filename)
submitter := NewSeleniumSubmitter()
problem := NewProblem(url, submitter)
result := problem.Submit(code)
```
