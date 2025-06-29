# 🌐 URL Shortener SaaS (Bitly-like) — Microservices Architecture

A scalable and production-ready **SaaS URL Shortener** platform inspired by Bitly. Built with **Go**, powered by **Redis**, **ClickHouse**, and **Kafka**, designed for **high throughput**, **analytics**, and **multi-tenant SaaS billing**.

---

## 🚀 Key Features

- 🔗 **Shorten Long URLs** with base62 or random aliases
- ⚡ **High-Performance Redirects** via Redis caching
- 📊 **Real-time Click Analytics** (device, geo, browser, time)
- 🌍 **GeoIP-based insights** (country, city)
- 🧑‍💼 **Multi-Tenant Support** (each user/org has their own namespace)
- 🧾 **Billing Integration** with Stripe (subscription plans)
- 📈 **Analytics Storage** in ClickHouse for scalable time-series insights
- 🔐 **Secure Auth** with JWT + OAuth2
- 🧱 **Built with Microservices** for maximum scalability
- 🛠️ Fully containerized with **Docker Compose** and ready for Kubernetes

---

## 🧩 Microservices Overview

| Service         | Description                               |
|----------------|-------------------------------------------|
| `api-gateway`   | Routes requests, rate limiting, JWT check |
| `auth-service`  | Login, register, refresh token, OAuth     |
| `user-service`  | Manages users, tenants, and plans         |
| `shortener-service` | Handles URL shortening logic          |
| `redirect-service`  | Redirects and logs click events       |
| `analytics-service` | Processes click events from Kafka     |
| `billing-service`   | Integrates Stripe for SaaS billing    |

---

## 🧱 Tech Stack

| Layer            | Technology                          |
|------------------|--------------------------------------|
| Language         | Golang                              |
| Caching          | Redis (Sentinel for HA)              |
| Database         | PostgreSQL                          |
| Analytics DB     | ClickHouse / TimescaleDB            |
| Message Queue    | Kafka (KRaft mode - no Zookeeper)    |
| Auth             | JWT, OAuth2                         |
| GeoIP            | MaxMind + `oschwald/geoip2`         |
| Billing          | Stripe (Webhooks & Checkout)        |
| Deployment       | Docker Compose → Kubernetes         |

---

## 📦 Infrastructure (via Docker Compose)

- PostgreSQL for persistent data
- Redis Master/Replica with Sentinel for HA
- ClickHouse for analytics logs
- Kafka in KRaft mode for event streaming

---

## 📊 Analytics

Click events are stored in ClickHouse and include:
- Timestamp
- IP address (resolved to geo)
- Device, OS, Browser
- Short URL code and tenant ID

This allows users to view:
- Click volume over time
- Geo distribution heatmaps
- Device/browser usage

---

## 🔐 SaaS Features

- Rate limiting per plan
- Monthly quota (clicks / URLs / API calls)
- Stripe-powered subscriptions
- Custom domain mapping per tenant
- Admin Dashboard for usage/management

---

## 🛣 Roadmap

| Phase | Features                                                      |
|-------|---------------------------------------------------------------|
| 1     | Core Shortener & Redirect Flow                                |
| 2     | Click Logging + GeoIP + Analytics                             |
| 3     | Auth + Multi-tenant Setup                                     |
| 4     | Stripe Billing & Subscription Plans                           |
| 5     | Frontend Dashboard + Management Tools                         |
| 6     | Deployment to Cloud / K8s with Auto Scaling                   |

---


---

## 👨‍💻 Contribution

Work in progress — initial architecture, service scaffolding and core flows are being implemented. Contributions and suggestions are welcome!

---

## 🗓 Development Status

> 📅 Initial setup in progress.  
> ✅ Redis Sentinel, Kafka KRaft, ClickHouse setup complete.  
> 🔜 Next: Implementing Shortener and Redirect services.

---

## 📄 License

MIT License.

