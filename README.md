# YouTube Video Fetcher


A full-stack application that fetches and stores YouTube videos based on search queries, providing a clean dashboard interface for viewing and filtering videos. The system includes an automated background service for continuous video fetching and a scalable API infrastructure.

## 🌟 Features

### Backend
- ⚡ Asynchronous YouTube API integration with 10-second refresh interval
- 🔄 Automatic API key rotation system
- 📊 PostgreSQL database with optimized indexing
- 🐳 Docker containerization
- ☁️ AWS deployment with RDS
- 📝 RESTful API with pagination support

### Frontend
- 📱 Responsive React dashboard
- 🎯 Advanced filtering capabilities
- ⏱️ Real-time sorting options
- 🔑 Dynamic API key management
- 🚀 Built with Vite for optimal performance
- Find the repository here https://github.com/RajendraKumarVesapogu/fampay-assignment-frontend.git

## 🛠️ Technology Stack

- **Backend**: Go
- **Frontend**: React + Vite
- **Database**: PostgreSQL
- **Deployment**: AWS, Docker
- **API**: YouTube Data API v3

## 🚀 Quick Start

### Backend Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/RajendraKumarVesapogu/fampay-assignment.git
   cd fampay-assignment
   ```

2. **Configure Environment Variables**
   Create a `.env` file in the root directory:
   ```env
   DB_HOST=your_database_host
   DB_PORT=your_database_port
   DB_USER=your_database_user
   DB_PASSWORD=your_database_password
   DB_NAME=your_database_name
   GIN_MODE="release"
   ALLOWED_ORIGINS="*"
   REDIS_URI="redis://localhost:6379"
   PORT=3000
   YOUTUBE_API_KEY1=your_youtube_api_key1
   YOUTUBE_API_KEY2=your_youtube_api_key2
   YOUTUBE_API_KEY3=your_youtube_api_key3
   ```

3. **Database Setup**
   - Execute the schema from `videos_schema.sql`
   - Ensure the table name is set to `videos`

4. **Install Dependencies and Run**
   ```bash
   go mod tidy
   go run main.go
   ```

### Frontend Setup

1. **Clone the Frontend Repository**
   ```bash
   git clone https://github.com/RajendraKumarVesapogu/fampay-assignment-frontend.git
   cd fampay-assignment-frontend
   ```

2. **Install Dependencies**
   ```bash
   npm install
   ```

3. **Start Development Server**
   ```bash
   npm run dev
   ```

   Access the application at `http://localhost:5173`


## 📡 API Documentation

### Base URL
```
http://3.108.83.52:3000/
```

### Endpoints

#### 1. Get Videos
```http
GET /videos
```

**Query Parameters**

| Parameter         | Type   | Required | Description                               |
|------------------|--------|----------|-------------------------------------------|
| sort_order       | string | Yes      | Sort order (asc/desc)                     |
| pagination_size  | int    | Yes      | Items per page (max 10)                   |
| pagination_page  | int    | Yes      | Page number                               |
| published_after  | string | No       | Filter by date (RFC 3339 format)          |

**Example Requests:**

<details>
<summary>cURL</summary>

```bash
curl 'http://3.108.83.52:3000/videos?sort_order=desc&pagination_size=10&pagination_page=1&published_after=2024-11-14T17:59:00Z' \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json'
```
</details>

<details>
<summary>JavaScript Fetch</summary>

```javascript
const params = new URLSearchParams({
  sort_order: 'desc',
  pagination_size: '10',
  pagination_page: '1',
  published_after: '2024-11-14T17:59:00Z'
});

fetch(`http://3.108.83.52:3000/videos?${params}`, {
  method: 'GET',
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json'
  }
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));
```
</details>

<details>
<summary>Postman</summary>

```plaintext
GET http://3.108.83.52:3000/videos?sort_order=desc&pagination_size=10&pagination_page=1&published_after=2024-11-14T17:59:00Z

Headers:
Accept: application/json
Content-Type: application/json
```
</details>

**Sample Response:**
```json
{
  "error": false,
  "response": {
    "videos": [
      {
        "VideoID": "1l_w5g7fbjA",
        "Title": "Sample Video Title",
        "Description": "Video description here...",
        "PublishedAt": "2024-11-02T09:13:49Z",
        "ThumbnailURL": "https://i.ytimg.com/vi/1l_w5g7fbjA/default.jpg",
        "ChannelTitle": "Channel Name",
        "ChannelID": "UC-crZTQNRzZgzyighTKF0nQ"
      }
    ]
  }
}
```

#### 2. Add API Key
```http
POST /videos/key
```

**Request Body**
```json
{
  "api_key": "your_youtube_api_key"
}
```

**Example Requests:**

<details>
<summary>cURL</summary>

```bash
curl 'http://3.108.83.52:3000/videos/key' \
  -X POST \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json' \
  --data-raw '{"api_key":"your_youtube_api_key"}'
```
</details>

<details>
<summary>JavaScript Fetch</summary>

```javascript
fetch('http://3.108.83.52:3000/videos/key', {
  method: 'POST',
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    api_key: 'your_youtube_api_key'
  })
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));
```
</details>

<details>
<summary>Postman</summary>

```plaintext
POST http://3.108.83.52:3000/videos/key

Headers:
Accept: application/json
Content-Type: application/json

Body (raw JSON):
{
  "api_key": "your_youtube_api_key"
}
```
</details>

**Sample Response:**
```json
{
  "error": false,
  "response": {
    "success": true
  }
}
```

### Testing with HTTPie
If you prefer using HTTPie, here are the equivalent commands:

```bash
# Get Videos
http GET 'http://3.108.83.52:3000/videos?sort_order=desc&pagination_size=10&pagination_page=1&published_after=2024-11-14T17:59:00Z'

# Add API Key
http POST http://3.108.83.52:3000/videos/key api_key=your_youtube_api_key
```













## 🌐 Live Deployments

- **Frontend Dashboard**: [https://main.d3s4vg61cjppvr.amplifyapp.com/](https://main.d3s4vg61cjppvr.amplifyapp.com/)
- **Backend API**: [http://3.108.83.52:3000/](http://3.108.83.52:3000/)

## ⚠️ Important Notes

- Enable insecure content in browser settings for the frontend dashboard
- Adjust the date filter if no videos are visible initially
- API keys are automatically rotated when quotas are exhausted. Add a new key to start fetching latest videos immediately,


## 👨‍💻 Author

Rajendra Kumar Vesapogu

## 🛠️ Tools refs used

- YouTube Data API
- AWS Services
- Stack Overflow, Go, React and Vite communities
- Claude and GPT
- Youtube and Official docs