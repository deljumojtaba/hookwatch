console.log("Initializing MongoDB...");

db = db.getSiblingDB("Eryaz"); // Replace with your database name

db.createCollection("configs"); // Replace with your collection name

db.configs.insertOne(
    {
        "merchant_id": 1,
        "name": "Eryaz",
        "config_version": "1.0",
        "db": {
            "product_db": {
                "type": "mssql",
                "version": "1.0",
                "host": "localhost",
                "port": "1433",
                "username": "sa",
                "password": "password123"
            },
            "order_db": {
                "type": "mssql",
                "version": "1.0",
                "host": "localhost",
                "port": "1433",
                "username": "sa",
                "password": "Password123"
            },
            "cart_db": {
                "type": "mssql",
                "version": "1.0",
                "host": "localhost",
                "port": "1433",
                "username": "sa",
                "password": "password123"
            }
        },
        "notifications": {
            "email": {
                "active": true,
                "smtp": {
                "host": "smtp.gmail.com",
                "port": "587",
                "username": "username",
                "password": "password"
                }
            }
        },
        "app": {
            "features": {
                "payment": {
                "active": true,
                "posList": [
                    "akbank",
                    "isbank"
                ]
                },
                "analytics": {
                "active": true,
                "type": "google"
                }
            }
        },
        "details": {
            "timezone": "Europe/Istanbul",
            "currency": "TRY",
            "language": "tr",
            "country": "TR"
        },
        "config": {
            "db": {
                "cart_db": {
                "password": "123"
                }
        }
        },
        created_at: new Date()
    });