# codeuploader
## 準備
go get github.com/go-sql-driver/mysql

# 実行方法
## TEST
go run src/main/main.go -c config.json

-cオプションを付けると設定ファイルを明示的に指定できる
-cオプションを付けなかった場合はmain.goファイルと同階層のconfig.jsonを探す

## 本番
./main -c config.json

-cオプションを付けると設定ファイルを明示的に指定できる
-cオプションを付けなかった場合はmain.goファイルと同階層のconfig.jsonを探す

