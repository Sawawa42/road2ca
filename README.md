# Road to Cyberagent

42Tokyo「Road to」CyberAgentカリキュラムリポジトリ

[README.mdの原本](README_original.md)

## 実行方法

**MySQL, Redis, SwaggerUIサーバの起動**

```
make up
```

**MySQL, Redis, SwaggerUIサーバの停止**

```
make down
```

**APIの起動**

```
make
```

**Dockerボリュームの削除**

```
make rmvolumes
```

## ER図

```mermaid
erDiagram
    USERS {
        BINARY(16) id PK "User ID (UUID)"
        VARCHAR(128) name
        INT highscore
        INT coin
        VARCHAR(128) token "UNIQUE"
        DATETIME createdAt
        DATETIME updatedAt
    }

    ITEMS {
        INT id PK "Item ID"
        VARCHAR(128) name
        TINYINT rarity
        INT weight
    }

    COLLECTIONS {
        BINARY(16) id PK "Collection ID (UUID)"
        BINARY(16) userId FK "User ID"
        INT itemId FK "Item ID"
    }

    SETTINGS {
        INT id PK "Setting ID"
        VARCHAR(128) name
        INT gachaCoinConsumption
        INT drawGachaMaxTimes
        INT getRankingLimit
        INT rewardCoin
        FLOAT rarity3Ratio
        FLOAT rarity2Ratio
        FLOAT rarity1Ratio
        DATETIME createdAt
        DATETIME updatedAt
    }

    USERS ||--o{ COLLECTIONS : "has"
    ITEMS ||--o{ COLLECTIONS : "is part of"

```
