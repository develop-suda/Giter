# Giter

セキュアなJWT認証を実装したWEBアプリケーション

## 📱 Giter概要

- 🎯 GitHubのコミット履歴を分析し、モチベーション向上につながる応援メッセージを表示
- 📊 コミット継続日数を正確に追跡・可視化
- 📧 当日のコミット内容を効率的にメール通知
- 👥 他ユーザーのコミット履歴も簡単に閲覧可能
- 🔐 シームレスなGitHubアカウント連携によるログイン機能
- 通常のメールアドレスからのログイン機能も実装

## 💾 データベース設計

- 🔗 ユーザーとGitHubリポジトリの効率的な紐付け
- 📊 userとrepositoryの最適化された二つのテーブル構成

## 📨 メール配信システム

- ✉️ ユーザー向けカスタマイズされたコミットメール
- ⏰ 0時に自動実行
- 📝 前日のコミット情報を確実にメール配信
- ⚡ 効率的なバッチ処理

### メール設定機能

- ⚙️ マイページで柔軟なメール設定が可能
- ✅ DBで管理された送信設定に基づく確実な配信

### システム構成

- 📦 機能別に最適化されたGiterとメール配信バッチの2リポジトリ構成

### 仕様
- publicリポジトリのみ
- 一つのリポジトリに対して⚪︎個ブランチ分の情報を取得する
- 一個のブランチに対して最新の⚪コミット分の情報を取得する
- [リポジトリ一覧取得API](https://docs.github.com/ja/rest/repos/repos?apiVersion=2022-11-28#list-repositories-for-a-user)