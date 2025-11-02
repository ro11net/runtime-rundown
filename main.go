package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	SecurityLevel string
	Message       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
}

type ConnectionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

var config Config

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Flux Training App</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container {
            background: white;
            border-radius: 10px;
            padding: 30px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
        }
        h1 {
            color: #333;
            border-bottom: 3px solid #667eea;
            padding-bottom: 10px;
        }
        .info-section {
            background: #f7f9fc;
            border-left: 4px solid #667eea;
            padding: 15px;
            margin: 20px 0;
            border-radius: 5px;
        }
        .info-item {
            margin: 10px 0;
            padding: 10px;
            background: white;
            border-radius: 5px;
        }
        .label {
            font-weight: bold;
            color: #555;
            display: inline-block;
            width: 150px;
        }
        .value {
            color: #333;
            font-family: 'Courier New', monospace;
            background: #e8f4f8;
            padding: 2px 8px;
            border-radius: 3px;
        }
        .security-high {
            background: #fee;
            border-left-color: #f44336;
        }
        .security-medium {
            background: #fff3cd;
            border-left-color: #ff9800;
        }
        .security-low {
            background: #e8f5e9;
            border-left-color: #4caf50;
        }
        button {
            background: #667eea;
            color: white;
            border: none;
            padding: 12px 30px;
            font-size: 16px;
            border-radius: 5px;
            cursor: pointer;
            transition: all 0.3s;
            margin: 20px 0;
        }
        button:hover {
            background: #5a67d8;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        button:active {
            transform: translateY(0);
        }
        #result {
            margin-top: 20px;
            padding: 15px;
            border-radius: 5px;
            display: none;
        }
        .success {
            background: #d4edda;
            border: 1px solid #c3e6cb;
            color: #155724;
        }
        .error {
            background: #f8d7da;
            border: 1px solid #f5c6cb;
            color: #721c24;
        }
        .timestamp {
            font-size: 12px;
            color: #666;
            margin-top: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Flux Training Application</h1>
        
        <div class="info-section security-{{.SecurityClass}}">
            <h2>Environment Configuration</h2>
            <div class="info-item">
                <span class="label">Security Level:</span>
                <span class="value">{{.SecurityLevel}}</span>
            </div>
            <div class="info-item">
                <span class="label">Message:</span>
                <span class="value">{{.Message}}</span>
            </div>
        </div>

        <div class="info-section">
            <h2>Database Configuration</h2>
            <div class="info-item">
                <span class="label">Host:</span>
                <span class="value">{{.DBHost}}</span>
            </div>
            <div class="info-item">
                <span class="label">Port:</span>
                <span class="value">{{.DBPort}}</span>
            </div>
            <div class="info-item">
                <span class="label">Database:</span>
                <span class="value">{{.DBName}}</span>
            </div>
            <div class="info-item">
                <span class="label">User:</span>
                <span class="value">{{.DBUser}}</span>
            </div>
        </div>

        <button onclick="testConnection()">🔌 Test Postgres Connection</button>
        
        <div id="result"></div>
    </div>

    <script>
        async function testConnection() {
            const resultDiv = document.getElementById('result');
            resultDiv.style.display = 'block';
            resultDiv.className = '';
            resultDiv.innerHTML = 'Testing connection...';
            
            try {
                const response = await fetch('/test-connection');
                const data = await response.json();
                
                if (data.success) {
                    resultDiv.className = 'success';
                    resultDiv.innerHTML = '✅ ' + data.message + '<div class="timestamp">Tested at: ' + data.time + '</div>';
                } else {
                    resultDiv.className = 'error';
                    resultDiv.innerHTML = '❌ ' + data.message + '<div class="timestamp">Tested at: ' + data.time + '</div>';
                }
            } catch (error) {
                resultDiv.className = 'error';
                resultDiv.innerHTML = '❌ Error testing connection: ' + error.message;
            }
        }
    </script>
</body>
</html>
`

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func loadConfig() {
	config = Config{
		SecurityLevel: getEnv("SECURITY_LEVEL", "LOW"),
		Message:       getEnv("MESSAGE", "Welcome to Flux Training!"),
		DBHost:        getEnv("DB_HOST", "postgres"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "postgres"),
	}
}

func getSecurityClass(level string) string {
	switch level {
	case "HIGH":
		return "high"
	case "MEDIUM":
		return "medium"
	default:
		return "low"
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("home").Parse(htmlTemplate))
	
	data := struct {
		Config
		SecurityClass string
	}{
		Config:        config,
		SecurityClass: getSecurityClass(config.SecurityLevel),
	}
	
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func testConnectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		json.NewEncoder(w).Encode(ConnectionResult{
			Success: false,
			Message: fmt.Sprintf("Failed to open connection: %v", err),
			Time:    time.Now().Format(time.RFC3339),
		})
		return
	}
	defer db.Close()
	
	// Set connection timeout
	db.SetConnMaxLifetime(time.Second * 5)
	
	err = db.Ping()
	if err != nil {
		json.NewEncoder(w).Encode(ConnectionResult{
			Success: false,
			Message: fmt.Sprintf("Failed to ping database: %v", err),
			Time:    time.Now().Format(time.RFC3339),
		})
		return
	}
	
	json.NewEncoder(w).Encode(ConnectionResult{
		Success: true,
		Message: "Successfully connected to PostgreSQL database!",
		Time:    time.Now().Format(time.RFC3339),
	})
}

func main() {
	loadConfig()
	
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/test-connection", testConnectionHandler)
	
	log.Printf("Starting Flux Training App on :8080")
	log.Printf("Security Level: %s", config.SecurityLevel)
	log.Printf("Message: %s", config.Message)
	log.Printf("Database Host: %s:%s", config.DBHost, config.DBPort)
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
