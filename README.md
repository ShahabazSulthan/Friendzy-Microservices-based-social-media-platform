# Friendzy - Microservice Based Social Media Platform

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Architecture](#architecture)
- [Microservices Overview](#microservices-overview)
  - [1. Auth Service](#1-auth-service)
  - [2. Post and Relation Service](#2-post-and-relation-service)
  - [3. Chat Service](#3-chat-service)
  - [4. Notification Service](#4-notification-service)
  - [5. API Gateway](#5-api-gateway)
- [Tech Stack](#tech-stack)
- [Algorithms Used](#algorithms-used)
- [Payment Integration](#payment-integration)
- [Database Schema](#database-schema)
- [Installation](#installation)
- [Usage](#usage)
- [API Routes](#api-routes)
- [API Testing](#api-testing)
- [Contributing](#contributing)
- [License](#license)

---

## Introduction

**Friendzy** is a scalable, microservice-based social media platform built using the **Go Fiber framework**. The platform is designed to support essential social media features like user authentication, posting, chatting, notifications, and more. It leverages modern technologies like **Kafka** for messaging, **WebSockets** for real-time chat, and **Redis** for caching. The platform also includes a unique monetization feature where users can pay to get a verified blue tick using **Razorpay**.

## Features

- **User Authentication** (Sign up, Login, OTP-based verification)
- **Posts and Feeds** (Create, Like, Comment, and Share posts)
- **Real-Time Chat** (One-to-one and Group Chat with WebSocket support)
- **Push Notifications** (For likes, comments, follows, and more)
- **Feed Algorithms**:
  - Popularity-Based Ranking
  - Social Graph & Recency Algorithm
  - Random Content Feed
- **Blue Tick Verification** (Users can get verified by paying a fee via Razorpay)
- **Microservices Architecture** for scalability and maintainability

## Architecture

The platform follows a **microservices architecture** with the following services:

- **Auth Service**: Handles user authentication and profile management.
- **Post and Relation Service**: Manages posts, likes, comments, and user relations (followers, following).
- **Chat Service**: Supports real-time messaging using WebSocket.
- **Notification Service**: Manages push notifications using Kafka.
- **API Gateway**: Acts as a single entry point to the microservices.

```
User -> API Gateway -> [Auth Service | Post Service | Chat Service | Notification Service]
                          |                  |              |                  |
                          -------------------- gRPC & Kafka --------------------
```

## Microservices Overview

### 1. Auth Service
- **Responsibilities**: User registration, login, OTP verification, JWT authentication, profile management, and payment-based blue tick verification.
- **Tech Stack**: Go, JWT, PostgreSQL, Razorpay.

### 2. Post and Relation Service
- **Responsibilities**: Create posts, like/unlike posts, comment on posts, manage followers/following, and implement feed algorithms.
- **Tech Stack**: Go, Kafka, Redis, PostgreSQL.

### 3. Chat Service
- **Responsibilities**: Supports one-to-one and group chats using WebSockets for real-time communication.
- **Tech Stack**: Go, WebSocket, gRPC, MongoDB.

### 4. Notification Service
- **Responsibilities**: Handles sending push notifications for user actions like likes, comments, and follows using Kafka.
- **Tech Stack**: Go Fiber, Kafka, gRPC, PostgreSQL.

### 5. API Gateway
- **Responsibilities**: Acts as a unified entry point for all microservices, handling request routing, authentication, and load balancing.
- **Tech Stack**: Go Fiber.

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Go Fiber
- **Message Queue**: Kafka
- **Real-Time Communication**: WebSocket
- **Cache**: Redis
- **Database**: PostgreSQL, MongoDB
- **Payment Gateway**: Razorpay
- **API Communication**: gRPC, REST

## Algorithms Used

### 1. Popularity-Based Ranking Algorithm
- Ranks posts based on likes, comments, and shares.
  
### 2. Social Graph and Recency Algorithm
- Prioritizes posts from users you follow and recent activities.

### 3. Random Content Feed Algorithm
- Shows random posts from users you do not follow to increase engagement.

## Payment Integration

- Users can get a **blue tick verification** badge by paying a fee of **₹1000**.
- Payment is processed using **Razorpay**.
- After successful payment, users receive a verified badge on their profile.

## Database Schema

### Example Tables:
- **Users Table**: Stores user information, profile details, and verification status.
- **Posts Table**: Stores post content, likes, comments, and timestamps.
- **Chats Table**: Stores chat messages and metadata.
- **Notifications Table**: Stores notification events and statuses.

## Installation

### Prerequisites
- [Go](https://golang.org/dl/) (v1.19+)
- [Docker](https://www.docker.com/)
- [Kafka](https://kafka.apache.org/quickstart)
- [Redis](https://redis.io/download)
- [PostgreSQL](https://www.postgresql.org/download)
- [Razorpay Account](https://razorpay.com/)

### Steps
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/ShahabazSulthan/Friendzy.git
   cd Friendzy
   ```

2. **Set up Environment Variables**:
   - Copy `.env.example` to `.env` and update the values.
   
3. **Build and Run Services**:
   ```bash
   docker-compose up --build
   ```

4. **Run Migrations**:
   ```bash
   go run ./cmd/migrate.go
   ```

5. **Access the Application**:
   - API Gateway: [http://localhost:8000](http://localhost:8000)

## API Routes

### Auth Service
- `POST /signup` - Register a new user
- `POST /verify` - OTP verification
- `POST /login` - User login
- `POST /forgotpassword` - Request password reset
- `PATCH /resetpassword` - Reset password
- `GET /profile` - Get user profile (Protected)

### Post Service
- `POST /post` - Create a new post
- `GET /post` - Get user posts
- `DELETE /post/:postid` - Delete a post
- `POST /post/like/:postid` - Like a post

### Chat Service
- `GET /chat/onetoonechats/:recipientid` - Get one-to-one chats
- `GET /chat/ws` - WebSocket connection for chat

### Notification Service
- `GET /notification` - Fetch user notifications

## API Testing

- The API is hosted at: [friendzy.shahabazsulthan.cloud](https://friendzy.shahabazsulthan.cloud/)
- Import the Postman collection from `docs/Friendzy.postman_collection.json`
- Use JWT tokens for protected routes

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a Pull Request.


