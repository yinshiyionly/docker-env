package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

const (
    weatherAPIURL = "https://api.weatherapi.com/v1/current.json"
)

func main() {
    // Create MCP server
    s := server.NewMCPServer(
        "MCP-demo",
        "1.0.0",
    )

    // Add tool
    // tool := mcp.NewTool("hello_world",
    //     mcp.WithDescription("Say hello to someone"),
    //     mcp.WithString("name",
    //         mcp.Required(),
    //         mcp.Description("Name of the person to greet"),
    //     ),
    // )

    // // Add tool handler
    // s.AddTool(tool, helloHandler)

    // Add weather tool
    weatherTool := mcp.NewTool("current_weather",
        mcp.WithDescription("The service retrieves the current weather for a specified city and returns only HTML and CSS code designed \n\t\t\tfor displaying weather information directly on a canvas. No additional comments should be included after the response. \n\t\t\tIt processes the city's name, gathers weather data such as temperature, humidity, wind speed, and general weather conditions, \n\t\t\tand then generates a structured HTML layout along with appropriate CSS styles. This HTML must be rendered visually on the canvas. \n\t\t\tAfter receiving the response, the text inside the HTML should be translated into the language of the request."),
        mcp.WithString("city",
            mcp.Required(),
            mcp.Description("The name of the city. This field is required and must be provided in English. \n\t\t\t\tOnly one city is allowed, and it must be the last one provided by the user."),
        ),
    )

    // Add weather tool handler
    s.AddTool(weatherTool, currentWeatherHandler)

    // Start the stdio server
    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    name, ok := request.Params.Arguments["name"].(string)
    if !ok {
        return nil, errors.New("name must be a string")
    }

    return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}

func currentWeatherHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Get city parameter
    city, ok := request.Params.Arguments["city"].(string)
    if !ok {
        return nil, errors.New("city must be a string")
    }

    // Get API key from environment variable
    apiKey := os.Getenv("WEATHER_API_KEY")
    if apiKey == "" {
        return nil, errors.New("WEATHER_API_KEY environment variable not set")
    }

    // Build API URL
    url := fmt.Sprintf("%s?q=%s&lang=zh&key=%s", weatherAPIURL, city, apiKey)

    // Make HTTP request
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch weather data: %v", err)
    }
    defer resp.Body.Close()

    // Check response status
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
    }

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    // Parse JSON
    var weatherData map[string]interface{}
    if err := json.Unmarshal(body, &weatherData); err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %v", err)
    }

    // Extract weather information
    location, hasLocation := weatherData["location"].(map[string]interface{})
    current, hasCurrent := weatherData["current"].(map[string]interface{})

    if !hasLocation || !hasCurrent {
        return nil, errors.New("invalid response format from weather API")
    }

    // Format HTML response
    html := generateWeatherHTML(location, current)
    
    return mcp.NewToolResultText(html), nil
}

// Helper function to generate HTML for weather display
func generateWeatherHTML(location, current map[string]interface{}) string {
    // Extract location data
    name := getStringValue(location, "name")
    country := getStringValue(location, "country")
    
    // Extract current weather data
    tempC := getFloatValue(current, "temp_c")
    tempF := getFloatValue(current, "temp_f")
    humidity := getIntValue(current, "humidity")
    windKph := getFloatValue(current, "wind_kph")
    windDir := getStringValue(current, "wind_dir")
    
    // Extract condition data if available
    var condition string
    var iconURL string
    if conditionData, ok := current["condition"].(map[string]interface{}); ok {
        condition = getStringValue(conditionData, "text")
        iconURL = getStringValue(conditionData, "icon")
        if iconURL != "" && iconURL[:2] == "//" {
            iconURL = "https:" + iconURL
        }
    }
    
    // Generate HTML
    html := `
    <!DOCTYPE html>
    <html>
    <head>
        <style>
            .weather-container {
                font-family: Arial, sans-serif;
                max-width: 500px;
                margin: 0 auto;
                padding: 20px;
                border-radius: 10px;
                box-shadow: 0 0 10px rgba(0,0,0,0.1);
                background: linear-gradient(to bottom right, #3498db, #2980b9);
                color: white;
            }
            .weather-header {
                display: flex;
                justify-content: space-between;
                align-items: center;
                margin-bottom: 20px;
            }
            .location {
                font-size: 24px;
                font-weight: bold;
            }
            .temperature {
                font-size: 48px;
                font-weight: bold;
                margin: 10px 0;
            }
            .condition {
                display: flex;
                align-items: center;
                margin: 10px 0;
                font-size: 18px;
            }
            .condition img {
                margin-right: 10px;
            }
            .details {
                display: grid;
                grid-template-columns: 1fr 1fr;
                gap: 10px;
                margin-top: 20px;
                background-color: rgba(255,255,255,0.1);
                padding: 15px;
                border-radius: 8px;
            }
            .detail-item {
                display: flex;
                flex-direction: column;
            }
            .detail-label {
                font-size: 14px;
                opacity: 0.8;
            }
            .detail-value {
                font-size: 18px;
                font-weight: bold;
            }
        </style>
    </head>
    <body>
        <div class="weather-container">
            <div class="weather-header">
                <div class="location">` + name + `, ` + country + `</div>
            </div>
            
            <div class="temperature">` + fmt.Sprintf("%.1f°C", tempC) + `</div>
            
            <div class="condition">
                ` + func() string {
                    if iconURL != "" {
                        return `<img src="` + iconURL + `" alt="Weather icon" width="50">`
                    }
                    return ""
                }() + `
                <span>` + condition + `</span>
            </div>
            
            <div class="details">
                <div class="detail-item">
                    <span class="detail-label">温度 (F)</span>
                    <span class="detail-value">` + fmt.Sprintf("%.1f°F", tempF) + `</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">湿度</span>
                    <span class="detail-value">` + fmt.Sprintf("%d%%", humidity) + `</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">风速</span>
                    <span class="detail-value">` + fmt.Sprintf("%.1f km/h", windKph) + `</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">风向</span>
                    <span class="detail-value">` + windDir + `</span>
                </div>
            </div>
        </div>
    </body>
    </html>
    `
    
    return html
}

// Helper functions to safely extract values from maps
func getStringValue(data map[string]interface{}, key string) string {
    if value, ok := data[key].(string); ok {
        return value
    }
    return ""
}

func getFloatValue(data map[string]interface{}, key string) float64 {
    switch value := data[key].(type) {
    case float64:
        return value
    case int:
        return float64(value)
    case float32:
        return float64(value)
    default:
        return 0
    }
}

func getIntValue(data map[string]interface{}, key string) int {
    switch value := data[key].(type) {
    case int:
        return value
    case float64:
        return int(value)
    default:
        return 0
    }
}
