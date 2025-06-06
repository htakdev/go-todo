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

## ローカル環境でのセットアップ

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

3. バックエンドの依存関係のインストール:
```bash
go mod download
```

4. フロントエンドの依存関係のインストール:
```bash
cd frontend
npm install
```

5. アプリケーションの起動:

バックエンドの起動:
```bash
go run .
```

フロントエンドの起動（別のターミナルで）:
```bash
cd frontend
npm run dev
```

6. アクセス方法

ブラウザで http://localhost:5173 にアクセスして、TODOアプリケーションを使用できます。

## 本番環境でのセットアップ

[ローカル環境でのセットアップ](#ローカル環境でのセットアップ)の内容をデプロイ先のサービスで置き換える必要があります。ただし「バックエンドの起動」や「フロントエンドの起動（別のターミナルで）」に関しては先にビルドすることをお勧めします。
