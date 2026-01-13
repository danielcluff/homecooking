package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/homecooking/backend/internal/config"
)

type AIService struct {
	cfg *config.Config
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		cfg: cfg,
	}
}

func (s *AIService) IsEnabled() bool {
	return s.cfg.AI.Enabled
}

func (s *AIService) GetConfig() *config.AIConfig {
	return &s.cfg.AI
}

type ExtractRecipeRequest struct {
	ImagePath string `json:"image_path"`
	ImageData []byte `json:"image_data,omitempty"`
	ImageType string `json:"image_type,omitempty"`
}

type ExtractRecipeResponse struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	MarkdownContent string   `json:"markdown_content"`
	PrepTimeMinutes *int     `json:"prep_time_minutes"`
	CookTimeMinutes *int     `json:"cook_time_minutes"`
	Servings        *int     `json:"servings"`
	Difficulty      string   `json:"difficulty"`
	Tags            []string `json:"tags"`
}

type EnhanceRecipeRequest struct {
	CurrentTitle       string `json:"current_title"`
	CurrentContent     string `json:"current_content"`
	EnhanceDescription bool   `json:"enhance_description"`
	ImproveStructure   bool   `json:"improve_structure"`
	AddTips            bool   `json:"add_tips"`
}

type EnhanceRecipeResponse struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	MarkdownContent string   `json:"markdown_content"`
	Tips            []string `json:"tips,omitempty"`
}

func (s *AIService) ExtractRecipe(req *ExtractRecipeRequest) (*ExtractRecipeResponse, error) {
	if !s.cfg.AI.Enabled {
		return nil, errors.New("AI features are not enabled")
	}

	switch s.cfg.AI.Provider {
	case "openai":
		return s.extractRecipeOpenAI(req)
	case "anthropic":
		return s.extractRecipeAnthropic(req)
	case "ollama", "local":
		return s.extractRecipeLocal(req)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", s.cfg.AI.Provider)
	}
}

func (s *AIService) EnhanceRecipe(req *EnhanceRecipeRequest) (*EnhanceRecipeResponse, error) {
	if !s.cfg.AI.Enabled {
		return nil, errors.New("AI features are not enabled")
	}

	switch s.cfg.AI.Provider {
	case "openai":
		return s.enhanceRecipeOpenAI(req)
	case "anthropic":
		return s.enhanceRecipeAnthropic(req)
	case "ollama", "local":
		return s.enhanceRecipeLocal(req)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", s.cfg.AI.Provider)
	}
}

func (s *AIService) extractRecipeOpenAI(req *ExtractRecipeRequest) (*ExtractRecipeResponse, error) {
	var imageData string
	if len(req.ImageData) > 0 {
		imageData = base64.StdEncoding.EncodeToString(req.ImageData)
		imageData = fmt.Sprintf("data:%s;base64,%s", req.ImageType, imageData)
	}

	prompt := `Extract the recipe from this image. Return a JSON object with the following fields:
- title: Recipe title (string)
- description: Brief description (string, optional)
- markdown_content: Full recipe content in markdown format with ## Ingredients and ## Instructions sections (string)
- prep_time_minutes: Preparation time in minutes (integer, optional)
- cook_time_minutes: Cooking time in minutes (integer, optional)
- servings: Number of servings (integer, optional)
- difficulty: Difficulty level (string, optional, one of: easy, medium, hard)
- tags: Array of relevant tags (array of strings, optional)

Only return the JSON object, nothing else.`

	messages := []map[string]interface{}{
		{
			"role": "user",
			"content": []interface{}{
				map[string]string{
					"type": "text",
					"text": prompt,
				},
			},
		},
	}

	if imageData != "" {
		messages[0]["content"] = []interface{}{
			map[string]string{
				"type": "text",
				"text": prompt,
			},
			map[string]interface{}{
				"type": "image_url",
				"image_url": map[string]string{
					"url": imageData,
				},
			},
		}
	}

	requestBody := map[string]interface{}{
		"model":    s.cfg.AI.Model,
		"messages": messages,
	}

	baseURL := "https://api.openai.com/v1/chat/completions"
	if s.cfg.AI.BaseURL != "" {
		baseURL = s.cfg.AI.BaseURL
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.cfg.AI.APIKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error: %s", string(respBody))
	}

	var openaiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBody, &openaiResponse); err != nil {
		return nil, err
	}

	if len(openaiResponse.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	content := openaiResponse.Choices[0].Message.Content

	var result ExtractRecipeResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	return &result, nil
}

func (s *AIService) extractRecipeAnthropic(req *ExtractRecipeRequest) (*ExtractRecipeResponse, error) {
	var imageData string
	if len(req.ImageData) > 0 {
		imageData = base64.StdEncoding.EncodeToString(req.ImageData)
	}

	prompt := `Extract the recipe from this image. Return a JSON object with the following fields:
- title: Recipe title (string)
- description: Brief description (string, optional)
- markdown_content: Full recipe content in markdown format with ## Ingredients and ## Instructions sections (string)
- prep_time_minutes: Preparation time in minutes (integer, optional)
- cook_time_minutes: Cooking time in minutes (integer, optional)
- servings: Number of servings (integer, optional)
- difficulty: Difficulty level (string, optional, one of: easy, medium, hard)
- tags: Array of relevant tags (array of strings, optional)

Only return the JSON object, nothing else.`

	messages := []map[string]interface{}{
		{
			"role": "user",
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": prompt,
				},
			},
		},
	}

	if imageData != "" {
		messages[0]["content"] = []map[string]interface{}{
			{
				"type": "text",
				"text": prompt,
			},
			{
				"type": "image",
				"source": map[string]string{
					"type":       "base64",
					"media_type": req.ImageType,
					"data":       imageData,
				},
			},
		}
	}

	requestBody := map[string]interface{}{
		"model":      s.cfg.AI.Model,
		"max_tokens": 4096,
		"messages":   messages,
	}

	baseURL := "https://api.anthropic.com/v1/messages"
	if s.cfg.AI.BaseURL != "" {
		baseURL = s.cfg.AI.BaseURL
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", s.cfg.AI.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Anthropic API error: %s", string(respBody))
	}

	var anthropicResponse struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(respBody, &anthropicResponse); err != nil {
		return nil, err
	}

	if len(anthropicResponse.Content) == 0 {
		return nil, errors.New("no response from Anthropic")
	}

	content := anthropicResponse.Content[0].Text

	var result ExtractRecipeResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic response: %w", err)
	}

	return &result, nil
}

func (s *AIService) extractRecipeLocal(req *ExtractRecipeRequest) (*ExtractRecipeResponse, error) {
	baseURL := "http://localhost:11434/api/chat"
	if s.cfg.AI.BaseURL != "" {
		baseURL = s.cfg.AI.BaseURL
	}

	prompt := `Extract the recipe from this image. Return a JSON object with the following fields:
- title: Recipe title (string)
- description: Brief description (string, optional)
- markdown_content: Full recipe content in markdown format with ## Ingredients and ## Instructions sections (string)
- prep_time_minutes: Preparation time in minutes (integer, optional)
- cook_time_minutes: Cooking time in minutes (integer, optional)
- servings: Number of servings (integer, optional)
- difficulty: Difficulty level (string, optional, one of: easy, medium, hard)
- tags: Array of relevant tags (array of strings, optional)

Only return the JSON object, nothing else.`

	requestBody := map[string]interface{}{
		"model": s.cfg.AI.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream": false,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API error: %s", string(respBody))
	}

	var ollamaResponse struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.Unmarshal(respBody, &ollamaResponse); err != nil {
		return nil, err
	}

	var result ExtractRecipeResponse
	if err := json.Unmarshal([]byte(ollamaResponse.Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Ollama response: %w", err)
	}

	return &result, nil
}

func (s *AIService) enhanceRecipeOpenAI(req *EnhanceRecipeRequest) (*EnhanceRecipeResponse, error) {
	prompt := `Improve this recipe. Return a JSON object with the following fields:
- title: Improved or original title (string)
- description: Enhanced description that makes the recipe more appealing (string)
- markdown_content: Improved recipe content with better formatting and clarity (string)
- tips: Array of cooking tips or variations (array of strings, optional)

Recipe to enhance:
Title: ` + req.CurrentTitle + `
Content:
` + req.CurrentContent + `

Only return the JSON object, nothing else.`

	requestBody := map[string]interface{}{
		"model": s.cfg.AI.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	baseURL := "https://api.openai.com/v1/chat/completions"
	if s.cfg.AI.BaseURL != "" {
		baseURL = s.cfg.AI.BaseURL
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.cfg.AI.APIKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error: %s", string(respBody))
	}

	var openaiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBody, &openaiResponse); err != nil {
		return nil, err
	}

	if len(openaiResponse.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	content := openaiResponse.Choices[0].Message.Content

	var result EnhanceRecipeResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	return &result, nil
}

func (s *AIService) enhanceRecipeAnthropic(req *EnhanceRecipeRequest) (*EnhanceRecipeResponse, error) {
	prompt := `Improve this recipe. Return a JSON object with the following fields:
- title: Improved or original title (string)
- description: Enhanced description that makes the recipe more appealing (string)
- markdown_content: Improved recipe content with better formatting and clarity (string)
- tips: Array of cooking tips or variations (array of strings, optional)

Recipe to enhance:
Title: ` + req.CurrentTitle + `
Content:
` + req.CurrentContent + `

Only return the JSON object, nothing else.`

	requestBody := map[string]interface{}{
		"model":      s.cfg.AI.Model,
		"max_tokens": 4096,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]string{
					{
						"type": "text",
						"text": prompt,
					},
				},
			},
		},
	}

	baseURL := "https://api.anthropic.com/v1/messages"
	if s.cfg.AI.BaseURL != "" {
		baseURL = s.cfg.AI.BaseURL
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", s.cfg.AI.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Anthropic API error: %s", string(respBody))
	}

	var anthropicResponse struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(respBody, &anthropicResponse); err != nil {
		return nil, err
	}

	if len(anthropicResponse.Content) == 0 {
		return nil, errors.New("no response from Anthropic")
	}

	content := anthropicResponse.Content[0].Text

	var result EnhanceRecipeResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic response: %w", err)
	}

	return &result, nil
}

func (s *AIService) enhanceRecipeLocal(req *EnhanceRecipeRequest) (*EnhanceRecipeResponse, error) {
	baseURL := "http://localhost:11434/api/chat"
	if s.cfg.AI.BaseURL != "" {
		baseURL = s.cfg.AI.BaseURL
	}

	prompt := `Improve this recipe. Return a JSON object with the following fields:
- title: Improved or original title (string)
- description: Enhanced description that makes the recipe more appealing (string)
- markdown_content: Improved recipe content with better formatting and clarity (string)
- tips: Array of cooking tips or variations (array of strings, optional)

Recipe to enhance:
Title: ` + req.CurrentTitle + `
Content:
` + req.CurrentContent + `

Only return the JSON object, nothing else.`

	requestBody := map[string]interface{}{
		"model": s.cfg.AI.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream": false,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API error: %s", string(respBody))
	}

	var ollamaResponse struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.Unmarshal(respBody, &ollamaResponse); err != nil {
		return nil, err
	}

	var result EnhanceRecipeResponse
	if err := json.Unmarshal([]byte(ollamaResponse.Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Ollama response: %w", err)
	}

	return &result, nil
}
