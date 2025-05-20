# Go TODOアプリケーション

シンプルなTODOアプリケーションです。Go言語で実装されており、PostgreSQLを使用してデータを永続化します。

## 機能

- TODOの追加
- TODOの一覧表示
- TODOの完了/未完了の切り替え
- TODOの削除

## 必要条件

- Go 1.16以上
- PostgreSQL 12以上

## セットアップ

1. PostgreSQLのインストールと設定:
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo service postgresql start
```

2. データベースの作成:
```bash
sudo -u postgres psql
postgres=# CREATE DATABASE todo;
postgres=# ALTER USER postgres WITH PASSWORD 'postgres';
postgres=# \q
```

3. 依存関係のインストール:
```bash
go mod download
```

4. アプリケーションの起動:
```bash
go run .
```

## 使用方法

ブラウザで http://localhost:8080 にアクセスして、TODOアプリケーションを使用できます。
