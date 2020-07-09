# acc

## Access AtCoder via command line

### Login

#### 1. GET

```shell script
curl -c cookie https://atcoder.jp/login | grep 'var csrfToken =' | sed 's/^.*var csrfToken = "\(.*\)"/\1/g' > csrfToken
```

- URL: `https://atcoder.jp/login`

| Option          | Description           |
| :-------------- | :-------------------- |
| `-c <filename>` | save Cookie to a file |

#### 2. POST

```shell script
curl -X POST -b cookie https://atcoder.jp/login -F "csrf_token=$(cat csrfToken)" -F "username=<username>" -F "password=<password>"
```

- Form Parameters
  - username
  - password
  - csrf_token

| Option            | Description                                               |
| :---------------- | :-------------------------------------------------------- |
| `-i`              | include HTTP Header in the output                         |
| `-b <filename>`   | use Cookie saved with `-c` option in the previous request |
| `-F <name=value>` | construct form of "Content-Type multipart/form-data"      |

### Submission

#### 1. GET

```shell script
curl -b cookie https://atcoder.jp/contests/abc173/tasks/abc173_a | grep 'var csrfToken =' | sed 's/^.*var csrfToken = "\(.*\)"/\1/g' > csrfToken
```

- URL: `https://atcoder.js/contests/<contest_name>/tasks/<task_name>`

#### 2. POST

```shell script
curl -X POST -b cookie https://atcoder.jp/contests/abc173/submit -F "data.TaskScreenName=<data.TaskScreenName>" -F "data.LanguageId=<data.LanguageId>" -F "sourceCode=<sourceCode>" -F "csrf_token=<csrf_token>"
```

- Form Parameters
  - data.TaskScreenName
  - data.LanguageId
  - sourceCode
  - csrf_token
