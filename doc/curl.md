# CURL を使って AtCoder にアクセスする

- ログイン
- 提出

## ログイン

### CSRF トークンを取得

```shell script
curl -c cookie https://atcoder.jp/login | grep 'var csrfToken =' | sed 's/^.*var csrfToken = "\(.*\)"/\1/g' > csrfToken
```

#### ログイン情報を POST

```shell script
curl -X POST -b cookie https://atcoder.jp/login -F "csrf_token=$(cat csrfToken)" -F "username=<username>" -F "password=<password>"
```

- フォームパラメータ
  - username
  - password
  - csrf_token

***

## 提出

### CSRF トークンを取得

```shell script
curl -b cookie https://atcoder.jp/contests/abc173/tasks/abc173_a | grep 'var csrfToken =' | sed 's/^.*var csrfToken = "\(.*\)"/\1/g' > csrfToken
```

### 解答を POST

```shell script
curl -X POST -b cookie https://atcoder.jp/contests/abc173/submit -F "data.TaskScreenName=<data.TaskScreenName>" -F "data.LanguageId=<data.LanguageId>" -F "sourceCode=<sourceCode>" -F "csrf_token=<csrf_token>"
```

- フォームパラメータ
  - data.TaskScreenName
  - data.LanguageId
  - sourceCode
  - csrf_token
